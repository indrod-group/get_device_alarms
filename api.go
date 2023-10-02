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
)

var client = &http.Client{}

func GetAlarmData(user User, currentTime, interval int64) ([]byte, error) {
	accessToken := os.Getenv("ACCESS_TOKEN")
	imei := user.Imei
	startTime := currentTime - interval
	endTime := currentTime
	url := fmt.Sprintf("https://open.iopgps.com/api/device/alarm?imei=%s&startTime=%d&endTime=%d", imei, startTime, endTime)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request for URL %s: %w", url, err)
	}

	req.Header.Add("AccessToken", accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request to URL %s: %w", url, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("error closing response body from URL %s: %s\n", url, closeErr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body from URL %s: %w", url, err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("IOGPS API Error | Status code: %d, Response: %s", resp.StatusCode, string(body))
	}

	if len(body) == 0 {
		return nil, fmt.Errorf("empty response body from URL %s", url)
	}

	log.Printf("IOGPS API | Status code: %d, Imei: %s\n", resp.StatusCode, imei)

	return body, nil
}

func saveAlarmInAPI(detail AlarmData) error {
	authToken := os.Getenv("AUTH_TOKEN")
	URL := os.Getenv("MY_API_URL")
	detailForPost := AlarmDataForPost(detail)
	jsonData, err := json.Marshal(detailForPost)
	if err != nil {
		return fmt.Errorf("error marshalling data for detail %+v: %w", detailForPost, err)
	}

	reqURL := URL + "details/"
	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request for URL %s: %w", reqURL, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Token "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request to URL %s: %w", reqURL, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Error closing response body from URL %s: %s\n", reqURL, closeErr)
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ACV API Error | Status code: %d for URL %s", resp.StatusCode, reqURL)
	}

	log.Printf("ACV API | Status code: %d\n", resp.StatusCode)
	return nil
}

func ProcessAlarmData(user User, data []byte) error {
	if len(data) == 0 {
		log.Printf("No data to process\n")
		return nil
	}

	var details ApiResponse
	err := json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("error unmarshalling data: %w", err)
	}

	var wgAlarms sync.WaitGroup

	for _, detail := range details.Details {
		wgAlarms.Add(1)
		go func(detail AlarmData) {
			defer wgAlarms.Done()
			CheckAndSendAlarm(user, detail)
		}(detail)

		err = saveAlarmInAPI(detail)
		if err != nil {
			log.Printf("Error saving data for detail %+v: %s", detail, err)
			continue
		}
	}

	wgAlarms.Wait()

	return nil
}
