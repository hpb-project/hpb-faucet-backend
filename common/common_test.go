package common

import (
	"fmt"
	"testing"
)

func TestGetCount(t *testing.T) {
	SetCount()
	for  {
		fmt.Println(GetCount())
	}
	for  {
		fmt.Println(GetCount())
	}
	for  {
		fmt.Println(GetCount())
	}
}
