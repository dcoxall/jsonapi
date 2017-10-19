package jsonapi

type sliceRelationship struct {
	Data []BaseNode `json:"data"`
}

// NewSliceRelationship will return a RelationshipManager capable of
// representing many nodes in a single relationship
func NewSliceRelationship() RelationshipManager {
	return &sliceRelationship{}
}

// Append will add a Node into the relationship
func (rel *sliceRelationship) Append(n BaseNode) {
	rel.Data = append(
		rel.Data,
		&referenceNode{
			ID:   n.GetID(),
			Type: n.GetType(),
		},
	)
	return
}
