package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	postgres "github.com/MaxKut3/BudgetBot/internal/useCases/repositories"

	"github.com/MaxKut3/BudgetBot/config"
	"github.com/MaxKut3/BudgetBot/internal/controller"
	"github.com/MaxKut3/BudgetBot/internal/useCases/useCaseCurrency"
	"github.com/MaxKut3/BudgetBot/pkg"
)

func Run(cfg *config.TgBotConfig) {

	cache := pkg.NewSimpleCache()
	currency := useCaseCurrency.NewCurrencyStr(cfg)
	rep := postgres.NewRepository(cfg)
	client := controller.NewTgBot(cfg, cache, currency, rep) // Ошибка с интерфейсом кэша

	go client.Run()

	//Graceful Shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signalChan

	log.Printf("%s signal caught", sig)

	client.Stop()
}
