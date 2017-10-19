package jsonapi

import "fmt"

type Cache interface {
	Add(Node) (bool, Node)
}

type coreCache struct {
	store map[string]Node
}

func NewCache() Cache {
	return &coreCache{
		store: make(map[string]Node),
	}
}

func (cache *coreCache) Add(n Node) (bool, Node) {
	key := fmt.Sprintf("%s=>%s", n.GetType(), n.GetID())
	if existing, ok := cache.store[key]; ok {
		return false, existing
	}
	cache.store[key] = n
	return true, n
}
