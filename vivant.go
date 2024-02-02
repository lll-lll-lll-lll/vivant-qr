package main

import (
	"encoding/base64"
	"encoding/hex"
	"strconv"
)

type VivantQR struct {
	cfg *Config
}

func (v *VivantQR) Generate() ([]string, error) {
	key, err := hex.DecodeString(v.cfg.SecretKey)
	if err != nil {
		return nil, err
	}
	iv, encrypted, err := Encrypt([]byte(vivant), key)
	if err != nil {
		return nil, err
	}
	encodedIV := base64.StdEncoding.EncodeToString(iv)
	encodedEncrypted := base64.StdEncoding.EncodeToString(encrypted)
	dummy := genDummy(16)
	return separate(encodedIV, encodedEncrypted, dummy, v.cfg.Order), nil
}

func separate(iv string, encryptedTxt string, dummy string, order int) []string {
	ivSepa := []string{iv[len(iv)/2:], iv[:len(iv)/2]}
	var orderStr string = strconv.Itoa(order)
	num := len(encryptedTxt) / 8
	var s = make([]string, 12)
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
	s[10] = dummy[len(dummy)/2:]
	s[11] = dummy[:len(dummy)/2]

	return s
}
