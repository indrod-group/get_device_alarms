package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"
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
const DEVICES_API_URL = "https://api.road-safety-ec.com/api/v1/devices/"

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
	var apiKey = os.Getenv("API_KEY")

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
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", apiKey))

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

// CleanAndValidateIMEI removes whitespace from the IMEI and checks that all characters are digits.
func CleanAndValidateIMEI(imei string) (string, error) {
	cleanIMEI := strings.ReplaceAll(imei, " ", "")

	for _, char := range cleanIMEI {
		if !unicode.IsDigit(char) {
			return "", errors.New("IMEI contains non-digit characters")
		}
	}

	return cleanIMEI, nil
}

const DEVICE_INFO_URL = "https://api.road-safety-ec.com/api/v1/devices/%s/"

func GetDeviceByImei(imei string) (*Device, error) {
	var apiKey = os.Getenv("API_KEY")

	cleanIMEI, err := CleanAndValidateIMEI(imei)
	if err != nil {
		return nil, fmt.Errorf("failed to clean and validate IMEI: %w", err)
	}

	client := &http.Client{}
	url := fmt.Sprintf(DEVICE_INFO_URL, cleanIMEI)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response status: %d", resp.StatusCode)
	}

	var device Device
	err = json.NewDecoder(resp.Body).Decode(&device)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &device, nil
}
