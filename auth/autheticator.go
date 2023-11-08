package auth

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

type Authenticator struct {
	accessToken string
	appID       string
	loginKey    string
	serviceURL  string
}

func InitAuthenticator() *Authenticator {
	return &Authenticator{
		accessToken: os.Getenv("ACCESS_TOKEN"),
		appID:       os.Getenv("APPID"),
		loginKey:    os.Getenv("LOGIN_KEY"),
		serviceURL:  os.Getenv("MY_API_URL"),
	}
}

func (a *Authenticator) createRequest() AuthRequest {
	timeNow := time.Now()
	currentTime := timeNow.Unix()
	signature := a.generateSignature(currentTime)

	authRequest := AuthRequest{
		Appid:     a.appID,
		Time:      currentTime,
		Signature: signature,
	}

	return authRequest
}

func (a *Authenticator) GetAccessToken() (string, error) {
	token, createdAt, err := a.readToken()
	if err != nil {
		return "", err
	}

	if token != "" && time.Now().Before(createdAt.Add(2*time.Hour-20*time.Minute)) {
		return token, nil
	}

	authRequest := a.createRequest()

	response, err := a.sendAuthRequest(authRequest)
	if err != nil {
		return "", err
	}

	authResponse, err := a.parseAuthResponse(response)
	if err != nil {
		return "", err
	}

	if authResponse.AccessToken != nil {
		a.writeToken(*authResponse.AccessToken, time.Now())
	}

	return a.accessToken, nil
}

func (a *Authenticator) sendAuthRequest(authRequest AuthRequest) (*http.Response, error) {
	authRequestBody, err := json.Marshal(authRequest)
	if err != nil {
		return nil, err
	}

	response, err := http.Post(a.serviceURL, "application/json", bytes.NewBuffer(authRequestBody))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return response, nil
}

func (a *Authenticator) parseAuthResponse(response *http.Response) (*AuthResponse, error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var authResponse AuthResponse
	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		return nil, err
	}

	if authResponse.Code == 0 {
		fmt.Println("Autenticación exitosa. Access Token:", authResponse.AccessToken)
	} else {
		fmt.Println("Error al autenticar. Código:", authResponse.Code, "Resultado:", authResponse.Result)
		return nil, fmt.Errorf("error al autenticar: %s", *authResponse.Result)
	}

	return &authResponse, nil
}

func (a *Authenticator) readToken() (string, time.Time, error) {
	token := os.Getenv("ACCESS_TOKEN")
	timeline := os.Getenv("TIME")

	unixTime, err := strconv.ParseInt(timeline, 10, 64)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("formato de archivo inválido")
	}
	time := time.Unix(unixTime, 0)

	return token, time, nil
}

func setTokenInEnv(token string, time string) {
	os.Setenv("ACCESS_TOKEN", token)
	os.Setenv("TIME", time)
}

func (a *Authenticator) writeToken(token string, createdAt time.Time) {
	a.accessToken = token
	time := strconv.FormatInt(createdAt.Unix(), 10)

	setTokenInEnv(token, time)

	envMap := map[string]string{
		"ACCESS_TOKEN": token,
		"TIME":         time,
	}

	err := godotenv.Write(envMap, ".env")
	if err != nil {
		logrus.Fatal("Error writing to .env file")
	}
}

func (a *Authenticator) generateSignature(time int64) string {
	hasher := md5.New()
	hasher.Write([]byte(a.loginKey))
	loginKeyHash := hex.EncodeToString(hasher.Sum(nil))
	hasher.Reset()
	hasher.Write([]byte(loginKeyHash + strconv.FormatInt(time, 10)))
	signature := hex.EncodeToString(hasher.Sum(nil))
	return signature
}
