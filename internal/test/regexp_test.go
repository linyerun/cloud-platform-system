package test

import (
	"fmt"
	"regexp"
	"testing"
)

func TestUse(t *testing.T) {
	re := regexp.MustCompile("^[a-zA-Z0-9]+$")
	matchString := re.MatchString("b704cd11893deaa086cc3e6626cd4d5d553哈哈哈4195442be14ded7c7262555194b70")
	fmt.Println(matchString)
}
