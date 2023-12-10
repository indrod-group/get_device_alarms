package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type MessageSender struct {
	next Handler
}

/*
SendMessage sends a WhatsApp message to multiple recipients using the Twilio API.

Inputs:
  - message (string): The content of the message to be sent.

Outputs:
  - None. The function only prints the message SID if the message is sent successfully or an error message if there is an error.

Example Usage:

	SendMessage("Hello, World!")

This code will send the message "Hello, World!" to the three WhatsApp numbers specified in the numbers array.
*/
func SendMessage(message string) {
	if discardMessage(message) {
		return
	}
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("TWILIO_ACCOUNT_SID"),
		Password: os.Getenv("TWILIO_AUTH_TOKEN"),
	})

	numbers := []string{
		"whatsapp:+593979368744",
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

func discardMessage(message string) bool {
	if message == "" {
		logrus.Warningf("Discarding message: %s", message)
		return true
	}
	return false
}

func (ms *MessageSender) Handle(data interface{}) (interface{}, error) {
	alarms, ok := data.([]Alarm)
	if !ok {
		err := fmt.Errorf("MessageSender.Handle: expected []Alarm, got %T", data)
		logrus.WithError(err).Error("Error in MessageSender.Handle")
		return nil, err
	}

	var filteredAlarms []Alarm
	for _, alarm := range alarms {
		if alarm.AlarmCode == "LOWVOT" || alarm.AlarmCode == "SOS" || alarm.AlarmCode == "REMOVE" {
			filteredAlarms = append(filteredAlarms, alarm)
		}
	}

	for _, alarm := range filteredAlarms {
		device, err := GetDeviceByImei(alarm.Imei)
		if err != nil {
			logrus.WithError(err).Error("Error getting device by IMEI")
			continue
		}
		if device == nil {
			logrus.Warning("Device is nil")
			continue
		}
		mb := NewMessageBuilder(device, &alarm)
		message := mb.BuildMessage()
		SendMessage(message)
	}

	if ms.next != nil {
		return ms.next.Handle(filteredAlarms)
	}
	return filteredAlarms, nil
}

func (ms *MessageSender) SetNext(next Handler) {
	ms.next = next
}
