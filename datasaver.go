package main

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

type DataSaver struct {
	next Handler
}

const MAX_ALARMS_TO_REGISTER = 25

func (ds *DataSaver) Handle(data interface{}) (interface{}, error) {
	alarms, ok := data.([]Alarm)
	if !ok {
		return nil, fmt.Errorf("DataSaver.Handle: expected []AlarmData, got %T", data)
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, MAX_ALARMS_TO_REGISTER)

	for _, alarm := range alarms {
		wg.Add(1)
		go func(alarm Alarm) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			err := alarm.CreateAlarm()
			if err != nil {
				// Log the alarm information
				logrus.WithFields(logrus.Fields{
					"error": err,
					"alarm": alarm,
				}).Warning("Error saving the alarm")
				return
			}
		}(alarm)
	}

	wg.Wait()

	if ds.next != nil {
		return ds.next.Handle(alarms)
	}
	return alarms, nil
}

func (ds *DataSaver) SetNext(next Handler) {
	ds.next = next
}
