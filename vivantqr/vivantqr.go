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
				Name:     "apikey",
				Value:    "",
				Usage:    "gemini api key",
				Required: true,
			},
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
			}, &cli.StringFlag{
				Name:  "secret",
				Value: "",
				Usage: "set secret value",
			},
		},
		Action: func(cCtx *cli.Context) error {
			if cCtx.Bool("read") {
				apiKey := cCtx.String("apikey")
				path := cCtx.String("file")
				cfg, err := NewCfg()
				if err != nil {
					return err
				}
				vivantQR := NewVivantQR(cfg)
				ctx := context.Background()
				ocrClient, err := NewOCRClient(ctx, apiKey, "gemini-pro-vision")
				if err != nil {
					return err
				}
				defer ocrClient.Close()
				content, err := ocrClient.Do(ctx, path)
				if err != nil {
					return err
				}
				formated := vivantQR.FormatRawData(content)
				r, err := vivantQR.Decode(formated)
				if err != nil {
					return err
				}
				fmt.Printf(r)
				return nil
			}
			if cCtx.Bool("write") {
				path := cCtx.String("file")
				secret := cCtx.String("secret")
				cfg, err := Refresh(secret)
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
				fmt.Println("secret key: ", cfg.SecretKey)
			}
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
