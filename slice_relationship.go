package jsonapi

type sliceRelationship struct {
	Data []ResourceIdentifier `json:"data"`
}

func newSliceRelationship() Relationship {
	return &sliceRelationship{}
}

func (rel *sliceRelationship) Append(n ResourceIdentifier) {
	rel.Data = append(
		rel.Data,
		&coreResIdentifier{
			ID:   n.GetID(),
			Type: n.GetType(),
		},
	)
	return
}
