package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func UnmarshalWhatsGPSAlarm(data []byte) (Alarm, error) {
	var r Alarm
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *WhatsGPSAlarmData) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type WhatsGPSAlarmData struct {
	Data  []WhatsGPSAlarm `json:"data"`
	Ret   int64           `json:"ret"`
	Total int64           `json:"total"`
}

type WhatsGPSAlarm struct {
	AlarmTime CustomTime `json:"alarmTime"`
	AlarmType int64      `json:"alarmType"`
	CarID     int64      `json:"carId"`
	Dir       int64      `json:"dir"`
	Lat       float64    `json:"lat"`
	Latc      float64    `json:"latc"`
	Lon       float64    `json:"lon"`
	Lonc      float64    `json:"lonc"`
	PointTime CustomTime `json:"pointTime"`
	PointType int64      `json:"pointType"`
	Remark    string     `json:"remark"`
	Speed     int64      `json:"speed"`
	UserName  string     `json:"userName"`
}

type CustomTime struct {
	time.Time
}

const ctLayout = "2006-01-02 15:04:05"

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(ctLayout, s)
	return
}

func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.UnixNano() == (time.Time{}).UnixNano() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(ctLayout))), nil
}
