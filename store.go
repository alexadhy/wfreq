package main

import (
	"github.com/patrickmn/go-cache"
	"sort"
	"time"
)

type store struct {
	c *cache.Cache
}

type pair struct {
	Key   string
	Value int
}

type pairList []pair

func (p pairList) Len() int           { return len(p) }
func (p pairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func newStore(exp time.Duration) *store {
	c := cache.New(exp, exp/2)
	return &store{c}
}

func (s *store) store(k string, v int) {
	s.c.Set(k, v, cache.DefaultExpiration)
}

func (s *store) load(k string) int {
	if v, found := s.c.Get(k); found {
		return v.(int)
	}
	return 0
}

func (s *store) loadAll(withSort bool) pairList {
	items := s.c.Items()
	var pairs []pair
	for k, v := range items {
		p := pair{Key: k, Value: v.Object.(int)}
		pairs = append(pairs, p)
	}
	pList := pairList(pairs)
	if withSort {
		sort.Sort(sort.Reverse(pList))
	}
	return pList
}
