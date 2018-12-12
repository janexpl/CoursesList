package controllers

import (
	"net/http"

	"github.com/janexpl/CoursesList/config"
	"github.com/janexpl/CoursesList/models"
)

type CertificatesController struct{}

func NewCertificatesController() *CertificatesController {
	return &CertificatesController{}
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
