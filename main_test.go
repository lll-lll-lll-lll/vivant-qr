package main

import (
	"fmt"
	"strconv"
	"testing"
)

func Test_Write(t *testing.T) {
	write([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"})
}

func Test(t *testing.T) {
	sp := []string{"0afj3GzarA==", "iGltBiiuE", "7OzXrHfAV", "Ng87OKrWA", "VU2nIJ8u/", "wENZ3ue1l", "WwDfA1pnu2ry", "46iPFZwf3", "Sw/YPoH2Y", "+hQye+C3h"}
	var correctNums = make([]string, 0, 11)
	for _, o := range "3814957206" {
		idx, _ := strconv.Atoi(string(o))
		fmt.Println(sp[idx])
		correctNums = append(correctNums, sp[idx])
	}
	fmt.Println("correctNums", correctNums)
	fmt.Println("coorectiv", "WwDfA1pnu2ry0afj3GzarA==")
}
