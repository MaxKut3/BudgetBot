package useCases

import (
	"sync"

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
	max := 0
	ch := make(chan int)

	var wg sync.WaitGroup
	wg.Add(len(c.Cfg.Providers))

	for key, provider := range c.Cfg.Providers {

		go func() {
			defer wg.Done()
			ch <- provider(cur, key)
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
