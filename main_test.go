package main

import (
	"fmt"
	"testing"
)

func Test_Decode(t *testing.T) {
	// encodedIV := base64.StdEncoding.EncodeToString([]byte("tests"))
	// fmt.Println(encodedIV)
	encodedString := bytesToOctalString([]byte("encodedIV"))
	fmt.Println("エンコードされた文字列:", encodedString)

	decodedBytes, err := octalStringToBytes(encodedString)
	if err != nil {
		fmt.Println("デコードエラー:", err)
		return
	}
	t.Log(decodedBytes)
}
