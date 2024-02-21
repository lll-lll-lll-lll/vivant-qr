package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	SecretKey   string
	SecretValue string
	Order       int
}

func Refresh(secretVal string) (*Config, error) {
	f, err := os.Create("env")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if _, err := f.Write([]byte(fmt.Sprintf("ORDER=%d\nSECRET_KEY=%d\nSECRET_VALUE=%s", genData(10), genData(32), secretVal))); err != nil {
		return nil, err
	}
	return NewCfg()
}

func NewCfg() (*Config, error) {
	var err error = errors.New("failed to initialize config")
	if e := godotenv.Load("env"); e != nil {
		return nil, err
	}
	order := os.Getenv("ORDER")
	if order == "" {
		return nil, err
	}
	key := os.Getenv("SECRET_KEY")
	if key == "" {
		return nil, err
	}
	val := os.Getenv("SECRET_VALUE")
	if val == "" {
		return nil, err
	}
	o, err := strconv.Atoi(order)
	if err != nil {
		return nil, err
	}
	return &Config{SecretKey: key, Order: o, SecretValue: val}, nil
}
