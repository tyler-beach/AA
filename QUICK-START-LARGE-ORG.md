# Quick Start Guide for Large GitHub Organizations

## 🚀 Optimized for skyeshanohan and Large Organizations

This application is now optimized for large-scale GitHub organizations with enhanced performance, batch processing, and memory management.

## ⚡ Quick Start Commands

### 1. Basic Large-Scale Processing
```bash
# Process all repositories with optimized defaults (generates pr-analysis.md)
./audit-ask --repos repositories-large-scale.yaml

# Or specify a custom markdown filename
./audit-ask --repos repositories-large-scale.yaml --output skyeshanohan-analysis.md
```

### 2. Recent Activity Analysis
```bash
# Last 6 months with high performance
./audit-ask --repos repositories-large-scale.yaml --start 2024-06-01 --workers 15 --page-size 250 --output recent-analysis.md
```

### 3. Batch Processing for Very Large Organizations
```bash
# Process in batches of 20 repositories
./audit-ask --repos repositories-large-scale.yaml --batch-size 20 --batch 1 --output batch1.md
./audit-ask --repos repositories-large-scale.yaml --batch-size 20 --batch 2 --output batch2.md
./audit-ask --repos repositories-large-scale.yaml --batch-size 20 --batch 3 --output batch3.md
```

### 4. Quick Overview Mode
```bash
# Get top 100 PRs per repository for quick overview
./audit-ask --repos repositories-large-scale.yaml --max-prs 100 --workers 15 --output overview.md
```

## 📊 Performance Expectations

| Organization Size | Repositories | Processing Time | Recommended Settings |
|------------------|--------------|-----------------|---------------------|
| Small | 5-10 | 1-3 minutes | Default settings |
| Medium | 10-50 | 3-10 minutes | `--workers 10 --page-size 200` |
| Large | 50-200 | 10-30 minutes | `--workers 15 --page-size 250` |
| Very Large | 200+ | 30+ minutes | `--workers 15 --page-size 250 --batch-size 20` |

## 🎯 Optimized Defaults

The application now uses large-scale optimized defaults:

- **Workers**: 10 (increased from 5)
- **Page Size**: 200 (increased from 100)
- **Memory Management**: Optimized for large datasets
- **Progress Reporting**: Enhanced with time estimates

## 📁 Configuration Files

### repositories-large-scale.yaml
```yaml
organization: "skyeshanohan"
repositories:
  - "repo1"
  - "repo2"
  - "repo3"
  # ... add all your repositories
```

## 🔧 Performance Tuning

### For Maximum Speed
```bash
./audit-ask --repos repos.yaml --workers 15 --page-size 250
```

### For Memory Efficiency
```bash
./audit-ask --repos repos.yaml --workers 5 --page-size 100
```

### For Very Large Datasets
```bash
./audit-ask --repos repos.yaml --workers 15 --page-size 250 --batch-size 20
```

## 📈 Monitoring Progress

The application provides enhanced progress reporting:

```
🚀 Large Dataset Mode: Fetching pull requests from 75 repositories using 15 workers...
📅 Start date filter: 2024-01-01
📅 End date filter: 2024-12-31
📄 Page size: 250 (optimized for large datasets)
⚡ Estimated processing time: 5-10 minutes

✅ skyeshanohan/repo1: Found 234 pull requests
✅ skyeshanohan/repo2: Found 156 pull requests
...

🎉 Large Dataset Processing Complete!
📊 Results: 75 repositories processed successfully, 0 failed
📈 Total PRs collected: 12,456
```

## 🚨 Troubleshooting

### Slow Processing
- Increase workers: `--workers 15`
- Increase page size: `--page-size 250`
- Use date filters: `--start 2024-06-01`

### Memory Issues
- Decrease page size: `--page-size 100`
- Decrease workers: `--workers 5`
- Use batch processing: `--batch-size 20`

### API Rate Limiting
- Decrease workers: `--workers 8`
- GitHub CLI handles rate limiting automatically

## 📋 Checklist for Large Organizations

- [ ] Update `repositories-large-scale.yaml` with your repository names
- [ ] Ensure GitHub CLI is authenticated: `gh auth status`
- [ ] Test with a small subset first
- [ ] Use appropriate performance settings for your organization size
- [ ] Save results to markdown file: `--output results.md` (default: pr-analysis.md)
- [ ] Consider batch processing for very large datasets

## 🎯 Recommended Workflow for skyeshanohan

1. **Test Run**: Start with a few repositories
   ```bash
   ./audit-ask --repos test-repos.yaml --max-prs 10
   ```

2. **Recent Analysis**: Focus on last 3 months
   ```bash
   ./audit-ask --repos repositories-large-scale.yaml --start 2024-09-01 --output recent.md
   ```

3. **Full Analysis**: Complete year analysis
   ```bash
   ./audit-ask --repos repositories-large-scale.yaml --start 2024-01-01 --output full-year.md
   ```

4. **Batch Processing**: For very large datasets
   ```bash
   # Process in batches of 25 repositories
   for i in {1..10}; do
     ./audit-ask --repos repositories-large-scale.yaml --batch-size 25 --batch $i --output batch$i.md
   done
   ```
