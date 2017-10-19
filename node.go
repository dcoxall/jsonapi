package jsonapi

type Node struct {
	ID            string        `json:"id,omitempty"`
	Type          string        `json:"type"`
	Attributes    Attributes    `json:"attributes,omitempty"`
	Relationships Relationships `json:"relationships,omitempty"`
}

type ReferenceNode struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type DataSetter interface {
	AppendData(*Node)
}

type Includer interface {
	AppendIncluded(*Node)
}

type JSONAPIResult interface {
	DataSetter
	Includer
}
