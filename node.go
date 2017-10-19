package jsonapi

type Attributes struct {
	storage map[string]interface{}
}

func (attrs *Attributes) Append(key string, value interface{}) {
	attrs.storage[key] = value
	return
}

type Node struct {
	ID            string                `json:"id,omitempty"`
	Type          string                `json:"type"`
	Attributes    Attributes            `json:"attributes,omitempty"`
	Relationships RelationshipsResponse `json:"relationships,omitempty"`
}

type ReferenceNode struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type RelationshipTracker interface {
	Append(*Node)
}

type RelationshipsResponse struct {
	Data map[string]RelationshipTracker
}

type SingleRelationship struct {
	Data *ReferenceNode `json:"data"`
}

type SliceRelationship struct {
	Data []*ReferenceNode `json:"data"`
}

type SingleResponse struct {
	Data     *Node   `json:"data,omitempty"`
	Included []*Node `json:"included,omitempty"`
}

type SliceResponse struct {
	Data     []*Node `json:"data,omitempty"`
	Included []*Node `json:"included,omitempty"`
}
