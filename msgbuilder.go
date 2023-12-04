package main

import (
	"fmt"
	"time"
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

func (mb *MessageBuilder) getAlarmAddress() string {
	const defaultLocation = "Ubicación desconocida\n"
	if mb.alarm.Lat == nil || mb.alarm.Lng == nil {
		return defaultLocation
	}
	if *mb.alarm.Lat == "" || *mb.alarm.Lng == "" {
		return defaultLocation
	}
	address := GetAddress(*mb.alarm.Lat, *mb.alarm.Lng)
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
	if mb.alarm.Lat == nil || mb.alarm.Lng == nil {
		return ""
	}
	if *mb.alarm.Lat == "" || *mb.alarm.Lng == "" {
		return ""
	}
	return fmt.Sprintf(googleMapsLinkBase, *mb.alarm.Lat, *mb.alarm.Lng)
}

func (mb *MessageBuilder) BuildMessage() string {
	localTime := unixToLocal(mb.alarm.Time)
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

func unixToLocal(unixTime int64) time.Time {
	loc, _ := time.LoadLocation("America/Guayaquil")
	return time.Unix(unixTime, 0).In(loc)
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
		case 11:
			return "⚡⚡ ALERTA DE CORTE DE CORRIENTE ⚡⚡"
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
