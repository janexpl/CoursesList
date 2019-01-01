package models

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/janexpl/CoursesList/logging"

	"github.com/janexpl/CoursesList/config"
)

type Course struct {
	ID            int64
	Name          string
	Symbol        string
	ExpiryTime    string
	CourseProgram json.RawMessage
	CertFrontpage string
}

type CourseProgram struct {
	Subject      string
	TheoryTime   string
	PracticeTime string
}

func (cr *Course) PutCourseJson(r *http.Request) error {
	json.NewDecoder(r.Body).Decode(&cr)
	_, err := config.DB.Exec("INSERT INTO courses(name,symbol,expirytime,courseprogram,certfrontpage) VALUES ($1,$2,$3,$4,$5)", cr.Name, cr.Symbol, cr.ExpiryTime, cr.CourseProgram, cr.CertFrontpage)
	if err != nil {
		errors.New("500. Internal Server Error." + err.Error())
	}
	return nil
}
func (cr *Course) PutCourse(r *http.Request) error {
	smb := r.FormValue("courseSymbol")
	course := Course{}
	err := config.DB.QueryRow("SELECT symbol FROM courses WHERE symbol=$1", smb).Scan(&course.Symbol)
	if course.Symbol != "" {
		return errors.New("Istnieje juz taki kurs.")
	}
	courseprograms := []CourseProgram{
		{
			Subject:      "Ogólne wiadomości z zakresu BHP",
			TheoryTime:   "3",
			PracticeTime: "0"},
		{
			Subject:      "Nowy temat",
			TheoryTime:   "3",
			PracticeTime: "1"},
	}
	bs, err := json.Marshal(courseprograms)
	if err != nil {
		logging.Error.Println(err.Error())
	}
	course.CourseProgram = bs
	course.Name = r.FormValue("courseName")
	course.Symbol = smb
	course.ExpiryTime = r.FormValue("courseExpT")
	if course.Symbol == "" || course.Name == "" {
		return errors.New("Proszę wypełnić wskazane pola")
	}
	_, err = config.DB.Exec("INSERT INTO courses(name,symbol,expirytime,courseprogram) VALUES ($1,$2,$3,$4)", course.Name, course.Symbol, course.ExpiryTime, course.CourseProgram)
	if err != nil {
		logging.Error.Println(err.Error())
		return errors.New("500. Internal Server Error." + err.Error())
	}
	return nil
}

func (cr *Course) AllCourses() ([]Course, error) {
	rows, err := config.DB.Query("SELECT id,name,symbol,expirytime,courseprogram FROM courses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	crs := []Course{}
	for rows.Next() {
		cr := Course{}
		err := rows.Scan(&cr.ID, &cr.Name, &cr.Symbol, &cr.ExpiryTime, &cr.CourseProgram)
		if err != nil {
			return nil, err
		}
		crs = append(crs, cr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return crs, nil
}
func (cr *Course) GetCourseWithId(r *http.Request) (Course, error) {
	id := r.FormValue("courseid")
	crs := Course{}
	err := config.DB.QueryRow("SELECT id,name,symbol,expirytime,courseprogram FROM courses WHERE id = $1", id).
		Scan(&crs.ID, &crs.Name, &crs.Symbol, &crs.ExpiryTime, &crs.CourseProgram)
	if err != nil {
		return crs, err
	}
	return crs, nil
}

func (cr *Course) DeleteCourse(r *http.Request) error {
	smb := r.FormValue("symbol")

	if smb == "" {
		return errors.New("400. Bad Request.")
	}

	_, err := config.DB.Exec("DELETE FROM courses WHERE symbol = $1", smb)
	if err != nil {
		return errors.New("500. Internal Server Error")
	}
	return nil
}

func (cr *Course) OneCourse(r *http.Request) (Course, error) {

	crs := Course{}
	smb := r.FormValue("symbol")
	if smb == "" {
		return crs, errors.New("400. Bad Request.")
	}

	row := config.DB.QueryRow("SELECT * FROM courses WHERE symbol = $1", smb)
	err := row.Scan(&crs.ID, &crs.Name, &crs.Symbol, &crs.ExpiryTime, &crs.CourseProgram, &crs.CertFrontpage)
	if err != nil {
		return crs, err
	}
	return crs, nil

}

func (cr *Course) UpdateCourse(r *http.Request) error {
	// course := Course{}
	// course.Name = r.FormValue("courseName")
	// course.Symbol = r.FormValue("courseSymbol")
	// course.ExpiryTime, _ = strconv.Atoi(r.FormValue("courseExpT"))

	// _, err := config.DB.Exec("UPDATE courses SET name=$1, symbol=$2, expirytime=$3 where symbol=$2", course.Name, course.Symbol, course.ExpiryTime)
	// if err != nil {
	// 	return err
	// }
	// return nil
	json.NewDecoder(r.Body).Decode(&cr)
	_, err := config.DB.Exec("UPDATE courses SET name=$1,symbol=$2,expirytime=$3,courseprogram=$4,certfrontpage=$5 WHERE symbol=$2", cr.Name, cr.Symbol, cr.ExpiryTime, cr.CourseProgram, cr.CertFrontpage)
	if err != nil {
		errors.New("500. Internal Server Error." + err.Error())
	}
	return nil
}