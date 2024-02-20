package main

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func Test(t *testing.T) {
	t.Run("", func(t *testing.T) {
		t.Setenv("APIKEY", " ")
		c, err := NewOCRClient(context.Background(), os.Getenv("APIKEY"), "gemini-pro-vision")
		if err != nil {
			t.Fatal(err)
		}
		content, err := c.Do(context.Background(), "./../sample.png")
		if err != nil {
			t.Fatal(err)
		}
		vqr := NewVivantQR(nil)
		data := vqr.FormatRawData(content)
		if len(data) != 11 {
			t.Fail()
		}
		fmt.Println(len(data))
	})
}
