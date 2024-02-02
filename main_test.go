package main

import (
	"fmt"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	b := []string{}
	s := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}
	for i := 2; i <= 12; i += 2 {
		a := strings.Join(s[i-2:i], " ")
		b = append(b, a)
	}
	fmt.Println(len(b))
}
