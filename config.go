package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Config struct {
	acvApiURL        string
	authToken        string
	appID            string
	loginKey         string
	twilioAccountSID string
	twilioAuthToken  string
}

// Set debug status to use the program in production
var DEBUG = os.Getenv("DEBUG") != "RENDER"

// Initialize the logging configuration based on the value of the DEBUG environment variable.
// If DEBUG is not set to "RENDER", the log output is directed to the standard output.
// Otherwise, a log file is created with a timestamped name in the "logs" folder, and the log output is directed to this file.
func initLog() {
	if !DEBUG {
		logrus.SetOutput(os.Stdout)
	} else {
		// Create the "logs" folder if it doesn't exist
		if _, err := os.Stat("logs"); os.IsNotExist(err) {
			os.Mkdir("logs", 0755)
		}

		// Generate a timestamped file name for the log file
		t := time.Now()
		fileName := fmt.Sprintf(
			"logs/log-%04d%02d%02d-%02d%02d%02d.log",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second(),
		)

		// Open the log file in write-only mode, creating it if it doesn't exist
		logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening log file: %v", err)
			return
		}

		// Configure Logrus to write to the log file
		logrus.SetOutput(logFile)

		// Set the Logrus format to JSON
		logrus.SetFormatter(&logrus.JSONFormatter{})

		// Set the Logrus level to Info
		logrus.SetLevel(logrus.InfoLevel)

		// Add a hook to send error-level logs to an external API (e.g., Loggly, Sentry, etc.)
		// logrus.AddHook(...)
	}
}
