package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"gospikes/salemspencer/ssdata"
)

var parallelFlag = flag.Bool("parallel", false, "use parallel DFS search instead of sequential")

func init() {
	flag.BoolVar(parallelFlag, "p", false, "shorthand for -parallel")
}

// subProblem is a work item: an SSSet state and the position from which
// to continue the left-to-right DFS.
type subProblem struct {
	ss    ssdata.SSSet
	start int
}

var (
	// parBestWeight is the globally-known best weight, updated atomically.
	// Workers read it lock-free for pruning; writes happen inside parMu.
	parBestWeight int64
	parMu         sync.Mutex
	parSets       map[ssdata.SSSet]bool
)

// searchP is the goroutine-safe DFS kernel. It mirrors search() but reads
// parBestWeight atomically for fast pruning and acquires parMu only when
// updating the shared result.
//
// Correctness note: any non-leaf node at weight W has children at weight W+1
// (each MoveLR increments Weight by exactly 1). If W is the final maximum,
// no child can exceed W, yet children are never pruned (pruning fires only
// when achievable < currentBest = W). This contradiction means non-leaves at
// the final maximum cannot exist, so every entry remaining in parSets at
// termination is a genuine DFS leaf.
func searchP(ss ssdata.SSSet, start int) {
	currentBest := int(atomic.LoadInt64(&parBestWeight))
	if ss.Weight+ss.Size-start+1 < currentBest {
		return
	}

	if ss.Weight >= currentBest {
		parMu.Lock()
		currentBest = int(atomic.LoadInt64(&parBestWeight))
		switch {
		case ss.Weight > currentBest:
			atomic.StoreInt64(&parBestWeight, int64(ss.Weight))
			parSets = make(map[ssdata.SSSet]bool)
			parSets[ss] = true
		case ss.Weight == currentBest:
			parSets[ss] = true
		}
		parMu.Unlock()
	}

	for i := start; i <= ss.Size; i++ {
		if ss.IsOpenAt(i) {
			next, _ := ss.MoveLR(i)
			searchP(next, i+1)
		}
	}
}

// generateSubProblems pre-generates all DFS states at depth 2 (two moves
// already applied), providing O(N²) independent tasks for better load
// distribution across workers. Depth-1 states with no valid second move
// are included directly so no work is missed.
func generateSubProblems(size int) []subProblem {
	initial := ssdata.NewSSSet(size)
	var problems []subProblem

	for i := 1; i <= size; i++ {
		if !initial.IsOpenAt(i) {
			continue
		}
		ss1, _ := initial.MoveLR(i)
		hasChild := false
		for j := i + 1; j <= size; j++ {
			if ss1.IsOpenAt(j) {
				ss2, _ := ss1.MoveLR(j)
				problems = append(problems, subProblem{ss2, j + 1})
				hasChild = true
			}
		}
		if !hasChild {
			// ss1 has no open positions past i; treat it as a leaf task.
			problems = append(problems, subProblem{ss1, i + 1})
		}
	}
	return problems
}

// findMaxSetsParallel runs the parallel search for sets of size N and prints
// one Markdown table row in the same format as findMaxSets.
func findMaxSetsParallel(size int, began time.Time) {
	parMu.Lock()
	atomic.StoreInt64(&parBestWeight, -1)
	parSets = make(map[ssdata.SSSet]bool)
	parMu.Unlock()

	problems := generateSubProblems(size)
	work := make(chan subProblem, len(problems))
	for _, p := range problems {
		work <- p
	}
	close(work)

	start := time.Now()
	numWorkers := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range work {
				searchP(p.ss, p.start)
			}
		}()
	}
	wg.Wait()

	ended := time.Now()
	parMu.Lock()
	weight := int(atomic.LoadInt64(&parBestWeight))
	count := len(parSets)
	parMu.Unlock()

	fmt.Printf(
		"%d | %d | %d | %v | %v\n",
		size, weight, count,
		ended.Sub(began), ended.Sub(start),
	)
}

// mainSearchParallel is the outer loop for the parallel search mode.
func mainSearchParallel(from, limit int) {
	fmt.Printf("Salem-Spencer Search (parallel Go implementation, %d workers)\n", runtime.GOMAXPROCS(0))
	fmt.Printf("N | Size | Count | Total time | Unit time\n")
	fmt.Printf(":-: | :-: | :-: | :-: | :-:\n")
	began := time.Now()
	for size := from; size <= limit; size++ {
		findMaxSetsParallel(size, began)
	}
}
