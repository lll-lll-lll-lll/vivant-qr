package main

import (
	"context"

	"github.com/otiai10/gosseract/v2"
)

type OCRTxt string

const defaultLanguage = "eng"

type OCRClient struct {
	c    *gosseract.Client
	lang string
}

func NewOCRClient() *OCRClient {
	client := gosseract.NewClient()
	return &OCRClient{c: client}
}

func (ocrC *OCRClient) Close() error {
	return ocrC.c.Close()
}

func (ocrC *OCRClient) Do(ctx context.Context, imgPath string) (OCRTxt, error) {
	if ocrC.lang == "" {
		ocrC.lang = defaultLanguage
	}
	if err := ocrC.c.SetImage(imgPath); err != nil {
		return "", err
	}
	if err := ocrC.c.SetVariable(gosseract.TESSEDIT_CHAR_BLACKLIST, "Â¢â€'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"); err != nil {
		return "", err
	}
	txt, err := ocrC.c.Text()
	if err != nil {
		return "", err
	}
	return OCRTxt(txt), nil
}

func (ocrC *OCRClient) Format(txt string) []string {
	return []string{}
}
