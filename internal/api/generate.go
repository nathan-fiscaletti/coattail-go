package api

import (
	"context"
	"os"
	"path/filepath"

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

	err = generator.GenerateUnits(ctx, path)
	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	log.Printf("Code generation completed successfully.\n")
}
