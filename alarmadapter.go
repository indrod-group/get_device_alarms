package main

import (
	"strconv"
)

type AlarmDataAdapter interface {
	ToAlarmRequest(data AlarmData) Alarm
}

type AlarmDataAdapterImpl struct{}

func (adapter AlarmDataAdapterImpl) ToAlarmRequest(data AlarmData) Alarm {
	return Alarm{
		Imei:         data.Imei,
		PositionType: data.PositionType,
		Lat:          data.Lat,
		Lng:          data.Lng,
		Time:         data.Time,
		AlarmCode:    data.AlarmCode,
		AlarmType:    data.AlarmType,
		Course:       data.Course,
		DeviceType:   data.DeviceType,
		Speed:        data.Speed,
	}
}

func ConvertAlarmDataToRequest(data AlarmData) Alarm {
	adapter := AlarmDataAdapterImpl{}
	return adapter.ToAlarmRequest(data)
}

type WhatsGPSAlarmDataAdapter interface {
	WhatsGPSToAlarmRequest(data WhatsGPSAlarm) Alarm
}

type WhatsGPSAlarmDataAdapterImpl struct{}

func getAlarmCode(alarmType int64) string {
	switch alarmType {
	case 1:
		return "SHAKE"
	case 2:
		return "POWEROFF"
	case 3:
		return "LOWVOT"
	case 4:
		return "SOS"
	case 5:
		return "OVERSPEED"
	case 6:
		return "FENCEOUT"
	case 7:
		return "REMOVE"
	case 8:
		return "LOWVOT"
	case 9:
		return "AREAOUT"
	case 10:
		return "REMOVE"
	case 11:
		return "REMOVE"
	case 12:
		return "MAGNETISM"
	case 13:
		return "REMOVECONTINUOUSLY"
	case 14:
		return "BLUETOOTH"
	case 15:
		return "SIGNALSHIELDING"
	case 16:
		return "PSEUDOBASESTATION"
	case 17:
		return "FENCEIN"
	case 18:
		return "FENCEIN"
	case 19:
		return "FENCEOUT"
	case 31:
		return "ACCON"
	case 32:
		return "ACCOFF"
	}
	return "UNKNOWN"
}

func getAlarmType(alarmType int64) int64 {
	switch alarmType {
	case 1:
		return 3
	case 2:
		return 92
	case 3:
		return 2
	case 4:
		return 99
	case 5:
		return 12
	case 6:
		return 16
	case 7:
		return 1
	case 8:
		return 2
	case 9:
		return 18
	case 10:
		return 11
	case 11:
		return 10
	case 12:
		return 103
	case 13:
		return 6
	case 14:
		return 102
	case 15:
		return 101
	case 16:
		return 14
	case 17:
		return 17
	case 18:
		return 17
	case 19:
		return 16
	case 31:
		return 44
	case 32:
		return 45
	}
	return alarmType
}

func getLatitude(lat float64) string {
	return strconv.FormatFloat(lat, 'f', 7, 64)
}

func getLongitude(lat float64) string {
	return strconv.FormatFloat(lat, 'f', 7, 64)
}

func (adapter WhatsGPSAlarmDataAdapterImpl) WhatsGPSToAlarmRequest(data WhatsGPSAlarm) Alarm {
	positionType := "GPS"
	lat := getLatitude(data.Lat)
	lon := getLongitude(data.Lon)
	alarmCode := getAlarmCode(data.AlarmType)
	alarmType := getAlarmType(data.AlarmType)
	course := int64(0)
	deviceType := int64(1)
	speed := data.Speed

	return Alarm{
		Imei:         strconv.FormatInt(data.CarID, 10),
		PositionType: &positionType,
		Lat:          &lat,
		Lng:          &lon,
		Time:         data.AlarmTime.Unix(),
		AlarmCode:    alarmCode,
		AlarmType:    alarmType,
		Course:       &course,
		DeviceType:   deviceType,
		Speed:        &speed,
	}
}

func ConvertWhatsGPSAlarmDataToRequest(data WhatsGPSAlarm) Alarm {
	adapter := WhatsGPSAlarmDataAdapterImpl{}
	return adapter.WhatsGPSToAlarmRequest(data)
}
