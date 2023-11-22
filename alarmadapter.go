package main

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
