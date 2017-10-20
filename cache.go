package jsonapi

import "fmt"

type Cache interface {
	Add(ResourceObject) (bool, ResourceObject)
}

type coreCache struct {
	store map[string]ResourceObject
}

func NewCache() Cache {
	return &coreCache{
		store: make(map[string]ResourceObject),
	}
}

func (cache *coreCache) Add(n ResourceObject) (bool, ResourceObject) {
	key := fmt.Sprintf("%s=>%s", n.GetType(), n.GetID())
	if existing, ok := cache.store[key]; ok {
		return false, existing
	}
	cache.store[key] = n
	return true, n
}
