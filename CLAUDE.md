# CLAUDE.md — Launchdata Codebase Guide

## Project Overview

**Launchdata** is a Go CLI tool that downloads, parses, and interactively browses rocket launch data scraped from Wikipedia. It converts Wikipedia spaceflight launch tables into structured JSON and provides a terminal UI (TUI) browser.

- **Language:** Go 1.18 (uses generics)
- **Data source:** Wikipedia spaceflight articles via the wikitable2json.com API
- **Key frameworks:** [Cobra](https://github.com/spf13/cobra) (CLI), [Bubble Tea](https://github.com/charmbracelet/bubbletea) (TUI)

---

## Repository Structure

```
launchdata/
├── main.go           # Entry point — creates and runs the root Cobra command
├── cmd/              # Cobra command definitions
│   ├── root.go       # Root command, persistent --dry-run flag, litter config
│   ├── cache.go      # `cache` subcommand: download and write JSON data
│   └── browse.go     # `browse` subcommand: open interactive TUI
├── parse/            # Core Wikipedia scraping and parsing logic
│   ├── parse.go      # Main parsing: RocketData, PayloadData, AllLaunchData structs
│   ├── url.go        # Wikipedia API URL construction (pre-2021 vs 2021+ logic)
│   ├── timestamp.go  # Timestamp parsing for Wikipedia's many date formats
│   └── testdata/     # Golden test files for parser tests
├── bubble/           # Bubble Tea TUI implementation
│   ├── bubble.go     # Main TUI model; keyboard bindings, viewport, detail view
│   ├── item.go       # MyItem wraps RocketData to implement list.Item
│   └── delegate.go   # ItemDelegate for custom list item rendering
├── list/             # Generic Bubble Tea list component (forked/custom)
│   └── list.go       # Model[I Item] — supports filtering, fuzzy search, pagination
├── jsonio/           # JSON file I/O and HTTP fetch helpers
├── config/           # Config struct (DryRun flag) loaded from Cobra flags
├── cli/              # Terminal utilities (ClearScreen)
├── slices/           # Generic slice utilities (Reverse, Insert, Delete)
├── data/             # Pre-generated JSON data files: launchdata-YYYY.json (1951–2022)
├── go.mod            # Module: github.com/vsinha/launchdata, Go 1.18
├── go.sum            # Dependency lock file
├── justfile          # Task runner (build, test, cache-all, smoke, test-accept)
└── .github/workflows/go.yml  # CI: build, test, regenerate data, auto-commit
```

---

## Development Commands

All common tasks are defined in the `justfile`. Run with [`just`](https://github.com/casey/just).

| Command | Description |
|---|---|
| `just build` | Run `go vet`, `staticcheck`, then `go build` |
| `just test` | Run all tests (`go test ./...`) |
| `just test-accept` | Update golden test files (`GOLDEN_UPDATE=true go test ./...`) |
| `just smoke` | Quick functional test: cache years 2020–2022 |
| `just cache-all` | Download and write all data files (1951–2022) |

**Direct Go commands:**
```bash
go build ./...               # Build
go test ./...                # Test
go run . browse 2022         # Run TUI browser for 2022
go run . cache -y 2022       # Cache single year
go run . cache -s 2020 -e 2022  # Cache a range of years
go run . cache all           # Cache all years (1951–2022)
go run . cache --dry-run -y 2022  # Dry run (no network, no writes)
```

---

## CLI Commands

```
launchdata [--dry-run]
  cache [-y year] [-s start] [-e end] [--output-dir dir]
  cache all [--output-dir dir]
  browse [year]
```

- `--dry-run`: Global flag; disables HTTP requests and file writes for safe testing.
- `--output-dir` / `-o`: Directory for output JSON (default: `./data`).

---

## Key Data Structures (`parse/parse.go`)

```go
type AllLaunchData struct {
    OrbitalFlights    []RocketData
    SuborbitalFlights []RocketData
}

type RocketData struct {
    Timestamp             TimeData
    Rocket                string
    FlightNumber          string
    LaunchSite            string
    LaunchServiceProvider string
    Notes                 string
    Payload               []PayloadData
}

type PayloadData struct {
    Payload  string
    Operator string
    Orbit    string
    Function string
    Decay    string
    Outcome  string
    Cubesat  bool
}
```

Output JSON is stored as `data/launchdata-YYYY.json` using this schema.

---

## Wikipedia Parsing Strategy (`parse/parse.go`)

Wikipedia launch tables have a two-level row structure:
- **Main rows** (non-empty column 1): One launch per row with rocket, site, provider, etc.
- **Payload rows** (empty column 1, indented): Associated payloads belonging to the preceding launch.
- **Note rows**: Same value repeated across all columns — filtered out.
- **Month headers / navigation rows**: Filtered by detecting arrows or non-date text.

**URL logic** (`parse/url.go`):
- Pre-2021: Single annual Wikipedia page `{year}_in_spaceflight`
- 2021 and later: Two pages split by half-year (January–June, July–December)

**Timestamp parsing** (`parse/timestamp.go`):
- Tries 7+ format strings to parse Wikipedia dates
- Handles `TBD` entries, "Early/Mid/Late" descriptors, and timezone variations
- `TimeData.LaunchedAlready()` returns whether the launch time is in the past

---

## Testing

Tests use Go's `testing` package with:
- [`testify`](https://github.com/stretchr/testify) — assertions (`require`, `assert`)
- [`go-golden`](https://github.com/jimeh/go-golden) — golden file tests for parser output
- [`go-cmp`](https://github.com/google/go-cmp) — deep comparison with custom transformers

**Golden file tests** (`parse/testdata/`):
- Test inputs are recorded Wikipedia API responses
- Test outputs are expected parsed structs
- Update with `GOLDEN_UPDATE=true go test ./...` (or `just test-accept`) when parser behavior intentionally changes

**Cubesat detection** is done by checking for Unicode characters `⚀` and `▫` in payload strings — not keywords.

---

## CI/CD (`.github/workflows/go.yml`)

Triggered on push/PR to `main`:
1. `go build -v ./...`
2. `go test -v ./...`
3. `go run . cache all --output-dir ./data` — regenerates all data files
4. Auto-commits generated data via `stefanzweifel/git-auto-commit-action@v4`

---

## Code Conventions

- Standard Go formatting (`gofmt`); enforced by `go vet` and `staticcheck`
- Exported types and functions are capitalized and documented
- `--dry-run` support flows through `config.Config` to `jsonio` and `parse` layers
- Generic slice utilities in `slices/` use Go 1.18+ type parameters (`[E any]`)
- The custom `list/` component uses generics (`Model[I Item]`) for type safety
- TODO comments exist in the codebase for future work — leave them in place unless addressing them

**Known TODOs in the code:**
- `parse/parse.go`: Add a command to load from file instead of HTTP
- `cmd/root.go`: Check the cache directory for what years are already cached
- `cli/cli.go`: Commented-out `less` pager integration

---

## External Dependencies

| Package | Purpose |
|---|---|
| `github.com/spf13/cobra` | CLI framework |
| `github.com/charmbracelet/bubbletea` | TUI framework |
| `github.com/charmbracelet/bubbles` | Bubble Tea UI components |
| `github.com/charmbracelet/lipgloss` | TUI styling |
| `github.com/sahilm/fuzzy` | Fuzzy search for list filtering |
| `github.com/deckarep/golang-set/v2` | Set data structure |
| `github.com/sanity-io/litter` | Pretty-printing Go structs |
| `github.com/stretchr/testify` | Test assertions |
| `github.com/jimeh/go-golden` | Golden file tests |
| `github.com/google/go-cmp` | Deep equality comparison |

---

## Environment Variables

| Variable | Purpose |
|---|---|
| `GOLDEN_UPDATE=true` | When set, updates golden test files instead of comparing against them |

No `.env` files — all runtime configuration is via CLI flags.
