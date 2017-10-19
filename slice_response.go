package jsonapi

type sliceResponse struct {
	Data     []*Node `json:"data,omitempty"`
	Included []*Node `json:"included,omitempty"`
}

// AppendData adds a Node into the main portion of the response
func (res *sliceResponse) AppendData(node *Node) {
	res.Data = append(res.Data, node)
	return
}

// AppendIncluded adds a Node into the included portion of the response
func (res *sliceResponse) AppendIncluded(node *Node) {
	res.Included = append(res.Included, node)
	return
}

// NewSliceResponse returns a JSONAPIResult that will represent the main
// portion of the response as a slice/array
func NewSliceResponse() JSONAPIResult {
	return &sliceResponse{
		Data:     make([]*Node, 0),
		Included: make([]*Node, 0),
	}
}
