package models

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/janexpl/CoursesList/config"
)

type Company struct {
	ID            int64
	Name          string
	Street        string
	City          string
	Zipcode       string
	Nip           string
	Email         string
	ContactPerson string
	TelephoneNo   string
	Note          string
	//Certificates []Cert
}

func (c *Company) PutCompany(r *http.Request) (int64, error) {

	nip := c.clearnip(r.FormValue("nip"))
	fmt.Println(nip)
	if !c.checkNip(nip) {
		return 0, errors.New("Błędny numer NIP")
	}
	company := Company{}
	err := config.DB.QueryRow("SELECT nip FROM companies WHERE nip=$1", nip).Scan(&company.Nip)
	if company.Nip != "" {
		return 0, errors.New("Istnieje juz firma o takim numerze nip.")
	}

	//company.ID = bson.NewObjectId()
	company.Name = r.FormValue("name")
	company.Street = r.FormValue("street")
	company.City = r.FormValue("city")
	company.Zipcode = r.FormValue("zipcode")
	company.Nip = nip
	company.Email = r.FormValue("email")
	company.ContactPerson = r.FormValue("contactperson")
	company.TelephoneNo = r.FormValue("telephone")
	company.Note = r.FormValue("note")
	var cid int64
	err = config.DB.QueryRow("INSERT INTO companies(name,street,city,zipcode,nip,email,contactperson,telephoneno,note) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id", company.Name, company.Street, company.City, company.Zipcode, company.Nip, company.Email, company.ContactPerson, company.TelephoneNo, company.Note).Scan(&cid)
	if err != nil {
		return 0, errors.New("500. Internal Server Error." + err.Error())
	}
	return cid, nil
}

func (c *Company) AllCompanies() ([]Company, error) {
	rows, err := config.DB.Query("SELECT * FROM companies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	companies := []Company{}
	for rows.Next() {
		cmp := Company{}
		err := rows.Scan(&cmp.ID, &cmp.Name, &cmp.Street, &cmp.City, &cmp.Zipcode, &cmp.Nip, &cmp.Email, &cmp.ContactPerson, &cmp.TelephoneNo, &cmp.Note)
		if err != nil {
			return nil, err
		}
		companies = append(companies, cmp)
	}
	if err != nil {
		return companies, err
	}

	return companies, nil
}

func (c *Company) DeleteCompany(r *http.Request) error {
	nip := r.FormValue("nip")

	if nip == "" {
		return errors.New("400. Bad Request.")
	}
	_, err := config.DB.Exec("DELETE FROM companies WHERE nip = $1", nip)

	if err != nil {
		return errors.New("500. Internal Server Error")
	}
	return nil
}
func (c *Company) GetCompanyWithId(id int) (Company, error) {

	cmp := Company{}

	err := config.DB.QueryRow("SELECT * FROM companies WHERE id = $1", id).
		Scan(&cmp.ID, &cmp.Name, &cmp.Street, &cmp.City, &cmp.Zipcode, &cmp.Nip, &cmp.Email, &cmp.ContactPerson, &cmp.TelephoneNo, &cmp.Note)
	if err != nil {
		return cmp, err
	}
	return cmp, nil

}
func (c *Company) OneCompany(r *http.Request) (Company, error) {

	cmp := Company{}
	nip := r.FormValue("nip")

	if nip == "" {
		return cmp, errors.New("400. Bad Request.")
	}
	row := config.DB.QueryRow("SELECT * FROM companies WHERE nip = $1", nip)
	err := row.Scan(&cmp.ID, &cmp.Name, &cmp.Street, &cmp.City, &cmp.Zipcode, &cmp.Nip, &cmp.Email, &cmp.ContactPerson, &cmp.TelephoneNo, &cmp.Note)

	if err != nil {
		return cmp, err
	}
	return cmp, nil

}

func (c *Company) UpdateCompany(r *http.Request) error {

	company := Company{}

	company.Name = r.FormValue("name")
	company.Street = r.FormValue("street")
	company.City = r.FormValue("city")
	company.Zipcode = r.FormValue("zipcode")
	company.Nip = r.FormValue("nip")
	company.Email = r.FormValue("email")
	company.ContactPerson = r.FormValue("contactperson")
	company.TelephoneNo = r.FormValue("telephone")
	company.Note = r.FormValue("note")

	_, err := config.DB.Exec("UPDATE companies SET name=$1, street=$2, city=$3, zipcode=$4,nip=$5,email=$6,contactperson=$7,telephoneno=$8,note=$9 where nip=$5", company.Name, company.Street, company.City, company.Zipcode, company.Nip, company.Email, company.ContactPerson, company.TelephoneNo, company.Note)
	if err != nil {
		return err
	}
	return nil
}

func (c *Company) checkNip(nip string) bool {

	nipnoSpaces := strings.Replace(nip, " ", "", -1)
	purenip := strings.Replace(nipnoSpaces, "-", "", -1)
	arrayOfLetters := strings.Split(purenip, "")
	if len(arrayOfLetters) != 10 {
		return false
	}
	nipA := []int{}

	for _, i := range arrayOfLetters {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		nipA = append(nipA, j)
	}
	wagi := []int{6, 5, 7, 2, 3, 4, 5, 6, 7}
	sumw := 0
	for i := 0; i < 9; i++ {
		sumw = sumw + nipA[i]*wagi[i]

	}
	cdigit := sumw % 11

	if cdigit == nipA[9] {
		return true
	} else {
		return false
	}

}

func (c *Company) clearnip(nip string) string {
	nipnoSpaces := strings.Replace(nip, " ", "", -1)

	return strings.Replace(nipnoSpaces, "-", "", -1)
}
