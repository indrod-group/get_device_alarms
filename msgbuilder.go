package main

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type MessageBuilder struct {
	device *Device
	alarm  *Alarm
}

func NewMessageBuilder(device *Device, alarm *Alarm) *MessageBuilder {
	return &MessageBuilder{
		device: device,
		alarm:  alarm,
	}
}

func unixToLocal(unixTime int64) (time.Time, error) {
	loc, err := time.LoadLocation("America/Guayaquil")
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load location: %w", err)
	}
	return time.Unix(unixTime, 0).In(loc), nil
}

func (mb *MessageBuilder) getCoordinates() (lat, lng string) {
	if mb.alarm.Lat != nil {
		lat = *mb.alarm.Lat
	}
	if mb.alarm.Lng != nil {
		lng = *mb.alarm.Lng
	}
	return lat, lng
}

func (mb *MessageBuilder) getAlarmAddress() string {
	const defaultLocation = "\nUbicación desconocida\n"
	lat, lng := mb.getCoordinates()
	if lat == "" || lng == "" {
		return defaultLocation
	}
	address := GetAddress(lat, lng)
	if address == nil {
		return defaultLocation
	}
	googleMapsLink := mb.getGoogleMapsLink()
	if googleMapsLink != "" {
		return mb.addDetail("Ubicación", *address) + "\nEnlace a Google Maps: " + googleMapsLink
	}
	return mb.addDetail("Ubicación", *address)
}

func (mb *MessageBuilder) getGoogleMapsLink() string {
	const googleMapsLinkBase = "https://www.google.com/maps/search/?api=1&query=%s,%s"
	lat, lng := mb.getCoordinates()
	if lat == "" || lng == "" {
		return ""
	}
	return fmt.Sprintf(googleMapsLinkBase, lat, lng)
}

func (mb *MessageBuilder) BuildMessage() string {
	localTime, err := unixToLocal(mb.alarm.Time)
	if err != nil {
		logrus.WithError(err).Error("Error converting unix time to local")
	}
	carOwner, licenseNumber, vin := mb.getUserDetails()
	alert := mb.getAlert()
	message := fmt.Sprintf("%s\nDatos del usuario:\nUsuario: %s", alert, mb.device.UserName)
	message += mb.addDetail("Propietario", carOwner)
	message += mb.addDetail("Placa del vehículo", licenseNumber)
	message += mb.addDetail("Vin", vin)
	message += fmt.Sprintf("\nHora de alarma: %s", localTime)
	message += mb.getAlarmAddress()
	return message
}

func (mb *MessageBuilder) getUserDetails() (carOwner, licenseNumber, vin string) {
	if mb.device.CarOwner != nil {
		carOwner = *mb.device.CarOwner
	}
	if mb.device.LicenseNumber != nil {
		licenseNumber = *mb.device.LicenseNumber
	}
	if mb.device.Vin != nil {
		vin = *mb.device.Vin
	}
	if licenseNumber == vin {
		vin = ""
	}
	return carOwner, licenseNumber, vin
}

func (mb *MessageBuilder) getAlert() string {
	switch mb.alarm.AlarmCode {
	case "SOS":
		return "🚨🚨 ALERTA DE SOS 🚨🚨"
	case "REMOVE":
		switch mb.alarm.AlarmType {
		case 1:
			return "🔧🔧 ALERTA DE DESMONTAJE 🔧🔧"
		case 10:
			return "💡💡 ALERTA DE SENSOR DE LUZ 💡💡"
		default:
			return "⚡⚡ ALERTA DE CORTE DE CORRIENTE ⚡⚡"
		}
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
