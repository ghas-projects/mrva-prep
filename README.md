# mrva-prep

A CLI tool that prepares an MRVA (Multi-Repository Variant Analysis) SQLite database for the Blazor WebAssembly reports UI. It creates indexes for faster queries, extracts dashboard aggregates into a lightweight JSON file for instant first-paint, and gzip-compresses the database for efficient browser delivery.

## Features

- **Index** - Creates indexes matching the query patterns used by the MRVA Reports UI, then runs `ANALYZE` and `VACUUM` to optimise the database.
- **Dashboard** - Runs aggregation queries and writes the results to a `dashboard.json` file so the UI can render immediately while the full database downloads in the background.
- **Compress** - Produces a `.gz` copy of the database using maximum gzip compression (level 9) for serving to the browser via `DecompressionStream`.

## Prerequisites

- [Go 1.25+](https://go.dev/doc/install)
- A CGO-enabled toolchain (required by the `go-sqlite3` driver)
- SQL Lite database created from `sql-transform` step

## Installation

```sh
go install github.com/ghas-projects/mrva-prep@latest
```

Or build from source:

```sh
git clone https://github.com/ghas-projects/mrva-prep.git
cd mrva-prep
go build -o mrva-prep .
```

## Usage

### Run all preparation steps

```sh
mrva-prep all --db mrva.db
```

This runs `index`, `dashboard`, and `compress` in sequence.

### Individual commands

```sh
# Create indexes on the database
mrva-prep index --db mrva.db

# Extract dashboard stats to JSON
mrva-prep dashboard --db mrva.db --output ./out

# Gzip-compress the database
mrva-prep compress --db mrva.db
```

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--db` | `-d` | `mrva.db` | Path to the SQLite database |
| `--output` | `-o` | `.` | Directory to write `dashboard.json` into (dashboard & all commands) |

## Project structure

```
cmd/           CLI command definitions (Cobra)
internal/
  dashboard/   Dashboard aggregation service
  index/       Index creation service
  models/      Shared data types (DashboardStats, etc.)
```

## Pipeline

`mrva-prep` is one component of the end-to-end MRVA reporting pipeline. The full flow is:

1. [**sarif-sql**](https://github.com/ghas-projects/sarif-sql) — Download SARIF artifacts and transform them into a normalized SQLite database.
2. **mrva-prep** (this repo) — Add query-optimized indexes, extract dashboard metrics, and compress the database.
3. [**mrva-reports**](https://github.com/advanced-security/mrva-reports) — Render the database as an interactive single-page dashboard in the browser.

A GitHub Actions workflow chains all three steps into a single automated run. See the [MRVA Documentation](https://github.com/advanced-security/mrva-documentation) for the full architecture and user guide.

## License

This project is licensed under the [MIT License](LICENSE).