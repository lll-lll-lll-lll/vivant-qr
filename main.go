package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const vivant = "https://www.netflix.com/jp/title/81726701"

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "write",
				Value: false,
				Usage: "",
			}, &cli.BoolFlag{
				Name:  "read",
				Value: false,
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
				formated := vivantQR.FormatDecode(content)
				r, err := vivantQR.Decode(formated)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf(r)
			}
			if cCtx.Bool("write") {
				cfg, err := Refresh()
				if err != nil {
					log.Fatal(err)
				}
				vivantQR := &VivantQR{cfg: cfg}
				encrypted, err := vivantQR.Encode()
				if err != nil {
					log.Fatal(err)
				}
				formated := vivantQR.FormatEncode(encrypted)
				if err := vivantQR.Output("./images/background.png", "./save.png", formated); err != nil {
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
