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

	deviceController.SetNext(requestGenerator)
	requestGenerator.SetNext(requestExecutor)
	requestExecutor.SetNext(dataSaver)

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
		request := "None"
		_, err := director.ProcessRequest(request)
		if err != nil {
			logrus.Println(err)
		}
		<-ticker.C
	}
}
