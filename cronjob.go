package main

import (
	"log"
	"sync"
	"time"
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
		log.Printf("Error getting alarm data: %s\n", err)
		return
	}

	err = ProcessAlarmData(c.user, apiData)
	if err != nil {
		log.Printf("Error processing alarm data: %s\n", err)
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
