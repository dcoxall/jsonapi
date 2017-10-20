package jsonapi

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
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

	var result Document

	if val.Kind() == reflect.Slice {
		result = newSliceDocument()
		for i := 0; i < val.Len(); i++ {
			result.AppendData(buildNode(cache, val.Index(i), result, false))
		}
	} else {
		result = newSingleDocument()
		result.AppendData(buildNode(cache, val, result, false))
	}

	err = encoder.Encode(result)

	return
}

func buildNode(cache Cache, val reflect.Value, includer Includer, addIncluded bool) ResourceObject {
	resource := NewResourceObject()

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
		// TODO: handle error
		resourceID, _ := formatID(value)
		resource.SetID(resourceID)
		resourceType, _ := tag.GetTypeIdentifier()
		resource.SetType(resourceType)
	}

	// at this stage we should have a cacheable node so let's do that
	if ok, existing := cache.Add(resource); !ok {
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
			resource.GetAttributes().Append(attrName, attribute(tag, value))
		} else if tag.IsRelation() {
			relationName, _ := tag.GetRelationName()
			if value.Kind() == reflect.Ptr {
				value = reflect.Indirect(value)
			}

			var relManager Relationship
			if value.Kind() == reflect.Slice {
				relManager = newSliceRelationship()
				for x := 0; x < value.Len(); x++ {
					relManager.Append(
						buildNode(cache, value.Index(x), includer, true),
					)
				}
			} else {
				relManager = newSingleRelationship()
				relManager.Append(buildNode(cache, value, includer, true))
			}
			resource.GetRelationships().Append(relationName, relManager)
		}
	}

	if addIncluded {
		includer.AppendIncluded(resource)
	}

	return resource
}

func attribute(tag Tag, val reflect.Value) interface{} {
	// TODO: handle types other than string
	return val.String()
}

// JSONAPI expects IDs to be a string so we need to convert
func formatID(value reflect.Value) (string, error) {
	if value.Kind() == reflect.Interface {
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10), nil
	case reflect.Float32:
		v, _ := value.Interface().(float32)
		return fmt.Sprintf("%G", v), nil
	case reflect.Float64:
		return fmt.Sprintf("%G", value.Float()), nil
	case reflect.Complex64:
		v, _ := value.Interface().(complex64)
		return fmt.Sprintf("%G", v), nil
	case reflect.Complex128:
		return fmt.Sprintf("%G", value.Complex()), nil
	case reflect.String:
		return value.String(), nil
	default:
		return "", fmt.Errorf("Unable to convert ID to string")
	}

}
