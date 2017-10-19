package jsonapi

import "encoding/json"

// Relationships represents a collection of multiple RelationshipManagers.
// With each manager representing a single relationship. Relationships should
// be used to represent all the numerous relationships for an object.
type Relationships interface {
	Append(string, RelationshipManager)
}

// RelationshipManager represents something that can append nodes to represent
// a particular relationship. The RelationshipManager only ever represents a
// single relationship.
type RelationshipManager interface {
	Append(BaseNode)
}

// NewRelationships will return a Relationships structure capable of
// representing all an objects relationships.
func NewRelationships() Relationships {
	return &relations{
		storage: make(map[string]RelationshipManager),
	}
}

type relations struct {
	storage map[string]RelationshipManager `json:"data"`
}

func (rels *relations) MarshalJSON() ([]byte, error) {
	return json.Marshal(rels.storage)
}

func (rel *relations) Append(name string, relManager RelationshipManager) {
	rel.storage[name] = relManager
	return
}
