# Audit Ask

A Go application that fetches pull requests from multiple GitHub repositories using the GitHub CLI with date filtering capabilities.

## Prerequisites

1. **Go 1.21 or later** - [Install Go](https://golang.org/doc/install)
2. **GitHub CLI (gh)** - [Install GitHub CLI](https://cli.github.com/)
3. **GitHub CLI Authentication** - Run `gh auth login` to authenticate

## Installation

1. Clone or download this project
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Build the application:
   ```bash
   go build -o audit-ask
   ```

## Configuration

The application supports two configuration formats:

### Single Organization Format (Recommended for skyeshanohan)

For repositories all within the same organization, use this simplified format:

```yaml
# repositories-single-org.yaml
organization: "skyeshanohan"
repositories:
  - "repo1"
  - "repo2"
  - "repo3"
  - "repo4"
  - "repo5"
```

### Full Format (For mixed organizations/users)

For repositories across different organizations or users:

```yaml
# repositories.yaml
repositories:
  - owner: "skyeshanohan"
    name: "repo1"
  - owner: "skyeshanohan"
    name: "repo2"
  - owner: "microsoft"
    name: "vscode"
  - owner: "facebook"
    name: "react"
```

## Usage

### Basic Usage

Fetch all pull requests from configured repositories:
```bash
./audit-ask
```

### With Date Filtering

Fetch pull requests created within a specific date range:
```bash
./audit-ask --start 2024-01-01 --end 2024-12-31
```

### With Custom Repositories File

Use a different repositories configuration file:
```bash
# Use the single-org format
./audit-ask --repos repositories-single-org.yaml

# Use the full format
./audit-ask --repos repositories.yaml
```

### Save Output to File

Save results to a file instead of displaying on screen:
```bash
./audit-ask --output results.txt
```

### Command Line Options

- `--repos, -r`: Path to repositories configuration file (default: repositories.yaml)
- `--start, -s`: Start date for filtering PRs (YYYY-MM-DD format)
- `--end, -e`: End date for filtering PRs (YYYY-MM-DD format)
- `--output, -o`: Output markdown file to write results (default: pr-analysis.md)
- `--workers, -w`: Maximum number of concurrent workers (default: 10 for large datasets)
- `--max-prs, -m`: Maximum PRs to fetch per repository (default: 0 = no limit)
- `--page-size, -p`: Number of PRs per page for pagination (default: 200 for large datasets)
- `--batch-size, -b`: Process repositories in batches (default: 0 = process all at once)
- `--batch, -n`: Batch number to process (used with --batch-size, default: 1)

### Examples

```bash
# Basic usage with default settings (5 workers, 100 PRs per page)
./audit-ask --repos repositories-single-org.yaml

# High-performance mode with more workers and larger page size
./audit-ask --repos repositories-single-org.yaml --workers 10 --page-size 200

# Fetch PRs from skyeshanohan org repositories (last 30 days)
./audit-ask --repos repositories-single-org.yaml --start 2024-11-01

# Fetch PRs from a specific month with performance tuning
./audit-ask --start 2024-10-01 --end 2024-10-31 --workers 8 --page-size 150

# Limit PRs per repository for faster processing
./audit-ask --start 2024-01-01 --max-prs 50 --workers 10

# Fetch PRs and save to file with optimized settings
./audit-ask --start 2024-01-01 --output skyeshanohan-pr-report.txt --workers 15 --page-size 250

# Process large organization in batches (for very large datasets)
./audit-ask --repos repositories-large-scale.yaml --batch-size 20 --batch 1 --output batch1.txt
./audit-ask --repos repositories-large-scale.yaml --batch-size 20 --batch 2 --output batch2.txt
```

## Output Format

The application generates a beautiful markdown report with the following features:

### 📊 Summary Section
- Total repositories and pull requests
- PR state breakdown
- Generation timestamp

### 📁 Repository Sections
- Organized by repository with clear headers
- PR count per repository
- PRs sorted by number (newest first)

### 🔗 Hyperlinked PRs
- PR IDs are hyperlinked to GitHub
- Direct links to view each PR
- Clean, readable format

### Sample Output Structure

```markdown
# Pull Request Analysis Report

**Generated:** 2024-01-15 14:30:00 MST

## Summary

- **Total Repositories:** 3
- **Total Pull Requests:** 25
- **PR States:**
  - merged: 18
  - open: 5
  - closed: 2

---

## 📁 skyeshanohan/repo1

**Total PRs:** 12

### [#123](https://github.com/skyeshanohan/repo1/pull/123) Fix authentication bug
- **State:** `merged`
- **Created:** 2024-01-15 10:30:00
- **Merged:** 2024-01-16 14:20:00
- **Link:** [View PR](https://github.com/skyeshanohan/repo1/pull/123)

### [#122](https://github.com/skyeshanohan/repo1/pull/122) Add new feature
- **State:** `open`
- **Created:** 2024-01-20 09:15:00
- **Link:** [View PR](https://github.com/skyeshanohan/repo1/pull/122)

---
```

## Performance Features

### Concurrent Processing
- **Worker Pool**: Process multiple repositories simultaneously using configurable worker threads
- **Pagination**: Efficiently fetch large numbers of PRs using GitHub CLI pagination
- **Memory Efficient**: Stream processing with configurable page sizes

### Performance Tuning
- **Workers**: Increase `--workers` for more concurrent repository processing (recommended: 5-15)
- **Page Size**: Larger `--page-size` reduces API calls but uses more memory (recommended: 100-250)
- **PR Limits**: Use `--max-prs` to limit PRs per repository for faster processing

### Recommended Settings
- **Small datasets** (< 10 repos): `--workers 5 --page-size 100`
- **Medium datasets** (10-50 repos): `--workers 10 --page-size 150`
- **Large datasets** (50-200 repos): `--workers 15 --page-size 200`
- **Very large datasets** (200+ repos): `--workers 15 --page-size 250 --batch-size 20`

## Notes

- The application fetches pull requests from the default branch of each repository
- All pull request states (open, closed, merged) are included
- Date filtering is based on the creation date of pull requests
- The GitHub CLI must be authenticated with appropriate permissions to access the repositories
- Rate limiting is handled by the GitHub CLI itself
- Concurrent processing significantly improves performance for multiple repositories

## Troubleshooting

1. **"GitHub CLI is not installed"**: Install GitHub CLI and ensure it's in your PATH
2. **"GitHub CLI is not authenticated"**: Run `gh auth login` to authenticate
3. **"Failed to fetch pull requests"**: Check if you have access to the repository and if the repository exists
4. **"No repositories found"**: Ensure your `repositories.yaml` file has the correct format and contains at least one repository
