package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tealeg/xlsx/v3"
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
	rootCmd.Flags().StringVarP(&startDate, "start", "s", "", "Start date for filtering PRs by merge date (YYYY-MM-DD format)")
	rootCmd.Flags().StringVarP(&endDate, "end", "e", "", "End date for filtering PRs by merge date (YYYY-MM-DD format)")
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

	// Collect all repositories (from both direct list and verticals) and deduplicate
	repositoryMap := make(map[string]Repository)
	
	// Add direct repositories
	for _, repo := range config.Repositories {
		key := fmt.Sprintf("%s/%s", repo.Owner, repo.Name)
		repositoryMap[key] = repo
	}
	
	// Add repositories from verticals
	for _, vertical := range config.Verticals {
		for _, repo := range vertical.Repositories {
			key := fmt.Sprintf("%s/%s", repo.Owner, repo.Name)
			repositoryMap[key] = repo
		}
	}
	
	// Convert map back to slice
	var allRepositories []Repository
	for _, repo := range repositoryMap {
		allRepositories = append(allRepositories, repo)
	}
	
	// Apply batch processing if specified
	repositoriesToProcess := allRepositories
	if batchSize > 0 {
		repositoriesToProcess = getBatchRepositories(allRepositories, batchSize, batchNumber)
		if len(repositoriesToProcess) == 0 {
			log.Fatalf("Batch %d is empty. Total repositories: %d, Batch size: %d", batchNumber, len(allRepositories), batchSize)
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
			fmt.Printf("📅 Start merge date filter: %s\n", filter.StartDate.Format("2006-01-02"))
		}
		if filter.EndDate != nil {
			fmt.Printf("📅 End merge date filter: %s\n", filter.EndDate.Format("2006-01-02"))
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
		Verticals  []string
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

		// Find the verticals for this repository
		verticals := findVerticalsForRepository(result.Repository, config)

		// Add repository info to each PR
		for _, pr := range result.PRs {
			allPRs = append(allPRs, struct {
				Repository string
				Verticals  []string
				PR         PullRequest
			}{
				Repository: result.Repository,
				Verticals:  verticals,
				PR:         pr,
			})
		}
	}

	fmt.Printf("\n🎉 Large Dataset Processing Complete!\n")
	fmt.Printf("📊 Results: %d repositories processed successfully, %d failed\n", successCount, errorCount)
	fmt.Printf("📈 Total PRs collected: %d\n", len(allPRs))

	// Output results
	outputResults(allPRs)
	
	// Output XLSX
	outputXLSX(allPRs)
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

// findVerticalsForRepository finds which verticals a repository belongs to
func findVerticalsForRepository(repoName string, config *RepositoriesConfig) []string {
	var verticals []string
	for _, vertical := range config.Verticals {
		for _, repo := range vertical.Repositories {
			if fmt.Sprintf("%s/%s", repo.Owner, repo.Name) == repoName {
				verticals = append(verticals, vertical.Name)
			}
		}
	}
	return verticals
}

func outputResults(allPRs []struct {
	Repository string
	Verticals  []string
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

	// Group PRs by repository - only include merged PRs
	repoPRs := make(map[string][]PullRequest)
	for _, item := range allPRs {
		// Only include PRs that are actually merged (have a merge date)
		if item.PR.MergedAt != nil && item.PR.State == "MERGED" {
			repoPRs[item.Repository] = append(repoPRs[item.Repository], item.PR)
		}
	}

	// Generate repository sections
	for repo, prs := range repoPRs {
		generateRepositorySection(output, repo, prs)
	}

}

func generateMarkdownHeader(output *os.File, allPRs []struct {
	Repository string
	Verticals  []string
	PR         PullRequest
}) {
	fmt.Fprintf(output, "# Merged Pull Request Analysis Report\n\n")
	
	// Generate timestamp
	fmt.Fprintf(output, "**Generated:** %s\n\n", time.Now().Format("2006-01-02 15:04:05 MST"))
	
	// Generate summary statistics - only count merged PRs
	repoCount := make(map[string]bool)
	mergedCount := 0
	
	for _, item := range allPRs {
		if item.PR.MergedAt != nil {
			repoCount[item.Repository] = true
			mergedCount++
		}
	}
	
	fmt.Fprintf(output, "## Summary\n\n")
	fmt.Fprintf(output, "- **Total Repositories with Merged PRs:** %d\n", len(repoCount))
	fmt.Fprintf(output, "- **Total Merged Pull Requests:** %d\n", mergedCount)
	
	fmt.Fprintf(output, "\n---\n\n")
}

func generateRepositorySection(output *os.File, repo string, prs []PullRequest) {
	// Filter for only merged PRs
	var mergedPRs []PullRequest
	for _, pr := range prs {
		if pr.MergedAt != nil {
			mergedPRs = append(mergedPRs, pr)
		}
	}
	
	// Skip repositories with no merged PRs
	if len(mergedPRs) == 0 {
		return
	}
	
	// Repository header
	fmt.Fprintf(output, "## %s\n\n", repo)
	
	// Sort PRs by number (descending)
	for i := 0; i < len(mergedPRs)-1; i++ {
		for j := i + 1; j < len(mergedPRs); j++ {
			if mergedPRs[i].Number < mergedPRs[j].Number {
				mergedPRs[i], mergedPRs[j] = mergedPRs[j], mergedPRs[i]
			}
		}
	}
	
	// Generate PR list
	for _, pr := range mergedPRs {
		generatePRMarkdown(output, repo, pr)
	}
	
	fmt.Fprintf(output, "\n---\n\n")
}

func generatePRMarkdown(output *os.File, repo string, pr PullRequest) {
	// Create GitHub PR URL
	prURL := fmt.Sprintf("https://github.com/%s/pull/%d", repo, pr.Number)
	
	// PR number with embedded URL
	fmt.Fprintf(output, "[#%d](%s)", pr.Number, prURL)
	
	// Author
	if pr.Author.Login != "" {
		fmt.Fprintf(output, " by **%s**", pr.Author.Login)
	}
	
	// Merge date (we know it's merged since we filtered for it)
	fmt.Fprintf(output, " - merged %s", pr.MergedAt.Format("2006-01-02"))
	
	fmt.Fprintf(output, "\n")
}

func outputXLSX(allPRs []struct {
	Repository string
	Verticals  []string
	PR         PullRequest
}) {
	// Create XLSX filename based on output file
	xlsxFile := strings.TrimSuffix(outputFile, ".md") + ".xlsx"
	
	// Create a new Excel file
	file := xlsx.NewFile()
	
	// Group PRs by repository with vertical info - only include merged PRs
	repoPRs := make(map[string]struct {
		Verticals []string
		PRs       []PullRequest
	})
	for _, item := range allPRs {
		// Only include PRs that are actually merged (have a merge date)
		if item.PR.MergedAt != nil && item.PR.State == "MERGED" {
			repoData := repoPRs[item.Repository]
			repoData.Verticals = item.Verticals
			repoData.PRs = append(repoData.PRs, item.PR)
			repoPRs[item.Repository] = repoData
		}
	}
	
	// Create a worksheet for each repository
	for repoName, repoData := range repoPRs {
		// Extract just the repository name (remove organization prefix)
		repoNameOnly := repoName
		if strings.Contains(repoName, "/") {
			parts := strings.Split(repoName, "/")
			repoNameOnly = parts[len(parts)-1] // Get the last part (repository name)
		}
		
		// Create worksheet name with vertical prefix
		var sheetName string
		if len(repoData.Verticals) > 0 {
			// Join multiple verticals with "/"
			verticalsStr := strings.Join(repoData.Verticals, "/")
			sheetName = fmt.Sprintf("%s - %s", verticalsStr, repoNameOnly)
		} else {
			sheetName = repoNameOnly
		}
		
		// Clean worksheet name for Excel (Excel has restrictions)
		sheetName = strings.ReplaceAll(sheetName, "/", "-")
		if len(sheetName) > 31 { // Excel worksheet name limit
			sheetName = sheetName[:31]
		}
		
		sheet, err := file.AddSheet(sheetName)
		if err != nil {
			log.Printf("Failed to create Excel sheet for %s: %v", repoName, err)
			continue
		}
		
		// Create header row
		headerRow := sheet.AddRow()
		headerRow.AddCell().SetString("PR_Number")
		headerRow.AddCell().SetString("Author")
		headerRow.AddCell().SetString("Merge_Date")
		
		// Add data rows
		for _, pr := range repoData.PRs {
			prURL := fmt.Sprintf("https://github.com/%s/pull/%d", repoName, pr.Number)
			author := pr.Author.Login
			if author == "" {
				author = "Unknown"
			}
			
			// Create row
			row := sheet.AddRow()
			
			// Create hyperlink cell for PR number
			prCell := row.AddCell()
			prCell.SetString(fmt.Sprintf("#%d", pr.Number))
			prCell.SetHyperlink(prURL, fmt.Sprintf("#%d", pr.Number), "")
			
			row.AddCell().SetString(author)
			row.AddCell().SetString(pr.MergedAt.Format("2006-01-02"))
		}
	}
	
	// Save the file
	err := file.Save(xlsxFile)
	if err != nil {
		log.Printf("Failed to save Excel file: %v", err)
		return
	}
	
	fmt.Printf("📊 Excel report generated: %s\n", xlsxFile)
}


