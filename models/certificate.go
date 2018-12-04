package models

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/janexpl/CoursesList/config"
)

type Certificate struct {
	ID              int64
	Date            time.Time
	Student         Student
	CourseDateStart time.Time
	CourseDateEnd   time.Time
	Registry        Registry
}

func (c *Certificate) AllCertificates() ([]Certificate, error) {

	rows, err := config.DB.Query(`SELECT
	certificates.id,
    certificates.date,
    students.firstname,
    students.secondname,
    students.lastname,
    students.birthdate,
    students.birthplace,
    students.pesel,
    certificates.coursedatestart,
    certificates.coursedateend,
    registries.number,
    registries.year,
    courses.symbol
FROM
    certificates
    INNER JOIN students ON students.id = certificates.student_id
    INNER JOIN registries ON certificates.registry_id= registries.id
    INNER JOIN courses ON courses.id=registries.course_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	certificates := []Certificate{}
	for rows.Next() {
		crt := Certificate{}
		err := rows.
			Scan(&crt.ID,
				&crt.Date,
				&crt.Student.Firstname, &crt.Student.Secondname, &crt.Student.Lastname, &crt.Student.Birthdate, &crt.Student.Birthplace, &crt.Student.Pesel,
				&crt.CourseDateStart, &crt.CourseDateEnd,
				&crt.Registry.Number, &crt.Registry.Year,
				&crt.Registry.Course.Symbol)

		if err != nil {
			return nil, err
		}
		certificates = append(certificates, crt)
	}

	if err != nil {
		return certificates, err
	}

	return certificates, nil
}

func (c *Certificate) AllCertificatesWithStudent(id int64) ([]Certificate, error) {
	rows, err := config.DB.Query(`SELECT
	certificates.id,
    certificates.date,
    students.firstname,
    students.secondname,
    students.lastname,
    students.birthdate,
    students.birthplace,
    students.pesel,
    certificates.coursedatestart,
    certificates.coursedateend,
    registries.number,
    registries.year,
	courses.symbol,
	courses.name
FROM
	certificates
	
    INNER JOIN students ON students.id = certificates.student_id
    INNER JOIN registries ON certificates.registry_id= registries.id
	INNER JOIN courses ON courses.id=registries.course_id
WHERE student_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	certificates := []Certificate{}
	for rows.Next() {
		crt := Certificate{}
		err := rows.
			Scan(&crt.ID,
				&crt.Date,
				&crt.Student.Firstname, &crt.Student.Secondname, &crt.Student.Lastname, &crt.Student.Birthdate, &crt.Student.Birthplace, &crt.Student.Pesel,
				&crt.CourseDateStart, &crt.CourseDateEnd,
				&crt.Registry.Number, &crt.Registry.Year,
				&crt.Registry.Course.Symbol, &crt.Registry.Course.Name)

		if err != nil {
			return nil, err
		}
		certificates = append(certificates, crt)
	}

	if err != nil {
		return certificates, err
	}

	return certificates, nil
}

func (c *Certificate) PutCertificate(r *http.Request) error {

	st := Student{}
	stid, _ := strconv.ParseInt(r.FormValue("student"), 0, 64)
	st, err := st.GetStudentWithId(stid)
	fmt.Println(st)
	cr := Certificate{}
	cr.Date, _ = time.Parse("2006-01-02", r.FormValue("certdate"))
	cr.Student = st
	cr.CourseDateStart, _ = time.Parse("2006-01-02", r.FormValue("startdate"))
	cr.CourseDateEnd, _ = time.Parse("2006-01-02", r.FormValue("enddate"))

	reg := Registry{}
	regid, err := reg.PutRegistry(r)

	_, err = config.DB.Exec("INSERT INTO certificates(date,student_id,coursedatestart,coursedateend,registry_id) VALUES ($1,$2,$3,$4,$5)", cr.Date, cr.Student.ID, cr.CourseDateStart, cr.CourseDateEnd, regid)

	if err != nil {
		return errors.New("500. Internal Server Error." + err.Error())
	}
	return nil
}

func (c *Certificate) GetCertificatesWithCourseId(id int64) ([]Certificate, error) {
	rows, err := config.DB.Query(`SELECT
	certificates.id,
    certificates.date,
    students.firstname,
    students.secondname,
    students.lastname,
    students.birthdate,
    students.birthplace,
    students.pesel,
    certificates.coursedatestart,
    certificates.coursedateend,
    registries.number,
    registries.year,
	courses.symbol,
	courses.name
FROM
	certificates
	
    INNER JOIN students ON students.id = certificates.student_id
    INNER JOIN registries ON certificates.registry_id= registries.id
	INNER JOIN courses ON courses.id=registries.course_id
WHERE course_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	certificates := []Certificate{}
	for rows.Next() {
		crt := Certificate{}
		err := rows.
			Scan(&crt.ID,
				&crt.Date,
				&crt.Student.Firstname, &crt.Student.Secondname, &crt.Student.Lastname, &crt.Student.Birthdate, &crt.Student.Birthplace, &crt.Student.Pesel,
				&crt.CourseDateStart, &crt.CourseDateEnd,
				&crt.Registry.Number, &crt.Registry.Year,
				&crt.Registry.Course.Symbol, &crt.Registry.Course.Name)

		if err != nil {
			return nil, err
		}
		certificates = append(certificates, crt)
	}

	if err != nil {
		return certificates, err
	}

	return certificates, nil
}

func (c *Certificate) DeleteCertificate(r *http.Request) error {
	cid, _ := strconv.ParseInt(r.FormValue("id"), 0, 64)
	fmt.Println("id", cid)
	var rid int64

	reg := Registry{}
	err := config.DB.QueryRow("SELECT registry_id FROM certificates WHERE id = $1", cid).Scan(&rid)
	fmt.Println("rid", rid)
	if err != nil {
		return err
	}

	err = reg.DeleteRegistryWithId(rid)

	if err != nil {
		return errors.New("500. Internal Server Error")
	}
	return nil
}
