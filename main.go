package main

import (
	"net/http"

	"github.com/janexpl/CoursesList/config"
	"github.com/janexpl/CoursesList/controllers"
	"github.com/janexpl/CoursesList/logging"
)

func main() {

	courses := controllers.NewCoursesController()

	http.HandleFunc("/", index)
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.HandleFunc("/courses/create/process", courses.CreateProcess)
	http.HandleFunc("/courses/create", courses.Create)
	http.HandleFunc("/courses/show", courses.Show)
	http.HandleFunc("/courses/delete/process", courses.DeleteProcess)
	http.HandleFunc("/courses", courses.Index)
	http.HandleFunc("/courses/update/process", courses.UpdateProcess)
	http.HandleFunc("/courses/update", courses.Update)
	http.HandleFunc("/courses/getcoursejson", courses.GetCourseJson)
	students := controllers.NewStudentsController()
	http.HandleFunc("/students/create/process", students.CreateProcess)
	http.HandleFunc("/students/create", students.Create)
	http.HandleFunc("/students/show", students.Show)
	http.HandleFunc("/students/delete/process", students.DeleteProcess)
	http.HandleFunc("/students", students.Index)
	http.HandleFunc("/students/getalljson", students.GetAllJson)
	http.HandleFunc("/students/update/process", students.UpdateProcess)
	http.HandleFunc("/students/update", students.Update)
	companies := controllers.NewCompaniesController()
	http.HandleFunc("/companies/create/process", companies.CreateProcess)
	http.HandleFunc("/companies/create", companies.Create)
	http.HandleFunc("/companies/show", companies.Show)
	http.HandleFunc("/companies/delete/process", companies.DeleteProcess)
	http.HandleFunc("/companies", companies.Index)
	http.HandleFunc("/companies/update/process", companies.UpdateProcess)
	http.HandleFunc("/companies/update", companies.Update)
	http.HandleFunc("/companies/getalljson", companies.GetAllJson)
	http.HandleFunc("/companies/addfrommodal", companies.CreateFromModal)
	certificates := controllers.NewCertificatesController()
	http.HandleFunc("/certificates", certificates.Index)
	http.HandleFunc("/certificates/create", certificates.Create)
	http.HandleFunc("/certificates/create/process", certificates.CreateProcess)
	http.HandleFunc("/certificates/delete/process", certificates.DeleteProcess)

	/*

		http.HandleFunc("/certificates/show", certificates.Show)


		http.HandleFunc("/certificates/update/process", certificates.UpdateProcess)
		http.HandleFunc("/certificates/update", certificates.Update)
		http.HandleFunc("/certificates/getalljson", certificates.GetAllJson)*/
	registries := controllers.NewRegistriesController()
	//http.HandleFunc("/certificates", certificates.Index)
	http.HandleFunc("/registries/getlastnumber", registries.GetLastNumberWithSymbol)
	logging.Info.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)

}
func index(w http.ResponseWriter, r *http.Request) {
	config.RenderTemplate(w, r, "index/home", nil)

}
