package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/janexpl/CoursesList/config"
	"github.com/janexpl/CoursesList/models"
)

type StudentsController struct{}

func NewStudentsController() *StudentsController {
	return &StudentsController{}
}
func (st *StudentsController) Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	s := models.Student{}
	students, err := s.AllStudents()
	fmt.Println(students)
	if err != nil {
		http.Error(w, http.StatusText(500)+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Data": students,
	}
	config.RenderTemplate(w, r, "students/students", data)
}

func (st *StudentsController) Create(w http.ResponseWriter, r *http.Request) {
	var flash string
	cp := models.Company{}
	companies, err := cp.AllCompanies()
	if err != nil {
		flash = err.Error()
	}
	flash = ""
	data := map[string]interface{}{
		"Data":  companies,
		"Flash": flash,
	}
	config.RenderTemplate(w, r, "students/create", data)
}

func (st *StudentsController) GetAllJson(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	std := models.Student{}
	students, err := std.AllStudents()

	if err != nil {
		http.Error(w, http.StatusText(500)+err.Error(), http.StatusInternalServerError)
		return
	}
	uj, err := json.Marshal(students)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(students)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	fmt.Fprintf(w, "%s\n", uj)

}

func (st *StudentsController) CreateProcess(w http.ResponseWriter, r *http.Request) {
	var flash string
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	s := models.Student{}
	err := s.PutStudent(r)
	if err != nil {

		flash = err.Error()
		config.RenderTemplate(w, r, "students/create", map[string]interface{}{
			"Data":  nil,
			"Flash": flash,
		})
		return
	} else {
		flash = "Zapisano kursanta poprawnie"

	}
	students, err := s.AllStudents()
	if err != nil {
		flash = err.Error()
	}
	data := map[string]interface{}{
		"Data":  students,
		"Flash": flash,
	}

	config.RenderTemplate(w, r, "students/students", data)
}

func (st *StudentsController) Show(w http.ResponseWriter, r *http.Request) {
	var flash string

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	s := models.Student{}
	student, err := s.OneStudent(r)

	certificate := models.Certificate{}
	crt, err := certificate.AllCertificatesWithStudent(student.ID)

	if err != nil {
		flash = err.Error()
	}
	flash = ""

	data := map[string]interface{}{
		"Data":  student,
		"Cert":  crt,
		"Flash": flash,
	}

	config.RenderTemplate(w, r, "students/show", data)

}

func (st *StudentsController) DeleteProcess(w http.ResponseWriter, r *http.Request) {
	var flash string
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	s := models.Student{}
	err := s.DeleteStudent(r)
	if err != nil {
		flash = err.Error()
	}
	students, err := s.AllStudents()
	if err != nil {
		flash = err.Error()
	}
	flash = "Kursant usunięty pomyślnie"
	data := map[string]interface{}{
		"Data":  students,
		"Flash": flash,
	}
	config.RenderTemplate(w, r, "students/students", data)
}

func (st *StudentsController) Update(w http.ResponseWriter, r *http.Request) {
	var flash string

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	s := models.Student{}
	student, err := s.OneStudent(r)

	if err != nil {
		flash = err.Error()
	}
	flash = ""

	data := map[string]interface{}{
		"Data":  student,
		"Flash": flash,
	}

	config.RenderTemplate(w, r, "students/update", data)
}

func (st *StudentsController) UpdateProcess(w http.ResponseWriter, r *http.Request) {
	var flash string
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	s := models.Student{}
	err := s.UpdateStudent(r)
	if err != nil {

		flash = err.Error()

	} else {
		flash = "Kursanta zapisano poprawnie"

	}
	students, err := s.AllStudents()
	if err != nil {
		flash = err.Error()
	}
	data := map[string]interface{}{
		"Data":  students,
		"Flash": flash,
	}

	config.RenderTemplate(w, r, "students/students", data)
}
