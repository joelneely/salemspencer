# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Purpose

A Go project to find all maximal Salem-Spencer sets (subsets of `{1..N}` containing no arithmetic progressions) for successive values of N, with the goal of beating the published record on [OEIS A262347](https://oeis.org/A262347). The compile-time ceiling is `LIMIT=150` in `ssdata/ssset.go`; the runtime search limit defaults to 75 and is controlled by the `-limit`/`-n` flag.

README.md contains four timing tables:
- **MacBook Pro (M2 Pro):** N=1–65, total run ~1h36m
- **Mac Studio (M3 Ultra):** N=1–70, total run ~5h56m (stopped before N=71 to stay within an 8-hour budget)
- **Mac Studio (M3 Ultra, hot-path opt.):** N=1–55, after eliminating redundant array copy/scan and per-node string allocation; unit times ~3× faster than the prior M3 Ultra run
- **Mac Studio (M3 Ultra, pruning + in-place undo):** N=1–55, after lazy tight-bound pruning and in-place undo (eliminating the per-step 151-byte duffcopy); a further ~25–27% speedup over the hot-path baseline

The M3 Ultra is consistently **17–21% faster** than the M2 Pro for this single-threaded search. Runtime grows roughly exponentially with N; each step is approximately 1.3–1.5× the previous. Current sequential baseline unit times (M3 Ultra, latest optimization):
- N=50: 4.5s | N=55: 17.0s

## Commands

```bash
# Build and run (sequential)
go run ssmain.go

# Build binary
go build -o salemspencer .

# Run sequential search (default: N=1..75)
./salemspencer

# Run with a custom limit (long and short forms)
./salemspencer -limit 50
./salemspencer -n 50

# Start from a value other than 1 (long and short forms)
./salemspencer -from 60
./salemspencer -f 60

# Run parallel search (uses runtime.GOMAXPROCS(0) goroutines)
./salemspencer -parallel
./salemspencer -p

# Combine flags
./salemspencer -p -f 60 -n 75

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

**`ssmain.go`** — Contains:
- `SearchResult` struct: `Weight` (best set size found so far) + `Sets` (`map[SSSet]bool` of all maximal sets at that size).
- `best`: package-level `SearchResult`, reset at the start of each `findMaxSets` call.
- `search(ss, start)`: recursive DFS. Prunes when `ss.Weight + ss.Size - start + 1 < best.Weight` (upper bound on achievable weight from this state is less than current best). Accumulates results in `best.Sets`.
- `findMaxSets(size, began)`: resets `best`, runs `search` for a single N, prints one Markdown table row.
- `limitFlag`: `flag.Int` for `-limit` (default 75); `-n` is a shorthand alias. Validated in `main()` to be in `[1, LIMIT]`.
- `fromFlag`: `flag.Int` for `-from` (default 1); `-f` is a shorthand alias. Validated to be in `[1, LIMIT]` and ≤ `limitFlag`.
- `mainSearch(from, limit int)`: prints the table header, then loops N=`from` to `limit` calling `findMaxSets`.
- `best.Sets` deduplicates maximal sets automatically because `SSSet` is a comparable value type usable as a map key.

**`ssparallel.go`** — Parallel alternative, selected by `-parallel` flag. Contains:
- `parallelFlag`: `flag.Bool` registered at package init; `main()` dispatches on it.
- `parBestWeight` (atomic `int64`): globally-known best weight, read lock-free in the hot path for pruning, written inside `parMu`.
- `parMu` / `parSets`: mutex-protected map of maximal sets at the current best weight.
- `searchP(ss, start)`: goroutine-safe DFS kernel mirroring `search()`. Reads `parBestWeight` atomically for pruning; acquires `parMu` only when `ss.Weight >= currentBest` to update the shared result.
- `generateSubProblems(size)`: pre-runs two DFS levels to produce O(N²) independent `subProblem` values, providing enough granularity for good load distribution across workers.
- `findMaxSetsParallel(size, began)`: fills a buffered channel with sub-problems, launches `runtime.GOMAXPROCS(0)` worker goroutines that drain it dynamically, waits with `sync.WaitGroup`, then prints one Markdown table row.
- `mainSearchParallel(from, limit int)`: outer loop equivalent of `mainSearch(from, limit)`.

## Key Design Decisions

- `SSSet` uses a fixed-size array `[MAXLENGTH]uint8` (not a slice) so it can be used as a map key for deduplication of maximal sets.
- `LIMIT=150` in `ssdata/ssset.go` is the compile-time ceiling; `MAXLENGTH = LIMIT + 1 = 151` (1-indexed arrays). The runtime search limit is set by `-limit`/`-n` (default 75, must be ≤ `LIMIT`).
- The search uses `MoveLR` (not `Move`) because the recursion always proceeds left-to-right. `Move` is retained for correctness testing.
- `TestMoves` in `ssdata/ssset_test.go` verifies that `Move` produces order-independent results (applying moves in different orders yields equal sets) while `MoveLR` is intentionally order-dependent (left-to-right only). Always run `go test ./ssdata/` before and after any changes to `Move` or `MoveLR`.
- Performance is the primary concern; the time to search grows roughly exponentially with N.

### Parallel search design
- Depth-2 pre-generation gives O(N²) sub-problems (~2800 for N=75), enough granularity for dynamic load distribution across many workers. Depth-1 (~75 items) would leave workers idle while one large sub-tree finishes.
- `parBestWeight` is monotonically non-decreasing. Stale atomic reads therefore give a lower (more conservative) pruning threshold — never over-prunes — so lock-free reads are safe.
- The lock (`parMu`) is only acquired when `ss.Weight >= currentBest`, which happens only at near-optimal nodes. Contention is low even with many goroutines.
- Non-leaf nodes are briefly added to `parSets` (matching the sequential behaviour) but are always superseded: every non-leaf has a child at `Weight+1`, which resets the map. Only true DFS leaves survive in the final result.
- Measured on M3 Ultra (28 workers): **~10.8× speedup** vs sequential at N=50 after hot-path optimization (6.2s sequential → 0.575s parallel). Pre-optimization baseline was ~16–17× (64s → ~4s); sequential improved ~10× while parallel improved ~7×, reducing parallelism efficiency from ~60% to ~39% of linear scaling across 28 workers. After the pruning + in-place undo optimization the sequential baseline at N=50 dropped to 4.5s; parallel speedup ratio not yet re-measured.

## Running Long Jobs on macOS

`timeout` is not available on macOS by default (produces exit code 127). Instead:
1. Build the binary first: `go build -o salemspencer .`
2. Run with the absolute path using `run_in_background=true` in Claude Code — output goes to the background task's capture file.
3. Stop manually with `TaskStop` when the desired time limit is reached.
4. The program prints each row as it completes, so partial output is usable even if the run is stopped early.

To stop a foreground run after a fixed number of rows, pipe through `head` (SIGPIPE terminates cleanly):
```bash
./salemspencer | head -$((K+3))         # capture through N=K, starting from 1
./salemspencer -f F | head -$((K-F+4))  # capture through N=K, starting from F
```
Prefer `-from`/`-n` flags to skip already-known values rather than piping through `head`.
