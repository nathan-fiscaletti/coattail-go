package authentication

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/nathan-fiscaletti/coattail-go/internal/keys"
)

const (
	secretKeyFile = "secret.key"
)

var (
	ErrAuthenticationNotFound = errors.New("authentication service not found in context")
	ErrInvalidToken           = errors.New("invalid token")
	ErrInvalidSource          = errors.New("invalid source")
	ErrInvalidSignature       = errors.New("invalid signature")
	ErrInvalidPermissions     = errors.New("invalid permissions")
	ErrTokenExpired           = errors.New("token expired")
)

type Service struct {
	secretKey []byte
}

func newService() (*Service, error) {
	service := &Service{}

	// Load or generate secret key
	err := service.loadSecretKey()
	if err != nil {
		return nil, err
	}

	return service, nil
}

// ContextWithService returns a context with the authentication service.
func ContextWithService(ctx context.Context) (context.Context, error) {
	auth, err := newService()
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, keys.AuthenticationKey, auth), nil
}

// GetService returns the authentication service from the context.
func GetService(ctx context.Context) (*Service, error) {
	auth, ok := ctx.Value(keys.AuthenticationKey).(*Service)
	if !ok {
		return nil, ErrAuthenticationNotFound
	}

	return auth, nil
}

// Issue issues a token with the provided claims.
func (s *Service) Issue(ctx context.Context, claims Claims) (*Token, error) {
	return NewToken(claims, s.secretKey)
}

// CreateToken creates a token with the provided claims and key.
func CreateToken(ctx context.Context, key []byte, claims Claims) (*Token, error) {
	return NewToken(claims, key)
}

// AuthenticationResult is the result of authenticating a token.
type AuthenticationResult struct {
	// Authenticated is true if the token was authenticated.
	Authenticated bool
	// Token is the token that was authenticated.
	Token *Token
}

// Authenticate authenticates a token.
func (s *Service) Authenticate(ctx context.Context, tokenStr string, source net.IP) (*AuthenticationResult, error) {
	token, err := NewTokenFromString(tokenStr)
	if err != nil {
		return nil, err
	}

	if err := token.VerifySignature(s.secretKey); err != nil {
		return nil, ErrInvalidToken
	}

	now := time.Now()
	if now.After(token.Expiry) {
		return nil, ErrTokenExpired
	}

	if !token.AuthorizedNetwork.Contains(source) {
		return nil, ErrInvalidSource
	}

	return &AuthenticationResult{
		Authenticated: true,
		Token:         token,
	}, nil
}

func (s *Service) loadSecretKey() error {
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
