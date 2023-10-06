package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type User struct {
	UserName      string  `json:"user_name"`
	Imei          string  `json:"imei"`
	LicenseNumber *string `json:"license_number,omitempty"`
	Vin           *string `json:"vin,omitempty"`
	CarOwner      *string `json:"car_owner,omitempty"`
	IsTracking    bool    `json:"is_tracking"`
}

var clientUsers = &http.Client{
	Timeout: time.Second * 10,
}

func GetUserFromApi() ([]User, error) {
	authToken := os.Getenv("AUTH_TOKEN")
	URL := os.Getenv("MY_API_URL")
	req, err := http.NewRequest("GET", URL+"user/?is_tracking=True", nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to create request")
		return nil, err
	}
	req.Header.Add("Authorization", "Token "+authToken)

	resp, err := clientUsers.Do(req)
	if err != nil {
		logrus.WithError(err).Error("The HTTP request failed")
		return nil, err
	}
	defer resp.Body.Close()

	var users []User
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		logrus.WithError(err).Error("Failed to decode response body")
		return nil, err
	}

	return users, nil
}
