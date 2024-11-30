# dailyare

Mark GitHub pull request notifications as read if the pull request has been merged.

dailyare helps keep your GitHub notifications tidy by automatically marking notifications as read for pull requests that have already been merged.

It directly integrates with GitHub CLI (gh) for authentication and API access.

## Installation

```bash
go install github.com/gkwa/dailyare@latest
```

Requires GitHub CLI to be installed and authenticated:

```bash
# Install GitHub CLI
brew install gh       # macOS
gh auth login        # Login to GitHub
```

## Usage

Basic usage:

```bash
# Mark notifications from last 7 days as read
dailyare

# Custom time period
dailyare --since 14d

# Bypass cache and fetch fresh data
dailyare --no-cache

# Increase logging verbosity
dailyare -v
dailyare -v -v

# JSON logging format
dailyare --log-format json
```

## How It Works

1. Fetches GitHub notifications for the configured time period

2. Filters for pull request notifications

3. Checks if each pull request has been merged

4. If merged, marks the notification as read

5. Maintains a local cache to avoid rechecking already processed notifications

## Performance

Uses caching by default to minimize API calls:

- Caches PR merge status
- Tracks already marked notifications
- Cache stored in `~/.dailyare/cache.json`

Override cache with `--no-cache` flag.

## Troubleshooting

Enable verbose logging to see what's happening:

```bash
# Show info logs
dailyare -v

# Show debug logs
dailyare -v -v

# JSON format for parsing
dailyare --log-format json
```

## Similar Tools

Shell script version using GitHub CLI:

```bash
#!/bin/bash
days=${1:-7}
since_date=$(date -v-${days}d +%Y-%m-%dT%H:%M:%SZ)

gh api "notifications?all=true&since=${since_date}" | \
  jq -r '.[] | select(.subject.type=="PullRequest")' | \
  while read -r n; do
    url=$(echo "$n" | jq -r .subject.url)
    id=$(echo "$n" | jq -r .id)
    if [[ $(gh api "$url" | jq '.merged') == "true" ]]; then
      gh api -X DELETE "notifications/threads/$id"
    fi
  done
```
