package jsonapi

import (
	"encoding/json"
	"io"
	"reflect"
)

type Encoder interface {
	Encode(interface{}) error
}

type coreEncoder struct {
	target io.Writer
}

func NewEncoder(writer io.Writer) Encoder {
	return &coreEncoder{
		target: writer,
	}
}

func (enc *coreEncoder) Encode(model interface{}) (err error) {
	// create a cache to store visited nodes
	cache := NewCache()
	encoder := json.NewEncoder(enc.target)

	// TODO: fix duplication in buildNode
	// resolve pointer if appropriate
	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	var result JSONAPIResult

	if val.Kind() == reflect.Slice {
		result = NewSliceResponse()
		for i := 0; i < val.Len(); i++ {
			result.AppendData(buildNode(cache, val.Index(i), result, false))
		}
	} else {
		result = NewSingleResponse()
		result.AppendData(buildNode(cache, val, result, false))
	}

	err = encoder.Encode(result)

	return
}

func buildNode(cache Cache, val reflect.Value, includer Includer, addIncluded bool) Node {
	node := &coreNode{
		Attributes:    NewAttributes(),
		Relationships: NewRelationships(),
	}

	// resolve pointer if appropriate
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		// at this stage we only want to find the id and type
		// but we also need to watch relationships so we don't
		// need to iterate over it all again
		tag := ParseTag(field.Tag.Get(TagName))
		if !tag.IsPrimary() {
			continue
		}

		// set the id and type for the node
		node.ID = value.String()
		nodeType, _ := tag.GetTypeIdentifier()
		node.Type = nodeType
	}

	// at this stage we should have a cacheable node so let's do that
	if ok, existing := cache.Add(node); !ok {
		// unable to add to cache so let's assume it was already added
		// therefore we don't need to traverse the objects relationships
		// and attributes. Instead we can just return the existing
		// reference
		return existing
	}

	// so the node was added to the cache. this should indicate that it
	// is the first time walking this node so we need to add attributes
	// and follow relationships. this does mean we have to re-iterate over
	// the fields
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		tag := ParseTag(field.Tag.Get(TagName))
		// this time we can skip the primary field
		if tag.IsPrimary() {
			continue
		} else if tag.IsAttribute() {
			attrName, _ := tag.GetAttributeName()
			node.GetAttributes().Append(attrName, attribute(tag, value))
		} else if tag.IsRelation() {
			relationName, _ := tag.GetRelationName()
			if value.Kind() == reflect.Ptr {
				value = reflect.Indirect(value)
			}

			var relManager RelationshipManager
			if value.Kind() == reflect.Slice {
				relManager = NewSliceRelationship()
				for x := 0; x < value.Len(); x++ {
					relManager.Append(
						buildNode(cache, value.Index(x), includer, true),
					)
				}
			} else {
				relManager = NewSingleRelationship()
				relManager.Append(buildNode(cache, value, includer, true))
			}
			node.GetRelationships().Append(relationName, relManager)
		}
	}

	if addIncluded {
		includer.AppendIncluded(node)
	}

	return node
}

func attribute(tag Tag, val reflect.Value) interface{} {
	// TODO: handle types other than string
	return val.String()
}
