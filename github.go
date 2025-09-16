package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

// GitHubClient handles GitHub CLI operations
type GitHubClient struct{}

// NewGitHubClient creates a new GitHub client
func NewGitHubClient() *GitHubClient {
	return &GitHubClient{}
}

// FetchPullRequests fetches pull requests for a repository using GitHub CLI
func (gc *GitHubClient) FetchPullRequests(owner, repo string, filter *PRFilter, workerConfig *WorkerConfig) ([]PullRequest, error) {
	// Build the GitHub CLI command
	cmd := exec.Command("gh", "pr", "list", 
		"--repo", fmt.Sprintf("%s/%s", owner, repo),
		"--state", "merged",
		"--json", "number,title,state,mergedAt,createdAt,author")

	// Set limit based on filter or use a reasonable default
	limit := 1000 // Default limit for GitHub CLI
	if filter != nil && filter.Limit > 0 {
		limit = filter.Limit
	}
	cmd.Args = append(cmd.Args, "--limit", strconv.Itoa(limit))

	// Add date filters if provided
	if filter != nil {
		if filter.StartDate != nil {
			cmd.Args = append(cmd.Args, "--created", fmt.Sprintf(">=%s", filter.StartDate.Format("2006-01-02")))
		}
		if filter.EndDate != nil {
			cmd.Args = append(cmd.Args, "--created", fmt.Sprintf("<=%s", filter.EndDate.Format("2006-01-02")))
		}
	}

	// Execute the command
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pull requests for %s/%s: %w", owner, repo, err)
	}

	// Parse the JSON output
	var prs []PullRequest
	if err := json.Unmarshal(output, &prs); err != nil {
		return nil, fmt.Errorf("failed to parse pull request data for %s/%s: %w", owner, repo, err)
	}

	return prs, nil
}

// CheckGitHubCLI checks if GitHub CLI is installed and authenticated
func (gc *GitHubClient) CheckGitHubCLI() error {
	// Check if gh command exists
	cmd := exec.Command("gh", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("GitHub CLI (gh) is not installed or not in PATH: %w", err)
	}

	// Check if authenticated
	cmd = exec.Command("gh", "auth", "status")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("GitHub CLI is not authenticated: %w", err)
	}

	// Check if output contains "Logged in"
	if !strings.Contains(string(output), "Logged in") {
		return fmt.Errorf("GitHub CLI is not authenticated")
	}

	return nil
}

// GetDefaultBranch gets the default branch for a repository
func (gc *GitHubClient) GetDefaultBranch(owner, repo string) (string, error) {
	cmd := exec.Command("gh", "repo", "view", fmt.Sprintf("%s/%s", owner, repo), "--json", "defaultBranchRef")
	
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get default branch for %s/%s: %w", owner, repo, err)
	}

	var result struct {
		DefaultBranchRef struct {
			Name string `json:"name"`
		} `json:"defaultBranchRef"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return "", fmt.Errorf("failed to parse default branch data: %w", err)
	}

	return result.DefaultBranchRef.Name, nil
}

// FetchPullRequestsConcurrent fetches pull requests from multiple repositories concurrently
func (gc *GitHubClient) FetchPullRequestsConcurrent(repositories []Repository, filter *PRFilter, workerConfig *WorkerConfig) []RepositoryResult {
	// Create channels for work distribution and results
	jobs := make(chan Repository, len(repositories))
	results := make(chan RepositoryResult, len(repositories))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < workerConfig.MaxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for repo := range jobs {
				prs, err := gc.FetchPullRequests(repo.Owner, repo.Name, filter, workerConfig)
				results <- RepositoryResult{
					Repository: fmt.Sprintf("%s/%s", repo.Owner, repo.Name),
					PRs:        prs,
					Error:      err,
				}
			}
		}()
	}

	// Send jobs to workers
	go func() {
		defer close(jobs)
		for _, repo := range repositories {
			jobs <- repo
		}
	}()

	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var allResults []RepositoryResult
	for result := range results {
		allResults = append(allResults, result)
	}

	return allResults
}
