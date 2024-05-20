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

// Don't change this value by any reason.
const MAX_REQUESTS_IN_WWT_API_PER_SECOND = 5

// This function handles a slice of URLs by sending requests to each URL and collecting the alarm data from the responses.
// It limits the number of concurrent requests to MAX_REQUEST_PER_SECOND using a semaphore channel.
// It uses a mutex to protect the shared slice of alarm data and a wait group to synchronize the goroutines.
// It passes the collected alarm data to the next handler in the chain, if any, or returns it as the final result.
func (re *RequestExecutor) Handle(data interface{}) (interface{}, error) {
	urls, ok := data.([]string)
	if !ok {
		return nil, fmt.Errorf("RequestExecutor.Handle: expected []string, got %T", data)
	}

	var allAlarmData []AlarmData
	var mutex sync.Mutex
	var wg sync.WaitGroup

	sem := make(chan struct{}, MAX_REQUESTS_IN_WWT_API_PER_SECOND)

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
