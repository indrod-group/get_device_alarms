package main

import (
	"encoding/json"
	"fmt"
	"io"
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

func GetUserFromApi() []User {
	authToken := os.Getenv("AUTH_TOKEN")
	URL := os.Getenv("MY_API_URL")
	req, _ := http.NewRequest("GET", URL+"user/?is_tracking=True", nil)
	req.Header.Add("Authorization", "Token "+authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := io.ReadAll(resp.Body)
		var users []User
		json.Unmarshal(data, &users)
		return users
	}
	return nil
}
