package test

import (
	"fmt"
	"testing"
)

func TestOne(t *testing.T) {
	// 2^6 - 1
	fmt.Printf("%b\n", 1<<41-1)
}
