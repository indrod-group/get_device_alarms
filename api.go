package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var client = &http.Client{}

func getNewUrl(imei string, startTime int64) string {
	endTime := time.Now().Unix() // Actualiza el endTime con el tiempo actual
	return fmt.Sprintf("https://open.iopgps.com/api/device/alarm?imei=%s&startTime=%d&endTime=%d", imei, startTime, endTime)
}

func createRequest(url, accessToken string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request for URL %s: %w", url, err)
	}
	req.Header.Add("AccessToken", accessToken)
	return req, nil
}

func doRequestWithRetry(req *http.Request, imei string, startTime int64, maxRetries int, baseDelay time.Duration) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i < maxRetries; i++ { // Número de intentos
		resp, err = client.Do(req)
		if err != nil {
			if strings.Contains(err.Error(), "i/o timeout") { // Comprueba si el error es un tiempo de espera
				log.Printf("Error making request to URL %s: %s\n", req.URL, err)
				delay := baseDelay * time.Duration(math.Pow(2, float64(i))) // Tiempo de espera antes del próximo intento
				time.Sleep(delay)
				url := getNewUrl(imei, startTime)                            // Obtiene una nueva URL con el tiempo actualizado
				req, err = createRequest(url, req.Header.Get("AccessToken")) // Crea una nueva solicitud con la nueva URL
				if err != nil {
					return nil, fmt.Errorf("error creating request for URL %s: %w", url, err)
				}
				continue
			}
			return nil, fmt.Errorf("error making request to URL %s: %w", req.URL, err) // Si el error no es un tiempo de espera, falla inmediatamente
		}
		break
	}

	if err != nil {
		return nil, fmt.Errorf("error making request to URL %s after retries: %w", req.URL, err)
	}

	return resp, nil
}

func GetAlarmData(user User, currentTime, interval int64) ([]byte, error) {
	accessToken := os.Getenv("ACCESS_TOKEN")
	imei := user.Imei
	startTime := currentTime - interval
	url := getNewUrl(imei, startTime)

	req, err := createRequest(url, accessToken)
	if err != nil {
		return nil, err
	}

	resp, err := doRequestWithRetry(req, imei, startTime, 7, 1*time.Second)
	if err != nil {
		return nil, err
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
