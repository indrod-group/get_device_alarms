package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func GetAlarmData(user User, currentTime, interval int64) ([]byte, error) {
	accessToken := os.Getenv("ACCESS_TOKEN")
	imei := user.Imei
	startTime := currentTime - interval
	endTime := currentTime
	url := fmt.Sprintf("https://open.iopgps.com/api/device/alarm?imei=%s&startTime=%d&endTime=%d", imei, startTime, endTime)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return nil, err
	}

	req.Header.Add("AccessToken", accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body: %s\n", err)
		return nil, err
	}

	if resp.StatusCode >= 400 {
		err = fmt.Errorf("IOGPS API Error | Status code: %d, Response: %s", resp.StatusCode, string(body))
		log.Println(err)
		return nil, err
	}

	if len(body) == 0 {
		err = fmt.Errorf("empty response body")
		log.Println(err)
		return nil, err
	}

	log.Printf("IOGPS API | Status code: %d, Imei: %s\n", resp.StatusCode, imei)

	return body, nil
}

func ProcessAlarmData(user User, data []byte) error {
	if len(data) == 0 {
		log.Printf("No data to process\n")
		return nil
	}

	authToken := os.Getenv("AUTH_TOKEN")
	URL := os.Getenv("MY_API_URL")
	var details ApiResponse
	err := json.Unmarshal(data, &details)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	for _, detail := range details.Details {
		detailForPost := AlarmDataForPost(detail)
		jsonData, err := json.Marshal(detailForPost)
		if err != nil {
			log.Printf("Error: %s", err)
			return err
		}

		client := &http.Client{}
		req, _ := http.NewRequest("POST", URL+"details/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", "Token "+authToken)

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error: %s", err)
			return err
		}

		log.Printf("ACV API | Status code: %d\n", resp.StatusCode)

		// Call the function to check and send alarm messages
		CheckAndSendAlarm(user, detail)
	}

	return nil
}
