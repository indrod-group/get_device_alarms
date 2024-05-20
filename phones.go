package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// UserPhoneNumbers is a slice of UserPhoneNumber.
type UserPhoneNumbers []UserPhoneNumber

// UnmarshalUserPhoneNumbers takes a byte slice and deserializes it into a UserPhoneNumbers.
// It returns a UserPhoneNumbers and an error if any occurred during deserialization.
func UnmarshalUserPhoneNumbers(data []byte) (UserPhoneNumbers, error) {
	var r UserPhoneNumbers
	err := json.Unmarshal(data, &r)
	return r, err
}

// Marshal serializes a UserPhoneNumbers into a byte slice.
// It returns a byte slice and an error if any occurred during serialization.
func (r *UserPhoneNumbers) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// UserPhoneNumber represents a user and their phone numbers.
// User is a UUID that uniquely identifies each user.
// PhoneNumbers is a slice of PhoneNumber that contains the user's phone numbers.
type UserPhoneNumber struct {
	User         string        `json:"user"`
	PhoneNumbers []PhoneNumber `json:"phone_numbers"`
}

// PhoneNumber represents a phone number.
// PhoneNumber is a string that contains the phone number.
type PhoneNumber struct {
	PhoneNumber string `json:"phone_number"`
}

func GetPhoneNumbersFromAPI(imei string) ([]string, error) {
	var apiKey = os.Getenv("API_KEY")
	url := fmt.Sprintf("https://api.road-safety-ec.com/api/v1/devices/%s/phones/", imei)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error making the request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error doing the request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading the response: %v", err)
	}

	var userPhoneNumbers []UserPhoneNumber
	err = json.Unmarshal(body, &userPhoneNumbers)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling the JSON: %v", err)
	}

	var phoneNumbers []string
	for _, userPhoneNumber := range userPhoneNumbers {
		for _, phoneNumber := range userPhoneNumber.PhoneNumbers {
			phoneNumbers = append(phoneNumbers, phoneNumber.PhoneNumber)
		}
	}

	return phoneNumbers, nil
}
