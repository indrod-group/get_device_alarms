package main

import "log"

// init is a function that runs before main and initializes the log and the token
func init() {
	initLog()
	var err error
	app.token, err = getAccessToken()
	if err != nil {
		log.Println("Error al obtener el token de acceso:", err)
	} else {
		log.Println("Token de acceso actualizado:", app.token)
	}
}

// app is an instance of the App structure
var app App

func main() {
	app.Run()
	select {} // Bucle infinito para evitar que main() termine
}
