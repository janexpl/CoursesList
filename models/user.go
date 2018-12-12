package models

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/janexpl/CoursesList/config"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64
	Email     string
	Password  []byte
	Firstname string
	Lastname  string
	Role      int
}

func (u *User) AllUsers() ([]User, error) {
	rows, err := config.DB.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	us := []User{}
	for rows.Next() {
		err = rows.Scan(&u.ID, &u.Email, &u.Password, &u.Firstname, &u.Lastname, &u.Role)
		if err != nil {
			return nil, err
		}
		us = append(us, *u)
	}
	return us, nil
}

func (u *User) OneUser(r *http.Request) (User, error) {
	id := r.FormValue("id")
	err := config.DB.QueryRow("SELECT * FROM users WHERE id=$1", id).Scan(&u.ID, &u.Email, &u.Password, &u.Firstname, &u.Lastname, &u.Role)
	if err != nil {
		return *u, err
	}
	return *u, nil
}
func (u *User) PutUser(r *http.Request) error {

	u.Email = r.FormValue("semail")
	u.Firstname = r.FormValue("sfirstname")
	u.Lastname = r.FormValue("slastname")

	u.Role = 0

	bpas, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("spassword")), bcrypt.MinCost)
	if err != nil {
		return err
	}
	u.Password = bpas
	var email string
	err = config.DB.QueryRow("SELECT email FROM users WHERE email = $1", u.Email).Scan(&email)
	if email != "" {
		return errors.New("Istnieje juz taki uzytkownik")
	}
	_, err = config.DB.Exec("INSERT INTO users(email, firstname, lastname, password, role) VALUES ($1, $2, $3, $4, $5)", u.Email, u.Firstname, u.Lastname, u.Password, u.Role)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetUser(email string) (User, error) {

	err := config.DB.QueryRow("SELECT * FROM users WHERE email = $1", email).Scan(&u.ID, &u.Email, &u.Password, &u.Firstname, &u.Lastname, &u.Role)

	if err != nil {
		return *u, err
	}
	return *u, nil

}

func (u *User) UpdateUser(r *http.Request) error {
	var query string
	id := r.FormValue("id")
	u.Email = r.FormValue("semail")
	u.Firstname = r.FormValue("sfirstname")
	u.Lastname = r.FormValue("slastname")
	if r.FormValue("srole") != "" {
		u.Role = 1
	} else {
		u.Role = 0
	}
	password := r.FormValue("spassword")
	if password != "" {
		bpas, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			return err
		}
		u.Password = bpas
	}
	fmt.Println(query)
	var email string
	err := config.DB.QueryRow("SELECT email FROM users WHERE email = $1 AND id <> $2", u.Email, id).Scan(&email)
	if email != "" {
		return errors.New("Istnieje juz taki uzytkownik")
	}
	if password != "" {
		_, err = config.DB.Exec("UPDATE users SET email=$1, password=$2, firstname=$3, lastname=$4, role=$5 WHERE id=$6", u.Email, u.Password, u.Firstname, u.Lastname, u.Role, id)
		if err != nil {
			return err
		}
	} else {
		_, err = config.DB.Exec("UPDATE users SET email=$1, firstname=$3, lastname=$4, role=$5 WHERE id=$6", u.Email, u.Firstname, u.Lastname, u.Role, id)
		if err != nil {
			return err
		}
	}

	return nil
}
