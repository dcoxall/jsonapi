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

func newAttributes() Attributes {
	return &attribs{
		storage: make(map[string]interface{}),
	}
}

func (attrs *attribs) Append(key string, value interface{}) {
	attrs.storage[key] = value
	return
}
