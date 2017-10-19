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
