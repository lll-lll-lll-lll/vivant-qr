package main

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type OCRTxt string

type OCRClient struct {
	c *genai.Client
}

func NewOCRClient(ctx context.Context, apiKey string) (*OCRClient, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &OCRClient{c: client}, nil
}

func (ocrC *OCRClient) Close() error {
	return ocrC.c.Close()
}

func (ocrC *OCRClient) Do(ctx context.Context, imgPath string) (OCRTxt, error) {
	return OCRTxt(""), nil
}
