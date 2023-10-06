package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// init is a function that runs before main and initializes the log and the token
func init() {
	initLog()
	var err error
	err = godotenv.Load(".env", ".env.development")
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	app.token, err = getAccessToken()
	if err != nil {
		logrus.Fatal("Error al obtener el token de acceso:", err)
	} else {
		logrus.Info("Token de acceso actualizado:", app.token)
	}
}

// app is an instance of the App structure
var app App

func main() {
	app.Run()
	select {} // Bucle infinito para evitar que main() termine
}
