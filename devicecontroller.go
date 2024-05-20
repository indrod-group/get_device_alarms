package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type DeviceController struct {
	next        Handler
	lastDevices []Device
}

func (dc *DeviceController) getDevices(queryParams map[string]string) ([]Device, error) {
	var apiKey = os.Getenv("API_KEY")

	// Create a new URL and set the raw query to the encoded query parameters
	u, _ := url.Parse(DEVICES_API_URL)
	q := u.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"url":   u.String(),
		}).Error("Error creating HTTP request")
		return nil, err
	}
	req.Header.Add("Authorization", "Token "+apiKey)

	// Create a context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Add the context to the request
	req = req.WithContext(ctx)

	// Create a new HTTP client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"url":   u.String(),
		}).Error("Error executing HTTP request")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
			"url":         u.String(),
		}).Error("Received non-200 response status")
		return nil, errors.New("received non-200 response status")
	}

	var devices []Device
	err = json.NewDecoder(resp.Body).Decode(&devices)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"url":   u.String(),
		}).Error("Error decoding response body")
		return nil, err
	}

	return devices, nil
}

func (dc *DeviceController) Handle(data interface{}) (interface{}, error) {
	queryParams, ok := data.(map[string]string)
	if !ok {
		return nil, errors.New("data is not of type map[string]string")
	}
	devices, err := dc.getDevices(queryParams)
	if err != nil {
		devices = dc.lastDevices
	}
	dc.lastDevices = devices
	if dc.next != nil {
		return dc.next.Handle(devices)
	}
	return devices, nil
}

func (dc *DeviceController) SetNext(next Handler) {
	dc.next = next
}
