package test

import (
	"fmt"
	"regexp"
	"testing"
)

func TestUse(t *testing.T) {
	re := regexp.MustCompile("^[a-zA-Z0-9]+$")
	fmt.Println(re.MatchString("b704cd11893deaa086cc3e6626cd4d5d553哈哈哈4195442be14ded7c7262555194b70"))

	fmt.Println(re.MatchString("e43291d6464a8f90dd653e28ba623c831fbde1a9460c64c252cc759bdd53539c"))
}
