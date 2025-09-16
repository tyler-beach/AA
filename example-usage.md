# Example Usage for skyeshanohan Organization

## Quick Start

1. **Update the repository list** in `repositories-single-org.yaml`:
   ```yaml
   organization: "skyeshanohan"
   repositories:
     - "your-actual-repo-1"
     - "your-actual-repo-2"
     - "your-actual-repo-3"
   ```

2. **Run the application**:
   ```bash
   # Fetch all PRs from skyeshanohan org repositories
   ./audit-ask --repos repositories-single-org.yaml
   
   # Fetch PRs from the last 30 days
   ./audit-ask --repos repositories-single-org.yaml --start 2024-11-01
   
   # Fetch PRs from a specific date range
   ./audit-ask --repos repositories-single-org.yaml --start 2024-01-01 --end 2024-12-31
   
   # Save results to a file
   ./audit-ask --repos repositories-single-org.yaml --output skyeshanohan-prs.txt
   ```

## Sample Output

```
Fetching pull requests from 3 repositories...
Start date filter: 2024-01-01
End date filter: 2024-12-31

Fetching PRs for skyeshanohan/repo1...
Found 15 pull requests
Fetching PRs for skyeshanohan/repo2...
Found 8 pull requests
Fetching PRs for skyeshanohan/repo3...
Found 12 pull requests

Pull Request Summary
===================

Total pull requests found: 35

Repository: skyeshanohan/repo1
PR #123: Fix authentication bug
State: merged
Created: 2024-03-15 10:30:00
Merged: 2024-03-16 14:20:00
---
Repository: skyeshanohan/repo1
PR #124: Add new feature
State: open
Created: 2024-03-20 09:15:00
---
...
```

## Prerequisites

Make sure you have:
1. GitHub CLI installed: `gh --version`
2. GitHub CLI authenticated: `gh auth status`
3. Access to the skyeshanohan organization repositories
