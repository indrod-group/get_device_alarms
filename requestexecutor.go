package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

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
		go func(url string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			req, err := createRequest(url)
			if err != nil {
				logrus.Println(err)
				return
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				logrus.Println(err)
				return
			}
			defer resp.Body.Close()

			var alarmResponse AlarmResponse
			err = json.NewDecoder(resp.Body).Decode(&alarmResponse)
			if err != nil {
				logrus.Println(err)
				return
			}

			mutex.Lock()
			allAlarmData = append(allAlarmData, alarmResponse.Details...)
			mutex.Unlock()
		}(url)
	}

	wg.Wait()

	if re.next != nil {
		return re.next.Handle(allAlarmData)
	}
	return allAlarmData, nil
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
