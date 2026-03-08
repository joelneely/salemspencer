// Working on a golang version of Salem-Spencer sets,
// with the goal of beating the number of values currently shown
// on https://oeis.org/A262347 (contains a link to n = 1..140)
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"gospikes/salemspencer/ssdata"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

//
// changing the sets to hash map to see whether performance improves
//

type SearchResult struct {
	Weight int
	Sets map[ssdata.SSSet] bool
}

var best SearchResult

func search(ss *ssdata.SSSet, start int) {
	remaining := ss.Size - start + 1
	needed := best.Weight - ss.Weight

	if remaining < needed {
		return
	}

	if remaining == needed {
		openCount := 0
		for k := start; k <= ss.Size; k++ {
			if ss.IsOpenAt(k) {
				openCount++
			}
		}
		if openCount < needed {
			return
		}
	}

	if ss.Weight == best.Weight {
		best.Sets[*ss] = true
	}

	if ss.Weight > best.Weight {
		best.Weight = ss.Weight
		best.Sets = make(map[ssdata.SSSet]bool)
		best.Sets[*ss] = true
	}

	for i := start; i <= ss.Size; i++ {
		if ss.IsOpenAt(i) {
			blocked := ss.ApplyMoveLR(i)
			search(ss, i+1)
			ss.UndoMoveLR(i, blocked)
		}
	}
}

func findMaxSets(size int, began time.Time) {
	best = SearchResult{-1, make(map[ssdata.SSSet] bool)}
	ss := ssdata.NewSSSet(size)
	start := time.Now()
	search(&ss, 1)
	ended := time.Now()
	fmt.Printf(
		"%d | %d | %d | %v | %v\n",
		size, best.Weight, len(best.Sets),
		ended.Sub(began), ended.Sub(start),
	)
}

func mainSearch(from, limit int) {
	fmt.Printf("Salem-Spencer Search (revised Go implementation)\n")
	fmt.Printf("N | Size | Count | Total time | Unit time\n")
	fmt.Printf(":-: | :-: | :-: | :-: | :-:\n")
	began := time.Now()
	for size := from; size <= limit; size++ {
		findMaxSets(size, began)
	}
}

var limitFlag = flag.Int("limit", 75, "search N up to this value (max compile-time LIMIT)")
var fromFlag = flag.Int("from", 1, "search N starting from this value (default 1)")

func init() {
	flag.IntVar(limitFlag, "n", 75, "shorthand for -limit")
	flag.IntVar(fromFlag, "f", 1, "shorthand for -from")
}

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not create CPU profile: %v\n", err)
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if *limitFlag < 1 || *limitFlag > ssdata.LIMIT {
		fmt.Fprintf(os.Stderr, "limit must be between 1 and %d\n", ssdata.LIMIT)
		os.Exit(1)
	}
	if *fromFlag < 1 || *fromFlag > ssdata.LIMIT {
		fmt.Fprintf(os.Stderr, "from must be between 1 and %d\n", ssdata.LIMIT)
		os.Exit(1)
	}
	if *fromFlag > *limitFlag {
		fmt.Fprintf(os.Stderr, "from (%d) must be at most limit (%d)\n", *fromFlag, *limitFlag)
		os.Exit(1)
	}
	if *parallelFlag {
		mainSearchParallel(*fromFlag, *limitFlag)
	} else {
		mainSearch(*fromFlag, *limitFlag)
	}
}
