package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadRepositories loads repositories from the YAML configuration file
// Supports multiple formats: single-org, single-org with verticals, and full format
func LoadRepositories(filename string) (*RepositoriesConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read repositories file %s: %w", filename, err)
	}

	// First try to parse as single-org with multiple verticals per repository format
	var singleOrgMultiVerticalConfig SingleOrgMultiVerticalConfig
	if err := yaml.Unmarshal(data, &singleOrgMultiVerticalConfig); err == nil && singleOrgMultiVerticalConfig.Organization != "" && len(singleOrgMultiVerticalConfig.Repositories) > 0 {
		// Convert single-org multi-vertical format to full format
		config := &RepositoriesConfig{
			Organization: singleOrgMultiVerticalConfig.Organization,
			Verticals:    []Vertical{},
		}
		
		// Group repositories by vertical
		verticalMap := make(map[string][]Repository)
		for _, repoWithVerticals := range singleOrgMultiVerticalConfig.Repositories {
			repo := Repository{
				Owner: singleOrgMultiVerticalConfig.Organization,
				Name:  repoWithVerticals.Name,
			}
			
			// Add repository to each of its verticals
			for _, verticalName := range repoWithVerticals.Verticals {
				verticalMap[verticalName] = append(verticalMap[verticalName], repo)
			}
		}
		
		// Convert map to Vertical slice
		for verticalName, repositories := range verticalMap {
			config.Verticals = append(config.Verticals, Vertical{
				Name:         verticalName,
				Repositories: repositories,
			})
		}
		
		return config, nil
	}

	// Try to parse as single-org with verticals format (one vertical per repository)
	var singleOrgVerticalConfig SingleOrgVerticalConfig
	if err := yaml.Unmarshal(data, &singleOrgVerticalConfig); err == nil && singleOrgVerticalConfig.Organization != "" && len(singleOrgVerticalConfig.Verticals) > 0 {
		// Convert single-org vertical format to full format
		config := &RepositoriesConfig{
			Organization: singleOrgVerticalConfig.Organization,
			Verticals:    make([]Vertical, len(singleOrgVerticalConfig.Verticals)),
		}
		
		for i, vertical := range singleOrgVerticalConfig.Verticals {
			config.Verticals[i] = Vertical{
				Name:         vertical.Name,
				Repositories: make([]Repository, len(vertical.Repositories)),
			}
			
			for j, repoName := range vertical.Repositories {
				config.Verticals[i].Repositories[j] = Repository{
					Owner: singleOrgVerticalConfig.Organization,
					Name:  repoName,
				}
			}
		}
		
		return config, nil
	}

	// Try to parse as single-org format (without verticals)
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

	if len(config.Repositories) == 0 && len(config.Verticals) == 0 {
		return nil, fmt.Errorf("no repositories found in configuration file")
	}

	return &config, nil
}
