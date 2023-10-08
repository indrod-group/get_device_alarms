package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// init is a function that runs before main and initializes the log and the token
func init() {
	initLog()
	err := godotenv.Load(".env", ".env.development")
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	app.config = Config{
		acvApiURL:        os.Getenv("MY_API_URL"),
		authToken:        os.Getenv("AUTH_TOKEN"),
		appID:            os.Getenv("APPID"),
		loginKey:         os.Getenv("LOGIN_KEY"),
		twilioAccountSID: os.Getenv("TWILIO_ACCOUNT_SID"),
		twilioAuthToken:  os.Getenv("TWILIO_AUTH_TOKEN"),
	}
}

// app is an instance of the App structure
var app App

func main() {
	app.Run()
	select {} // Bucle infinito para evitar que main() termine
}
