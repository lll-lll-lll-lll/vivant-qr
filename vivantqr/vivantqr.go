package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file",
				Value:    "",
				Usage:    "generated image save file path",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "write",
				Value: false,
				Usage: "generate image file",
			}, &cli.BoolFlag{
				Name:  "read",
				Value: false,
				Usage: "read generated image file",
			},
		},
		Action: func(cCtx *cli.Context) error {
			if cCtx.Bool("read") {
				path := cCtx.String("file")
				cfg, err := NewCfg()
				if err != nil {
					return err
				}
				vivantQR := NewVivantQR(cfg)
				ocrClient := NewOCRClient()
				defer ocrClient.Close()
				content, err := ocrClient.Do(context.TODO(), path)
				if err != nil {
					return err
				}
				formated := vivantQR.DecodeRawData(content)
				r, err := vivantQR.Decode(formated)
				if err != nil {
					return err
				}
				fmt.Printf(r)
				return nil
			}
			if cCtx.Bool("write") {
				path := cCtx.String("file")
				cfg, err := Refresh()
				if err != nil {
					return err
				}
				vivantQR := NewVivantQR(cfg)
				encrypted, err := vivantQR.Encode()
				if err != nil {
					return err
				}
				formated := vivantQR.EncodeRawData(encrypted)
				if err := vivantQR.Output(path, formated); err != nil {
					return err
				}
				fmt.Println("order: ", cfg.Order)
				fmt.Println("secret value: ", cfg.SecretValue)
			}
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
