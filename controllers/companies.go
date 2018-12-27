package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/janexpl/CoursesList/config"
	"github.com/janexpl/CoursesList/logging"
	"github.com/janexpl/CoursesList/models"
)

type CompaniesController struct{}

func NewCompaniesController() *CompaniesController {
	return &CompaniesController{}
}
func (c *CompaniesController) GetAllJson(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	cp := models.Company{}
	companies, err := cp.AllCompanies()

	if err != nil {
		http.Error(w, http.StatusText(500)+err.Error(), http.StatusInternalServerError)
		return
	}
	uj, err := json.Marshal(companies)
	if err != nil {
		logging.Error.Println(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	fmt.Fprintf(w, "%s\n", uj)

}
func (c *CompaniesController) Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	cp := models.Company{}
	companies, err := cp.AllCompanies()

	if err != nil {
		http.Error(w, http.StatusText(500)+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Data": companies,
	}
	config.RenderTemplate(w, r, "companies/companies", data)
}

func (c *CompaniesController) Create(w http.ResponseWriter, r *http.Request) {
	config.RenderTemplate(w, r, "companies/create", nil)
}

func (c *CompaniesController) CreateFromModal(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	cp := models.Company{}

	id, err := cp.PutCompany(r)
	if err != nil {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}
	uj, err := json.Marshal(id)
	if err != nil {
		logging.Error.Println(err.Error())

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	fmt.Fprintf(w, "%s\n", uj)
	//http.Redirect(w, r, "#", http.StatusMovedPermanently)

}

func (c *CompaniesController) CreateProcess(w http.ResponseWriter, r *http.Request) {
	var flash string
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	cp := models.Company{}
	_, err := cp.PutCompany(r)
	if err != nil {
		flash = err.Error()
		config.RenderTemplate(w, r, "companies/create", map[string]interface{}{
			"Data":  nil,
			"Flash": flash,
		})
		return
	} else {
		flash = "Dane zapisano poprawnie"

	}
	companies, err := cp.AllCompanies()
	if err != nil {
		flash = err.Error()
	}
	data := map[string]interface{}{
		"Data":  companies,
		"Flash": flash,
	}

	config.RenderTemplate(w, r, "companies/companies", data)
}

func (c *CompaniesController) Show(w http.ResponseWriter, r *http.Request) {
	var flash string

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	cp := models.Company{}
	company, err := cp.OneCompany(r)
	st := models.Student{}

	students, err := st.AllStudentsWithCompany(company)

	if err != nil {
		flash = err.Error()
	}
	flash = ""

	data := map[string]interface{}{
		"Data":  company,
		"Data1": students,
		"Flash": flash,
	}

	config.RenderTemplate(w, r, "companies/show", data)

}

func (c *CompaniesController) DeleteProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	cp := models.Company{}
	err := cp.DeleteCompany(r)
	if err != nil {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}
	companies, err := cp.AllCompanies()
	if err != nil {
		http.Error(w, http.StatusText(500)+err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Data":  companies,
		"Flash": "Dane usunięte poprawnie.",
	}
	config.RenderTemplate(w, r, "companies/companies", data)
}

func (c *CompaniesController) Update(w http.ResponseWriter, r *http.Request) {
	var flash string

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	cp := models.Company{}
	company, err := cp.OneCompany(r)

	if err != nil {
		flash = err.Error()
	}
	flash = ""

	data := map[string]interface{}{
		"Data":  company,
		"Flash": flash,
	}

	config.RenderTemplate(w, r, "companies/update", data)
}

func (c *CompaniesController) UpdateProcess(w http.ResponseWriter, r *http.Request) {
	var flash string
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	cp := models.Company{}
	err := cp.UpdateCompany(r)
	if err != nil {

		flash = err.Error()

	} else {
		flash = "Dane zapisano poprawnie"

	}
	companies, err := cp.AllCompanies()
	if err != nil {
		flash = err.Error()
	}
	data := map[string]interface{}{
		"Data":  companies,
		"Flash": flash,
	}

	config.RenderTemplate(w, r, "companies/companies", data)
}
