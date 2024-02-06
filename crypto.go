package main

import (
	"crypto/aes"
	crand "crypto/rand"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func write(s []string) {
	f, err := os.Create("writed.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err := f.Write([]byte(strings.Join(s, " "))); err != nil {
		log.Fatal(err)
	}
}

func genOrder(num int64) int {
	rand.New(rand.NewSource(num))
	var res string
	for _, v := range rand.Perm(int(num)) {
		res += strconv.Itoa(v)
	}
	a, _ := strconv.Atoi(res)
	return a
}

func genDummy(num int64) string {
	rand.New(rand.NewSource(num))

	charset := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	result := make([]byte, num)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}
func genSecret(num int64) string {
	rand.New(rand.NewSource(num))

	charset := "0123456789"

	result := make([]byte, num)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}

func genIV() ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := crand.Read(iv); err != nil {
		return nil, err
	}
	return iv, nil
}
