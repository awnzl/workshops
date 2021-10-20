package api

import (
	"calendar/internal/app"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Auth struct {
	Username string
	Password string
}

type User struct {
	login    string
	timezone string
}

type Event struct {
	id          string
	title       string
	description string
	time        string
	timezone    string
	duration    int32
	notes       []string
}

type API struct {
	app *app.App
	log *log.Logger
}

func New(a *app.App, l *log.Logger) *API {
	return &API{
		app: a,
		log: l,
	}
}

const eventID = "id"

func (a *API) RegisterHandlers(router *mux.Router, mwFuncs ...mux.MiddlewareFunc) {
	router.HandleFunc("/", a.main)
	router.HandleFunc("/login", a.login)
	router.HandleFunc("/logout", a.logout)
	router.HandleFunc("/api/user", a.user)
	router.HandleFunc("/api/events", a.events)
	router.HandleFunc(fmt.Sprintf("/api/event/{%v}", eventID), a.event)

	router.Use(mwFuncs...)
}

func (a *API) main(w http.ResponseWriter, r *http.Request) {
	// todo: what should I return here?
	resp := struct{ Calendar string }{ Calendar: "main end-point" }

	b, err := json.Marshal(resp)
	if err != nil {
		a.log.Println(err)
		a.writeError("internal server error", http.StatusInternalServerError, w)
		return
	}

	if err := a.writeResponse(b, w); err != nil {
		a.log.Println("response writing error", err)
	}
}

func (a *API) login(w http.ResponseWriter, r *http.Request) {
	var auth Auth

	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		a.log.Println(err)
		a.writeError("failed to read request body", http.StatusBadRequest, w)
		return
	}
	defer r.Body.Close()

	err = a.app.Login(auth.Username, auth.Password)
	if err != nil {
		a.log.Println(err)
		a.writeError("authorization fail", http.StatusUnauthorized, w)
		return
	}

	if err := a.writeResponse([]byte{}, w); err != nil {
		a.log.Println(err.Error())
	}
}

func (a *API) logout(w http.ResponseWriter, r *http.Request) {
	//todo: set content type?
}

func (a *API) user(w http.ResponseWriter, r *http.Request) {

}

func (a *API) events(w http.ResponseWriter, r *http.Request) {

}

func (a *API) event(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

    id, ok := vars[eventID]
    if !ok {
        a.log.Println(eventID, "is missing in parameters")
    }

	a.log.Println(`id := `, id) //todo: remove
}

func (a*API) writeResponse(b []byte, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}

func (a *API) writeError(msg string, status int, w http.ResponseWriter) {
	w.WriteHeader(status)

	resp := struct{	Error string }{ Error: msg }

	b, err := json.Marshal(resp)
	if err != nil {
		a.log.Println("failed to marshal", err)
		return
	}

	if _, err := w.Write(b); err != nil {
		a.log.Println("failed to write response", err)
	}
}

