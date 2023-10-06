package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type AuthRequest struct {
	Appid     string `json:"appid"`
	Time      int64  `json:"time"`
	Signature string `json:"signature"`
}

type AuthResponse struct {
	Code        int    `json:"code"`
	Expiresin   int    `json:"expiresIn"`
	AccessToken string `json:"accessToken"`
	Result      string `json:"result"`
}

const serviceURL string = "https://open.iopgps.com/api/auth"

func getAccessToken() (string, error) {
	appid := os.Getenv("APPID")
	loginKey := os.Getenv("LOGIN_KEY")
	token, createdAt, err := readToken()
	if err != nil {
		return "", err
	}

	if token != "" && time.Now().Before(createdAt.Add(2*time.Hour-20*time.Minute)) {
		return token, nil
	}

	timeNow := time.Now()
	currentTime := timeNow.Unix()
	signature := generateSignature(loginKey, currentTime)

	authRequest := AuthRequest{
		Appid:     appid,
		Time:      currentTime,
		Signature: signature,
	}

	authRequestBody, err := json.Marshal(authRequest)
	if err != nil {
		return "", err
	}

	response, err := http.Post(serviceURL, "application/json", bytes.NewBuffer(authRequestBody))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var authResponse AuthResponse
	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		return "", err
	}

	if authResponse.Code == 0 {
		fmt.Println("Autenticación exitosa. Access Token:", authResponse.AccessToken)
	} else {
		fmt.Println("Error al autenticar. Código:", authResponse.Code, "Resultado:", authResponse.Result)
		return "", fmt.Errorf("error al autenticar: %s", authResponse.Result)
	}

	writeToken(authResponse.AccessToken, timeNow)

	return authResponse.AccessToken, nil
}

func readToken() (string, time.Time, error) {
	token := os.Getenv("ACCESS_TOKEN")
	timeline := os.Getenv("TIME")

	unixTime, err := strconv.ParseInt(timeline, 10, 64)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("formato de archivo inválido")
	}
	time := time.Unix(unixTime, 0)

	return token, time, nil
}

func writeToken(token string, createdAt time.Time) {
	time := strconv.FormatInt(createdAt.Unix(), 10)
	os.Setenv("ACCESS_TOKEN", token)
	os.Setenv("TIME", time)
	envMap := map[string]string{
		"ACCESS_TOKEN": token,
		"TIME":         time,
	}

	err := godotenv.Write(envMap, ".env")
	if err != nil {
		logrus.Fatal("Error writing to .env file")
	}
}

func generateSignature(loginKey string, time int64) string {
	hasher := md5.New()
	hasher.Write([]byte(loginKey))
	loginKeyHash := hex.EncodeToString(hasher.Sum(nil))

	hasher.Reset()
	hasher.Write([]byte(loginKeyHash + strconv.FormatInt(time, 10)))
	signature := hex.EncodeToString(hasher.Sum(nil))

	return signature
}
