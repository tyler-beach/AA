package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	reposFile      string
	startDate      string
	endDate        string
	outputFile     string
	maxWorkers     int
	maxPRsPerRepo  int
	pageSize       int
	batchSize      int
	batchNumber    int
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "audit-ask",
		Short: "Fetch pull requests from GitHub repositories",
		Long:  "A tool to fetch pull requests from multiple GitHub repositories using GitHub CLI with date filtering",
		Run:   run,
	}

	rootCmd.Flags().StringVarP(&reposFile, "repos", "r", "repositories.yaml", "Path to repositories configuration file")
	rootCmd.Flags().StringVarP(&startDate, "start", "s", "", "Start date for filtering PRs (YYYY-MM-DD format)")
	rootCmd.Flags().StringVarP(&endDate, "end", "e", "", "End date for filtering PRs (YYYY-MM-DD format)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "pr-analysis.md", "Output markdown file to write results (default: pr-analysis.md)")
	rootCmd.Flags().IntVarP(&maxWorkers, "workers", "w", 10, "Maximum number of concurrent workers (default: 10 for large datasets)")
	rootCmd.Flags().IntVarP(&maxPRsPerRepo, "max-prs", "m", 0, "Maximum PRs to fetch per repository (0 = no limit)")
	rootCmd.Flags().IntVarP(&pageSize, "page-size", "p", 200, "Number of PRs per page for pagination (default: 200 for large datasets)")
	rootCmd.Flags().IntVarP(&batchSize, "batch-size", "b", 0, "Process repositories in batches (0 = process all at once)")
	rootCmd.Flags().IntVarP(&batchNumber, "batch", "n", 1, "Batch number to process (used with --batch-size)")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	// Load repositories configuration
	config, err := LoadRepositories(reposFile)
	if err != nil {
		log.Fatalf("Failed to load repositories: %v", err)
	}

	// Apply batch processing if specified
	repositoriesToProcess := config.Repositories
	if batchSize > 0 {
		repositoriesToProcess = getBatchRepositories(config.Repositories, batchSize, batchNumber)
		if len(repositoriesToProcess) == 0 {
			log.Fatalf("Batch %d is empty. Total repositories: %d, Batch size: %d", batchNumber, len(config.Repositories), batchSize)
		}
	}

	// Parse date filters
	filter, err := parseDateFilter()
	if err != nil {
		log.Fatalf("Failed to parse date filter: %v", err)
	}

	// Create worker configuration
	workerConfig := &WorkerConfig{
		MaxWorkers:    maxWorkers,
		MaxPRsPerRepo: maxPRsPerRepo,
		PageSize:      pageSize,
	}

	// Initialize GitHub client
	githubClient := NewGitHubClient()

	// Check if GitHub CLI is available and authenticated
	if err := githubClient.CheckGitHubCLI(); err != nil {
		log.Fatalf("GitHub CLI check failed: %v", err)
	}

	if batchSize > 0 {
		totalBatches := (len(config.Repositories) + batchSize - 1) / batchSize
		fmt.Printf("🚀 Large Dataset Mode: Processing batch %d/%d (%d repositories) using %d workers...\n", 
			batchNumber, totalBatches, len(repositoriesToProcess), workerConfig.MaxWorkers)
	} else {
		fmt.Printf("🚀 Large Dataset Mode: Fetching pull requests from %d repositories using %d workers...\n", 
			len(repositoriesToProcess), workerConfig.MaxWorkers)
	}
	if filter != nil {
		if filter.StartDate != nil {
			fmt.Printf("📅 Start date filter: %s\n", filter.StartDate.Format("2006-01-02"))
		}
		if filter.EndDate != nil {
			fmt.Printf("📅 End date filter: %s\n", filter.EndDate.Format("2006-01-02"))
		}
	}
	if workerConfig.MaxPRsPerRepo > 0 {
		fmt.Printf("🔢 Max PRs per repository: %d\n", workerConfig.MaxPRsPerRepo)
	}
	fmt.Printf("📄 Page size: %d (optimized for large datasets)\n", workerConfig.PageSize)
	fmt.Printf("⚡ Estimated processing time: %d-%.0f minutes\n", len(config.Repositories)/workerConfig.MaxWorkers, float64(len(config.Repositories))/float64(workerConfig.MaxWorkers)*2)
	fmt.Println()

	// Fetch pull requests concurrently
	results := githubClient.FetchPullRequestsConcurrent(repositoriesToProcess, filter, workerConfig)

	// Process results
	var allPRs []struct {
		Repository string
		PR         PullRequest
	}

	successCount := 0
	errorCount := 0

	for _, result := range results {
		if result.Error != nil {
			log.Printf("❌ Warning: Failed to fetch PRs for %s: %v", result.Repository, result.Error)
			errorCount++
			continue
		}

		successCount++
		fmt.Printf("✅ %s: Found %d pull requests\n", result.Repository, len(result.PRs))

		// Add repository info to each PR
		for _, pr := range result.PRs {
			allPRs = append(allPRs, struct {
				Repository string
				PR         PullRequest
			}{
				Repository: result.Repository,
				PR:         pr,
			})
		}
	}

	fmt.Printf("\n🎉 Large Dataset Processing Complete!\n")
	fmt.Printf("📊 Results: %d repositories processed successfully, %d failed\n", successCount, errorCount)
	fmt.Printf("📈 Total PRs collected: %d\n", len(allPRs))

	// Output results
	outputResults(allPRs)
}

