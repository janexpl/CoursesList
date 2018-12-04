package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/janexpl/CoursesList/config"
	"github.com/janexpl/CoursesList/models"
)

type CoursesController struct{}

func NewCoursesController() *CoursesController {
	return &CoursesController{}
}
func (cr *CoursesController) Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	course := models.Course{}

	crs, err := course.AllCourses()
	if err != nil {
		http.Error(w, http.StatusText(500)+err.Error(), http.StatusInternalServerError)
		return
	}
	for _, z := range crs {
		fmt.Println(z.Name)

	}

	data := map[string]interface{}{
		"Data": crs,
	}
	config.RenderTemplate(w, r, "courses/courses", data)
}

func (cr *CoursesController) Create(w http.ResponseWriter, r *http.Request) {
	config.RenderTemplate(w, r, "courses/create", nil)
}
func (cr *CoursesController) GetCourseJson(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	crs := models.Course{}
	course, err := crs.GetCourseWithId(r)

	if err != nil {
		http.Error(w, http.StatusText(500)+err.Error(), http.StatusInternalServerError)
		return
	}
	uj, err := json.Marshal(course)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(course)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	fmt.Fprintf(w, "%s\n", uj)

}
func (cr *CoursesController) CreateProcess(w http.ResponseWriter, r *http.Request) {
	var flash string
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	course := models.Course{}
	err := course.PutCourse(r)
	if err != nil {

		flash = err.Error()
		config.RenderTemplate(w, r, "courses/create", map[string]interface{}{
			"Data":  nil,
			"Flash": flash,
		})
		return

	} else {
		flash = "Kurs zapisano poprawnie"

	}
	crs, err := course.AllCourses()
	if err != nil {
		flash = err.Error()
	}
	data := map[string]interface{}{
		"Data":  crs,
		"Flash": flash,
	}

	config.RenderTemplate(w, r, "courses/courses", data)
}

func (cr *CoursesController) Show(w http.ResponseWriter, r *http.Request) {
	var flash string

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	crs := models.Course{}
	crs, err := crs.OneCourse(r)

	cert := models.Certificate{}
	certs, err := cert.GetCertificatesWithCourseId(crs.ID)

	if err != nil {
		flash = err.Error()
	}
	flash = ""

	data := map[string]interface{}{
		"Course": crs,
		"Cert":   certs,
		"Flash":  flash,
	}

	config.RenderTemplate(w, r, "courses/show", data)

}

func (cr *CoursesController) DeleteProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	course := models.Course{}
	err := course.DeleteCourse(r)
	if err != nil {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}
	crs, err := course.AllCourses()
	if err != nil {
		http.Error(w, http.StatusText(500)+err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Data":  crs,
		"Flash": "Kurs usunięty pomyślnie.",
	}
	config.RenderTemplate(w, r, "courses/courses", data)
}

func (cr *CoursesController) Update(w http.ResponseWriter, r *http.Request) {
	var flash string

	fmt.Println("Update")
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	course := models.Course{}
	crs, err := course.OneCourse(r)

	if err != nil {
		flash = err.Error()
	}
	flash = ""
	data := map[string]interface{}{
		"Data":  crs,
		"Flash": flash,
	}
	fmt.Println(crs)
	config.RenderTemplate(w, r, "courses/update", data)
}

func (cr *CoursesController) UpdateProcess(w http.ResponseWriter, r *http.Request) {
	var flash string
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	course := models.Course{}
	err := course.UpdateCourse(r)
	if err != nil {

		flash = err.Error()

	} else {
		flash = "Kurs zapisano poprawnie"

	}
	crs, err := course.AllCourses()
	if err != nil {
		flash = err.Error()
	}
	data := map[string]interface{}{
		"Data":  crs,
		"Flash": flash,
	}

	config.RenderTemplate(w, r, "courses/courses", data)
}
