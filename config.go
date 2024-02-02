package main

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	SecretKey string
	Order     int
}

func NewCfg() (*Config, error) {
	var err error = errors.New("failed to initialize config")
	if e := godotenv.Load(); e != nil {
		return nil, err
	}
	order := os.Getenv("ORDER")
	if order == "" {
		return nil, err
	}
	key := os.Getenv("SECRET_KEY")
	if key == "" {
		log.Fatal("secret key is nil")
	}
	o, err := strconv.Atoi(order)
	if err != nil {
		return nil, err
	}
	return &Config{SecretKey: key, Order: o}, nil
}
