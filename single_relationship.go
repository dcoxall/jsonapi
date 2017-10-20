package jsonapi

type singleRelationship struct {
	Data ResourceIdentifier `json:"data"`
}

func newSingleRelationship() Relationship {
	return &singleRelationship{}
}

func (rel *singleRelationship) Append(n ResourceIdentifier) {
	rel.Data = &coreResIdentifier{
		ID:   n.GetID(),
		Type: n.GetType(),
	}
	return
}
