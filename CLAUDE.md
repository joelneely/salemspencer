# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Purpose

A Go project to find all maximal Salem-Spencer sets (subsets of `{1..N}` containing no arithmetic progressions) for successive values of N, with the goal of beating the published record on [OEIS A262347](https://oeis.org/A262347). The current implementation handles N up to 75 (set via `LIMIT` in `ssdata/ssset.go`).

README.md contains two timing tables:
- **MacBook Pro (M2 Pro):** N=1–65, total run ~1h36m
- **Mac Studio (M3 Ultra):** N=1–70, total run ~5h56m (stopped before N=71 to stay within an 8-hour budget)

The M3 Ultra is consistently **17–21% faster** than the M2 Pro for this single-threaded search. Runtime grows roughly exponentially with N; each step is approximately 1.3–1.5× the previous.

## Commands

```bash
# Build and run
go run ssmain.go

# Build binary
go build -o salemspencer .

# Run all tests
go test ./...

# Run tests in the ssdata package (with verbose output)
go test -v ./ssdata/

# Run a specific test
go test -v ./ssdata/ -run TestMoves
```

Output is formatted as a Markdown table. When adding a new hardware run, append a new section to README.md in the same format as the existing tables, with an additional **vs [previous hardware]** column showing unit-time percentage change for rows where unit time exceeds one second.

## Architecture

The project has two layers:

**`ssdata/ssset.go`** — The core data structure (`SSSet`), used as a value type (struct with a fixed-size array `[MAXLENGTH]uint8`). Using an array (not a slice) enables `SSSet` to be used as a map key. Each element in `Data` is one of three states: `OPEN` (0), `BLOCKED` (1), or `CLOSED` (2). `Weight` counts how many positions are `CLOSED` (i.e., the size of the set built so far).

Two move methods exist:
- `Move(i)` — general-case; checks both directions from position `i` for arithmetic progression conflicts.
- `MoveLR(i)` — optimized for left-to-right traversal only; skips checking positions to the right of `i` (they haven't been visited yet), giving a significant performance advantage over `Move`.

**`ssmain.go`** — Recursive depth-first search (`search`) over all subsets, pruning branches where the remaining positions cannot improve on `best.Weight`. Maximal sets are accumulated in `best.Sets` (a `map[SSSet]bool`), which inherently deduplicates results because `SSSet` is a comparable value type.

## Key Design Decisions

- `SSSet` uses a fixed-size array `[MAXLENGTH]uint8` (not a slice) so it can be used as a map key for deduplication of maximal sets.
- `LIMIT` in `ssdata/ssset.go` controls the maximum N searched. `MAXLENGTH = LIMIT + 1` (1-indexed arrays).
- The search uses `MoveLR` (not `Move`) because the recursion always proceeds left-to-right. `Move` is retained for correctness testing.
- Performance is the primary concern; the time to search grows roughly exponentially with N.
