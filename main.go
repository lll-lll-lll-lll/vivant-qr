package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

const vivant = "https://www.netflix.com/jp/title/81726701"

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "refresh",
				Value: false,
				Usage: "",
			}, &cli.BoolFlag{
				Name:  "write",
				Value: false,
				Usage: "",
			}, &cli.BoolFlag{
				Name:  "read",
				Value: false,
				Usage: "",
			},
			&cli.StringFlag{
				Name:  "file",
				Value: "",
				Usage: "",
			},
		},
		Action: func(cCtx *cli.Context) error {
			if cCtx.Bool("read") {
				if err := godotenv.Load(); err != nil {
					log.Fatal("Error loading .env file")
				}
				order := os.Getenv("ORDER")
				if order == "" {
					log.Fatal("order is nil. refresh order")
				}
				k := os.Getenv("SECRET_KEY")
				if k == "" {
					log.Fatal("secret key is nil")
				}
				key, err := hex.DecodeString(k)
				if err != nil {
					log.Fatal(err)
				}
				path := cCtx.Args().Get(2)
				if path == "" {
					log.Fatal("file not set. --file {file path}")
				}
				f, err := os.ReadFile(path)
				if err != nil {
					log.Fatal(err)
				}
				content := string(f)
				splittedCttSlice := strings.Split(content, " ")
				var correctNums = make([]string, 0, 11)
				for _, o := range order {
					idx, _ := strconv.Atoi(string(o))
					correctNums = append(correctNums, splittedCttSlice[idx])
				}
				iv := correctNums[len(correctNums)-1] + correctNums[len(correctNums)-2]
				ivByte, err := base64.StdEncoding.DecodeString(iv)
				if err != nil {
					log.Fatal(err)
				}
				notDummyEncrypted := strings.Join(correctNums[:8], "")
				nde, err := base64.StdEncoding.DecodeString(notDummyEncrypted)
				if err != nil {
					log.Fatal(err)
				}
				decrypted, _ := decrypt(nde, key, ivByte)
				fmt.Printf("Decrypted: %s\n", decrypted)
			}
			if cCtx.Bool("write") {
				if err := godotenv.Load(); err != nil {
					log.Fatal("Error loading .env file")
				}
				order := os.Getenv("ORDER")
				if order == "" {
					log.Fatal("order is nil. refresh order")
				}
				k := os.Getenv("SECRET_KEY")
				if k == "" {
					log.Fatal("secret key is nil")
				}
				key, _ := hex.DecodeString(k)
				iv, encrypted, _ := Encrypt([]byte(vivant), key)
				encodedIV := base64.StdEncoding.EncodeToString(iv)
				fmt.Println("encdoeIV", encodedIV)
				encodedEncrypted := base64.StdEncoding.EncodeToString(encrypted)
				fmt.Println("encode", encodedEncrypted)
				dummy := genDummy(12)
				o, err := strconv.Atoi(order)
				if err != nil {
					log.Fatal(err)
				}
				result := separate(encodedIV, encodedEncrypted, dummy, o)
				write(result)
			}
			if cCtx.Bool("refresh") {
				if err := godotenv.Load(); err != nil {
					log.Fatal("Error loading .env file")
				}
				f, err := os.Create(".env")
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()
				if _, err := f.Write([]byte(fmt.Sprintf("ORDER=%d", genOrder(10)))); err != nil {
					log.Fatal(err)
				}
			}
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
