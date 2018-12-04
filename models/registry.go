package models

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/janexpl/CoursesList/config"
)

type Registry struct {
	ID     int64
	Course Course
	Year   int
	Number int
}

func (reg *Registry) GetLastNumber(r *http.Request) (int, error) {
	var number int
	id := r.FormValue("id")
	year := r.FormValue("year")
	err := config.DB.QueryRow("SELECT MAX(number) FROM registries WHERE course_id = $1 and year = $2", id, year).Scan(&number)
	fmt.Println(number)
	if err != nil {
		if number == 0 {
			number = 0
		}

	}
	return number, nil

}

func (reg *Registry) PutRegistry(r *http.Request) (int, error) {
	rg := Registry{}
	cr := Course{}
	cr, err := cr.GetCourseWithId(r)
	rg.Number, _ = strconv.Atoi(r.FormValue("number"))
	rg.Year, _ = strconv.Atoi(r.FormValue("year"))
	rg.Course = cr
	id := 0
	err = config.DB.QueryRow("INSERT INTO registries(course_id,year,number) VALUES ($1,$2,$3) RETURNING id", rg.Course.ID, rg.Year, rg.Number).Scan(&id)

	fmt.Println(id)
	if err != nil {
		return id, errors.New("500. Internal Server Error." + err.Error())
	}
	return id, nil

}

func (reg *Registry) DeleteRegistryWithId(id int64) error {

	_, err := config.DB.Exec("DELETE FROM registries WHERE id = $1", id)

	if err != nil {
		return errors.New("500. Internal Server Error")
	}
	return nil
}
