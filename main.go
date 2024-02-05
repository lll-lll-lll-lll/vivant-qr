package main

import (
	"context"
	"fmt"
	"log"
	"os"
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
				vivantQR := &VivantQR{cfg: cfg}
				ocrClient := NewOCRClient()
				defer ocrClient.c.Close()
				content, err := ocrClient.Do(context.TODO(), "./save.png")
				if err != nil {
					log.Fatal(err)
				}
				splitedCnt := strings.Split(string(content), " ")
				var decoded []string
				for _, v := range splitedCnt {
					v = strings.ReplaceAll(v, "\n", " ")
					spliteV := strings.Split(v, " ")
					if len(spliteV) == 2 {
						for i := 0; i < len(spliteV); i++ {
							decoded = append(decoded, spliteV[i])
						}
						continue
					}
					decoded = append(decoded, v)
				}
				r, err := vivantQR.Decrypt(decoded)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf(string(r))
			}
			if cCtx.Bool("write") {
				cfg, err := NewCfg()
				if err != nil {
					log.Fatal(err)
				}
				vivantQR := &VivantQR{cfg: cfg}
				result, err := vivantQR.Encrpto()
				if err != nil {
					log.Fatal(err)
				}

				var res []string
				for i := 2; i <= 12; i += 2 {
					a := strings.Join(result[i-2:i], " ")
					res = append(res, a)
				}
				if err := vivantQR.Output("./images/background.png", "./save.png", res); err != nil {
					log.Fatal(err)
				}
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

func bytesToOctalString(bytes []byte) string {
	var octalString string

	for _, b := range bytes {
		// 8進数で3桁になるようにフォーマット
		octalString += fmt.Sprintf("%03o", b)
	}

	return octalString
}

func octalStringToBytes(octalString string) ([]byte, error) {
	var bytes []byte

	for i := 0; i < len(octalString); i += 3 {
		// 残りの文字が3文字未満の場合の対処
		endIndex := i + 3
		if endIndex > len(octalString) {
			endIndex = len(octalString)
		}

		octalByte := octalString[i:endIndex]
		var b int
		_, err := fmt.Sscanf(octalByte, "%03o", &b)
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, byte(b))
	}

	return bytes, nil
}
