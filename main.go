package main

import (
	"net/http"

	"github.com/janexpl/CoursesList/config"
	"github.com/janexpl/CoursesList/controllers"
	"github.com/janexpl/CoursesList/logging"
)

func main() {

	http.HandleFunc("/", index)
	courses := controllers.NewCoursesController()

	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	http.Handle("/node_modules/", http.StripPrefix("/node_modules", http.FileServer(http.Dir("./node_modules"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.HandleFunc("/courses/create/process", authorized(courses.CreateProcess))
	http.HandleFunc("/courses/create", authorized(courses.Create))
	http.HandleFunc("/courses/show", authorized(courses.Show))
	http.HandleFunc("/courses/delete/process", authorized(courses.DeleteProcess))
	http.HandleFunc("/courses", authorized(courses.Index))
	http.HandleFunc("/courses/json", authorized(courses.HandleJson))
	http.HandleFunc("/courses/createjson", authorized(courses.Createjson))
	http.HandleFunc("/courses/update/process", authorized(courses.UpdateProcess))
	http.HandleFunc("/courses/update", authorized(courses.Update))
	http.HandleFunc("/courses/getcoursejson", authorized(courses.GetCourseJson))
	students := controllers.NewStudentsController()
	http.HandleFunc("/students/create/process", authorized(students.CreateProcess))
	http.HandleFunc("/students/create", authorized(students.Create))
	http.HandleFunc("/students/show", authorized(students.Show))
	http.HandleFunc("/students/delete/process", authorized(students.DeleteProcess))
	http.HandleFunc("/students", authorized(students.Index))
	http.HandleFunc("/students/getalljson", authorized(students.GetAllJson))
	http.HandleFunc("/students/update/process", authorized(students.UpdateProcess))
	http.HandleFunc("/students/update", authorized(students.Update))
	http.HandleFunc("/students/addfrommodal", authorized(students.CreateFromModal))
	companies := controllers.NewCompaniesController()
	http.HandleFunc("/companies/create/process", authorized(companies.CreateProcess))
	http.HandleFunc("/companies/create", authorized(companies.Create))
	http.HandleFunc("/companies/show", authorized(companies.Show))
	http.HandleFunc("/companies/delete/process", authorized(companies.DeleteProcess))
	http.HandleFunc("/companies", authorized(companies.Index))
	http.HandleFunc("/companies/update/process", authorized(companies.UpdateProcess))
	http.HandleFunc("/companies/update", authorized(companies.Update))
	http.HandleFunc("/companies/getalljson", authorized(companies.GetAllJson))
	http.HandleFunc("/companies/addfrommodal", authorized(companies.CreateFromModal))
	certificates := controllers.NewCertificatesController()
	http.HandleFunc("/certificates", authorized(certificates.Index))
	http.HandleFunc("/certificates/create", authorized(certificates.Create))
	http.HandleFunc("/certificates/update", authorized(certificates.Update))
	http.HandleFunc("/certificates/create/process", authorized(certificates.CreateProcess))
	http.HandleFunc("/certificates/delete/process", authorized(certificates.DeleteProcess))
	http.HandleFunc("/certificates/print", authorized(certificates.Print))
	users := controllers.NewUsersController()
	http.HandleFunc("/users/signup", users.Signup)
	//http.HandleFunc("/users/login/process", users.Login)

	http.HandleFunc("/users/create", authorized(users.Create))
	http.HandleFunc("/users/login", users.LoginForm)
	http.HandleFunc("/users/login/process", users.LoginProcess)
	http.HandleFunc("/users/logout", authorized(users.Logout))
	http.HandleFunc("/users", authorized(users.Index))
	http.HandleFunc("/users/json", authorized(users.HandleJson))
	http.HandleFunc("/users/update", authorized(users.Update))
	http.HandleFunc("/users/update/process", authorized(users.UpdateProcess))
	registries := controllers.NewRegistriesController()

	http.HandleFunc("/registries/getlastnumber", registries.GetLastNumberWithSymbol)
	logging.Info.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)

}
func index(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/certificates", http.StatusSeeOther)

}

func authorized(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// code before
		if !config.AlreadyLoggedIn(w, r) {
			logging.Trace.Println("authorized")
			http.Redirect(w, r, "/users/login", http.StatusSeeOther)
			return
		}
		h.ServeHTTP(w, r)
		// code after
	})
}
