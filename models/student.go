package models

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/janexpl/CoursesList/config"
)

type Student struct {
	ID            int64
	Firstname     string
	Lastname      string
	Secondname    string
	Birthdate     time.Time
	Birthplace    string
	Pesel         string
	AddressStreet string
	AddressCity   string
	TelephoneNo   string
	Company       Company
	AddressZip    string
	//Certificates []Cert
}

func (s *Student) PutStudent(r *http.Request) (int64, error) {
	// if !s.checkPesel(r.FormValue("pesel")) {
	// 	return 0, errors.New("Błędny numer pesel")
	// }
	student := Student{}
	pesel := strings.Replace(r.FormValue("pesel"), " ", "", -1)
	err := config.DB.QueryRow("SELECT pesel FROM students WHERE pesel=$1", pesel).Scan(&student.Pesel)
	if student.Pesel != "" {
		return 0, errors.New("Istnieje juz kursant o takim samym numerze pesel")
	}
	var cpid int
	cp := Company{}
	if r.FormValue("company") == "" {
		cpid = 0
	} else {
		cpid, _ = strconv.Atoi(r.FormValue("company"))
		cp, err = cp.GetCompanyWithId(cpid)
		if err != nil {
			return 0, err
		}
	}
	//course.ID = bson.NewObjectId()
	student.Firstname = r.FormValue("firstname")
	student.Lastname = r.FormValue("lastname")
	student.Secondname = r.FormValue("secondname")
	student.Birthdate, _ = time.Parse("2006-01-02", r.FormValue("birthdate"))
	student.Birthplace = r.FormValue("birthplace")
	student.Pesel = pesel
	student.AddressStreet = r.FormValue("adstreet")
	student.AddressCity = r.FormValue("adcity")
	student.AddressZip = r.FormValue("zip")
	student.TelephoneNo = r.FormValue("telephone")
	student.Company = cp
	var sid int64
	err = config.DB.QueryRow("INSERT INTO students(firstname,lastname,secondname,birthdate,birthplace,pesel,addressstreet,addresscity,telephoneno,company_id,addresszip) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id", student.Firstname, student.Lastname, student.Secondname, student.Birthdate, student.Birthplace, student.Pesel, student.AddressStreet, student.AddressCity, student.TelephoneNo, student.Company.ID, student.AddressZip).Scan(&sid)

	if err != nil {
		return 0, errors.New("500. Internal Server Error." + err.Error())
	}
	return sid, nil
}
func (s *Student) AllStudentsWithCompany(id int64) ([]Student, error) {

	rows, err := config.DB.Query(`SELECT
	students.id,
	students.firstname,
	students.lastname,
	students.secondname,
	students.birthdate,
	students.birthplace,
	students.pesel,
	students.addressstreet,
	students.addresscity,
	students.telephoneno,
	companies.id,
	companies.name,
	students.addresszip
FROM 
	students 
	LEFT OUTER JOIN companies ON students.company_id = companies.id
WHERE students.company_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	students := []Student{}
	for rows.Next() {
		st := Student{}
		err := rows.Scan(&st.ID, &st.Firstname, &st.Lastname, &st.Secondname, &st.Birthdate, &st.Birthplace, &st.Pesel, &st.AddressStreet, &st.AddressCity, &st.TelephoneNo, &st.Company.ID, &st.Company.Name, &st.AddressZip)
		if err != nil {
			return nil, err
		}
		students = append(students, st)
	}

	if err != nil {
		return students, err
	}

	return students, nil
}

func (s *Student) GetStudentWithId(id int64) (Student, error) {
	st := Student{}
	var cpid sql.NullInt64
	var cpname sql.NullString
	err := config.DB.QueryRow(`SELECT
		students.id,
		students.firstname,
		students.lastname,
		students.secondname,
		students.birthdate,
		students.birthplace,
		students.pesel,
		students.addressstreet,
		students.addresscity,
		students.telephoneno,
		companies.id,
		companies.name,
		students.addresszip
	FROM 
		students 
		LEFT OUTER JOIN companies ON students.company_id = companies.id 
	WHERE students.id=$1`, id).
		Scan(&st.ID, &st.Firstname, &st.Lastname, &st.Secondname, &st.Birthdate, &st.Birthplace, &st.Pesel, &st.AddressStreet, &st.AddressCity, &st.TelephoneNo, &cpid, &cpname, &st.AddressZip)
	if cpid.Valid {
		st.Company.ID = cpid.Int64
	} else {
		st.Company.ID = 0
	}
	if cpname.Valid {
		st.Company.Name = cpname.String
	} else {
		st.Company.Name = ""
	}
	if err != nil {
		return st, err
	}
	return st, nil
}

func (s *Student) AllStudents() ([]Student, error) {

	rows, err := config.DB.Query(`SELECT
	students.id,
	students.firstname,
	students.lastname,
	students.secondname,
	students.birthdate,
	students.birthplace,
	students.pesel,
	students.addressstreet,
	students.addresscity,
	students.telephoneno,
	companies.id,
	companies.name,
	students.addresszip
	
FROM 
	students 
	LEFT OUTER JOIN companies ON students.company_id = companies.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	students := []Student{}
	for rows.Next() {
		st := Student{}
		var cpid sql.NullInt64
		var cpname sql.NullString

		err := rows.Scan(&st.ID, &st.Firstname, &st.Lastname, &st.Secondname, &st.Birthdate, &st.Birthplace, &st.Pesel, &st.AddressStreet, &st.AddressCity, &st.TelephoneNo, &cpid, &cpname, &st.AddressZip)
		if cpid.Valid {
			st.Company.ID = cpid.Int64
		} else {
			st.Company.ID = 0
		}
		if cpname.Valid {
			st.Company.Name = cpname.String
		} else {
			st.Company.Name = ""
		}

		if err != nil {
			return nil, err
		}
		students = append(students, st)

	}

	if err != nil {
		return students, err
	}

	return students, nil
}

func (s *Student) DeleteStudent(r *http.Request) error {
	pesel := r.FormValue("pesel")

	if pesel == "" {
		return errors.New("Brak pesel")
	}
	peseli, _ := strconv.ParseInt(pesel, 0, 64)
	_, err := config.DB.Exec("DELETE FROM students WHERE pesel = $1", peseli)
	return err
}

func (s *Student) OneStudent(r *http.Request) (Student, error) {

	st := Student{}
	var cpid sql.NullInt64
	var cpname sql.NullString
	pesel := r.FormValue("pesel")
	if pesel == "" {
		return st, errors.New("400. Bad Request.")
	}

	err := config.DB.QueryRow(`SELECT
		students.id,
		students.firstname,
		students.lastname,
		students.secondname,
		students.birthdate,
		students.birthplace,
		students.pesel,
		students.addressstreet,
		students.addresscity,
		students.telephoneno,
		companies.id,
		companies.name,
		students.addresszip
	FROM 
		students 
		LEFT OUTER JOIN companies ON students.company_id = companies.id 
	WHERE pesel=$1`, pesel).
		Scan(&st.ID, &st.Firstname, &st.Lastname, &st.Secondname, &st.Birthdate, &st.Birthplace, &st.Pesel, &st.AddressStreet, &st.AddressCity, &st.TelephoneNo, &cpid, &cpname, &st.AddressZip)
	if cpid.Valid {
		st.Company.ID = cpid.Int64
	} else {
		st.Company.ID = 0
	}
	if cpname.Valid {
		st.Company.Name = cpname.String
	} else {
		st.Company.Name = ""
	}

	if err != nil {
		return st, err
	}

	return st, nil

}

func (s *Student) UpdateStudent(r *http.Request) error {
	cp := Company{}
	var cpid int
	if r.FormValue("company") == "" {
		cpid = 0
	} else {
		cpid, _ = strconv.Atoi(r.FormValue("company"))

	}
	cp, err := cp.GetCompanyWithId(cpid)
	if err != nil {
		return err
	}
	student := Student{}
	student.Firstname = r.FormValue("firstname")
	student.Secondname = r.FormValue("secondname")
	student.Lastname = r.FormValue("lastname")
	student.Birthdate, _ = time.Parse("2006-01-02", r.FormValue("birthdate"))
	student.Birthplace = r.FormValue("birthplace")
	student.Pesel = r.FormValue("pesel")
	student.AddressStreet = r.FormValue("adstreet")
	student.AddressCity = r.FormValue("adcity")
	student.AddressZip = r.FormValue("zip")
	student.TelephoneNo = r.FormValue("telephone")
	student.Company = cp
	_, err = config.DB.Exec("UPDATE students SET firstname=$1,lastname=$2,secondname=$3,birthdate=$4,birthplace=$5,pesel=$6,addressstreet=$7,addresscity=$8,telephoneno=$9,company_id=$10,addresszip=$11 where pesel=$6", student.Firstname, student.Lastname, student.Secondname, student.Birthdate, student.Birthplace, student.Pesel, student.AddressStreet, student.AddressCity, student.TelephoneNo, student.Company.ID, student.AddressZip)
	if err != nil {
		return err
	}
	return nil
}

func (s *Student) checkPesel(pesels string) bool {

	arrayOfLetters := strings.Split(pesels, "")
	if len(arrayOfLetters) != 11 {
		return false
	}
	peselA := []int{}

	for _, i := range arrayOfLetters {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		peselA = append(peselA, j)
	}
	wagi := []int{1, 3, 7, 9, 1, 3, 7, 9, 1, 3}
	sumw := 0
	for i := 0; i < 10; i++ {
		sumw = sumw + peselA[i]*wagi[i]

	}
	cdigit := (10 - (sumw % 10)) % 10

	if cdigit == peselA[10] {
		return true
	} else {
		return false
	}

}
