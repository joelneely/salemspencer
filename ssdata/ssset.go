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
	LIMIT     = 150
	MAXLENGTH = LIMIT + 1
)

type SSSet struct {
	Size, Weight int
	Data [MAXLENGTH]uint8
}

func NewSSSet(size int) SSSet {
	return SSSet{Size: size}
}

func (this *SSSet) Equals(that SSSet) bool {
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

func (this *SSSet) IsClosedAt(move int) bool {
	return this.Data[move] == CLOSED
}

func (this *SSSet) IsOpenAt(move int) bool {
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

func (this *SSSet) Move(move int) (SSSet, bool) {
	if !this.IsOpenAt(move) {
		return *this, false
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
	return SSSet{this.Size, this.Weight + 1, dd}, true
}

//
// This special-case version of Move only considers left-to-right
// traversals, eliminating tests for CLOSED positions to the right
// of the current move position.
//

func (this *SSSet) MoveLR(move int) (SSSet, bool) {
	if !this.IsOpenAt(move) {
		return *this, false
	}
	dd := this.Data
	for i, j := move - 1, move + 1; 1 <= i && j <= this.Size; i, j = i - 1, j + 1 {
		if this.IsClosedAt(i) {
			dd[j] = BLOCKED
		}
	}
	dd[move] = CLOSED
	return SSSet{this.Size, this.Weight + 1, dd}, true
}

// ApplyMoveLR is the in-place counterpart to MoveLR. It mutates the receiver,
// marking position move as CLOSED and any newly-blocked positions as BLOCKED,
// and incrementing Weight. It returns the list of positions changed from OPEN
// to BLOCKED so the caller can undo them. The caller must guarantee move is OPEN.
func (this *SSSet) ApplyMoveLR(move int) []int {
	var blocked []int
	for i, j := move-1, move+1; 1 <= i && j <= this.Size; i, j = i-1, j+1 {
		if this.IsClosedAt(i) && this.Data[j] == OPEN {
			this.Data[j] = BLOCKED
			blocked = append(blocked, j)
		}
	}
	this.Data[move] = CLOSED
	this.Weight++
	return blocked
}

// UndoMoveLR reverses an ApplyMoveLR call, restoring the receiver to the
// state it was in before that call.
func (this *SSSet) UndoMoveLR(move int, blocked []int) {
	this.Data[move] = OPEN
	this.Weight--
	for _, j := range blocked {
		this.Data[j] = OPEN
	}
}

func (s SSSet) String() string {
	return fmt.Sprintf("{%d %d %v}", s.Size, s.Weight, s.Data[1:s.Size+1])
}
