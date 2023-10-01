package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type User struct {
	UserName      string  `json:"user_name"`
	Imei          string  `json:"imei"`
	LicenseNumber *string `json:"license_number,omitempty"`
	Vin           *string `json:"vin,omitempty"`
	CarOwner      *string `json:"car_owner,omitempty"`
	IsTracking    bool    `json:"is_tracking"`
}

func GetUserFromApi() ([]User, error) {
	authToken := os.Getenv("AUTH_TOKEN")
	URL := os.Getenv("MY_API_URL")
	req, err := http.NewRequest("GET", URL+"user/?is_tracking=True", nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return nil, err
	}
	req.Header.Add("Authorization", "Token "+authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("The HTTP request failed with error %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, err
	}

	var users []User
	err = json.Unmarshal(data, &users)
	if err != nil {
		log.Printf("Failed to unmarshal response body: %v", err)
		return nil, err
	}

	return users, nil
}
