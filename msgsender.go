package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type MessageSender struct {
	next Handler
	mu   sync.Mutex
}

/*
SendMessage sends a WhatsApp message to multiple recipients using the Twilio API.

Inputs:
  - message (string): The content of the message to be sent.
  - imei (string): The IMEI of the device whose associated
    users' phone numbers will be retrieved.

Outputs:
  - None. The function only prints the message SID
    if the message is sent successfully or
    an error message if there is an error.

Example Usage:

	SendMessage("Hello, World!", "86541232122351")

This code will send the message "Hello, World!"
to the WhatsApp numbers associated with the device specified by the IMEI.
*/
func SendMessage(message, imei string) {
	if discardMessage(message) {
		return
	}
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("TWILIO_ACCOUNT_SID"),
		Password: os.Getenv("TWILIO_AUTH_TOKEN"),
	})

	numbers, err := GetPhoneNumbersFromAPI(imei)
	if err != nil {
		logrus.WithError(err).Error("Error retrieving phone numbers")
		return
	}

	for _, number := range numbers {
		params := &api.CreateMessageParams{}
		params.SetFrom("whatsapp:+14155238886")
		params.SetBody(message)
		formatedNumber := fmt.Sprintf("whatsapp:%s", number)
		params.SetTo(formatedNumber)

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
		if alarm.AlarmCode == "SOS" || alarm.AlarmCode == "LOWVOT" || alarm.AlarmCode == "REMOVE" {
			ms.mu.Lock()
			filteredAlarms = append(filteredAlarms, alarm)
			ms.mu.Unlock()
		}
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 25)

	for _, alarm := range filteredAlarms {
		wg.Add(1)
		go func(alarm Alarm) {
			defer wg.Done()
			sem <- struct{}{}        // Acquire a token
			defer func() { <-sem }() // Release the token back into the pool

			device, err := GetDeviceByImei(alarm.Imei)
			if err != nil {
				logrus.WithError(err).Error("Error getting device by IMEI")
				return
			}
			if device == nil {
				logrus.Warning("Device is nil")
				return
			}
			mb := NewMessageBuilder(device, &alarm)
			message := mb.BuildMessage()
			SendMessage(message, alarm.Imei)
		}(alarm)
	}

	wg.Wait()

	if ms.next != nil {
		return ms.next.Handle(filteredAlarms)
	}
	return filteredAlarms, nil
}

func (ms *MessageSender) SetNext(next Handler) {
	ms.next = next
}
