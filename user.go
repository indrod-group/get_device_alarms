package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type User struct {
	ID            int64   `json:"id"`
	AccountID     int64   `json:"account_id"`
	AccountName   string  `json:"account_name"`
	UserName      string  `json:"user_name"`
	Imei          string  `json:"imei"`
	LicenseNumber *string `json:"license_number,omitempty"`
	Vin           *string `json:"vin,omitempty"`
	CarOwner      *string `json:"car_owner,omitempty"`
	IsTracking    bool    `json:"is_tracking"`
}

func GetUserFromApi() []User {
	authToken := os.Getenv("AUTH_TOKEN")
	URL := os.Getenv("MY_API_URL")
	req, _ := http.NewRequest("GET", URL+"user/?is_tracking=True", nil)
	req.Header.Add("Authorization", "Token "+authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := io.ReadAll(resp.Body)
		var users []User
		json.Unmarshal(data, &users)
		return users
	}
	return nil
}

func (user *User) GetSOSMessage(time int64) string {
	localTime := UnixToLocal(time)
	carOwner := "desconocido"
	licenseNumber := "desconocido"
	if user.CarOwner != nil {
		carOwner = *user.CarOwner
	}
	if user.LicenseNumber != nil {
		licenseNumber = *user.LicenseNumber
	}
	if user.UserName == carOwner {
		return fmt.Sprintf("🚨🚨 ALERTA DE SOS 🚨🚨\nDatos del usuario:\nUsuario y propietario: %s\nPlaca del vehículo: %s\nHora de alarma: %s", user.UserName, licenseNumber, localTime)
	}
	return fmt.Sprintf("🚨🚨 ALERTA DE SOS 🚨🚨\nDatos del usuario:\nUsuario: %s\nPropietario: %s\nPlaca del vehículo: %s\nHora de alarma: %s", user.UserName, carOwner, licenseNumber, localTime)
}

func (user *User) GetRemoveMessage(time int64) string {
	localTime := UnixToLocal(time)
	carOwner := "desconocido"
	licenseNumber := "desconocido"
	if user.CarOwner != nil {
		carOwner = *user.CarOwner
	}
	if user.LicenseNumber != nil {
		licenseNumber = *user.LicenseNumber
	}
	if user.UserName == carOwner {
		return fmt.Sprintf("⚡⚡ ALERTA DE CORTE DE CORRIENTE ⚡⚡\nDatos del usuario:\nUsuario y propietario: %s\nPlaca del vehículo: %s\nHora de alarma: %s", user.UserName, licenseNumber, localTime)
	}
	return fmt.Sprintf("⚡⚡ ALERTA DE CORTE DE CORRIENTE ⚡⚡\nDatos del usuario:\nUsuario: %s\nPropietario: %s\nPlaca del vehículo: %s\nHora de alarma: %s", user.UserName, carOwner, licenseNumber, localTime)
}

func (user *User) GetLowvotMessage(time int64) string {
	localTime := UnixToLocal(time)
	carOwner := "desconocido"
	licenseNumber := "desconocido"
	if user.CarOwner != nil {
		carOwner = *user.CarOwner
	}
	if user.LicenseNumber != nil {
		licenseNumber = *user.LicenseNumber
	}
	if user.UserName == carOwner {
		return fmt.Sprintf("⚡⚡ ALERTA DE CORRIENTE BAJA ⚡⚡\nDatos del usuario:\nUsuario y propietario: %s\nPlaca del vehículo: %s\nHora de alarma: %s", user.UserName, licenseNumber, localTime)
	}
	return fmt.Sprintf("⚡⚡ ALERTA DE CORRIENTE BAJA ⚡⚡\nDatos del usuario:\nUsuario: %s\nPropietario: %s\nPlaca del vehículo: %s\nHora de alarma: %s", user.UserName, carOwner, licenseNumber, localTime)
}

func UnixToLocal(unixTime int64) time.Time {
	loc, _ := time.LoadLocation("America/Guayaquil")
	return time.Unix(unixTime, 0).In(loc)
}
