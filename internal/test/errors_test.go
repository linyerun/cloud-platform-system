package test

import (
	"github.com/pkg/errors"
	"testing"
)

func TestUseWrap(t *testing.T) {
	//err := errors.WithStack(GetError())
	//err := errors.Wrap(GetError(), "222")
	err := errors.New("111")
	println(err)
}

func GetError() error {
	return errors.New("111")
}
