package main

import (
	"github.com/MaxKut3/BudgetBot/config"
	"github.com/MaxKut3/BudgetBot/internal/app"
)

func main() {

	cfg := config.NewTgBotConfig()
	app.Run(cfg)

}
