package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"

	"calendar/internal/app"
	"calendar/internal/helpers"
	"calendar/internal/middleware"
	"calendar/internal/models"
)

type API struct {
	app *app.App
	log *log.Logger
}

type loginResponse struct {
	Token string `json:"token"`
}

func New(a *app.App, l *log.Logger) *API {
	return &API{
		app: a,
		log: l,
	}
}

const eventID = "id"

func (a *API) RegisterHandlers(router *mux.Router) {
	chain := alice.New(middleware.Authorization(a.log, a.app.AuthApp()))

	handler := func(endpoint http.HandlerFunc) http.Handler {
		return chain.Then(http.HandlerFunc(endpoint))
	}

	router.HandleFunc("/login", a.login).Methods("POST")
	router.Handle("/logout", handler(a.logout))
	router.Handle("/api/user", handler(a.user)).Methods("PUT")
	router.Handle("/api/events", handler(a.events)).Methods("GET", "POST")
	router.Handle(fmt.Sprintf("/api/event/{%v}", eventID), handler(a.event)).Methods("GET", "PUT")

	router.Use(middleware.Logger(a.log))
}

func (a *API) login(w http.ResponseWriter, r *http.Request) {
	var auth models.Auth

	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		a.writeError(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	token, err := a.app.Login(auth.Username, auth.Password)
	if err != nil {
		a.writeError(w, err, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(loginResponse{ Token: token })
	if err != nil {
		a.writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) logout(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(helpers.CtxValKey("username"))

	usr, ok := username.(string)
	if !ok {
		a.writeError(w, fmt.Errorf("incorrect context value"), http.StatusInternalServerError)
		return
	}

	a.app.Logout(usr)

	w.WriteHeader(http.StatusOK)
}

// update user's timezone
// can be done only for logged in users
func (a *API) user(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		a.writeError(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = a.app.User(user)
	if err != nil {
		a.writeError(w, err, http.StatusInternalServerError)
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
		a.writeError(w, err, http.StatusInternalServerError)
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
		a.writeError(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = a.app.AddEvents(events)
	if err != nil {
		a.writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) getEvent(w http.ResponseWriter, id string) {
	event, err := a.app.Event(id)
	if err != nil {
		a.writeError(w, err, http.StatusInternalServerError)
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
		a.writeError(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = a.app.UpdateEvent(event)
	if err != nil {
		a.writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}


/* Helpers */

func (a *API) writeError(w http.ResponseWriter, err error, status int) {
	a.log.Println(err)
	http.Error(w, err.Error(), status)
}
