package jsonapi_test

import (
	"bytes"
	"fmt"
	"github.com/dcoxall/jsonapi"
)

type Author struct {
	ID        uint   `jsonapi:"primary,author"`
	FirstName string `jsonapi:"attr,first_name"`
	LastName  string `jsonapi:"attr,last_name"`
	Age       uint   `jsonapi:"-"`
}

type Book struct {
	ID     uint    `jsonapi:"primary,book"`
	Name   string  `jsonapi:"attr,name"`
	Author *Author `jsonapi:"relation,author"`
}

func Example() {
	book := &Book{
		ID:   1234,
		Name: "JSON:API in Go",
		Author: &Author{
			ID:        9876,
			FirstName: "Jane",
			LastName:  "Doe",
			Age:       42,
		},
	}

	buffer := &bytes.Buffer{}
	encoder := jsonapi.NewIndentEncoder(buffer, "", "    ")
	if err := encoder.Encode(book); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(buffer.String())
	}
	// Output:
	// {
	//     "data": {
	//         "id": "1234",
	//         "type": "book",
	//         "attributes": {
	//             "name": "JSON:API in Go"
	//         },
	//         "relationships": {
	//             "author": {
	//                 "data": {
	//                     "id": "9876",
	//                     "type": "author"
	//                 }
	//             }
	//         }
	//     },
	//     "included": [
	//         {
	//             "id": "9876",
	//             "type": "author",
	//             "attributes": {
	//                 "first_name": "Jane",
	//                 "last_name": "Doe"
	//             }
	//         }
	//     ]
	// }
}
