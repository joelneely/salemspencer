// Working on a golang version of Salem-Spencer sets,
// with the goal of beating the number of values currently shown
// on https://oeis.org/A262347 (contains a link to n = 1..140)
package main

import (
	"fmt"
	"time"
	"gospikes/salemspencer/ssdata"
)

//
// changing the sets to hash map to see whether performance improves
//

type SearchResult struct {
	Weight int
	Sets map[ssdata.SSSet] bool
}

var best SearchResult

func search(ss ssdata.SSSet, start int, prefix string) {

	if ss.Weight + ss.Size - start + 1 < best.Weight {
		return
	}

	if ss.Weight == best.Weight {
		best.Sets[ss] = true
	}

	if ss.Weight > best.Weight {
		best.Weight = ss.Weight
		best.Sets = make(map[ssdata.SSSet] bool)
		best.Sets[ss] = true
	}

	for i := start; i <= ss.Size; i++ {
		if ss.IsOpenAt(i) {
			next, _ := ss.MoveLR(i)
			search(next, i + 1, prefix + "\t")
		}
	}
}

func findMaxSets(size int, began time.Time) {
	best = SearchResult{-1, make(map[ssdata.SSSet] bool)}
	ss := ssdata.NewSSSet(size)
	start := time.Now()
	search(ss, 1, "\t\t\t\t\t")
	ended := time.Now()
	fmt.Printf(
		"%d | %d | %d | %v | %v\n",
		size, best.Weight, len(best.Sets),
		ended.Sub(began), ended.Sub(start),
	)
}

func mainSearch() {
	fmt.Printf("Salem-Spencer Search (revised Go implementation)\n")
	fmt.Printf("N | Size | Count | Total time | Unit time\n")
	fmt.Printf(":-: | :-: | :-: | :-: | :-:\n")
	began := time.Now()
	for size := 1; size <=ssdata.LIMIT; size++ {
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
