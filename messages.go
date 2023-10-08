package main

import (
	"github.com/sirupsen/logrus"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

/*
SendMessage sends a WhatsApp message to multiple recipients using the Twilio API.

Inputs:
  - message (string): The content of the message to be sent.

Outputs:
  - None. The function only prints the message SID if the message is sent successfully or an error message if there is an error.

Example Usage:

	SendMessage("Hello, World!")

This code will send the message "Hello, World!" to the three WhatsApp numbers specified in the `numbers` array.
*/
func SendMessage(message string) {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: app.config.twilioAccountSID,
		Password: app.config.twilioAuthToken,
	})

	numbers := []string{
		"whatsapp:+593979368744",
		"whatsapp:+593984924265",
		"whatsapp:+593987129357",
	}

	for _, number := range numbers {
		params := &api.CreateMessageParams{}
		params.SetFrom("whatsapp:+14155238886")
		params.SetBody(message)
		params.SetTo(number)

		resp, err := client.Api.CreateMessage(params)
		if err != nil {
			logrus.WithError(err).Error("Error sending message")
			continue
		}

		if resp.Sid != nil {
			logrus.Printf("Message sent successfully, SID: %s\n", *resp.Sid)
		} else {
			logrus.Warningf("Message sent successfully, but no SID returned")
		}
	}
}

// Function to check the alarm type and send a message if necessary
func CheckAndSendAlarm(user User, detail AlarmData) {
	alarmCodes := []string{"SOS", "REMOVE", "LOWVOT"}
	for _, alarmCode := range alarmCodes {
		if detail.AlarmCode == alarmCode {
			mb := NewMessageBuilder(&user, alarmCode, detail.AlarmTime)
			message := mb.BuildMessage()
			SendMessage(message)
		}
	}
}
