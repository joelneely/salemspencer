// checkout functions to verify that Move can handle moves in any order,
// while MoveLR works properly only for left-to-right traversal
//
package ssdata

import (
	"fmt"
	"testing"
//	"time"
//	"gospikes/salemspencer/ssdata"
)

func show(caption string, s0, s1 SSSet) {
	fmt.Printf("%12s %v : %v -> %v / %v\n", caption, s0, s1, s0.Equals(s1), s1.Equals(s0))
}

func checkSSSetEquality(t *testing.T, s0, s1 SSSet, expected bool) {
	var description string
	if expected {
		description = "equality"
	} else {
		description = "difference"
	}
	if s0.Equals(s1) != expected {
		t.Errorf("expected %s! %v : %v -> %v", description, s0, s1, s0.Equals(s1))
	}
}

func expectEqualSSSets(t *testing.T, caption string, s0, s1 SSSet) {
	show(caption, s0, s1)
	checkSSSetEquality(t, s0, s1, true)
	checkSSSetEquality(t, s1, s0, true)
}

func expectDifferentSSSets(t *testing.T, caption string, s0, s1 SSSet) {
	show(caption, s0, s1)
	checkSSSetEquality(t, s0, s1, false)
	checkSSSetEquality(t, s1, s0, false)
}

func expectWeight(t *testing.T, s SSSet, expected int) {
	if s.Weight != expected {
		t.Errorf("expected weight of %d but observed %d in %v", expected, s.Weight, s)
	}
}

func TestMoves(t *testing.T) {
	sFull := NewSSSet(7)
	sLtoR := NewSSSet(7)
	sAlso := NewSSSet(7)

	fmt.Printf("\nnew instances should be equal\n")

	expectEqualSSSets(t, "sFull:sLtoR", sFull, sLtoR)
	expectEqualSSSets(t, "sFull:sAlso", sFull, sAlso)
	expectWeight(t, sFull, 0)
	expectWeight(t, sLtoR, 0)
	expectWeight(t, sAlso, 0)

	fmt.Printf("\nfirst move doesn't show difference between full and L-to-R moves\n")

	sLtoR, _ = sLtoR.MoveLR(3)
	sFull, _ = sFull.Move(3)
	expectEqualSSSets(t, "sFull:sLtoR", sFull, sLtoR)
	expectWeight(t, sLtoR, 1)
	expectWeight(t, sFull, 1)

	fmt.Printf("\nsecond move does show difference between full and L-to-R moves\n")

	sLtoR, _ = sLtoR.MoveLR(5)
	sFull, _ = sFull.Move(5)
	expectDifferentSSSets(t, "sFull:sLtoR", sFull, sLtoR)
	expectWeight(t, sLtoR, 2)
	expectWeight(t, sFull, 2)

	fmt.Printf("\npartial move sequence yields different sets\n")

	expectDifferentSSSets(t, "sFull:sAlso", sFull, sAlso)

	sAlso, _ = sAlso.Move(5)
	expectDifferentSSSets(t, "sFull:sAlso", sFull, sAlso)
	expectWeight(t, sAlso, 1)

	fmt.Printf("\nequivalent move sequence yields equal sets for full moves\n")

	sAlso, _ = sAlso.Move(3)
	expectEqualSSSets(t, "sFull:sAlso", sFull, sAlso)
	expectWeight(t, sAlso, 2)
}
