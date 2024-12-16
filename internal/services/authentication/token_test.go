package authentication_test

import (
	"net"
	"testing"
	"time"

	"github.com/nathan-fiscaletti/coattail-go/internal/services/authentication"
	"github.com/nathan-fiscaletti/coattail-go/internal/services/permission"
)

func TestToken(t *testing.T) {
	key := []byte("test")

	_, ipnet, _ := net.ParseCIDR("127.0.0.1/32")
	expiry, _ := time.Parse(time.RFC3339, "2022-01-01T00:00:00Z")
	tokenData := authentication.Claims{
		AuthorizedNetwork: *ipnet,
		Permitted:         permission.PermissionMask(permission.All),
		Expiry:            expiry,
	}

	token, err := authentication.NewToken(tokenData, key)
	if err != nil {
		t.Fatal(err)
	}

	if err := token.VerifySignature(key); err != nil {
		t.Fatal(err)
	}

	val, err := token.MarshalString()
	if err != nil {
		t.Fatal(err)
	}

	if token.AuthorizedNetwork.String() != tokenData.AuthorizedNetwork.String() {
		t.Errorf("expected %s, got %s", tokenData.AuthorizedNetwork, token.AuthorizedNetwork)
	}

	if !token.Expiry.Equal(tokenData.Expiry) {
		t.Errorf("expected %s, got %s", tokenData.Expiry, token.Expiry)
	}

	if token.Permitted != tokenData.Permitted {
		t.Errorf("expected %d, got %d", tokenData.Permitted, token.Permitted)
	}

	expectedToken := "hLFBdXRob3JpemVkTmV0d29ya4KiSVDECTEyNy4wLjAuMaRNYXNrxAT/////qVBlcm1pdHRlZNIAAAAHrkF1dGhvcml6YXRpb25zwKZFeHBpcnnW/2HPmYA=.F5490Sc4vGG9JPDg8m4fpMaWkWCRKolKZhS3Cqc9DN0="
	if val != expectedToken {
		t.Errorf("got %s, want %s", val, expectedToken)
	}

	token2, err := authentication.NewTokenFromString(val)
	if err != nil {
		t.Fatal(err)
	}

	if err := token2.VerifySignature(key); err != nil {
		t.Fatal(err)
	}

	if token2.AuthorizedNetwork.String() != tokenData.AuthorizedNetwork.String() {
		t.Errorf("expected %s, got %s", tokenData.AuthorizedNetwork, token.AuthorizedNetwork)
	}

	if !token2.Expiry.Equal(tokenData.Expiry) {
		t.Errorf("expected %s, got %s", tokenData.Expiry, token.Expiry)
	}

	if token2.Permitted != tokenData.Permitted {
		t.Errorf("expected %d, got %d", tokenData.Permitted, token.Permitted)
	}
}
