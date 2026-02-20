# umami-cli

CLI for self-hosted Umami analytics using the Umami API.

## Quick start

```
# Build
make build

# Set endpoint (required)
export UMAMI_URL="https://analytics.example.com"

# Login and save token
./bin/umami auth login --username you --password secret

# Verify token
./bin/umami auth verify

# List websites
./bin/umami websites list

# List teams
./bin/umami teams list

# List websites for a team
./bin/umami teams websites <team-id>

# Analytics examples
./bin/umami analytics active <website-id>
./bin/umami analytics stats <website-id> --start-at 1704067200000 --end-at 1706745600000
./bin/umami analytics pageviews <website-id> --start-at 1704067200000 --end-at 1706745600000 --unit day
./bin/umami analytics metrics <website-id> --start-at 1704067200000 --end-at 1706745600000 --type path --limit 100
./bin/umami analytics metrics-expanded <website-id> --start-at 1704067200000 --end-at 1706745600000 --type referrer --limit 100
./bin/umami analytics events-series <website-id> --start-at 1704067200000 --end-at 1706745600000 --unit day
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
umami auth login --username <user> --password <pass>
umami auth verify

umami websites list

umami teams list
umami teams websites <team-id>

umami analytics active <website-id>
umami analytics events-series <website-id> [--start-at <ms>] [--end-at <ms>] [--unit <unit>] [--timezone <tz>] [filters]
umami analytics metrics <website-id> --type <type> [--start-at <ms>] [--end-at <ms>] [--limit <n>] [--offset <n>] [filters]
umami analytics metrics-expanded <website-id> --type <type> [--start-at <ms>] [--end-at <ms>] [--limit <n>] [--offset <n>] [filters]
umami analytics pageviews <website-id> [--start-at <ms>] [--end-at <ms>] [--unit <unit>] [--timezone <tz>] [--compare <prev|yoy>] [filters]
umami analytics stats <website-id> [--start-at <ms>] [--end-at <ms>] [filters]
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
cmd/umami           // entrypoint
internal/cmd        // CLI commands
internal/client     // HTTP client
internal/config     // config loading/saving
internal/out        // output helpers
```
