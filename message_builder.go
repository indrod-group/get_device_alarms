package main

import (
	"fmt"
	"time"
)

type MessageBuilder struct {
	user        *User
	messageType string
	time        int64
}

func NewMessageBuilder(user *User, messageType string, time int64) *MessageBuilder {
	return &MessageBuilder{
		user:        user,
		messageType: messageType,
		time:        time,
	}
}

func (mb *MessageBuilder) BuildMessage() string {
	localTime := unixToLocal(mb.time)
	carOwner, licenseNumber, vin := mb.getUserDetails()
	alert := mb.getAlert()
	message := fmt.Sprintf("%s\nDatos del usuario:\nUsuario: %s", alert, mb.user.UserName)
	message += mb.addDetail("Propietario", carOwner)
	message += mb.addDetail("Placa del vehículo", licenseNumber)
	message += mb.addDetail("Vin", vin)
	message += fmt.Sprintf("\nHora de alarma: %s", localTime)
	return message
}

func unixToLocal(unixTime int64) time.Time {
	loc, _ := time.LoadLocation("America/Guayaquil")
	return time.Unix(unixTime, 0).In(loc)
}

func (mb *MessageBuilder) getUserDetails() (carOwner, licenseNumber, vin string) {
	if mb.user.CarOwner != nil {
		carOwner = *mb.user.CarOwner
	}
	if mb.user.LicenseNumber != nil {
		licenseNumber = *mb.user.LicenseNumber
	}
	if mb.user.Vin != nil {
		vin = *mb.user.Vin
	}
	if licenseNumber == vin {
		vin = "" // Si Vin es igual a LicenseNumber, se omite Vin
	}
	return carOwner, licenseNumber, vin
}

func (mb *MessageBuilder) getAlert() string {
	switch mb.messageType {
	case "SOS":
		return "🚨🚨 ALERTA DE SOS 🚨🚨"
	case "REMOVE":
		return "⚡⚡ ALERTA DE CORTE DE CORRIENTE ⚡⚡"
	case "LOWVOT":
		return "⚡⚡ ALERTA DE CORRIENTE BAJA ⚡⚡"
	default:
		return "ALERTA DESCONOCIDA"
	}
}

func (mb *MessageBuilder) addDetail(label, value string) string {
	if value != "" {
		return fmt.Sprintf("\n%s: %s", label, value)
	}
	return ""
}
