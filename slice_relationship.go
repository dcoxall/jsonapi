package jsonapi

type sliceRelationship struct {
	Data []*ReferenceNode `json:"data"`
}

// NewSliceRelationship will return a RelationshipManager capable of
// representing many nodes in a single relationship
func NewSliceRelationship() RelationshipManager {
	return &sliceRelationship{}
}

// Append will add a Node into the relationship
func (rel *sliceRelationship) Append(node *Node) {
	rel.Data = append(
		rel.Data,
		&ReferenceNode{
			ID:   node.ID,
			Type: node.Type,
		},
	)
	return
}
