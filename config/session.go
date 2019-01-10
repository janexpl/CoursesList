package config

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/janexpl/CoursesList/logging"

	"github.com/gobuffalo/uuid"
)

type SessionsController struct {
}

func NewSessionsController() *SessionsController {

	return &SessionsController{}
}

type loggedUser struct {
	Email string
	Admin int
}
type sessions struct {
	un string

	lastActivity time.Time
}

var dbUsers = map[string]loggedUser{}
var dbSessions = map[string]sessions{} // session ID, session
var dbSessionsCleaned time.Time

const sessionLength int = 2000

// var dbUsers = map[string]models.User{} // user ID, user
func SetFlash(w http.ResponseWriter, r *http.Request, m []byte) {
	c := &http.Cookie{
		Name:  "flash",
		Value: encode(m),
	}
	http.SetCookie(w, c)
}
func getFlash(w http.ResponseWriter, r *http.Request) (string, error) {
	c, err := r.Cookie("flash")
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return "", nil
		default:
			return "", err
		}
	}
	value, err := decode(c.Value)
	if err != nil {
		return "", err
	}
	dc := &http.Cookie{Name: "flash", MaxAge: -1, Path: "/users/login", Expires: time.Unix(1, 0)}
	http.SetCookie(w, dc)
	st := fmt.Sprintf("%s", value)
	return st, nil
}
func AlreadyLoggedIn(w http.ResponseWriter, req *http.Request) bool {

	c, err := req.Cookie("session")
	if err != nil {

		return false
	}
	se, ok := dbSessions[c.Value]
	if ok {
		se.lastActivity = time.Now()
		dbSessions[c.Value] = se
	}
	// refresh session
	c.MaxAge = sessionLength
	c.Path = "/"
	http.SetCookie(w, c)

	return true
}

func GetLoggedUser(req *http.Request) (loggedUser, error) {

	c, err := req.Cookie("session")

	if err != nil {
		logging.Error.Println(err.Error())
		return loggedUser{}, err
	}
	se, _ := dbUsers[c.Value]

	return se, nil
}

func NewSession(w http.ResponseWriter, r *http.Request, email string, admin int) {
	sID, _ := uuid.NewV4()
	c := &http.Cookie{
		Name:  "session",
		Value: sID.String(),
	}
	c.MaxAge = sessionLength
	c.Path = "/"

	http.SetCookie(w, c)
	dbUsers[c.Value] = loggedUser{email, admin}
	dbSessions[c.Value] = sessions{email, time.Now()}

}

func DeleteSession(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("session")
	// delete the session
	delete(dbSessions, c.Value)
	delete(dbUsers, c.Value)
	// remove the cookie
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}

	http.SetCookie(w, c)
	if time.Now().Sub(dbSessionsCleaned) > (time.Second * 30) {
		go cleanSessions()
	}

}
func cleanSessions() {

	for k, v := range dbSessions {
		if time.Now().Sub(v.lastActivity) > (time.Second * 30) {
			delete(dbSessions, k)
		}
	}
	dbSessionsCleaned = time.Now()

}

func encode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func decode(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}
