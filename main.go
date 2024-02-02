package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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
				cfg, err := NewCfg()
				if err != nil {
					log.Fatal(err)
				}
				key, err := hex.DecodeString(cfg.SecretKey)
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
				for _, o := range strconv.Itoa(cfg.Order) {
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
				cfg, err := NewCfg()
				if err != nil {
					log.Fatal(err)
				}
				vivantQR := &VivantQR{cfg: cfg}
				result, err := vivantQR.Generate()
				if err != nil {
					log.Fatal(err)
				}
				write(result)
			}
			if cCtx.Bool("refresh") {
				f, err := os.Create(".env")
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()
				if _, err := f.Write([]byte(fmt.Sprintf("ORDER=%d\nSECRET_KEY=%s", genOrder(10), genSecret(32)))); err != nil {
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
