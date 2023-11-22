package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Alarm struct {
	Imei         string  `json:"device_imei"`
	PositionType *string `json:"position_type,omitempty"`
	Lat          *string `json:"lat,omitempty"`
	Lng          *string `json:"lng,omitempty"`
	Time         int64   `json:"time"`
	AlarmCode    string  `json:"alarm_code"`
	AlarmType    int64   `json:"alarm_type"`
	Course       *int64  `json:"course,omitempty"`
	DeviceType   int64   `json:"device_type"`
	Speed        *int64  `json:"speed,omitempty"`
}

const ALARMS_API_URL = "http://127.0.0.1:8001/api/v1/alarms/"

func (a *Alarm) CreateAlarm() error {
	var API_KEY = os.Getenv("API_KEY")

	jsonAlarm, err := json.Marshal(a)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", ALARMS_API_URL, bytes.NewBuffer(jsonAlarm))
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

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAlreadyReported {
		return fmt.Errorf("failed to create alarm: %v", resp.Status)
	}

	return nil
}
