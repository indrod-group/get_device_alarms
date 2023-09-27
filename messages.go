package main

import (
	"log"
	"os"

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
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
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
			log.Printf("Error sending message: %s\n", err.Error())
			continue
		}

		if resp.Sid != nil {
			log.Printf("Message sent successfully, SID: %s\n", *resp.Sid)
		} else {
			log.Println("Message sent successfully, but no SID returned")
		}
	}
}
