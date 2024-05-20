package main

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Director manages a chain of responsibility pattern for handling requests.
type Director struct {
	first Handler // first points to the first handler in the chain.
}

// directorInstance holds a singleton instance of Director.
var directorInstance *Director
var once sync.Once

// GetDirectorInstance returns a singleton instance of Director, creating it if necessary.
func GetDirectorInstance() *Director {
	once.Do(func() {
		directorInstance = &Director{}
	})
	return directorInstance
}

// BuildChain constructs the chain of responsibility for handling requests.
func (d *Director) BuildChain() {
	deviceController := &DeviceController{}
	requestGenerator := &RequestGenerator{}
	requestExecutor := &RequestExecutor{}
	dataSaver := &DataSaver{}
	messageSender := &MessageSender{}

	// Sets the next handler for each component in the chain.
	deviceController.SetNext(requestGenerator)
	requestGenerator.SetNext(requestExecutor)
	requestExecutor.SetNext(dataSaver)
	dataSaver.SetNext(messageSender)

	d.first = deviceController
}

// ProcessRequest processes a request through the chain of handlers.
// It returns an error if any handler fails.
func (d *Director) ProcessRequest(request interface{}) (interface{}, error) {
	if d.first != nil {
		return d.first.Handle(request)
	}
	return request, nil
}

var trackingAlarmsStarted bool    // trackingAlarmsStarted indicates whether alarm tracking has started.
var trackingAlarmsLock sync.Mutex // trackingAlarmsLock provides a mutex for controlling access to trackingAlarmsStarted.

// InitiateTrackingAlarms starts the alarm tracking process, ensuring it runs only once.
func InitiateTrackingAlarms() {
	trackingAlarmsLock.Lock()
	defer trackingAlarmsLock.Unlock()

	if trackingAlarmsStarted {
		return // If already started, do nothing.
	}

	trackingAlarmsStarted = true
	director := GetDirectorInstance()
	director.BuildChain()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	running := make(chan bool, 1) // running controls the execution overlap of process invocations.
	for range ticker.C {
		select {
		case running <- true:
			go func() {
				defer func() { <-running }()
				queryParams := map[string]string{"is_tracking_alarms": "true"}
				_, err := director.ProcessRequest(queryParams)
				if err != nil {
					logrus.Println(err)
				}
			}()
		default:
			// Skips this tick if the previous process is still running.
			logrus.Warn("Skipping this tick as the previous one is still processing")
		}
	}
}
