package jsonapi

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"time"
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

type MultiAuthorAuthor struct {
	ID       string             `jsonapi:"primary,author"`
	FullName string             `jsonapi:"attr,name"`
	Books    []*MultiAuthorBook `jsonapi:"relation,books"`
}

type MultiAuthorBook struct {
	ID      string               `jsonapi:"primary,book"`
	Name    string               `jsonapi:"attr,name"`
	Authors []*MultiAuthorAuthor `jsonapi:"relation,authors"`
}

type VariableIDType struct {
	ID interface{} `jsonapi:"primary,vari-id"`
}

type PlainStruct struct {
	Time  time.Time `json:"time"`
	Place string    `json:"place"`
}

type CompositeStructA struct {
	CompositeStructB
	ID           string `jsonapi:"primary,comp-a"`
	StringField  string `jsonapi:"attr,string"`
	IntField     int    `jsonapi:"attr,integer"`
	IgnoredField string `jsonapi:"-"`
}

type CompositeStructB struct {
	BooleanField bool        `json:"boolean"`
	StructField  PlainStruct `json:"object"`
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

	expected, err := ioutil.ReadFile("testdata/basic_serialization.json")
	assert.NoError(t, err)
	assert.JSONEq(t, string(expected), buffer.String())
}

func TestSliceSerialization(t *testing.T) {
	books := []*Book{
		{
			ID:   "book-123",
			Name: "Awesome Book 1",
		},
		{
			ID:   "book-234",
			Name: "Awesome Book 2",
		},
	}
	books[0].Author = &Author{
		ID:       "author-123",
		FullName: "Jane Doe",
		Book:     books[0], // circular
	}
	books[1].Author = &Author{
		ID:       "author-234",
		FullName: "John Doe",
		Book:     books[1], // circular
	}
	buffer := &bytes.Buffer{}
	encoder := NewEncoder(buffer)

	err := encoder.Encode(books)
	assert.NoError(t, err)

	expected, err := ioutil.ReadFile("testdata/slice_serialization.json")
	assert.NoError(t, err)
	assert.JSONEq(t, string(expected), buffer.String())
}

func TestSliceRelationshipSerialization(t *testing.T) {
	book := &MultiAuthorBook{
		ID:   "book-123",
		Name: "Awesome Book",
		Authors: []*MultiAuthorAuthor{
			{
				ID:       "author-123",
				FullName: "Jane Doe",
				Books:    make([]*MultiAuthorBook, 0),
			},
			{
				ID:       "author-234",
				FullName: "John Doe",
				Books:    make([]*MultiAuthorBook, 0),
			},
		},
	}
	book.Authors[0].Books = append(book.Authors[0].Books, book) // circular
	book.Authors[1].Books = append(book.Authors[1].Books, book) // circular
	buffer := &bytes.Buffer{}
	encoder := NewEncoder(buffer)

	err := encoder.Encode(book)
	assert.NoError(t, err)

	expected, err := ioutil.ReadFile(
		"testdata/slice_relationship_serialization.json",
	)
	assert.NoError(t, err)
	assert.JSONEq(t, string(expected), buffer.String())
}

func TestIDSerialization(t *testing.T) {
	tests := []struct {
		typeName    string
		idValue     interface{}
		expectation string
	}{
		{"float32", float32(1.23), `{"data":{"type":"vari-id","id":"1.23"}}`},
		{"float64", float64(1.23), `{"data":{"type":"vari-id","id":"1.23"}}`},
		{"int", int(123), `{"data":{"type":"vari-id","id":"123"}}`},
		{"int8", int8(123), `{"data":{"type":"vari-id","id":"123"}}`},
		{"int16", int16(123), `{"data":{"type":"vari-id","id":"123"}}`},
		{"int32", int32(123), `{"data":{"type":"vari-id","id":"123"}}`},
		{"int64", int64(123), `{"data":{"type":"vari-id","id":"123"}}`},
		{"uint", uint(123), `{"data":{"type":"vari-id","id":"123"}}`},
		{"uint8", uint8(123), `{"data":{"type":"vari-id","id":"123"}}`},
		{"uint16", uint16(123), `{"data":{"type":"vari-id","id":"123"}}`},
		{"uint32", uint32(123), `{"data":{"type":"vari-id","id":"123"}}`},
		{"uint64", uint64(123), `{"data":{"type":"vari-id","id":"123"}}`},
		{"complex32", complex(float32(1.23), float32(0)), `{"data":{"type":"vari-id","id":"(1.23+0i)"}}`},
		{"complex64", complex(float64(1.23), float64(0)), `{"data":{"type":"vari-id","id":"(1.23+0i)"}}`},
		{"string", "foo", `{"data":{"type":"vari-id","id":"foo"}}`},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Serializing IDs with a type of %s", test.typeName), func(t *testing.T) {
			obj := VariableIDType{ID: test.idValue}
			buffer := &bytes.Buffer{}
			encoder := NewEncoder(buffer)
			assert.NoError(t, encoder.Encode(obj))
			assert.JSONEq(t, test.expectation, buffer.String())
		})
	}
}

func TestCompositeStructSerialization(t *testing.T) {
	obj := &CompositeStructA{
		ID:           "composite-123",
		StringField:  "this is a string",
		IntField:     -123,
		IgnoredField: "YOU CAN'T SEE ME",
		CompositeStructB: CompositeStructB{
			BooleanField: true,
			StructField: PlainStruct{
				Time:  time.Date(2017, time.October, 20, 14, 43, 0, 0, time.FixedZone("BST", 3600)),
				Place: "Hertfordshire",
			},
		},
	}
	buffer := &bytes.Buffer{}
	encoder := NewEncoder(buffer)

	err := encoder.Encode(obj)
	assert.NoError(t, err)

	expected, err := ioutil.ReadFile(
		"testdata/composite_serialization.json",
	)
	assert.NoError(t, err)
	assert.JSONEq(t, string(expected), buffer.String())
}
