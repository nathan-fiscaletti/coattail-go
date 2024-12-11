package generator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nathan-fiscaletti/coattail-go/internal/generator/templates"
	"gopkg.in/yaml.v3"
)

type actionsYaml struct {
	Actions []templates.ActionTemplateData `yaml:"actions"`
}

func GenerateActions(destination string, actionsYamlFile string) error {
	yamlData, err := os.ReadFile(actionsYamlFile)
	if err != nil {
		return fmt.Errorf("failed to read actions yaml file %s: %s", actionsYamlFile, err)
	}

	var actions actionsYaml
	err = yaml.Unmarshal(yamlData, &actions)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml file: %s: %s", actionsYamlFile, err)
	}

	logger := log.Default()

	logger.Printf("Found %d actions in %v\n", len(actions.Actions), actionsYamlFile)

	for _, action := range actions.Actions {
		actionFileName := fmt.Sprintf("action.%s.go", strings.ToLower(action.Name))

		actionPath := filepath.Join(destination, actionFileName)
		if _, err := os.Stat(filepath.Join(destination, actionFileName)); err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("failed to stat action file %s: %s", actionFileName, err)
			}

			logger.Printf("Generating action: %s\n", actionPath)

			actionTemplate := templates.ActionTemplateData{
				Name:       action.Name,
				InputType:  action.InputType,
				OutputType: action.OutputType,
			}

			if err := templates.NewActionTemplate(actionTemplate).Fill(destination); err != nil {
				return fmt.Errorf("failed to generate action: %w", err)
			}

			continue
		}

		logger.Printf("Skipping action: %s (exists)\n", actionPath)
	}

	// Generate the action registry
	actionRegistryTemplate := templates.ActionRegistryTemplateData{
		Actions: actions.Actions,
	}

	if err := templates.NewActionRegistryTemplate(actionRegistryTemplate).Fill(filepath.Join(destination, "..")); err != nil {
		return fmt.Errorf("failed to generate action registry: %w", err)
	}

	return nil
}
