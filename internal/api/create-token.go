package api

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/nathan-fiscaletti/coattail-go/internal/logging"
	"github.com/nathan-fiscaletti/coattail-go/internal/services/authentication"
	"github.com/nathan-fiscaletti/coattail-go/internal/services/permission"
	"github.com/nathan-fiscaletti/coattail-go/internal/util"
)

func CreateToken(keyfile, network string, perm int, expiryStr string) string {
	ctx, err := util.CreateServiceContext(context.Background())
	if err != nil {
		panic(err)
	}

	log, err := logging.GetLogger(ctx)
	if err != nil {
		panic(err)
	}

	_, ipnet, err := net.ParseCIDR(network)
	if err != nil {
		log.Printf("Error: failed to parse network: %s\n", err)
		os.Exit(1)
	}

	expiry := time.Now().Add(time.Hour * 24)
	if expiryStr != "" {
		expiryVal, err := time.Parse(time.RFC3339, expiryStr)
		if err != nil {
			log.Printf("Error: failed to parse expiry: %s\n", err)
			os.Exit(1)
		}
		expiry = expiryVal
	}

	log.Println("Creating token with the following parameters:")

	log.Println()
	log.Printf("  Keyfile:    %s\n", keyfile)
	log.Printf("  Network:    %s\n", network)
	log.Printf("  Permission: %s (%d)\n", permission.GetPermissions(int32(perm)).String(), perm)
	log.Printf("  Expiry:     %s\n", expiry.String())
	log.Println()

	// make sure the keyfile exists
	if _, err := os.Stat(keyfile); os.IsNotExist(err) {
		log.Printf("Error: keyfile does not exist.\n")
		os.Exit(1)
	}

	// Read the keyfile into a byte slice
	key, err := os.ReadFile(keyfile)
	if err != nil {
		log.Printf("Error: failed to read keyfile: %s\n", err)
		os.Exit(1)
	}

	token, err := authentication.CreateToken(ctx, key, authentication.Claims{
		AuthorizedNetwork: *ipnet,
		Permitted:         int32(perm),
		Expiry:            expiry,
	})
	if err != nil {
		log.Printf("Error: failed to create token: %s\n", err)
		os.Exit(1)
	}

	log.Println("Token created successfully.")
	log.Println()
	log.Printf("  Token:      %s\n", token.String())

	return token.String()
}
