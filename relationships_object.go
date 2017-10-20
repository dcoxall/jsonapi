package jsonapi

import "encoding/json"

// RelationshipsObject represents a collection of multiple Relationships. Relationships should
// be used to represent all the numerous relationships for an object.
type RelationshipsObject interface {
	Append(string, Relationship)
}

// Relationship represents something that can append resource identifiers to represent
// a particular relationship. The Relationship only ever represents a single relationship.
type Relationship interface {
	Append(ResourceIdentifier)
}

func newRelationships() RelationshipsObject {
	return &relations{
		storage: make(map[string]Relationship),
	}
}

type relations struct {
	storage map[string]Relationship
}

func (rels *relations) MarshalJSON() ([]byte, error) {
	return json.Marshal(rels.storage)
}

func (rel *relations) Append(name string, relManager Relationship) {
	rel.storage[name] = relManager
	return
}
