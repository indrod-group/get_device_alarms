package main

import (
	"fmt"
	"sync"
)

type DataSaver struct {
	next Handler
}

func (ds *DataSaver) Handle(data interface{}) (interface{}, error) {
	alarmData, ok := data.([]AlarmData)
	if !ok {
		return nil, fmt.Errorf("DataSaver.Handle: expected []AlarmData, got %T", data)
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex
	var allAlarms []Alarm
	sem := make(chan struct{}, 50)

	for _, data := range alarmData {
		wg.Add(1)
		go func(data AlarmData) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			// Convert AlarmData to Alarm
			alarm := ConvertAlarmDataToRequest(data)

			// Create the alarm
			err := alarm.CreateAlarm()
			if err != nil {
				fmt.Println(err)
				return
			}

			mutex.Lock()
			allAlarms = append(allAlarms, alarm)
			mutex.Unlock()
		}(data)
	}

	wg.Wait()

	if ds.next != nil {
		return ds.next.Handle(allAlarms)
	}
	return allAlarms, nil
}

func (ds *DataSaver) SetNext(next Handler) {
	ds.next = next
}
