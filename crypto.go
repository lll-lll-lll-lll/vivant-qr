package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
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

func separate(iv string, encryptedTxt string, dummy string, order int) []string {
	ivSepa := []string{iv[len(iv)/2:], iv[:len(iv)/2]}
	var orderStr string = strconv.Itoa(order)
	num := len(encryptedTxt) / 8
	var s = make([]string, 10)
	for i, v := range orderStr[:len(orderStr)-2] {
		idx, _ := strconv.Atoi(string(v))
		startIndex := i * num
		endIndex := (i + 1) * num
		s[idx] = encryptedTxt[startIndex:endIndex]
	}
	for i, v := range orderStr[len(orderStr)-2:] {
		idx, _ := strconv.Atoi(string(v))
		s[idx] = ivSepa[i]
	}

	return s
}

func genOrder(num int) int {
	var res string
	for _, v := range rand.Perm(num) {
		res += strconv.Itoa(v)
	}
	a, _ := strconv.Atoi(res)
	return a
}

func genDummy(num int64) string {
	rand.New(rand.NewSource(num))

	charset := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ=^|"

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
func pkcs7Pad(data []byte) []byte {
	length := aes.BlockSize - (len(data) % aes.BlockSize)
	trailing := bytes.Repeat([]byte{byte(length)}, length)
	return append(data, trailing...)
}
func Encrypt(data []byte, key []byte) (iv []byte, encrypted []byte, err error) {
	iv, err = genIV()
	if err != nil {
		return nil, nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	padded := pkcs7Pad(data)
	encrypted = make([]byte, len(padded))
	cbcEncrypter := cipher.NewCBCEncrypter(block, iv)
	cbcEncrypter.CryptBlocks(encrypted, padded)
	return iv, encrypted, nil
}

func pkcs7Unpad(data []byte) []byte {
	dataLength := len(data)
	padLength := int(data[dataLength-1])
	return data[:dataLength-padLength]
}

// decrypted, _ := decrypt(encrypted, key, iv)
// fmt.Printf("Decrypted: %s\n", decrypted)
func decrypt(data []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	decrypted := make([]byte, len(data))
	cbcDecrypter := cipher.NewCBCDecrypter(block, iv)
	cbcDecrypter.CryptBlocks(decrypted, data)
	return pkcs7Unpad(decrypted), nil
}
