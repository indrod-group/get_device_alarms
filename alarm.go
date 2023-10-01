package main

type AlarmData struct {
	Code         int64   `json:"code"`
	PositionType *string `json:"positionType,omitempty"`
	Imei         string  `json:"imei"`
	Lat          *string `json:"lat,omitempty"`
	Lng          *string `json:"lng,omitempty"`
	Time         int64   `json:"time"`
	Speed        *int64  `json:"speed,omitempty"`
	Course       *int64  `json:"course,omitempty"`
	AlarmCode    string  `json:"alarmCode"`
	AlarmTime    int64   `json:"alarmTime"`
	DeviceType   int64   `json:"deviceType"`
	AlarmType    int64   `json:"alarmType"`
}

type AlarmDataForPost struct {
	Code         int64   `json:"code"`
	PositionType *string `json:"position_type,omitempty"`
	Imei         string  `json:"imei"`
	Lat          *string `json:"lat,omitempty"`
	Lng          *string `json:"lng,omitempty"`
	Time         int64   `json:"time"`
	Speed        *int64  `json:"speed,omitempty"`
	Course       *int64  `json:"course,omitempty"`
	AlarmCode    string  `json:"alarm_code"`
	AlarmTime    int64   `json:"alarm_time"`
	DeviceType   int64   `json:"device_type"`
	AlarmType    int64   `json:"alarm_type"`
}

type ApiResponse struct {
	Code    int64       `json:"code"`
	Details []AlarmData `json:"details"`
}
