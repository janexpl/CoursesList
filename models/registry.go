package models

import (
	"errors"
	"net/http"
	"strconv"
	"time"

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
	if err != nil {
		if number == 0 {
			number = 0
		}

	}
	return number, nil

}
func checkDateIsValid(id int64, year int, number int, date time.Time) (bool, error) {
	type Result struct {
		date   time.Time
		number int
	}
	ress := []*Result{}

	rows, err := config.DB.Query(`
		SELECT
    	certificates. "date",
    	registries. "number"
	FROM
    	certificates
    	INNER JOIN registries ON certificates.registry_id = registries.id
    	INNER JOIN courses ON registries.course_id = courses.id
   		WHERE registries."year" = $1 and courses.id = $2
    ORDER BY
        "number"`, year, id)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		res := new(Result)
		err := rows.Scan(&res.date, &res.number)
		if err != nil {
			return false, err
		}
		ress = append(ress, res)
	}
	min := -1
	max := len(ress)
	for i := range ress {
		if ress[i].number <= number {
			if min <= i {
				min = i
			}
		}
		if ress[i].number >= number {
			if max >= i {
				max = i
			}
		}
	}

	if min == -1 {
		res := new(Result)
		res.number = number
		res.date = date
		ress = append([]*Result{res}, ress...)
		min++
	}
	if max == len(ress) {
		res := new(Result)
		res.number = number
		res.date = date
		ress = append(ress, res)
		//max++
	}

	mindate := ress[min].date.Unix()
	maxdate := ress[max].date.Unix()

	if date.Unix() < mindate || date.Unix() > maxdate {

		return false, err
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
func (reg *Registry) PutRegistry(r *http.Request) (int, error) {

	cr := Course{}
	cr, err := cr.GetCourseWithId(r)
	reg.Number, _ = strconv.Atoi(r.FormValue("number"))
	reg.Year, _ = strconv.Atoi(r.FormValue("year"))
	date, _ := time.Parse("2006-01-02", r.FormValue("certdate"))
	reg.Course = cr
	id := 0
	// sprawdzanie daty poprzedniego numeru
	ok, err := checkDateIsValid(reg.Course.ID, reg.Year, reg.Number, date)
	if err != nil {
		return 0, err
	}
	if ok {
		err = config.DB.QueryRow("INSERT INTO registries(course_id,year,number) VALUES ($1,$2,$3) RETURNING id", reg.Course.ID, reg.Year, reg.Number).Scan(&id)
		if err != nil {
			return id, errors.New("500. Internal Server Error." + err.Error())
		}
		return id, nil
	} else {
		return 0, err
	}

}

func (reg *Registry) DeleteRegistryWithId(id int64) error {

	_, err := config.DB.Exec("DELETE FROM registries WHERE id = $1", id)

	if err != nil {
		return errors.New("500. Internal Server Error")
	}
	return nil
}
