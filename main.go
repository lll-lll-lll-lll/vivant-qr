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
					return err
				}
				vivantQR := &VivantQR{cfg: cfg}
				ocrClient := NewOCRClient()
				defer ocrClient.c.Close()
				content, err := ocrClient.Do(context.TODO(), "./save.png")
				if err != nil {
					return err
				}
				formated := vivantQR.FormatDecode(content)
				r, err := vivantQR.Decode(formated)
				if err != nil {
					return err
				}
				fmt.Printf(r)
				return nil
			}
			if cCtx.Bool("write") {
				cfg, err := Refresh()
				if err != nil {
					return err
				}
				vivantQR := &VivantQR{cfg: cfg}
				encrypted, err := vivantQR.Encode()
				if err != nil {
					return err
				}
				formated := vivantQR.FormatEncode(encrypted)
				if err := vivantQR.Output("./images/background.png", "./save.png", formated); err != nil {
					return err
				}
			}
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
