package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

type DeviceController struct {
	next        Handler
	lastDevices []Device
}

const TRACKED_DEVICES_API_URL = "http://127.0.0.1:8001/api/v1/devices/?is_tracking_alarms=true"

func (dc *DeviceController) getDevices() ([]Device, error) {
	var apiKey = os.Getenv("API_KEY")
	req, err := http.NewRequest("GET", TRACKED_DEVICES_API_URL, nil)
	if err != nil {
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
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("received non-200 response status")
	}

	var devices []Device
	err = json.NewDecoder(resp.Body).Decode(&devices)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func (dc *DeviceController) Handle(data interface{}) (interface{}, error) {
	devices, err := dc.getDevices()
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
