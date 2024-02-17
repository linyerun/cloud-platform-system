package test

import (
	"cloud-platform-system/internal/utils"
	"fmt"
	"testing"
)

func TestSnowFlakeGenerate(t *testing.T) {
	m := make(map[string]struct{})
	for i := 0; i < 10; i++ {
		str := utils.GetSnowFlakeIdAndBase64()
		m[str] = struct{}{}
		fmt.Println(str)
	}
	fmt.Println(len(m))
}

func TestInt64To8Byte(t *testing.T) {
	id := utils.GetSnowFlakeId()
	Int64Print(id)
	fmt.Println()
	buf := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		buf[i] = byte(id)
		id >>= 8
	}
	for i := 0; i < 8; i++ {
		BytePrint(buf[i])
	}
	fmt.Println()
}

func Int64Print(x int64) {
	for i := 63; i >= 0; i-- {
		fmt.Print(x >> i & 1)
	}
}

func BytePrint(x byte) {
	for i := 7; i >= 0; i-- {
		fmt.Print(x >> i & 1)
	}
}
