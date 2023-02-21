package config

import (
	"log"
	"os"

	"github.com/jackc/pgx"

	"github.com/joho/godotenv"
)

type TgBotConfig struct {
	Key       string
	Providers map[string]string
	Connect   pgx.ConnConfig
}

func NewTgBotConfig() *TgBotConfig {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	return &TgBotConfig{
		Key: os.Getenv("KEY"),

		Providers: map[string]string{
			"https://api.apilayer.com/fixer/convert?to=RUB&from=%s&amount=1":              os.Getenv("FIXER"),
			"https://api.apilayer.com/exchangerates_data/convert?to=RUB&from=%s&amount=1": os.Getenv("EXCHANGERATES"),
		},

		Connect: pgx.ConnConfig{
			Host:     "localhost",
			Port:     5434,
			User:     "postgres",
			Password: os.Getenv("DBKEY"),
			Database: "postgres",
		},
	}

}
