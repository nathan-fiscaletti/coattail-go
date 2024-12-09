package authentication

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/nathan-fiscaletti/coattail-go/internal/keys"
	"github.com/nathan-fiscaletti/coattail-go/internal/logging"
)

const (
	secretKeyFile = "secret.key"
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

func newService(ctx context.Context) (*Service, error) {
	service := &Service{}

	// Load or generate secret key
	err := service.loadSecretKey(ctx)
	if err != nil {
		return nil, err
	}

	// Issue a single token
	// TODO: Remove this, for debugging only.
	_, ipnet, err := net.ParseCIDR("127.0.0.1/32")
	if err != nil {
		return nil, err
	}
	token, err := service.Issue(ctx, ipnet)
	if err != nil {
		return nil, err
	}
	if logger, err := logging.GetLogger(ctx); err == nil {
		logger.Printf("created dev token: %s", token)
	}

	return service, nil
}

func ContextWithService(ctx context.Context) (context.Context, error) {
	auth, err := newService(ctx)
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

func (s *Service) Issue(ctx context.Context, ipNet *net.IPNet) (string, error) {
	h := hmac.New(sha256.New, s.secretKey)
	h.Write([]byte(ipNet.String()))
	signature := h.Sum(nil)

	token := fmt.Sprintf("%s;%s", ipNet.String(), base64.StdEncoding.EncodeToString(signature))
	return token, nil
}

func (s *Service) Authenticate(ctx context.Context, token string, source net.IP) (bool, error) {
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

func (s *Service) loadSecretKey(ctx context.Context) error {
	// check if the `secret.key` file exists
	if _, err := os.Stat(secretKeyFile); err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		// if the file doesn't exist, generate a new key
		key, err := s.generateKey(2048)
		if err != nil {
			return err
		}

		// write the key to a file
		secretFile, err := os.Create(secretKeyFile)
		if err != nil {
			return err
		}
		defer secretFile.Close()

		_, err = secretFile.Write(key)
		if err != nil {
			return err
		}
	}

	// Load the secret key
	secretFile, err := os.Open(secretKeyFile)
	if err != nil {
		return err
	}

	s.secretKey, err = io.ReadAll(secretFile)
	if err != nil {
		return err
	}

	return nil
}

// GenerateHMACKey generates a random key of the specified length for use with HMAC
func (s *Service) generateKey(length int) ([]byte, error) {
	// Ensure the key length is valid
	if length <= 0 {
		return nil, fmt.Errorf("key length must be greater than 0")
	}

	// Create a byte slice to hold the key
	key := make([]byte, length)

	// Fill the byte slice with secure random bytes
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	return key, nil
}
