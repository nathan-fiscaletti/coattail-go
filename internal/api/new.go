package api

import (
	"context"
	"os"
	"path/filepath"

	"github.com/nathan-fiscaletti/coattail-go/internal/generator"
	"github.com/nathan-fiscaletti/coattail-go/internal/logging"
	"github.com/nathan-fiscaletti/coattail-go/internal/util"
)

func CreateCoattailInstance(destination string, packageName string) {
	ctx, err := util.CreateServiceContext(context.Background())
	if err != nil {
		panic(err)
	}

	log, err := logging.GetLogger(ctx)
	if err != nil {
		panic(err)
	}

	destination, err = filepath.Abs(destination)
	if err != nil {
		panic(err)
	}

	log.Printf("Creating new Coattail instance at: %v\n", destination)

	err = generator.GenerateNewMod(destination, packageName)
	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	log.Printf("Successfully created new Coattail instance at: %v\n", destination)
}
