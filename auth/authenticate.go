package auth

type Authenticate interface {
	GetAccessToken() string
	InitiateTokenRenewal(stopChan <-chan struct{})
}
