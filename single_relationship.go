package jsonapi

type singleRelationship struct {
	Data BaseNode `json:"data"`
}

// NewSingleRelationship will return a RelationshipManager capable of
// representing a single Node relationship
func NewSingleRelationship() RelationshipManager {
	return &singleRelationship{}
}

// Append will set the Node as the relationships target
func (rel *singleRelationship) Append(n BaseNode) {
	rel.Data = &referenceNode{
		ID:   n.GetID(),
		Type: n.GetType(),
	}
	return
}
