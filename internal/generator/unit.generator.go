package generator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nathan-fiscaletti/coattail-go/internal/generator/templates"
	"github.com/nathan-fiscaletti/coattail-go/internal/logging"
	"gopkg.in/yaml.v3"
)

type actionsYaml struct {
	Actions []templates.ActionTemplateData `yaml:"actions"`
}

type receiversYaml struct {
	Receivers []templates.ReceiverTemplateData `yaml:"receivers"`
}

func GenerateUnits(ctx context.Context, root string) error {
	log, err := logging.GetLogger(ctx)
	if err != nil {
		return err
	}

	// ====== ACTIONS ======

	actionsYamlFile := filepath.Join(root, "actions.yaml")
	actionsYamlData, err := os.ReadFile(actionsYamlFile)
	if err != nil {
		return fmt.Errorf("failed to read actions yaml file %s: %s", actionsYamlFile, err)
	}

	var actions actionsYaml
	err = yaml.Unmarshal(actionsYamlData, &actions)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml file: %s: %s", actionsYamlFile, err)
	}

	log.Printf("Found %d actions in %v\n", len(actions.Actions), actionsYamlFile)

	actionsDir := filepath.Join(root, "internal", "actions")
	for _, action := range actions.Actions {
		actionFileName := fmt.Sprintf("action.%s.go", strings.ToLower(action.Name))

		actionPath := filepath.Join(actionsDir, actionFileName)
		if _, err := os.Stat(actionPath); err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("failed to stat action file %s: %s", actionFileName, err)
			}

			log.Printf("Generating action: %s\n", actionPath)

			if err := templates.NewActionTemplate(action).Fill(actionsDir); err != nil {
				return fmt.Errorf("failed to generate action: %w", err)
			}

			continue
		}

		log.Printf("Skipping action: %s (exists)\n", actionPath)
	}

	// ====== RECEIVERS ======

	receiversYamlFile := filepath.Join(root, "receivers.yaml")
	receiversYamlData, err := os.ReadFile(receiversYamlFile)
	if err != nil {
		return fmt.Errorf("failed to read receivers yaml file %s: %s", receiversYamlFile, err)
	}

	var receivers receiversYaml
	err = yaml.Unmarshal(receiversYamlData, &receivers)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml file: %s: %s", receiversYamlFile, err)
	}

	log.Printf("Found %d receivers in %v\n", len(receivers.Receivers), receiversYamlFile)

	receiversDir := filepath.Join(root, "internal", "receivers")
	for _, receiver := range receivers.Receivers {
		receiverFileName := fmt.Sprintf("receiver.%s.go", strings.ToLower(receiver.Name))

		receiverPath := filepath.Join(receiversDir, receiverFileName)
		if _, err := os.Stat(receiverPath); err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("failed to stat receiver file %s: %s", receiverFileName, err)
			}

			log.Printf("Generating receiver: %s\n", receiverPath)

			if err := templates.NewReceiverTemplate(receiver).Fill(receiversDir); err != nil {
				return fmt.Errorf("failed to generate receiver: %w", err)
			}

			continue
		}

		log.Printf("Skipping receiver: %s (exists)\n", receiverPath)
	}

	// ====== REGISTRY ======

	log.Printf("Generating registry\n")

	registryTemplate := templates.RegistryTemplateData{
		Actions:   actions.Actions,
		Receivers: receivers.Receivers,
	}

	if err := templates.NewRegistryTemplate(registryTemplate).Fill(filepath.Join(root, "internal")); err != nil {
		return fmt.Errorf("failed to generate registry: %w", err)
	}

	log.Printf("Registry generated successfully.\n")

	// ====== SDK ======

	log.Printf("Generating SDK\n")

	sdkTemplate := templates.SdkTemplateData{
		Actions: actions.Actions,
	}

	if err := templates.NewSdkTemplate(sdkTemplate).Fill(filepath.Join(root, "pkg", "sdk")); err != nil {
		return fmt.Errorf("failed to generate sdk: %w", err)
	}

	log.Printf("SDK generated successfully.\n")

	return nil
}
