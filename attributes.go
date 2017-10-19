package jsonapi

import "encoding/json"

// Attributes represents a key value collection that can be appended to.
type Attributes interface {
	Append(string, interface{})
}

type attribs struct {
	storage map[string]interface{}
}

func (attrs *attribs) MarshalJSON() ([]byte, error) {
	return json.Marshal(attrs.storage)
}

// NewAttributes will return an Attributes structure that can be used to
// represent an objects attributes.
func NewAttributes() Attributes {
	return &attribs{
		storage: make(map[string]interface{}),
	}
}

// Append will add a given key and its value into the key value store.
func (attrs *attribs) Append(key string, value interface{}) {
	attrs.storage[key] = value
	return
}
