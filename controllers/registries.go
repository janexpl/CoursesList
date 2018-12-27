package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/janexpl/CoursesList/logging"
	"github.com/janexpl/CoursesList/models"
)

type RegistriesController struct{}

func NewRegistriesController() *RegistriesController {
	return &RegistriesController{}
}

func (c *RegistriesController) GetLastNumberWithSymbol(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	rg := models.Registry{}
	number, err := rg.GetLastNumber(r)

	if err != nil {
		http.Error(w, http.StatusText(500)+err.Error(), http.StatusInternalServerError)
		return
	}
	uj, err := json.Marshal(number)
	if err != nil {
		logging.Error.Println(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	fmt.Fprintf(w, "%s\n", uj)

}
