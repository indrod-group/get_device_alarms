package main

import (
	"os"
	"os/signal"
	"syscall"

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

	app.config = Config{
		acvApiURL:        os.Getenv("MY_API_URL"),
		authToken:        os.Getenv("AUTH_TOKEN"),
		twilioAccountSID: os.Getenv("TWILIO_ACCOUNT_SID"),
		twilioAuthToken:  os.Getenv("TWILIO_AUTH_TOKEN"),
	}

	authenticator = auth.InitAuthenticator()
}

// app is an instance of the App structure
var app App
var authenticator *auth.Authenticator

func main() {
	stopChan := make(chan struct{})
	go authenticator.InitiateTokenRenewal(stopChan)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for range signalChan {
			close(stopChan)
			os.Exit(0)
		}
	}()

	app.Run()

	select {}
}
