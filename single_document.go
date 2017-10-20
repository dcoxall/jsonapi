package jsonapi

type singleDocument struct {
	Data     ResourceObject   `json:"data,omitempty"`
	Included []ResourceObject `json:"included,omitempty"`
}

func newSingleDocument() Document {
	return &singleDocument{
		Included: make([]ResourceObject, 0),
	}
}

func (res *singleDocument) AppendIncluded(n ResourceObject) {
	res.Included = append(res.Included, n)
	return
}

func (res *singleDocument) AppendData(n ResourceObject) {
	res.Data = n
	return
}
