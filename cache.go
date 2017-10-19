package jsonapi

import "fmt"

type Cache interface {
	Add(*Node) (bool, *Node)
}

type coreCache struct {
	store map[string]*Node
}

func NewCache() Cache {
	return &coreCache{
		store: make(map[string]*Node),
	}
}

func (cache *coreCache) Add(node *Node) (bool, *Node) {
	fmt.Printf("Attempting add for %#v\n", node)
	key := fmt.Sprintf("%s=>%s", node.Type, node.ID)
	if existing, ok := cache.store[key]; ok {
		fmt.Printf("Already present %#v", existing)
		return false, existing
	}
	cache.store[key] = node
	return true, node
}
