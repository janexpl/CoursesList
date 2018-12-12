package users

import (
	"net/http"

	"github.com/janexpl/CoursesList/controllers"
)

func RequestUser(r *http.Request) (string, error) {
	c := controllers.NewSessionsController()
	u := c.CreateNewSession()
	un, err := u.GetSessionUser(r)
	return un, err

}
