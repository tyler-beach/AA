package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadRepositories loads repositories from the YAML configuration file
// Supports both full format (owner/name pairs) and single-org format (organization + repo names)
func LoadRepositories(filename string) (*RepositoriesConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read repositories file %s: %w", filename, err)
	}

	// First try to parse as single-org format
	var singleOrgConfig SingleOrgConfig
	if err := yaml.Unmarshal(data, &singleOrgConfig); err == nil && singleOrgConfig.Organization != "" {
		// Convert single-org format to full format
		config := &RepositoriesConfig{
			Organization: singleOrgConfig.Organization,
			Repositories: make([]Repository, len(singleOrgConfig.Repositories)),
		}
		
		for i, repoName := range singleOrgConfig.Repositories {
			config.Repositories[i] = Repository{
				Owner: singleOrgConfig.Organization,
				Name:  repoName,
			}
		}
		
		if len(config.Repositories) == 0 {
			return nil, fmt.Errorf("no repositories found in configuration file")
		}
		
		return config, nil
	}

	// Fall back to full format
	var config RepositoriesConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse repositories file: %w", err)
	}

	if len(config.Repositories) == 0 {
		return nil, fmt.Errorf("no repositories found in configuration file")
	}

	return &config, nil
}
