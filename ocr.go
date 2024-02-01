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

func (ocrC *OCRClient) Do(ctx context.Context, imgPath string) (OCRTxt, error) {
	if ocrC.lang == "" {
		ocrC.lang = defaultLanguage
	}
	if err := ocrC.c.SetImage(imgPath); err != nil {
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
