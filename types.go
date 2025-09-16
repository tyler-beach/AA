package main

import "time"

// Repository represents a GitHub repository
// Owner can be either a GitHub username or organization name
type Repository struct {
	Owner string `yaml:"owner"` // GitHub username or organization name
	Name  string `yaml:"name"`  // Repository name
}

// Vertical represents a business vertical with its repositories
type Vertical struct {
	Name         string       `yaml:"name"`         // Vertical name (e.g., "Provider", "Payer")
	Repositories []Repository `yaml:"repositories"` // Repositories under this vertical
}

// RepositoriesConfig represents the configuration file structure
type RepositoriesConfig struct {
	Organization string       `yaml:"organization,omitempty"` // Optional: for single-org configs
	Repositories []Repository `yaml:"repositories"`
	Verticals    []Vertical   `yaml:"verticals,omitempty"` // Optional: for vertical-based configs
}

// SingleOrgConfig represents a simplified configuration for a single organization
type SingleOrgConfig struct {
	Organization string   `yaml:"organization"`
	Repositories []string `yaml:"repositories"`
}

// RepositoryWithVerticals represents a repository with its associated verticals
type RepositoryWithVerticals struct {
	Name      string   `yaml:"name"`
	Verticals []string `yaml:"verticals"`
}

// SingleOrgVerticalConfig represents a simplified configuration with verticals
type SingleOrgVerticalConfig struct {
	Organization string   `yaml:"organization"`
	Verticals    []struct {
		Name         string   `yaml:"name"`
		Repositories []string `yaml:"repositories"`
	} `yaml:"verticals"`
}

// SingleOrgMultiVerticalConfig represents a configuration where repositories can have multiple verticals
type SingleOrgMultiVerticalConfig struct {
	Organization string                    `yaml:"organization"`
	Repositories []RepositoryWithVerticals `yaml:"repositories"`
}

// PullRequest represents a GitHub pull request
type PullRequest struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"`
	MergedAt *time.Time `json:"mergedAt,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	Author struct {
		Login string `json:"login"`
	} `json:"author"`
}

// PRFilter represents filtering options for pull requests
type PRFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	Limit     int // Maximum number of PRs to fetch per repository (0 = no limit)
}

// WorkerConfig represents configuration for concurrent processing
type WorkerConfig struct {
	MaxWorkers    int // Maximum number of concurrent workers
	MaxPRsPerRepo int // Maximum PRs to fetch per repository
	PageSize      int // Number of PRs per page for pagination
}

// RepositoryResult represents the result of processing a single repository
type RepositoryResult struct {
	Repository string
	PRs        []PullRequest
	Error      error
}
