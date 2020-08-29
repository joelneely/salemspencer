// separating the core data structure into its own package
package ssdata

import (
	"fmt"
)

const (
	OPEN    = 0
	BLOCKED = 1
	CLOSED  = 2
)

const (
	LIMIT = 75
	MAXLENGTH = LIMIT + 1
)

type SSSet struct {
	Size, Weight int
	Data [MAXLENGTH]uint8
}

func MakeSSSet(d [MAXLENGTH]uint8, size int) SSSet {
	dd := d
	wt := 0
	for _, state := range dd {
		if state == CLOSED {
			wt++
		}
	}
	return SSSet{size, wt, dd}
}

func NewSSSet(size int) SSSet {
	empty := [MAXLENGTH]uint8{}
	return MakeSSSet(empty, size)
}

// TODO: examine the effect of modifying the methods to take a pointer
// in other words, will that eliminate copying and speed up the code?

func (this SSSet) Equals(that SSSet) bool {
	if this.Size != that.Size {
		return false;
	}
	if this.Weight != that.Weight {
		return false;
	}
	for i := 1; i <= this.Size; i++ {
		if this.Data[i] != that.Data[i] {
			return false;
		}
	}
	return true;
}

func (this SSSet) IsClosedAt(move int) bool {
	return this.Data[move] == CLOSED
}

func (this SSSet) IsOpenAt(move int) bool {
	return this.Data[move] == OPEN
}

func checkBlock(dd *[MAXLENGTH]uint8, i, j int) {
	if dd[i] == CLOSED && dd[j] != CLOSED {
		dd[j] = BLOCKED
	}
}

// The general-case Move method looks both ways from the
// current move position, checking for CLOSED positions that
// cause other positions to become BLOCKED.

func (this SSSet) Move(move int) (SSSet, bool) {
	if !this.IsOpenAt(move) {
		return this, false
	}
	dd := this.Data
	//
	// a closed position on one side of move closes its mirror image
	//
	for i, j := move - 1, move + 1; 1 <= i && j <= this.Size; i, j = i - 1, j + 1 {
		checkBlock(&dd, i, j)
		checkBlock(&dd, j, i)
	}
	//
	// check for triples to the left of move
	//
	for i, j := move - 2, move - 1; 1 <= i; i, j = i - 2, j - 1 {
		checkBlock(&dd, i, j)
		checkBlock(&dd, j, i)
	}
	//
	// check for triples to the right of move
	//
	for i, j := move + 1, move + 2; j <= this.Size; i, j = i + 1, j + 2 {
		checkBlock(&dd, i, j)
		checkBlock(&dd, j, i)
	}
	dd[move] = CLOSED
	return MakeSSSet(dd, this.Size), true
}

//
// This special-case version of Move only considers left-to-right
// traversals, eliminating tests for CLOSED positions to the right
// of the current move position.
//

func (this SSSet) MoveLR(move int) (SSSet, bool) {
	if !this.IsOpenAt(move) {
		return this, false
	}
	dd := this.Data
	for i, j := move - 1, move + 1; 1 <= i && j <= this.Size; i, j = i - 1, j + 1 {
		if this.IsClosedAt(i) {
			dd[j] = BLOCKED
		}
	}
	dd[move] = CLOSED
	return MakeSSSet(dd, this.Size), true
}

func (s SSSet) String() string {
	return fmt.Sprintf("{%d %d %v}", s.Size, s.Weight, s.Data[1:s.Size+1])
}
