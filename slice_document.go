package jsonapi

type sliceDocument struct {
	Data     []ResourceObject `json:"data,omitempty"`
	Included []ResourceObject `json:"included,omitempty"`
}

func (res *sliceDocument) AppendData(n ResourceObject) {
	res.Data = append(res.Data, n)
	return
}

func (res *sliceDocument) AppendIncluded(n ResourceObject) {
	res.Included = append(res.Included, n)
	return
}

func newSliceDocument() Document {
	return &sliceDocument{
		Data:     make([]ResourceObject, 0),
		Included: make([]ResourceObject, 0),
	}
}
