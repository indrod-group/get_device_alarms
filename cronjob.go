package main

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type CronJob struct {
	user        User
	currentTime int64
	interval    int64
}

func NewCronJob(user User, startTime int64, interval int64) *CronJob {
	return &CronJob{
		user:        user,
		currentTime: startTime,
		interval:    interval,
	}
}

func (c *CronJob) Start(sem *sync.WaitGroup) {
	defer sem.Done()

	apiData, err := GetAlarmData(c.user, c.currentTime, c.interval)
	if err != nil {
		logrus.WithError(err).Error("Error getting alarm data")
		return
	}

	err = ProcessAlarmData(c.user, apiData)
	if err != nil {
		logrus.WithError(err).Error("Error processing alarm data")
		return
	}

	c.currentTime = time.Now().Unix()
}

func ProcessUsers(users []User, sem chan bool, wg *sync.WaitGroup) {
	for _, user := range users {
		wg.Add(1)
		go func(user User) {
			sem <- true
			defer func() { <-sem }()
			cronJob := NewCronJob(user, time.Now().Unix(), 60)
			cronJob.Start(wg)
		}(user)
	}
}
