package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var DEBUG = os.Getenv("DEBUG") != "RENDER"

func initLog() {
	if !DEBUG {
		log.SetOutput(os.Stdout)
	} else {
		// Crear la carpeta logs si no existe
		if _, err := os.Stat("logs"); os.IsNotExist(err) {
			os.Mkdir("logs", 0755)
		}

		// Crear un nombre de archivo con la fecha y hora actual
		t := time.Now()
		fileName := fmt.Sprintf(
			"logs/log-%04d%02d%02d-%02d%02d%02d.log",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second(),
		)

		// Abrir el archivo de log
		logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening log file: %v", err)
			return
		}

		// Configurar el paquete log para escribir en el archivo de log
		log.SetOutput(logFile)
	}
}

func processUsers(users []User, sem chan bool, wg *sync.WaitGroup) {
	for _, user := range users {
		wg.Add(1)
		go func(user User) {
			sem <- true
			defer func() { <-sem }()
			cronJob := NewCronJob(user, time.Now().Unix(), 60)
			cronJob.Run(wg)
		}(user)
	}
}

func main() {
	initLog()

	var token string
	var err error

	tokenTicker := time.NewTicker(10 * time.Minute)
	go func() {
		for {
			token, err = getAccessToken()
			if err != nil {
				log.Println("Error al obtener el token de acceso:", err)
			} else {
				log.Println("Token de acceso actualizado:", token)
			}
			<-tokenTicker.C
		}
	}()

	time.Sleep(5 * time.Second)

	sem := make(chan bool, 10)
	var wg sync.WaitGroup

	users := GetUserFromApi()
	if users == nil {
		log.Println("Error al obtener los usuarios")
	}

	userTicker := time.NewTicker(5 * time.Minute)
	go func() {
		for range userTicker.C {
			users = GetUserFromApi()
			if users == nil {
				log.Println("Error al obtener los usuarios")
			}
		}
	}()

	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C {
			if users == nil {
				continue
			}
			for i := 0; i < len(users); i += 9 {
				end := i + 9
				if end > len(users) {
					end = len(users)
				}
				processUsers(users[i:end], sem, &wg)
				wg.Wait()
			}
		}
	}()

	select {} // Bucle infinito para evitar que main() termine
}
