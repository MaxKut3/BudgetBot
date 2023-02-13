package pkg

import (
	"time"

	"github.com/MaxKut3/BudgetBot/internal/Message"
)

type Currency interface {
	GetValue(cur string) int
}

type Cache interface {
	Get(msg *Message.Message) (int, bool)
	Set(msg *Message.Message, c Currency)
}

type structCache struct {
	val  int
	time time.Time
}

type simpleCache struct {
	cache map[string]structCache
}

func NewSimpleCache() *simpleCache {
	return &simpleCache{
		cache: make(map[string]structCache),
	}
}

func (s *simpleCache) Get(msg *Message.Message) (int, bool) {
	c, ok := s.cache[msg.Cur]
	if !ok {
		return 0, ok
	}
	return c.val, ok
}

func (s *simpleCache) Set(msg *Message.Message, c Currency) {
	s.cache[msg.Cur] = structCache{
		val:  c.GetValue(msg.Cur),
		time: time.Now(),
	}
}
