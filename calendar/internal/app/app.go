package app

import (
	"log"
)

type App struct {
	log *log.Logger
}

func New(l *log.Logger) *App {
	return &App{
		log: l,
	}
}

// todo: remove
var temporaryStorage = make(map[string]string)

func (a *App) Login(usr, psw string) error {
	a.log.Println("app.Login()", usr, psw)

	if _, ok := temporaryStorage[usr]; !ok {
		temporaryStorage[usr] = psw
	}

	return nil
}
