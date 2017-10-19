package jsonapi

type Node interface {
	BaseNode
	SetID(string)
	SetType(string)
	GetAttributes() Attributes
	GetRelationships() Relationships
}

type BaseNode interface {
	GetID() string
	GetType() string
}

type coreNode struct {
	ID            string        `json:"id,omitempty"`
	Type          string        `json:"type"`
	Attributes    Attributes    `json:"attributes,omitempty"`
	Relationships Relationships `json:"relationships,omitempty"`
}

func (n *coreNode) GetID() string {
	return n.ID
}

func (n *coreNode) SetID(id string) {
	n.ID = id
	return
}

func (n *coreNode) GetType() string {
	return n.Type
}

func (n *coreNode) SetType(t string) {
	n.Type = t
	return
}

func (n *coreNode) GetAttributes() Attributes {
	return n.Attributes
}

func (n *coreNode) GetRelationships() Relationships {
	return n.Relationships
}

type referenceNode struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func (n *referenceNode) GetID() string {
	return n.ID
}

func (n *referenceNode) GetType() string {
	return n.Type
}

type DataSetter interface {
	AppendData(Node)
}

type Includer interface {
	AppendIncluded(Node)
}

type JSONAPIResult interface {
	DataSetter
	Includer
}
