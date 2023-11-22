package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Device represents a device with various properties.
type Device struct {
	Imei             string  `json:"imei"`
	UserName         string  `json:"user_name"`
	CarOwner         *string `json:"car_owner"`
	LicenseNumber    *string `json:"license_number"`
	Vin              *string `json:"vin"`
	IsTrackingAlarms bool    `json:"is_tracking_alarms"`
	LastTimeTracked  int64   `json:"last_time_tracked"`
}

const DEVICE_ALARM_URL = "https://open.iopgps.com/api/device/alarm?imei=%s&startTime=%d&endTime=%d"
const TWENTY_FOUR_HOURS_IN_SECONDS = 86400
const DEVICES_API_URL = "http://127.0.0.1:8001/api/v1/devices/"

func (d *Device) GenerateURL() string {
	endTime := time.Now().Unix()
	var startTime int64
	if d.LastTimeTracked == 0 {
		startTime = endTime - TWENTY_FOUR_HOURS_IN_SECONDS
	} else {
		startTime = d.LastTimeTracked
	}
	d.LastTimeTracked = endTime
	return fmt.Sprintf(DEVICE_ALARM_URL, d.Imei, startTime, endTime)
}

func (d *Device) UpdateDevice() error {
	var API_KEY = os.Getenv("API_KEY")

	jsonDevice, err := json.Marshal(d)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", DEVICES_API_URL, bytes.NewBuffer(jsonDevice))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", API_KEY))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to update device: %v", resp.Status)
	}

	return nil
}
