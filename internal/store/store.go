package store

import (
	"sort"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/alexadhy/wfreq/internal/model"
)

// Store is the storage for the entire application
type Store struct {
	c *cache.Cache
}

func New(exp time.Duration) *Store {
	c := cache.New(exp, exp/2)
	return &Store{c}
}

// Store stores key value of type (string, int) to cache
func (s *Store) Store(k string, v int) {
	s.c.Set(k, v, cache.DefaultExpiration)
}

// Load a value from the store
func (s *Store) Load(k string) int {
	if v, found := s.c.Get(k); found {
		return v.(int)
	}
	return 0
}

// LoadAll loads all value from the store, and sort it optionally
func (s *Store) LoadAll(withSort bool, limit int) model.PairList {
	items := s.c.Items()
	var pairs []model.Pair
	for k, v := range items {
		p := model.Pair{Key: k, Value: v.Object.(int)}
		pairs = append(pairs, p)
	}
	pList := model.PairList(pairs)
	if withSort {
		sort.Sort(sort.Reverse(pList))
	}
	if limit > 0 {
		if len(pList) >= limit {
			pList = pList[:limit-1]
		}
	}
	return pList
}

// Clear clears all value from the store
func (s *Store) Clear() {
	s.c.Flush()
}
