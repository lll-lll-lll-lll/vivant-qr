package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"vivant-qr/vivantqr"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file",
				Value:    "",
				Usage:    "",
				Required: true,
			},
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
				path := cCtx.String("file")
				cfg, err := vivantqr.NewCfg()
				if err != nil {
					return err
				}
				vivantQR := vivantqr.NewVivantQR(cfg)
				ocrClient := vivantqr.NewOCRClient()
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
				cfg, err := vivantqr.Refresh()
				if err != nil {
					return err
				}
				vivantQR := vivantqr.NewVivantQR(cfg)
				encrypted, err := vivantQR.Encode()
				if err != nil {
					return err
				}
				formated := vivantQR.EncodeRawData(encrypted)
				if err := vivantQR.Output(path, formated); err != nil {
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
