package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/janexpl/CoursesList/config"
	"github.com/janexpl/CoursesList/logging"
	"github.com/janexpl/CoursesList/models"
)

type CertificatesController struct{}

func NewCertificatesController() *CertificatesController {
	return &CertificatesController{}
}

func (crt *CertificatesController) Print(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	c := models.Certificate{}
	cert, err := c.GetCertificate(r)
	if err != nil {
		logging.Error.Println(err.Error())
	}
	bck := []models.CourseProgram{}
	err = json.Unmarshal(cert.Registry.Course.CourseProgram, &bck)
	front := parseHtml(cert)
	front = `<!doctype html><html><head><meta charset="utf-8"><title>ZAŚWIADCZENIE</title></head><body>` + front
	back := `
	<p style="page-break-before: always;">
	<table> <thead>
	<tr>
	  <th>Lp.</th>
	  <th>Temat szkolenia</th>
	  <th>Liczba godzin zajęć teoretycznych (wykładów)</th>
	  <th>Liczba godzin  zajęć praktycznych (ćwiczeń)</th>
	</tr>
  </thead>`
	var i int = 1
	for _, row := range bck {
		rw := fmt.Sprintf("<tr><td>%v</td><td>%v</td><td class='hour'>%v</td><td class='hour'>%v</td></tr>", i, row.Subject, row.TheoryTime, row.PracticeTime)
		back = back + rw
		i++
		fmt.Println(row)
	}
	back = back + `</table></body><style>
	table {
		border-collapse: collapse;
	  }
	  .hour {
		  text-align: center;
	  }
	  table, th, td {
		padding: 15px;
		border: 1px solid black;

	  }
	  h2{
	  font-size: 50px;	
		}
	p {
	font-size: 20px;
		}</style></html>`
	fmt.Println(back)

	page := front + back
	//Printing certificate

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.Grayscale.Set(false)

	// Create a new input page from an URL
	page1 := wkhtmltopdf.NewPageReader(strings.NewReader(page))

	// Set options for this page
	//page1.FooterRight.Set("[page]")
	//page1.FooterFontSize.Set(10)
	page1.Zoom.Set(0.95)

	//page2.FooterRight.Set("[page]")
	//page2.FooterFontSize.Set(10)

	// Add to document
	pdfg.AddPage(page1)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	// Write buffer contents to file on disk
	//err = pdfg.WriteFile("./simplesample.pdf")
	buf := pdfg.Buffer()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
	// Output: Done

	buf.WriteTo(w)
}

func (crt *CertificatesController) Index(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	c := models.Certificate{}
	certificates, err := c.AllCertificates()
	if err != nil {
		http.Error(w, http.StatusText(500)+err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Data": certificates,
	}
	config.RenderTemplate(w, r, "certificates/certificates", data)
}

func (crt *CertificatesController) Create(w http.ResponseWriter, r *http.Request) {
	var flash string
	crs := models.Course{}
	courses, err := crs.AllCourses()
	if err != nil {
		flash = err.Error()
	}
	flash = ""
	data := map[string]interface{}{
		"Course": courses,
		"Flash":  flash,
	}
	config.RenderTemplate(w, r, "certificates/create", data)
}

func (crt *CertificatesController) Update(w http.ResponseWriter, r *http.Request) {
	var flash string

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	cp := models.Certificate{}
	cert, err := cp.GetCertificate(r)
	crs := models.Course{}
	courses, err := crs.AllCourses()
	if err != nil {
		flash = err.Error()
	}
	flash = ""

	data := map[string]interface{}{
		"Course": courses,
		"Data":   cert,
		"Flash":  flash,
	}

	config.RenderTemplate(w, r, "certificates/update", data)
}
func (crt *CertificatesController) CreateProcess(w http.ResponseWriter, r *http.Request) {
	var flash string
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	cr := models.Certificate{}
	err := cr.PutCertificate(r)
	if err != nil {

		flash = err.Error()
		config.RenderTemplate(w, r, "certificates/create", map[string]interface{}{
			"Data":  nil,
			"Flash": flash,
		})
		return
	} else {
		flash = "Zapisano kursanta poprawnie"

	}
	crs, err := cr.AllCertificates()
	if err != nil {
		flash = err.Error()
	}
	data := map[string]interface{}{
		"Data":  crs,
		"Flash": flash,
	}

	config.RenderTemplate(w, r, "certificates/certificates", data)
}

func (crt *CertificatesController) DeleteProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	cr := models.Certificate{}
	err := cr.DeleteCertificate(r)
	if err != nil {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	crs, err := cr.AllCertificates()
	if err != nil {
		http.Error(w, http.StatusText(500)+err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Data":  crs,
		"Flash": "Dane usunięte poprawnie.",
	}
	config.RenderTemplate(w, r, "certificates/certificates", data)
}
func parseHtml(c models.Certificate) string {
	htmls := c.Registry.Course.CertFrontpage
	tags := make(map[string]string)
	tagRe, _ := regexp.Compile(`(?U){{.*}}`)
	capList := tagRe.FindAllStringIndex(htmls, -1)
	for _, cap := range capList {
		tag := htmls[cap[0]+2 : cap[1]-2]
		tag = strings.Replace(tag, " ", "", -1)

		switch tag {
		case "imie":
			tags[htmls[cap[0]:cap[1]]] = c.Student.Firstname
		case "drugie_imie":
			tags[htmls[cap[0]:cap[1]]] = c.Student.Secondname
		case "nazwisko":
			tags[htmls[cap[0]:cap[1]]] = c.Student.Lastname
		case "data_urodzenia":
			tags[htmls[cap[0]:cap[1]]] = (c.Student.Birthdate).Format("02.01.2006")
		case "miejsce_urodzenia":
			tags[htmls[cap[0]:cap[1]]] = c.Student.Birthplace
		case "nazwa_kursu":
			tags[htmls[cap[0]:cap[1]]] = c.Registry.Course.Name
		case "data_rozpoczecia":
			tags[htmls[cap[0]:cap[1]]] = (c.CourseDateStart).Format("02.01.2006")
		case "data_zakonczenia":
			tags[htmls[cap[0]:cap[1]]] = (c.CourseDateEnd).Format("02.01.2006")
		case "data_wystawienia":
			tags[htmls[cap[0]:cap[1]]] = (c.Date).Format("02.01.2006")
		case "numer_zaswiadczenia":
			tags[htmls[cap[0]:cap[1]]] = fmt.Sprintf("%v/%v/%v", c.Registry.Number, c.Registry.Course.Symbol, c.Registry.Year)
		default:
			tags[htmls[cap[0]:cap[1]]] = ""
		}

	}
	for key, tag := range tags {
		htmls = strings.Replace(htmls, key, tag, -1)
	}
	return htmls
}
