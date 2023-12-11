package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type RequestExecutor struct {
	next Handler
}

func (re *RequestExecutor) Handle(data interface{}) (interface{}, error) {
	urls, ok := data.([]string)
	if !ok {
		return nil, fmt.Errorf("RequestExecutor.Handle: expected []string, got %T", data)
	}

	var allAlarmData []AlarmData
	var mutex sync.Mutex
	var wg sync.WaitGroup

	sem := make(chan struct{}, 10)

	for _, url := range urls {
		wg.Add(1)
		go re.processURL(url, &allAlarmData, &mutex, &wg, sem)
	}

	wg.Wait()

	if re.next != nil {
		return re.next.Handle(allAlarmData)
	}
	return allAlarmData, nil
}

func (re *RequestExecutor) processURL(url string, allAlarmData *[]AlarmData, mutex *sync.Mutex, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()

	sem <- struct{}{}
	defer func() { <-sem }()

	req, err := createRequest(url)
	if err != nil {
		logrus.Warning(err)
		return
	}

	// Create a context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Add the context to the request
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Warning(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.WithFields(logrus.Fields{
			"status": resp.StatusCode,
			"url":    url,
		}).Warning("Received non-200 HTTP status")
		return
	}

	var alarmResponse AlarmResponse
	err = json.NewDecoder(resp.Body).Decode(&alarmResponse)
	if err != nil {
		// Log the alarm information even if there is an error
		logrus.Info("Alarm details: ", alarmResponse.Details)
		logrus.WithFields(logrus.Fields{
			"error": err,
			"url":   url,
		}).Warning("Error decoding the response body")
		return
	}
	// Check if alarmResponse.Details is empty
	if len(alarmResponse.Details) == 0 {
		return
	}

	mutex.Lock()
	*allAlarmData = append(*allAlarmData, alarmResponse.Details...)
	mutex.Unlock()
}

func (re *RequestExecutor) SetNext(next Handler) {
	re.next = next
}

func createRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request for URL %s: %w", url, err)
	}

	token := os.Getenv("ACCESS_TOKEN")

	req.Header.Add("AccessToken", token)
	return req, nil
}
