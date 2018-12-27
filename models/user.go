package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/janexpl/CoursesList/config"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Password  []byte `json:"password"`
	SPassword string `json:"spassword"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Role      int    `json:"role"`
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
func (u *User) PutUser(r *http.Request) (User, error) {
	if strings.Contains(r.Header.Get("Content-Type"), "json") {
		json.NewDecoder(r.Body).Decode(&u)
		bpas, err := bcrypt.GenerateFromPassword([]byte(u.SPassword), bcrypt.MinCost)
		if err != nil {
			return *u, err
		}
		u.Password = bpas
		
	} else {
		u.Email = r.FormValue("semail")
		u.Firstname = r.FormValue("sfirstname")
		u.Lastname = r.FormValue("slastname")
		if r.FormValue("srole") != "" {
			u.Role = 1
		} else {
			u.Role = 0
		}
		bpas, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("spassword")), bcrypt.MinCost)
		if err != nil {
			return *u, err
		}
		u.Password = bpas
	}
	var email string
	err := config.DB.QueryRow("SELECT email FROM users WHERE email = $1", u.Email).Scan(&email)
	if email != "" {
		return *u, errors.New("Istnieje juz taki uzytkownik")
	}
	err = config.DB.QueryRow("INSERT INTO users(email, firstname, lastname, password, role) VALUES ($1, $2, $3, $4, $5) RETURNING id", u.Email, u.Firstname, u.Lastname, u.Password, u.Role).Scan(&u.ID)
	fmt.Println(*u)
	if err != nil {
		return *u, err
	}

	return *u, nil
}

func (u *User) GetUser(email string) (User, error) {

	err := config.DB.QueryRow("SELECT * FROM users WHERE email = $1", email).Scan(&u.ID, &u.Email, &u.Password, &u.Firstname, &u.Lastname, &u.Role)

	if err != nil {
		return *u, err
	}
	return *u, nil

}
func (u *User) DeleteUser(r *http.Request) error {
	id := r.FormValue("id")
	_, err := config.DB.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}
func (u *User) UpdateUser(r *http.Request) error {
	var password string
	if strings.Contains(r.Header.Get("Content-Type"), "json") {

		json.NewDecoder(r.Body).Decode(&u)
		cost, _ := bcrypt.Cost([]byte(u.SPassword))
		fmt.Println(cost)
		fmt.Println(*u)
		if cost == 0 {
			fmt.Println("Zmiana hasła")
			bpas, err := bcrypt.GenerateFromPassword([]byte(u.SPassword), bcrypt.MinCost)
			if err != nil {
				return err
			}
			u.Password = bpas
		}
		fmt.Println(*u)
	} else {
		u.ID, _ = strconv.ParseInt(r.FormValue("id"), 0, 64)
		u.Email = r.FormValue("semail")
		u.Firstname = r.FormValue("sfirstname")
		u.Lastname = r.FormValue("slastname")
		if r.FormValue("srole") != "" {
			u.Role = 1
		} else {
			u.Role = 0
		}
		password = r.FormValue("spassword")
		cost, _ := bcrypt.Cost([]byte(password))
		fmt.Println(cost)
		if cost == 0 {
			bpas, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				return err
			}
			u.Password = bpas
		}
	}
	//var email string
	// err := config.DB.QueryRow("SELECT email FROM users WHERE email = $1 AND id <> $2", u.Email, u.ID).Scan(&email)
	// if email != "" {
	// 	return errors.New("Istnieje juz taki uzytkownik")
	// }
	fmt.Println(*u)
	_, err := config.DB.Exec("UPDATE users SET email=$1, password=$2, firstname=$3, lastname=$4, role=$5 WHERE id=$6", u.Email, u.Password, u.Firstname, u.Lastname, u.Role, u.ID)
	if err != nil {
		return err
	}
	return nil
}
