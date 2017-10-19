package jsonapi

type singleResponse struct {
	Data     Node   `json:"data,omitempty"`
	Included []Node `json:"included,omitempty"`
}

// NewSingleResponse returns a JSONAPIResult that can handle a single node
// to represent the main response object
func NewSingleResponse() JSONAPIResult {
	return &singleResponse{
		Included: make([]Node, 0),
	}
}

// AppendIncluded will add a node into the included portion of the response
func (res *singleResponse) AppendIncluded(n Node) {
	res.Included = append(res.Included, n)
	return
}

// AppendData will set the node as the main target of the response
func (res *singleResponse) AppendData(n Node) {
	res.Data = n
	return
}
