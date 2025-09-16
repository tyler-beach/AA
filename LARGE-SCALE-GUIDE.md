# Large-Scale GitHub Organization Processing Guide

## Overview

This guide is specifically designed for processing large GitHub organizations like skyeshanohan with many repositories and thousands of pull requests.

## Optimized Defaults for Large Datasets

The application now uses optimized defaults for large-scale processing:

- **Workers**: 10 (increased from 5)
- **Page Size**: 200 (increased from 100)
- **Memory Management**: Optimized for large datasets
- **Progress Reporting**: Enhanced with emojis and time estimates

## Large-Scale Configuration

### 1. Repository Configuration

Use the `repositories-large-scale.yaml` format:

```yaml
organization: "skyeshanohan"
repositories:
  - "repo1"
  - "repo2"
  - "repo3"
  # ... add all your repositories
```

### 2. Performance Settings

For large organizations (50+ repositories):

```bash
# Maximum performance (recommended for large orgs)
./audit-ask --repos repositories-large-scale.yaml --workers 15 --page-size 250

# Balanced performance (good for most large orgs)
./audit-ask --repos repositories-large-scale.yaml --workers 10 --page-size 200

# Memory-constrained (if you have limited RAM)
./audit-ask --repos repositories-large-scale.yaml --workers 8 --page-size 150
```

## Processing Strategies

### Strategy 1: Full Analysis (All PRs)
```bash
# Process all PRs from all repositories
./audit-ask --repos repositories-large-scale.yaml --output full-analysis.txt
```

### Strategy 2: Recent Analysis (Last 6 months)
```bash
# Focus on recent activity
./audit-ask --repos repositories-large-scale.yaml --start 2024-06-01 --output recent-analysis.txt
```

### Strategy 3: Limited Analysis (Top PRs per repo)
```bash
# Get top 100 PRs per repository for quick overview
./audit-ask --repos repositories-large-scale.yaml --max-prs 100 --output overview.txt
```

### Strategy 4: Time-Bounded Analysis
```bash
# Specific time period with high performance
./audit-ask --repos repositories-large-scale.yaml --start 2024-01-01 --end 2024-12-31 --workers 15 --page-size 250 --output yearly-analysis.txt
```

## Performance Expectations

### Small Organization (5-10 repos)
- **Processing Time**: 1-3 minutes
- **Memory Usage**: Low
- **Recommended Settings**: Default settings work well

### Medium Organization (10-50 repos)
- **Processing Time**: 3-10 minutes
- **Memory Usage**: Medium
- **Recommended Settings**: `--workers 10 --page-size 200`

### Large Organization (50-200 repos)
- **Processing Time**: 10-30 minutes
- **Memory Usage**: High
- **Recommended Settings**: `--workers 15 --page-size 250`

### Very Large Organization (200+ repos)
- **Processing Time**: 30+ minutes
- **Memory Usage**: Very High
- **Recommended Settings**: `--workers 15 --page-size 250 --max-prs 500`

## Memory Management

### For Large Datasets
- Use larger page sizes (200-250) to reduce API calls
- Monitor memory usage during processing
- Consider using `--max-prs` to limit data per repository

### Memory-Optimized Settings
```bash
# Lower memory usage
./audit-ask --repos repos.yaml --workers 5 --page-size 100 --max-prs 200
```

## Monitoring and Progress

The application provides enhanced progress reporting for large datasets:

```
🚀 Large Dataset Mode: Fetching pull requests from 75 repositories using 15 workers...
📅 Start date filter: 2024-01-01
📅 End date filter: 2024-12-31
📄 Page size: 250 (optimized for large datasets)
⚡ Estimated processing time: 5-10 minutes

✅ skyeshanohan/repo1: Found 234 pull requests
✅ skyeshanohan/repo2: Found 156 pull requests
✅ skyeshanohan/repo3: Found 89 pull requests
...

🎉 Large Dataset Processing Complete!
📊 Results: 75 repositories processed successfully, 0 failed
📈 Total PRs collected: 12,456
```

## Troubleshooting Large Datasets

### Slow Processing
1. Increase workers: `--workers 15`
2. Increase page size: `--page-size 250`
3. Use date filters to reduce data: `--start 2024-06-01`

### Memory Issues
1. Decrease page size: `--page-size 100`
2. Decrease workers: `--workers 5`
3. Limit PRs per repo: `--max-prs 200`

### API Rate Limiting
1. Decrease workers: `--workers 8`
2. GitHub CLI handles rate limiting automatically
3. Consider running during off-peak hours

## Best Practices for Large Organizations

1. **Start Small**: Test with a few repositories first
2. **Use Date Filters**: Focus on relevant time periods
3. **Monitor Resources**: Watch CPU and memory usage
4. **Save Results**: Always use `--output` for large datasets
5. **Batch Processing**: Consider processing in smaller batches if needed

## Example Commands for skyeshanohan

```bash
# Quick overview (last 3 months, top 50 PRs per repo)
./audit-ask --repos repositories-large-scale.yaml --start 2024-09-01 --max-prs 50 --output overview.txt

# Full year analysis with maximum performance
./audit-ask --repos repositories-large-scale.yaml --start 2024-01-01 --end 2024-12-31 --workers 15 --page-size 250 --output full-year.txt

# Recent activity analysis
./audit-ask --repos repositories-large-scale.yaml --start 2024-10-01 --workers 12 --page-size 200 --output recent.txt
```
