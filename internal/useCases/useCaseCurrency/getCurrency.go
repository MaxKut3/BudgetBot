package useCaseCurrency

import (
	"sync"

	"github.com/MaxKut3/BudgetBot/internal/useCases/provider"

	"github.com/MaxKut3/BudgetBot/config"
)

type Currency interface {
	GetValue(cur string) int
}

type currencyStr struct {
	Cfg *config.TgBotConfig
}

func NewCurrencyStr(cfg *config.TgBotConfig) *currencyStr {
	return &currencyStr{
		Cfg: cfg,
	}
}

func (c *currencyStr) GetValue(cur string) int {
	ch := make(chan int)

	var wg sync.WaitGroup
	wg.Add(len(c.Cfg.Providers))

	max := 0

	for url, key := range c.Cfg.Providers {

		prov := provider.NewProvider(url, cur, key)

		go func() {
			defer wg.Done()

			ch <- provider.Sender(prov)

		}()
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	for v := range ch {

		if v > max {
			max = v
		}
	}

	return max
}
