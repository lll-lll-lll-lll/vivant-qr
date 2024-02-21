package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const prompt = "load numbers.Remove line breaks and replace them with spaces"

type OCRClient struct {
	c *genai.Client
	m *genai.GenerativeModel
}

func NewOCRClient(ctx context.Context, apiKey, model string) (*OCRClient, error) {
	c, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	m := c.GenerativeModel(model)
	return &OCRClient{c: c, m: m}, nil
}

func (ocrC *OCRClient) Close() error {
	return ocrC.c.Close()
}

func (ocrC *OCRClient) Do(ctx context.Context, imgPath string) (string, error) {
	var out string
	img, err := os.ReadFile(imgPath)
	if err != nil {
		return "", err
	}
	res, err := ocrC.m.GenerateContent(ctx, genai.Text(prompt), genai.ImageData("png", img))
	if err != nil {
		return "", err
	}
	for _, cand := range res.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				out += fmt.Sprint(part)
			}
		}
	}
	return out, nil
}
