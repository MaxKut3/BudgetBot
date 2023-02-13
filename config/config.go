package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type GetCurRate func(cur, key string) int

type TgBotConfig struct {
	Key       string
	Providers map[string]GetCurRate
}

func NewTgBotConfig() *TgBotConfig {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	return &TgBotConfig{
		Key: os.Getenv("KEY"),
		Providers: map[string]GetCurRate{
			"coinGAteAPI":              coinGateAPI,
			os.Getenv("FIXER"):         fixerAPI,
			os.Getenv("EXCHANGERATES"): exchangeratesAPI,
		},
	}
}
