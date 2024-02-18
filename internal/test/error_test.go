package test

import (
	"fmt"
	"github.com/pkg/errors"
	"testing"
)

func TestErrorIs(t *testing.T) {
	err01 := errors.New("aaa")
	err02 := errors.Wrap(err01, "bbb")

	fmt.Println(errors.Is(err01, err02))
	fmt.Println(errors.Is(err02, err01))
}
