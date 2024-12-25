package api

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/nathan-fiscaletti/coattail-go/internal/generator"
	"github.com/nathan-fiscaletti/coattail-go/internal/logging"
	"github.com/nathan-fiscaletti/coattail-go/internal/util"
)

func RunGeneration(path string) {
	ctx, err := util.CreateServiceContext(context.Background())
	if err != nil {
		panic(err)
	}

	log, err := logging.GetLogger(ctx)
	if err != nil {
		panic(err)
	}

	// extract package name from go.mod
	goModPath := filepath.Join(path, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		log.Printf("Error: go.mod does not exist. Are you in a coattail instance?\n")
		os.Exit(1)
	}

	goModData, err := os.ReadFile(goModPath)
	if err != nil {
		log.Printf("Error: failed to read go.mod file: %s\n", err)
		os.Exit(1)
	}

	goModData = bytes.Split(goModData, []byte("\n"))[0]
	packageName := strings.TrimPrefix(strings.TrimSpace(string(goModData)), "module ")

	// make sure that host-config.yaml exists
	hostConfigPath := filepath.Join(path, "host-config.yaml")
	if _, err := os.Stat(hostConfigPath); os.IsNotExist(err) {
		log.Printf("Error: host-config.yaml does not exist. Are you in a coattail instance?\n")
		os.Exit(1)
	}

	// make sure actions.yaml exists
	actionsYamlPath := filepath.Join(path, "/actions.yaml")
	if _, err := os.Stat(actionsYamlPath); os.IsNotExist(err) {
		log.Printf("Error: actions.yaml does not exist. Are you in a coattail instance?\n")
		os.Exit(1)
	}

	err = generator.GenerateUnits(ctx, path, packageName)
	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	log.Printf("Code generation completed successfully.\n")
}
