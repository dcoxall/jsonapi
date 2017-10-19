package jsonapi

type singleRelationship struct {
	Data *ReferenceNode `json:"data"`
}

// NewSingleRelationship will return a RelationshipManager capable of
// representing a single Node relationship
func NewSingleRelationship() RelationshipManager {
	return &singleRelationship{}
}

// Append will set the Node as the relationships target
func (rel *singleRelationship) Append(node *Node) {
	rel.Data = &ReferenceNode{
		ID:   node.ID,
		Type: node.Type,
	}
	return
}
