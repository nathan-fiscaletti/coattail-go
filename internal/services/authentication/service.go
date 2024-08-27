package authentication

import (
	"context"

	"github.com/nathan-fiscaletti/coattail-go/internal/services"
)

type Service struct{}

func newService() *Service {
	return &Service{}
}

func (s *Service) Authenticate(token string) string {
	return token
}

func ContextWithService(ctx context.Context) context.Context {
	return context.WithValue(ctx, services.AuthenticationServiceKey, newService())
}

func GetService(ctx context.Context) *Service {
	if v := ctx.Value(services.AuthenticationServiceKey); v != nil {
		if s, ok := v.(*Service); ok {
			return s
		}
	}

	return nil
}
