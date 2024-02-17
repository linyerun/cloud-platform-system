package test

import (
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/utils"
	"fmt"
	"testing"
	"time"
)

func TestJson(t *testing.T) {
	user := &models.User{Name: "哈哈哈"}
	obj, err := utils.NewDefaultTokenObject(user, time.Second*10)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(obj)
	data, err := obj.MarshalBinary()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
	obj.UnmarshalBinary(data)
	fmt.Println(data)
}
