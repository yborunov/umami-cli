# umami-cli

CLI for self-hosted Umami analytics using the Umami API.

### Install (Homebrew)
```
brew install yborunov/tap/umami-cli
```

## Quick start
```
# Set endpoint (required)
export UMAMI_URL="https://analytics.example.com"

# Login and save token
umami-cli auth login --username you --password secret

# List websites
umami-cli websites list
```

### Configure Umami URL
```
export UMAMI_URL="https://analytics.example.com"
```

### Authenticate

```
# Login and save token
umami-cli auth login --username you --password secret

# Verify token
umami-cli auth verify
```

### Commands

```
# List websites
umami-cli websites list

# List teams
umami-cli teams list

# List websites for a team
umami-cli teams websites <team-id>

# Analytics examples
umami-cli analytics active <website-id>
umami-cli analytics stats <website-id> --start-at 1704067200000 --end-at 1706745600000
umami-cli analytics pageviews <website-id> --start-at 1704067200000 --end-at 1706745600000 --unit day
umami-cli analytics metrics <website-id> --start-at 1704067200000 --end-at 1706745600000 --type path --limit 100
umami-cli analytics metrics-expanded <website-id> --start-at 1704067200000 --end-at 1706745600000 --type referrer --limit 100
umami-cli analytics events-series <website-id> --start-at 1704067200000 --end-at 1706745600000 --unit day
```

## Manual build

```
make build
```

## Installation

### Install with Go

```
go install github.com/yborunov/umami-cli/cmd/umami-cli@latest
```

### Install from source

```
git clone https://github.com/yborunov/umami-cli.git
cd umami-cli
make build
./bin/umami-cli
```

## Configuration

The CLI stores the API token at:

- macOS/Linux: `~/.config/umami-cli/config.json`

Environment variables:

- `UMAMI_URL` (required) – Umami base URL (the CLI appends `/api`)
- `UMAMI_USERNAME` – default username for `auth login`
- `UMAMI_PASSWORD` – default password for `auth login`
- `UMAMI_TOKEN` – override stored token

## Commands

```
umami-cli auth login --username <user> --password <pass>
umami-cli auth verify

umami-cli websites list

umami-cli teams list
umami-cli teams websites <team-id>

umami-cli analytics active <website-id>
umami-cli analytics events-series <website-id> [--start-at <ms>] [--end-at <ms>] [--unit <unit>] [--timezone <tz>] [filters]
umami-cli analytics metrics <website-id> --type <type> [--start-at <ms>] [--end-at <ms>] [--limit <n>] [--offset <n>] [filters]
umami-cli analytics metrics-expanded <website-id> --type <type> [--start-at <ms>] [--end-at <ms>] [--limit <n>] [--offset <n>] [filters]
umami-cli analytics pageviews <website-id> [--start-at <ms>] [--end-at <ms>] [--unit <unit>] [--timezone <tz>] [--compare <prev|yoy>] [filters]
umami-cli analytics stats <website-id> [--start-at <ms>] [--end-at <ms>] [filters]
```

Common analytics flags:

- `--start-at` and `--end-at` are optional and default to the last 24 hours (milliseconds since epoch).
- `--unit` supports `year`, `month`, `day`, `hour`, `minute`.
- Filters: `--path` `--referrer` `--title` `--query` `--browser` `--os` `--device` `--country` `--region` `--city` `--hostname` `--tag` `--distinct-id` `--segment` `--cohort`
- Metric types: `path` `entry` `exit` `title` `query` `referrer` `channel` `domain` `country` `region` `city` `browser` `os` `device` `language` `screen` `event` `hostname` `tag` `distinctId`

## Notes

- This CLI uses `/api/auth/login` to obtain a token and then sends it as a Bearer token for subsequent requests.
- It targets the same endpoints as the official Umami API client docs: https://umami.is/docs/api/api-client

## Project layout

```
cmd/umami-cli       // entrypoint
internal/cmd        // CLI commands
internal/client     // HTTP client
internal/config     // config loading/saving
internal/out        // output helpers
```
