// Starting to work on the golang version of salem-spencer sets,
// with the goal of beating the number of values currently shown
// on https://oeis.org/A262347 (contains a link to n = 1..140)
package main

import (
	"fmt"
	"time"
	"gospikes/salemspencer/ssdata"
)

const (
	OPEN    = 0
	BLOCKED = 1
	CLOSED  = 2
)

const (
	LIMIT = 45
	MAXLENGTH = LIMIT + 1
)

type SearchResult struct {
	Weight int
	Sets []ssdata.SSSet
}

var best SearchResult

func intMax(m, n int) int {
	if m >= n {
		return m
	} else {
		return n
	}
}

func search(ss ssdata.SSSet, start int, prefix string) {
	// fmt.Printf("%s%d in %v ; %v ", prefix, start, ss, best)
	if ss.Weight + ss.Size - start + 1 < best.Weight {
		return
	}
	if ss.Weight == best.Weight {
		found := false
		for _, s := range best.Sets {
			if ss.Equals(s) {
				found = false
				break
			}
		}
		if !found {
			best.Sets = append(best.Sets, ss)
		}
		// fmt.Printf("=,%5v -> %v", found, best)
	}
	if ss.Weight > best.Weight {
		best.Weight = ss.Weight
		best.Sets = append(best.Sets[:0], ss)
		// fmt.Printf(">       -> %v", best)
	}
	// fmt.Printf("\n")
	for i := start; i <= ss.Size; i++ {
		if ss.IsOpenAt(i) {
			next, _ := ss.MoveLR(i)
			search(next, i + 1, prefix + "\t")
		}
	}
	// fmt.Printf("%s <== %v\n", prefix, best)
}

func findMaxSets(size int, began time.Time) {
	best = SearchResult{-1, []ssdata.SSSet{}}
	ss := ssdata.NewSSSet(size)
//	start := time.Now()
	search(ss, 1, "\t\t\t\t\t")
	ended := time.Now()
	// fmt.Printf(
	// 	"%4d\t%4d\t%5d\t%v\t%v\n",
	// 	size, best.Weight, len(best.Sets),
	// 	ended.Sub(began), ended.Sub(start),
	// )
	fmt.Printf(
		"%d | %d | %d | %v\n",
		size, best.Weight, len(best.Sets),
		ended.Sub(began),
	)
}

func mainSearch() {
	fmt.Printf("Salem-Spencer Search (first Go implementation)\n")
	fmt.Printf("N | Size | Count | Total Time\n")
	fmt.Printf(":_: | :-: | :-: | :_:\n")
	began := time.Now()
	for size := 1; size <=LIMIT; size++ {
		findMaxSets(size, began)
	}
	// fmt.Printf("\nMaximal sets for %d\n", LIMIT)
	// for _, example := range best.Sets {
	// 	fmt.Printf("\t%v\n", example)
	// }
}

func main() {
	mainSearch()
}
