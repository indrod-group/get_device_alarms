package main

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Director struct {
	first Handler
}

var directorInstance *Director
var once sync.Once

func GetDirectorInstance() *Director {
	once.Do(func() {
		directorInstance = &Director{}
	})
	return directorInstance
}

func (d *Director) BuildChain() {
	deviceController := &DeviceController{}
	requestGenerator := &RequestGenerator{}
	requestExecutor := &RequestExecutor{}
	dataSaver := &DataSaver{}
	messageSender := &MessageSender{}

	deviceController.SetNext(requestGenerator)
	requestGenerator.SetNext(requestExecutor)
	requestExecutor.SetNext(dataSaver)
	dataSaver.SetNext(messageSender)

	d.first = deviceController
}

func (d *Director) ProcessRequest(request interface{}) (interface{}, error) {
	if d.first != nil {
		return d.first.Handle(request)
	}
	return request, nil
}

func InitiateTrackingAlarms() {
	director := GetDirectorInstance()
	director.BuildChain()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		queryParams := map[string]string{
			"is_tracking_alarms": "true",
		}
		_, err := director.ProcessRequest(queryParams)
		if err != nil {
			logrus.Println(err)
		}
		<-ticker.C
	}
}
