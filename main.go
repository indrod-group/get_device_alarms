package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/indrod-group/get_device_alarms/auth"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// init performs initial setup before the main function executes.
// It initializes logging and loads environment variables from .env files.
// It also initializes the authenticator for handling API authentication.
func init() {
	initLog()
	err := godotenv.Load(".env", ".env.development")
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}
	authenticator = auth.InitAuthenticator()
}

// authenticator holds an instance of the Authenticator which manages authentication.
var authenticator *auth.Authenticator

// main sets up signal handling and starts background goroutines.
// It listens for OS termination signals to gracefully shut down the application.
func main() {
	// Initiates token renewal and alarm tracking in separate goroutines.
	go authenticator.InitiateTokenRenewal()
	go InitiateTrackingAlarms()

	// Creates a channel to receive operating system signals.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Blocks until a signal is received.
	<-sigChan
	// Logs the reception of the signal and exits the program.
	logrus.Info("Program terminated by interrupt signal")
	os.Exit(0)
}
