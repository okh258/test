package base

import (
	"fmt"
	"testing"
)

func TestDe(t *testing.T) {
	m := make(map[string]int)
	fmt.Printf("len(m): %v\n", len(m))
	m["1"] = 1
	m["2"] = 2
	m["3"] = 3
	fmt.Printf("len(m): %v\n", len(m))
}
