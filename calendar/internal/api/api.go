package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"calendar/internal/app"
	"calendar/internal/models"
)


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
	router.HandleFunc("/login", a.login).Methods("POST")
	router.HandleFunc("/logout", a.logout)
	router.HandleFunc("/api/user", a.user).Methods("PUT")
	router.HandleFunc("/api/events", a.events).Methods("GET", "POST")
	router.HandleFunc(fmt.Sprintf("/api/event/{%v}", eventID), a.event).Methods("GET", "PUT")

	router.Use(mwFuncs...)
}

func (a *API) login(w http.ResponseWriter, r *http.Request) {
	var auth models.Auth

	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		a.writeError(err, http.StatusBadRequest, w)
		return
	}
	defer r.Body.Close()

	err = a.app.Login(auth.Username, auth.Password)
	if err != nil {
		a.writeError(err, http.StatusUnauthorized, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) logout(w http.ResponseWriter, r *http.Request) {
	a.app.Logout()

	w.WriteHeader(http.StatusOK)
}

// update user's timezone
// can be done only for logged in users
func (a *API) user(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		a.writeError(err, http.StatusBadRequest, w)
		return
	}
	defer r.Body.Close()

	//TODO: find the correct way to check authorization header before loading Body
	if a.app.IsAuthorized_DebugPurposeSolution(user.Login) {
		a.writeError(err, http.StatusUnauthorized, w)
		return
	}

	err = a.app.User(user)
	if err != nil {
		a.writeError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// get all events or create events
func (a *API) events(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.writeAllEvents(w)
	case http.MethodPost:
		a.addEvents(w, r)
	}
}

func (a *API) event(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

    id, ok := vars[eventID]
    if !ok {
        a.log.Println(eventID, "is missing in parameters")
    }

	switch r.Method {
	case http.MethodGet:
		a.getEvent(w, id)
	case http.MethodPut:
		a.updateEvent(w, r)
	}
}

func (a *API) writeAllEvents(w http.ResponseWriter) {
	events, err := a.app.Events()
	if err != nil {
		a.writeError(err, http.StatusInternalServerError, w)
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		a.log.Println(err)
	}
}

func (a *API) addEvents(w http.ResponseWriter, r *http.Request) {
	var events []models.Event

	err := json.NewDecoder(r.Body).Decode(&events)
	if err != nil {
		a.writeError(err, http.StatusBadRequest, w)
		return
	}
	defer r.Body.Close()

	err = a.app.AddEvents(events)
	if err != nil {
		a.writeError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) getEvent(w http.ResponseWriter, id string) {
	event, err := a.app.Event(id)
	if err != nil {
		a.writeError(err, http.StatusInternalServerError, w)
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		a.log.Println(err)
	}
}

func (a *API) updateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		a.writeError(err, http.StatusBadRequest, w)
		return
	}
	defer r.Body.Close()

	err = a.app.UpdateEvent(event)
	if err != nil {
		a.writeError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}


/* Helpers */

func (a *API) writeError(err error, status int, w http.ResponseWriter) {
	a.log.Println(err)
	http.Error(w, err.Error(), status)
}
