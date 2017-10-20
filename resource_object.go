package jsonapi

type DataSetter interface {
	AppendData(ResourceObject)
}

type Includer interface {
	AppendIncluded(ResourceObject)
}

type Document interface {
	DataSetter
	Includer
}

type ResourceObject interface {
	ResourceIdentifier
	SetID(string)
	SetType(string)
	GetAttributes() Attributes
	GetRelationships() RelationshipsObject
}

type ResourceIdentifier interface {
	GetID() string
	GetType() string
}

type coreResourceObject struct {
	ID                  string              `json:"id,omitempty"`
	Type                string              `json:"type"`
	Attributes          Attributes          `json:"attributes,omitempty"`
	RelationshipsObject RelationshipsObject `json:"relationships,omitempty"`
}

func NewResourceObject() ResourceObject {
	return &coreResourceObject{
		Attributes:          newAttributes(),
		RelationshipsObject: newRelationships(),
	}
}

func (n *coreResourceObject) GetID() string {
	return n.ID
}

func (n *coreResourceObject) SetID(id string) {
	n.ID = id
	return
}

func (n *coreResourceObject) GetType() string {
	return n.Type
}

func (n *coreResourceObject) SetType(t string) {
	n.Type = t
	return
}

func (n *coreResourceObject) GetAttributes() Attributes {
	return n.Attributes
}

func (n *coreResourceObject) GetRelationships() RelationshipsObject {
	return n.RelationshipsObject
}

type coreResIdentifier struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func (n *coreResIdentifier) GetID() string {
	return n.ID
}

func (n *coreResIdentifier) GetType() string {
	return n.Type
}
