package pkg

import (
	"sync"
	"time"

	"github.com/MaxKut3/BudgetBot/internal/models"
)

type Currency interface {
	GetValue(cur string) int
}

type Cache interface {
	Get(msg *models.Message) (int, bool)
	Set(msg *models.Message, c Currency)
}

type structCache struct {
	val  int
	time time.Time
}

type simpleCache struct {
	cache map[string]structCache
	mu    sync.RWMutex
}

func NewSimpleCache() *simpleCache {
	return &simpleCache{
		cache: make(map[string]structCache),
	}
}

func (s *simpleCache) Get(msg *models.Message) (int, bool) {
	c, ok := s.cache[msg.Cur]
	if !ok {
		return 0, ok
	}
	return c.val, ok
}

func (s *simpleCache) Set(msg *models.Message, c Currency) {
	s.cache[msg.Cur] = structCache{
		val:  c.GetValue(msg.Cur),
		time: time.Now(),
	}
}
