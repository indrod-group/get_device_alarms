package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func getNewUrl(imei string, startTime int64) string {
	endTime := time.Now().Unix() // Actualiza el endTime con el tiempo actual
	return fmt.Sprintf("https://open.iopgps.com/api/device/alarm?imei=%s&startTime=%d&endTime=%d", imei, startTime, endTime)
}

func createRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request for URL %s: %w", url, err)
	}
	req.Header.Add("AccessToken", app.accessToken)
	return req, nil
}

var client = &http.Client{
	Timeout: time.Second * 10,
}

func doRequestWithRetry(req *http.Request, imei string, startTime int64, maxRetries int, baseDelay time.Duration) (*http.Response, error) {
	for i := 0; i < maxRetries; i++ {
		resp, err := client.Do(req)
		if err != nil {
			if resp != nil {
				resp.Body.Close()
			}
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				logrus.Errorf("Timeout Error making request to URL %s: %s\n", req.URL, err)
				delay := baseDelay * time.Duration(math.Pow(2, float64(i)))
				jitter := time.Duration(rand.Int63n(int64(delay)))
				sleepTime := delay + jitter
				time.Sleep(sleepTime)
				uri := getNewUrl(imei, startTime)
				req.URL, err = url.Parse(uri)
				if err != nil {
					return nil, fmt.Errorf("error creating request for URL %s: %w", uri, err)
				}
				continue
			}
			return nil, fmt.Errorf("error making request to URL %s: %w", req.URL, err)
		}
		return resp, nil
	}
	return nil, fmt.Errorf("error making request to URL %s after retries", req.URL)
}

func readResponseBody(resp *http.Response) ([]byte, error) {
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logrus.Errorf("error closing response body from URL %s: %s\n", resp.Request.URL, closeErr)
		}
	}()

	var buf bytes.Buffer
	_, err := io.Copy(&buf, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body from URL %s: %w", resp.Request.URL, err)
	}
	body := buf.Bytes()
	return body, nil
}

func GetAlarmData(user User, currentTime, interval int64) ([]byte, error) {
	imei := user.Imei
	startTime := currentTime - interval
	uri := getNewUrl(imei, startTime)

	req, err := createRequest(uri)
	if err != nil {
		return nil, err
	}

	resp, respErr := doRequestWithRetry(req, imei, startTime, 7, 1*time.Second)
	if respErr != nil {
		return nil, respErr
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, err
	}

	body, err := readResponseBody(resp)
	if len(body) == 0 {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"status_code": resp.StatusCode,
		"imei":        imei,
	}).Info("Successfully retrieved alarm data from IOGPS API")

	return body, nil
}

func saveAlarmInAPI(detail AlarmData) error {
	detailForPost := AlarmDataForPost(detail)
	jsonData, err := json.Marshal(detailForPost)
	if err != nil {
		return fmt.Errorf("error marshalling data for detail %+v: %w", detailForPost, err)
	}

	reqURL := app.config.acvApiURL + "details/"
	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request for URL %s: %w", reqURL, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Token "+app.config.authToken)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request to URL %s: %w", reqURL, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logrus.Errorf("Error closing response body from URL %s: %s\n", reqURL, closeErr)
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ACV API Error | Status code: %d for URL %s", resp.StatusCode, reqURL)
	}

	logrus.Printf("ACV API | Status code: %d\n", resp.StatusCode)
	return nil
}

func ProcessAlarmData(user User, data []byte) error {
	if len(data) == 0 {
		logrus.Warning("No data to process\n")
		return nil
	}

	var details ApiResponse
	err := json.Unmarshal(data, &details)
	if err != nil {
		return fmt.Errorf("error unmarshalling data: %w", err)
	}

	var wgAlarms sync.WaitGroup
	errChan := make(chan error, 1) // A channel to hold the first error we get

	for _, detail := range details.Details {
		wgAlarms.Add(1)
		go func(detail AlarmData) {
			defer wgAlarms.Done()
			CheckAndSendAlarm(user, detail)
		}(detail)

		wgAlarms.Add(1)
		go func(detail AlarmData) {
			defer wgAlarms.Done()
			if err := saveAlarmInAPI(detail); err != nil {
				select {
				case errChan <- err:
				default:
				}
			}
		}(detail)
	}

	wgAlarms.Wait()

	select {
	case err := <-errChan:
		return fmt.Errorf("error saving data: %w", err)
	default:
		return nil
	}
}
