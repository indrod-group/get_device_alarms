package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type AlarmData struct {
	Code         int64   `json:"code"`
	PositionType *string `json:"positionType,omitempty"`
	Imei         string  `json:"imei"`
	Lat          *string `json:"lat,omitempty"`
	Lng          *string `json:"lng,omitempty"`
	Time         int64   `json:"time"`
	Speed        *int64  `json:"speed,omitempty"`
	Course       *int64  `json:"course,omitempty"`
	AlarmCode    string  `json:"alarmCode"`
	AlarmTime    int64   `json:"alarmTime"`
	DeviceType   int64   `json:"deviceType"`
	AlarmType    int64   `json:"alarmType"`
}

type AlarmDataForPost struct {
	Code         int64   `json:"code"`
	PositionType *string `json:"position_type,omitempty"`
	Imei         string  `json:"imei"`
	Lat          *string `json:"lat,omitempty"`
	Lng          *string `json:"lng,omitempty"`
	Time         int64   `json:"time"`
	Speed        *int64  `json:"speed,omitempty"`
	Course       *int64  `json:"course,omitempty"`
	AlarmCode    string  `json:"alarm_code"`
	AlarmTime    int64   `json:"alarm_time"`
	DeviceType   int64   `json:"device_type"`
	AlarmType    int64   `json:"alarm_type"`
}

type ApiResponse struct {
	Code    int64       `json:"code"`
	Details []AlarmData `json:"details"`
}

type CronJob struct {
	user        User
	currentTime int64
	interval    int64
}

func NewCronJob(user User, startTime int64, interval int64) *CronJob {
	return &CronJob{
		user:        user,
		currentTime: startTime,
		interval:    interval,
	}
}

func (c *CronJob) GetAlarmData() []byte {
	accessToken := os.Getenv("ACCESS_TOKEN")
	imei := c.user.Imei
	startTime := c.currentTime - c.interval
	endTime := c.currentTime
	url := fmt.Sprintf("https://open.iopgps.com/api/device/alarm?imei=%s&startTime=%d&endTime=%d", imei, startTime, endTime)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return nil
	}

	req.Header.Add("AccessToken", accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body: %s\n", err)
		return nil
	}

	if resp.StatusCode >= 400 {
		log.Printf("IOGPS API Error | Status code: %d, Response: %s\n", resp.StatusCode, string(body))
		return nil
	}

	if len(body) == 0 {
		log.Printf("Empty response body\n")
		return nil
	}

	log.Printf("IOGPS API | Status code: %d, Imei: %s\n", resp.StatusCode, imei)

	c.currentTime = time.Now().Unix()

	return body
}

func (c *CronJob) ProcessAlarmData(data []byte) {
	if len(data) == 0 {
		log.Printf("No data to process\n")
		return
	}

	authToken := os.Getenv("AUTH_TOKEN")
	URL := os.Getenv("MY_API_URL")
	var details ApiResponse
	err := json.Unmarshal(data, &details)
	if err != nil {
		log.Fatalf("Error: %s", err)
		return
	}

	for _, detail := range details.Details {
		detailForPost := AlarmDataForPost(detail)
		jsonData, err := json.Marshal(detailForPost)
		if err != nil {
			log.Fatalf("Error: %s", err)
			return
		}

		req, _ := http.NewRequest("POST", URL+"details/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", "Token "+authToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error: %s", err)
			return
		}

		log.Printf("ACV API | Status code: %d\n", resp.StatusCode)

		alarmCodes := []string{"SOS", "REMOVE", "LOWVOT"}
		for _, alarmCode := range alarmCodes {
			if detail.AlarmCode == alarmCode {
				mb := NewMessageBuilder(&c.user, alarmCode, detail.AlarmTime)
				message := mb.BuildMessage()
				SendMessage(message)
			}
		}

	}
}

func (c *CronJob) Run(sem *sync.WaitGroup) {
	defer sem.Done()

	apiData := c.GetAlarmData()
	c.ProcessAlarmData(apiData)
}
