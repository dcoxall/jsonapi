JSON:API
========

[![Build Status](https://travis-ci.org/dcoxall/jsonapi.svg?branch=dev)](https://travis-ci.org/dcoxall/jsonapi)

A Go JSON:API serialization and deserialization library.

Overview
--------

This JSON API library was built to support circular references when serializing
and deserializing Go structures. Due to the dual phase approach to avoid stack
overflow errors then it is unlikely to be the most performant.

Usage
-----

```go
import (
    "bytes"
    "fmt"
    "github.com/dcoxall/jsonapi"
)

type Author struct {
    ID        string `jsonapi:"primary,author"`
    FirstName string `jsonapi:"attr,first_name"`
    LastName  string `jsonapi:"attr,last_name"`
}

type Book struct {
    ID     string  `jsonapi:"primary,book"`
    Name   string  `jsonapi:"attr,name"`
    Author *Author `jsonapi:"relation,author"`
}

func ExampleEncode() {
    book := &Book{
        ID:     "book-123",
        Name:   "JSON Book",
        Author: &Author{
            ID:        "author-123",
            FirstName: "Jane",
            LastName:  "Doe",
        },
    }
    buffer := new(bytes.Buffer)
    encoder := jsonapi.NewEncoder(buffer)
    if err := encoder.Encode(book); err != nil {
        fmt.Println(err.Error())
        return
    }
    fmt.Println(buffer.String())
}
```

The above would produce...

```json
{
  "data": {
    "id": "book-123",
    "type": "book",
    "attributes": {
      "name": "JSON Book"
    },
    "relationships": {
      "author": {
        "data": {
          "id": "author-123",
          "type": "author"
        }
      }
    }
  },
  "included": [
    {
      "id": "author-123",
      "type": "author",
      "attributes": {
        "first_name": "Jane",
        "last_name": "Doe"
      }
    }
  ]
}
```
