/*
Copyright © 2024 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Chaos represents the structure of each Chaos in the profile.yaml
type Chaos struct {
	Systems   []string `yaml:"systems"`   // Systems can be a list of strings
	Severity  []string `yaml:"severity"`  // Severity can be a list of strings
	Operation string   `yaml:"operation"` // Chaos operation (e.g., "delete", "add", etc.)
	Resource  string   `yaml:"resource"`  // Resource typ
	Namespace string   `yaml:"namespace"` // Namespace
	Count     int      `yaml:"count"`     // Namespace
}

// Configuration represents the overall structure of the YAML configuration
type Configuration struct {
	ChaosEvents map[string]Chaos `yaml:"chaosEvents"`
}

// Load configuration from the profile.yaml file
func loadConfiguration(filepath string) (Configuration, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return Configuration{}, fmt.Errorf("failed to read YAML file: %v", err)
	}

	var config Configuration
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Configuration{}, fmt.Errorf("failed to parse YAML: %v", err)
	}

	return config, nil
}
