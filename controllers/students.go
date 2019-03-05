package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/janexpl/CoursesList/config"
	"github.com/janexpl/CoursesList/logging"
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
		logging.Error.Println(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	fmt.Fprintf(w, "%s\n", uj)

}

func (st *StudentsController) CreateProcess(w http.ResponseWriter, r *http.Request) {
	//var flash string
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	s := models.Student{}
	_, err := s.PutStudent(r)
	if err != nil {
		config.SetFlash(w, r, []byte(err.Error()))
		http.Redirect(w, r, "/students/create", http.StatusSeeOther)
		// config.RenderTemplate(w, r, "students/create", map[string]interface{}{
		// 	"Data":  nil,
		// 	"Flash": flash,
		// })
		return
	} else {
		config.SetFlash(w, r, []byte("Kursanta zapisano poprawnie"))

	}
	// students, err := s.AllStudents()
	// if err != nil {
	// 	flash = err.Error()
	// }
	// data := map[string]interface{}{
	// 	"Data":  students,
	// 	"Flash": flash,
	// }
	http.Redirect(w, r, "/students", http.StatusSeeOther)
	//config.RenderTemplate(w, r, "students/students", data)
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
	//var flash string
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	s := models.Student{}
	err := s.DeleteStudent(r)
	if err != nil {
		config.SetFlash(w, r, []byte(err.Error()))
		//flash = err.Error()
	}
	// students, err := s.AllStudents()
	// if err != nil {
	// 	flash = err.Error()
	// }
	// flash = "Kursant usunięty pomyślnie"
	// data := map[string]interface{}{
	// 	"Data":  students,
	// 	"Flash": flash,
	// }
	config.SetFlash(w, r, []byte("Kursant kursant usunięty poprawnie"))
	http.Redirect(w, r, "/students", http.StatusSeeOther)
	//config.RenderTemplate(w, r, "students/students", data)
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
func (st *StudentsController) CreateFromModal(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	std := models.Student{}

	sid, err := std.PutStudent(r)
	if err != nil {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	uj, err := json.Marshal(sid)
	if err != nil {
		logging.Error.Println(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	fmt.Fprintf(w, "%s\n", uj)

}
func (st *StudentsController) UpdateProcess(w http.ResponseWriter, r *http.Request) {
	//var flash string
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	s := models.Student{}
	err := s.UpdateStudent(r)
	if err != nil {
		config.SetFlash(w, r, []byte(err.Error()))
		http.Redirect(w, r, "/students/update", http.StatusSeeOther)
	} else {
		config.SetFlash(w, r, []byte("Kursanta zapisano poprawnie"))
	}
	// students, err := s.AllStudents()
	// if err != nil {
	// 	flash = err.Error()
	// }
	// data := map[string]interface{}{
	// 	"Data":  students,
	// 	"Flash": flash,
	// }
	http.Redirect(w, r, "students/students", http.StatusSeeOther)
	//config.RenderTemplate(w, r, "students/students", data)
}
