package jsonapi

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

type Author struct {
	ID       string `jsonapi:"primary,author"`
	FullName string `jsonapi:"attr,name"`
	Book     *Book  `jsonapi:"relation,book"`
}

type Book struct {
	ID     string  `jsonapi:"primary,book"`
	Name   string  `jsonapi:"attr,name"`
	Author *Author `jsonapi:"relation,author"`
}

func TestBasicSerialization(t *testing.T) {
	book := &Book{
		ID:   "book-123",
		Name: "Awesome Book",
		Author: &Author{
			ID:       "author-123",
			FullName: "Jane Doe",
		},
	}
	book.Author.Book = book // make it circular
	buffer := &bytes.Buffer{}
	encoder := NewEncoder(buffer)

	err := encoder.Encode(book)
	assert.NoError(t, err)

	var expected []byte
	expected, err = ioutil.ReadFile("testdata/basic_serialization.json")
	assert.NoError(t, err)
	assert.JSONEq(t, string(expected), buffer.String())
}
