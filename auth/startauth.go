package auth

import (
	"time"

	"github.com/sirupsen/logrus"
)

func (a *Authenticator) InitiateTokenRenewal() {
	tokenTicker := time.NewTicker(10 * time.Minute)
	defer tokenTicker.Stop()
	for {
		accessToken, err := a.GetAccessToken()
		if err != nil {
			logrus.WithError(err).Error("Error al obtener el token de acceso")
			continue
		}
		logrus.Println("Token de acceso actualizado:", accessToken)
		<-tokenTicker.C
	}
}
