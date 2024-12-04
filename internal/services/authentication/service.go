package authentication

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/nathan-fiscaletti/coattail-go/internal/keys"
)

var (
	ErrAuthenticationNotFound = errors.New("authentication service not found in context")
	ErrInvalidToken           = errors.New("invalid token")
	ErrInvalidSource          = errors.New("invalid source")
	ErrInvalidSignature       = errors.New("invalid signature")
)

type Service struct {
	secretKey []byte
}

func newService() (*Service, error) {
	return &Service{
		secretKey: []byte("secret"),
	}, nil
}

func ContextWithService(ctx context.Context) (context.Context, error) {
	auth, err := newService()
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, keys.AuthenticationKey, auth), nil
}

func GetService(ctx context.Context) (*Service, error) {
	auth, ok := ctx.Value(keys.AuthenticationKey).(*Service)
	if !ok {
		return nil, ErrAuthenticationNotFound
	}

	return auth, nil
}

func (s *Service) Issue(ipNet net.IPNet) (string, error) {
	h := hmac.New(sha256.New, s.secretKey)
	h.Write([]byte(ipNet.String()))
	signature := h.Sum(nil)

	token := fmt.Sprintf("%s;%s", ipNet.String(), base64.StdEncoding.EncodeToString(signature))
	return token, nil
}

func (s *Service) Authenticate(token string, source net.IP) (bool, error) {
	parts := strings.Split(token, ";")
	if len(parts) != 2 {
		return false, ErrInvalidToken
	}

	_, ipNet, err := net.ParseCIDR(parts[0])
	if err != nil {
		return false, ErrInvalidToken
	}

	if !ipNet.Contains(source) {
		return false, ErrInvalidSource
	}

	signatureB64 := parts[1]
	signature, err := base64.StdEncoding.DecodeString(signatureB64)
	if err != nil {
		return false, ErrInvalidSignature
	}

	h := hmac.New(sha256.New, s.secretKey)
	h.Write([]byte(ipNet.String()))
	expectedSignature := h.Sum(nil)

	if !hmac.Equal(expectedSignature, []byte(signature)) {
		return false, ErrInvalidSignature
	}

	return true, nil
}