func parseDateFilter() (*PRFilter, error) {
	var filter *PRFilter

	if startDate != "" || endDate != "" || maxPRsPerRepo > 0 {
		filter = &PRFilter{}

		if startDate != "" {
			parsed, err := time.Parse("2006-01-02", startDate)
			if err != nil {
				return nil, fmt.Errorf("invalid start date format: %v", err)
			}
			filter.StartDate = &parsed
		}

		if endDate != "" {
			parsed, err := time.Parse("2006-01-02", endDate)
			if err != nil {
				return nil, fmt.Errorf("invalid end date format: %v", err)
			}
			filter.EndDate = &parsed
		}

		if maxPRsPerRepo > 0 {
			filter.Limit = maxPRsPerRepo
		}
	}

	return filter, nil
}

// getBatchRepositories returns a slice of repositories for the specified batch
func getBatchRepositories(repositories []Repository, batchSize, batchNumber int) []Repository {
	start := (batchNumber - 1) * batchSize
	end := start + batchSize
	
	if start >= len(repositories) {
		return []Repository{}
	}
	
	if end > len(repositories) {
		end = len(repositories)
	}
	
	return repositories[start:end]
}

func outputResults(allPRs []struct {
	Repository string
	PR         PullRequest
}) {
	var output *os.File
	var err error

	// Always create a markdown file
	output, err = os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer output.Close()

	// Generate markdown header
	generateMarkdownHeader(output, allPRs)

	if len(allPRs) == 0 {
		fmt.Fprintf(output, "No pull requests found matching the criteria.\n")
		return
	}

	// Group PRs by repository
	repoPRs := make(map[string][]PullRequest)
	for _, item := range allPRs {
		repoPRs[item.Repository] = append(repoPRs[item.Repository], item.PR)
	}

	// Generate repository sections
	for repo, prs := range repoPRs {
		generateRepositorySection(output, repo, prs)
	}

	// Generate footer
	generateMarkdownFooter(output, allPRs)
}

func generateMarkdownHeader(output *os.File, allPRs []struct {
	Repository string
	PR         PullRequest
}) {
	fmt.Fprintf(output, "# Pull Request Analysis Report\n\n")
	
	// Generate timestamp
	fmt.Fprintf(output, "**Generated:** %s\n\n", time.Now().Format("2006-01-02 15:04:05 MST"))
	
	// Generate summary statistics
	repoCount := make(map[string]bool)
	stateCount := make(map[string]int)
	
	for _, item := range allPRs {
		repoCount[item.Repository] = true
		stateCount[item.PR.State]++
	}
	
	fmt.Fprintf(output, "## Summary\n\n")
	fmt.Fprintf(output, "- **Total Repositories:** %d\n", len(repoCount))
	fmt.Fprintf(output, "- **Total Pull Requests:** %d\n", len(allPRs))
	
	if len(stateCount) > 0 {
		fmt.Fprintf(output, "- **PR States:**\n")
		for state, count := range stateCount {
			fmt.Fprintf(output, "  - %s: %d\n", state, count)
		}
	}
	
	fmt.Fprintf(output, "\n---\n\n")
}

func generateRepositorySection(output *os.File, repo string, prs []PullRequest) {
	// Repository header
	fmt.Fprintf(output, "## 📁 %s\n\n", repo)
	fmt.Fprintf(output, "**Total PRs:** %d\n\n", len(prs))
	
	// Sort PRs by number (descending)
	for i := 0; i < len(prs)-1; i++ {
		for j := i + 1; j < len(prs); j++ {
			if prs[i].Number < prs[j].Number {
				prs[i], prs[j] = prs[j], prs[i]
			}
		}
	}
	
	// Generate PR list
	for _, pr := range prs {
		generatePRMarkdown(output, repo, pr)
	}
	
	fmt.Fprintf(output, "\n---\n\n")
}

func generatePRMarkdown(output *os.File, repo string, pr PullRequest) {
	// Create GitHub PR URL
	prURL := fmt.Sprintf("https://github.com/%s/pull/%d", repo, pr.Number)
	
	// PR header with hyperlinked ID
	fmt.Fprintf(output, "### [#%d](%s) %s\n", pr.Number, prURL, pr.Title)
	
	// PR details
	fmt.Fprintf(output, "- **State:** `%s`\n", pr.State)
	fmt.Fprintf(output, "- **Created:** %s\n", pr.CreatedAt.Format("2006-01-02 15:04:05"))
	
	if pr.MergedAt != nil {
		fmt.Fprintf(output, "- **Merged:** %s\n", pr.MergedAt.Format("2006-01-02 15:04:05"))
	}
	
	// Add direct link
	fmt.Fprintf(output, "- **Link:** [View PR](%s)\n", prURL)
	
	fmt.Fprintf(output, "\n")
}

func generateMarkdownFooter(output *os.File, allPRs []struct {
	Repository string
	PR         PullRequest
}) {
	fmt.Fprintf(output, "---\n\n")
	fmt.Fprintf(output, "## 📊 Analysis Complete\n\n")
	fmt.Fprintf(output, "This report was generated by [audit-ask](https://github.com/your-org/audit-ask) - a tool for analyzing GitHub pull requests across multiple repositories.\n\n")
	
	// Add generation info
	fmt.Fprintf(output, "**Report generated on:** %s\n", time.Now().Format("2006-01-02 15:04:05 MST"))
	fmt.Fprintf(output, "**Total processing time:** See console output above\n")
}

