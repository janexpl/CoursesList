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
	for _, t := range layout.Templates() {
		fmt.Println(t.Name())
	}
	templates = template.Must(
		template.
			New("t").Funcs(layoutFuncs).
			ParseGlob("templates/*/*.gohtml"))
	//	lms, err := filepath.Glob("templates/*.gohtml")
	/* if err != nil {
		log.Fatal(err)
	}
	for _, lm := range lms {
		b, err := ioutil.ReadFile(lm)
		if err != nil {
			log.Fatal(err)
		}
		name := strings.TrimPrefix(lm, trimPrefix)
		template.Must(LYT.New(name).Parse(string(b)))
	}*/

}

func RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}

	//data["CurrentUser"] = RequestUser(r)
	//data["Flash"] = r.URL.Query().Get("flash")
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

	err := layoutClone.Execute(w, data)

	if err != nil {
		http.Error(
			w,
			fmt.Sprintf(errorTemplate, name, err),
			http.StatusInternalServerError,
		)
	}

}
