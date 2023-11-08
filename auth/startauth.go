package auth

import (
	"time"

	"github.com/sirupsen/logrus"
)

func (a *Authenticator) InitiateTokenRenewal(stopChan <-chan struct{}) {
	tokenTicker := time.NewTicker(10 * time.Minute)
	defer tokenTicker.Stop()

	for {
		select {
		case <-tokenTicker.C:
			accessToken, err := a.GetAccessToken()
			if err != nil {
				logrus.WithError(err).Error("Error al obtener el token de acceso")
				continue
			}
			logrus.Println("Token de acceso actualizado:", accessToken)
		case <-stopChan:
			return
		}
	}
}
