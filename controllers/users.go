package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/janexpl/CoursesList/logging"

	"github.com/janexpl/CoursesList/config"
	"github.com/janexpl/CoursesList/models"
	"golang.org/x/crypto/bcrypt"
)

type UsersController struct{}

func NewUsersController() *UsersController {
	return &UsersController{}
}
func (u *UsersController) HandleJson(w http.ResponseWriter, r *http.Request) {
	us := models.User{}
	switch r.Method {
	case "GET":
		usrs, err := us.AllUsers()
		if err != nil {
			http.Error(w, http.StatusText(500)+err.Error(), http.StatusInternalServerError)
			return
		}
		uj, err := json.Marshal(usrs)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // 200
		fmt.Fprintf(w, "%s\n", uj)
	case "POST":
		fmt.Println("POST")

		us, err := us.PutUser(r)
		if err != nil {
			http.Error(w, http.StatusText(406), http.StatusNotAcceptable)
			return
		}
		fmt.Println(us)
		uj, err := json.Marshal(us)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // 201
		fmt.Fprintf(w, "%s\n", uj)
	case "DELETE":
		err := us.DeleteUser(r)
		if err != nil {
			http.Error(w, http.StatusText(406), http.StatusNotAcceptable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // 201
	case "PUT":
		err := us.UpdateUser(r)
		if err != nil {
			http.Error(w, http.StatusText(406), http.StatusNotAcceptable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // 201

	}
}
func (u *UsersController) Signup(w http.ResponseWriter, r *http.Request) {
	var flash string
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	user := models.User{}
	_, err := user.PutUser(r)
	if err != nil {
		flash = err.Error()
		config.RenderTemplate(w, r, "users/login", map[string]interface{}{

			"Flash": flash,
		})
		return
	} else {
		flash = "Uzytkownika dodano poprawnie."

	}
	data := map[string]interface{}{
		"Flash": flash,
	}

	config.RenderTemplate(w, r, "users/login", data)

}
func (u *UsersController) Create(w http.ResponseWriter, r *http.Request) {
	config.RenderTemplate(w, r, "users/create", nil)
}
func (u *UsersController) LoginForm(w http.ResponseWriter, r *http.Request) {

	if config.AlreadyLoggedIn(w, r) {

		http.Redirect(w, r, "/certificates", http.StatusSeeOther)
		return
	}
	config.RenderTemplate(w, r, "users/login", nil)
}

func (u *UsersController) LoginProcess(w http.ResponseWriter, r *http.Request) {

	if config.AlreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/certificates", http.StatusSeeOther)
		logging.Trace.Println("Sprawdzzam")
		return
	}
	//var flash string
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	user := models.User{}

	user, err := user.GetUser(r.FormValue("email"))
	p := r.FormValue("password")
	if err != nil {
		logging.Error.Println("Hasło albo uzytkownik błędne")
		config.SetFlash(w, r, []byte("Błędny uzytkownik lub hasło"))
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(p))
	if err != nil {
		logging.Error.Println("Hasło albo uzytkownik błędne")
		config.SetFlash(w, r, []byte("Bledny uzytkownik lub haslo"))
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	logging.Trace.Println("Logowanie")
	config.NewSession(w, r, user.Email, user.Role)
	http.Redirect(w, r, "/certificates", http.StatusSeeOther)
	//	fmt.Println(dbSessions)
	//	config.RenderTemplate(w, r, "users/login", nil)
}
func (u *UsersController) Logout(w http.ResponseWriter, r *http.Request) {

	config.DeleteSession(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (u *UsersController) Index(w http.ResponseWriter, r *http.Request) {
	var flash string
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	// us := models.User{}
	// uss, err := us.AllUsers()
	// if err != nil {
	// 	flash = err.Error()
	// 	return
	// }
	data := map[string]interface{}{
		"Flash": flash,
	}
	config.RenderTemplate(w, r, "users/users", data)
}

func (u *UsersController) Update(w http.ResponseWriter, r *http.Request) {
	var flash string
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	us := models.User{}
	us, err := us.OneUser(r)
	if err != nil {
		flash = err.Error()
	}
	data := map[string]interface{}{
		"Data":  us,
		"Flash": flash,
	}
	config.RenderTemplate(w, r, "users/update", data)
}

func (u *UsersController) UpdateProcess(w http.ResponseWriter, r *http.Request) {
	var flash string
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	us := models.User{}
	err := us.UpdateUser(r)
	if err != nil {

		flash = err.Error()

	} else {
		flash = "Zmiany zapisano poprawnie"

	}
	uss, err := us.AllUsers()
	if err != nil {
		flash = err.Error()
	}
	data := map[string]interface{}{
		"Data":  uss,
		"Flash": flash,
	}

	config.RenderTemplate(w, r, "users/users", data)
}
