package logo

import (
	"fmt"
	"math/rand"
)

var _ LogoServicer = new(ServiceLogo)

type ServiceLogo struct {
	apiKey  string
	baseUri string
}

func NewServiceLogo(apiKey string, baseUri string) *ServiceLogo {
	return &ServiceLogo{
		apiKey:  apiKey,
		baseUri: baseUri,
	}
}

func (s *ServiceLogo) generateRandomString() string {
	//generate a random string
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, 8)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (s *ServiceLogo) GenerateLogo() (string, error) {
	return fmt.Sprintf("%s/%s", s.baseUri, s.generateRandomString()), nil
}
