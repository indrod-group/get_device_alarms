package main

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type App struct {
	sem    chan bool
	wg     sync.WaitGroup
	users  []User
	config Config
}

// Run is a method that runs the main logic of the application
func (app *App) Run() {
	app.sem = make(chan bool, 10)

	originalUsers := app.users

	userTicker := time.NewTicker(5 * time.Minute)
	go func() {
		for {
			newUsers, err := GetUserFromApi()
			if err != nil {
				logrus.WithError(err).Error("Error al obtener los usuarios")
				app.users = originalUsers
			} else {
				app.users = newUsers
			}
			<-userTicker.C
		}
	}()

	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for {
			if app.users == nil {
				continue
			}
			for i := 0; i < len(app.users); i += 9 {
				end := i + 9
				if end > len(app.users) {
					end = len(app.users)
				}
				ProcessUsers(app.users[i:end], app.sem, &app.wg)
				app.wg.Wait()
			}
			<-ticker.C
		}
	}()
}
