package config

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type Layout string

var layout *template.Template

const (
	trimPrefix        = "templates/"
	Home       Layout = "home"
	Ajax       Layout = "ajax"
	User       Layout = "user"
)

var (
	layoutFuncs = template.FuncMap{
		"yield": func() (string, error) {
			return "", fmt.Errorf("yield called unexpectedly")
		},
		"dateformat":  dateformat,
		"formdatefmt": formdatefmt,
	}
)

func dateformat(t time.Time) string {
	return t.Format("02.01.2006")
}
func formdatefmt(t time.Time) string {

	return t.Format("2006-01-02")
}

var errorTemplate = `
<html>
	<body>
		<h1>Error rendering template %s</h1>
		<p>%s</p>
	</body>
</html>
`
var templates *template.Template

func init() {
	layout = template.Must(
		template.
			New("layout.gohtml").
			Funcs(layoutFuncs).
			ParseFiles("templates/layout.gohtml", "templates/header.gohtml"))

	templates = template.Must(
		template.
			New("t").Funcs(layoutFuncs).
			ParseGlob("templates/*/*.gohtml"))

}

func RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	un, err := GetLoggedUser(r)
	
	data["CurrentUser"] = un
	if data["Flash"] == "" {
		data["Flash"], _ = getFlash(w, r)
	}

	templates = template.Must(
		template.
			New("t").Funcs(layoutFuncs).
			ParseGlob("templates/*/*.gohtml"))

	funcs := template.FuncMap{
		"yield": func() (template.HTML, error) {
			buf := bytes.NewBuffer(nil)
			err := templates.ExecuteTemplate(buf, name, data)
			return template.HTML(buf.String()), err
		},
		"dateformat": dateformat,
	}

	layoutClone, _ := layout.Clone()

	layoutClone.Funcs(funcs)

	err = layoutClone.Execute(w, data)

	if err != nil {
		http.Error(
			w,
			fmt.Sprintf(errorTemplate, name, err),
			http.StatusInternalServerError,
		)
	}

}
