package authentication

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/vmihailenco/msgpack/v5"
)

var (
	ErrMalformedToken = errors.New("malformed token")
)

// Token is a token that can be used to authenticate a peer.
type Token struct {
	// Claims is the claims of the token.
	Claims
	// Signature is the signature of the token.
	Signature []byte
}

func (t *Token) String() string {
	res, _ := t.MarshalString()
	return res
}

// NewToken creates a new token with the provided claims and key.
func NewToken(data Claims, key []byte) (*Token, error) {
	payload, err := msgpack.Marshal(data)
	if err != nil {
		return nil, err
	}

	h := hmac.New(sha256.New, key)
	_, err = h.Write(payload)
	if err != nil {
		return nil, err
	}
	signature := h.Sum(nil)

	return &Token{
		Claims:    data,
		Signature: signature,
	}, nil
}

// NewTokenFromString creates a new token from the provided string.
func NewTokenFromString(data string) (*Token, error) {
	parts := strings.Split(data, ".")

	if len(parts) != 2 {
		return nil, ErrMalformedToken
	}

	claimsData, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode claims: %w", err)
	}

	signature, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature: %w", err)
	}

	// unmarshal the payload to TokenData
	var claims Claims
	err = msgpack.Unmarshal(claimsData, &claims)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	return &Token{
		Claims:    claims,
		Signature: signature,
	}, nil
}

// MarshalString marshals the token into a string.
func (s *Token) MarshalString() (string, error) {
	payload, err := msgpack.Marshal(s.Claims)
	if err != nil {
		return "", err
	}

	dataEncoded := base64.StdEncoding.EncodeToString(payload)
	signatureEncoded := base64.StdEncoding.EncodeToString(s.Signature[:])

	return dataEncoded + "." + signatureEncoded, nil
}

// VerifySignature verifies the signature of the token.
func (s *Token) VerifySignature(key []byte) error {
	payload, err := msgpack.Marshal(s.Claims)
	if err != nil {
		return err
	}

	h := hmac.New(sha256.New, key)
	_, err = h.Write(payload)
	if err != nil {
		return err
	}
	expectedSignature := h.Sum(nil)

	if subtle.ConstantTimeCompare(s.Signature, expectedSignature) != 1 {
		return errors.New("invalid signature for token")
	}

	return nil
}
