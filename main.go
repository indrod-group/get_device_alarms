package main

import (
	"alarm_notifications/auth"

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

	authenticator = auth.InitAuthenticator()

}

// app is an instance of the App structure
var authenticator *auth.Authenticator

func main() {

	go authenticator.InitiateTokenRenewal()
	go InitiateTrackingAlarms()

	select {}
}
