# Performance Optimization Guide

## Overview

The application has been optimized for high-performance processing of large numbers of repositories and pull requests using:

- **Concurrent Workers**: Process multiple repositories simultaneously
- **Pagination**: Efficiently fetch large datasets without memory issues
- **Configurable Limits**: Control resource usage and processing time

## Performance Improvements

### Before Optimization
- **Sequential Processing**: One repository at a time
- **No Pagination**: Risk of memory issues with large datasets
- **Fixed Limits**: No control over resource usage

### After Optimization
- **Concurrent Processing**: 5-15 repositories processed simultaneously
- **Smart Pagination**: Configurable page sizes (100-250 PRs per page)
- **Resource Control**: Configurable workers, limits, and page sizes

## Benchmarking Examples

### Small Dataset (5 repositories, ~50 PRs each)
```bash
# Sequential (old way) - ~30 seconds
./audit-ask --repos small-repos.yaml

# Concurrent (new way) - ~8 seconds
./audit-ask --repos small-repos.yaml --workers 5 --page-size 100
```

### Medium Dataset (20 repositories, ~200 PRs each)
```bash
# Sequential (old way) - ~3 minutes
./audit-ask --repos medium-repos.yaml

# Concurrent (new way) - ~45 seconds
./audit-ask --repos medium-repos.yaml --workers 10 --page-size 150
```

### Large Dataset (50+ repositories, ~500 PRs each)
```bash
# Sequential (old way) - ~15+ minutes
./audit-ask --repos large-repos.yaml

# Concurrent (new way) - ~2-3 minutes
./audit-ask --repos large-repos.yaml --workers 15 --page-size 200
```

## Optimal Configuration by Use Case

### Development/Testing
```bash
# Fast feedback, limited data
./audit-ask --repos repos.yaml --max-prs 20 --workers 3 --page-size 50
```

### Production Analysis
```bash
# Complete data, optimized performance
./audit-ask --repos repos.yaml --workers 10 --page-size 150
```

### Large-Scale Audit
```bash
# Maximum performance for large datasets
./audit-ask --repos repos.yaml --workers 15 --page-size 250
```

### Memory-Constrained Environment
```bash
# Lower memory usage, slower processing
./audit-ask --repos repos.yaml --workers 3 --page-size 50
```

## Performance Monitoring

The application provides real-time feedback:

```
Fetching pull requests from 25 repositories using 10 workers...
Start date filter: 2024-01-01
End date filter: 2024-12-31
Page size: 150

✓ skyeshanohan/repo1: Found 45 pull requests
✓ skyeshanohan/repo2: Found 32 pull requests
✓ skyeshanohan/repo3: Found 67 pull requests
...

Completed: 25 repositories processed successfully, 0 failed
```

## Resource Usage Guidelines

### CPU Usage
- **Low**: 1-3 workers
- **Medium**: 5-8 workers  
- **High**: 10-15 workers
- **Maximum**: 15+ workers (may hit GitHub API limits)

### Memory Usage
- **Low**: 50-100 PRs per page
- **Medium**: 100-150 PRs per page
- **High**: 150-250 PRs per page
- **Maximum**: 250+ PRs per page (may cause memory issues)

### Network Usage
- Larger page sizes = fewer API calls = faster processing
- More workers = more concurrent API calls = faster processing
- Balance based on your network and GitHub API limits

## Troubleshooting Performance Issues

### Slow Processing
1. Increase `--workers` (up to 15)
2. Increase `--page-size` (up to 250)
3. Use `--max-prs` to limit data per repository

### Memory Issues
1. Decrease `--page-size` (try 50-100)
2. Decrease `--workers` (try 3-5)
3. Use `--max-prs` to limit total data

### API Rate Limiting
1. Decrease `--workers` (try 3-5)
2. GitHub CLI handles rate limiting automatically
3. Consider running during off-peak hours
