package main

import (
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
)

type RequestGenerator struct {
	next Handler
}

const MAX_DEVICES_FOR_UPDATE = 22

func (rg *RequestGenerator) Handle(data interface{}) (interface{}, error) {
	devices, ok := data.([]Device)
	if !ok {
		logrus.Println("Error: Unable to cast data to []Device")
		return nil, errors.New("unable to cast data to []Device")
	}

	urls := make([]string, len(devices))
	var wg sync.WaitGroup

	sem := make(chan struct{}, MAX_DEVICES_FOR_UPDATE)

	for i, device := range devices {
		wg.Add(1)
		go func(i int, device Device) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			urls[i] = device.GenerateURL()
			device.UpdateDevice()
		}(i, device)
	}

	wg.Wait()

	if rg.next != nil {
		return rg.next.Handle(urls)
	}
	return urls, nil
}

func (rg *RequestGenerator) SetNext(next Handler) {
	rg.next = next
}
