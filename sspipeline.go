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

var pipelineFlag = flag.Bool("pipeline", false, "use pipelined parallel DFS: N+1 starts on idle cores while N straggles")

func init() {
	flag.BoolVar(pipelineFlag, "pp", false, "shorthand for -pipeline")
}

// pipeJobState holds all per-N search state for the pipelined search.
// Keeping state per-N (rather than in globals) allows N and N+1 to run
// concurrently with fully isolated bestWeight, result sets, and WaitGroups.
type pipeJobState struct {
	size        int
	began       time.Time // start of the whole run (for total time column)
	started     time.Time // when this N's sub-problems were submitted
	bestWeight  int64     // accessed atomically; monotonically non-decreasing
	mu          sync.Mutex
	sets        map[ssdata.SSSet]bool
	wg          sync.WaitGroup
}

// pipeWorkItem pairs a sub-problem with the pipeJobState it belongs to,
// so workers in the shared pool know which N's state to update.
type pipeWorkItem struct {
	job   *pipeJobState
	ss    ssdata.SSSet
	start int
}

// searchPipe is the goroutine-safe DFS kernel for the pipelined search.
// It reads job.bestWeight atomically for fast pruning and acquires job.mu
// only when updating the shared result set.
//
// Correctness: stale atomic reads give a conservative (lower) pruning
// threshold, never over-pruning. Per-N isolation means N's bestWeight
// updates never affect N+1's pruning decisions.
func searchPipe(job *pipeJobState, ss *ssdata.SSSet, start int) {
	currentBest := int(atomic.LoadInt64(&job.bestWeight))
	remaining := ss.Size - start + 1
	needed := currentBest - ss.Weight

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

	if ss.Weight >= currentBest {
		job.mu.Lock()
		currentBest = int(atomic.LoadInt64(&job.bestWeight))
		switch {
		case ss.Weight > currentBest:
			atomic.StoreInt64(&job.bestWeight, int64(ss.Weight))
			job.sets = make(map[ssdata.SSSet]bool)
			job.sets[*ss] = true
		case ss.Weight == currentBest:
			job.sets[*ss] = true
		}
		job.mu.Unlock()
	}

	for i := start; i <= ss.Size; i++ {
		if ss.IsOpenAt(i) {
			blocked := ss.ApplyMoveLR(i)
			searchPipe(job, ss, i+1)
			ss.UndoMoveLR(i, blocked)
		}
	}
}

// mainSearchPipeline runs the pipelined parallel search for all N in [from, limit].
//
// Unlike the basic parallel mode (-p), which completes N entirely before starting
// N+1, this mode submits N+1's sub-problems to the shared worker pool without
// waiting for N to finish. Cores that go idle when N's work queue is exhausted
// (but a straggler sub-tree is still running) immediately pick up N+1's work.
//
// Pipeline design:
//   - A persistent pool of workers drains a shared channel of pipeWorkItems.
//   - An enqueuer goroutine submits sub-problems for each N in sequence and
//     moves on immediately, so N+1's items enter the pool while N straggles.
//   - The main goroutine waits for each N's WaitGroup in order before printing,
//     so output is always sequential regardless of overlap.
func mainSearchPipeline(from, limit int) {
	numWorkers := runtime.GOMAXPROCS(0)

	// Buffer large enough to hold two full N values of sub-problems so the
	// enqueuer can run ahead and have N+1's work ready when N's queue drains.
	work := make(chan pipeWorkItem, limit*limit+numWorkers)

	// Start persistent worker pool. Workers serve all N values from the shared
	// channel and exit when the channel is closed after all jobs complete.
	var poolDone sync.WaitGroup
	poolDone.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		go func() {
			defer poolDone.Done()
			for item := range work {
				searchPipe(item.job, &item.ss, item.start)
				item.job.wg.Done()
			}
		}()
	}

	fmt.Printf("Salem-Spencer Search (pipelined parallel Go implementation, %d workers)\n", numWorkers)
	fmt.Printf("N | Size | Count | Total time | Unit time\n")
	fmt.Printf(":-: | :-: | :-: | :-: | :-:\n")

	began := time.Now()

	// Buffered to hold all job handles so the enqueuer never blocks here.
	jobs := make(chan *pipeJobState, limit-from+1)

	// Enqueuer: submit sub-problems for each N without waiting for completion,
	// so N+1's work enters the pool while N's straggler is still running.
	go func() {
		for size := from; size <= limit; size++ {
			job := &pipeJobState{
				size:  size,
				began: began,
				sets:  make(map[ssdata.SSSet]bool),
			}
			atomic.StoreInt64(&job.bestWeight, -1)

			problems := generateSubProblems(size)
			job.wg.Add(len(problems))
			job.started = time.Now()
			for _, p := range problems {
				work <- pipeWorkItem{job, p.ss, p.start}
			}
			jobs <- job
		}
		close(jobs)
	}()

	// Result collector: wait for each N to complete (in order) then print.
	for job := range jobs {
		job.wg.Wait()
		ended := time.Now()

		job.mu.Lock()
		weight := int(atomic.LoadInt64(&job.bestWeight))
		count := len(job.sets)
		job.mu.Unlock()

		fmt.Printf("%d | %d | %d | %v | %v\n",
			job.size, weight, count,
			ended.Sub(job.began), ended.Sub(job.started))
	}

	// All jobs complete; drain workers and exit cleanly.
	close(work)
	poolDone.Wait()
}
