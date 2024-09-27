package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"calendar/internal/app"
	"calendar/internal/helpers"
	"calendar/internal/middleware"
	"calendar/internal/models"
)

type API struct {
	app *app.App
	val middleware.JWTValidator
	log *log.Logger
}

type loginResponse struct {
	Token string `json:"token"`
}

func New(a *app.App, v middleware.JWTValidator, l *log.Logger) *API {
	return &API{
		app: a,
		val: v,
		log: l,
	}
}

const eventID = "id"

func (a *API) RegisterHandlers(router *mux.Router) {
	chain := alice.New(middleware.Authorization(a.log, a.val))

	handler := func(endpoint http.HandlerFunc) http.Handler {
		return chain.Then(http.HandlerFunc(endpoint))
	}

	router.HandleFunc("/login", a.login).Methods("POST")
	router.Handle("/logout", handler(a.logout))
	router.Handle("/api/user", handler(a.user)).Methods("PUT")
	router.Handle("/api/events", handler(a.events)).Methods("GET", "POST")
	router.Handle(fmt.Sprintf("/api/event/{%v}", eventID), handler(a.event)).Methods("GET", "PUT")

	router.Use(middleware.Logger(a.log))

	router.Handle("/metrics", promhttp.Handler())
}

func (a *API) login(w http.ResponseWriter, r *http.Request) {
	var auth models.Auth

	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		a.log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	token, err := a.app.Login(auth.Username, auth.Password)
	if err != nil {
		a.log.Println(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(loginResponse{ Token: token })
	if err != nil {
		a.log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) logout(w http.ResponseWriter, r *http.Request) {
	username, ok := a.username(w, r)
	if !ok {
		return
	}

	a.app.Logout(username)

	w.WriteHeader(http.StatusOK)
}

// update user's timezone
func (a *API) user(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		a.log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = a.app.User(user)
	if err != nil {
		a.log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// get all events or create events
func (a *API) events(w http.ResponseWriter, r *http.Request) {
	username, ok := a.username(w, r)
	if !ok {
		return
	}

	switch r.Method {
	case http.MethodGet:
		a.writeAllEvents(username, w)
	case http.MethodPost:
		a.addEvents(username, w, r)
	}
}

func (a *API) event(w http.ResponseWriter, r *http.Request) {
	username, ok := a.username(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
    id, ok := vars[eventID]
    if !ok {
		err := fmt.Errorf("%v is missing in request parameters", eventID)
		a.log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
    }

	switch r.Method {
	case http.MethodGet:
		a.writeEvent(w, username, id)
	case http.MethodPut:
		a.updateEvent(w, r, username, id)
	}
}

func (a *API) writeAllEvents(username string, w http.ResponseWriter) {
	events, err := a.app.Events(username)
	if err != nil {
		a.log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		a.log.Println(err)
	}
}

func (a *API) addEvents(username string, w http.ResponseWriter, r *http.Request) {
	var events []models.Event

	err := json.NewDecoder(r.Body).Decode(&events)
	if err != nil {
		a.log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = a.app.AddEvents(username, events)
	if err != nil {
		a.log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) writeEvent(w http.ResponseWriter, username, id string) {
	event, err := a.app.Event(username, id)
	if err != nil {
		a.log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		a.log.Println(err)
	}
}

func (a *API) updateEvent(w http.ResponseWriter, r *http.Request, usr, eventID string) {
	var event models.Event

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		a.log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = a.app.UpdateEvent(usr, eventID, event)
	if err != nil {
		a.log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

/* Helpers */
func (a *API) username(w http.ResponseWriter, r *http.Request) (username string, ok bool) {
	val := r.Context().Value(helpers.CtxValKey("username"))

	username, ok = val.(string)
	if !ok {
		err := fmt.Errorf("incorrect context value")
		a.log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}
