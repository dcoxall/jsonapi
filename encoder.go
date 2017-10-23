package jsonapi

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"reflect"
	"strconv"
)

type Encoder interface {
	Encode(interface{}) error
}

type coreEncoder struct {
	encoder *json.Encoder
}

func NewEncoder(writer io.Writer) Encoder {
	return &coreEncoder{
		encoder: json.NewEncoder(writer),
	}
}

func NewIndentEncoder(writer io.Writer, prefix, indent string) Encoder {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent(prefix, indent)
	return &coreEncoder{
		encoder: encoder,
	}
}

func (enc *coreEncoder) Encode(model interface{}) error {
	cache := NewCache()

	// TODO: fix duplication in buildNode
	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	var result Document

	if val.Kind() == reflect.Slice {
		result = newSliceDocument()
		for i := 0; i < val.Len(); i++ {
			resource, err := buildResource(cache, val.Index(i), result, false)
			if err != nil {
				return errors.Wrap(err, "unable to serialize resource")
			}
			result.AppendData(resource)
		}
	} else {
		result = newSingleDocument()
		resource, err := buildResource(cache, val, result, false)
		if err != nil {
			return errors.Wrap(err, "unable to serialize resource")
		}
		result.AppendData(resource)
	}

	return enc.encoder.Encode(result)
}

func buildResource(cache Cache, val reflect.Value, includer Includer, addIncluded bool) (ResourceObject, error) {
	resource := newResourceObject()

	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		tag, err := ParseTag(field)
		if err != nil {
			return nil, errors.Wrap(err, "invalid jsonapi struct tag")
		}
		if !tag.IsPrimary() {
			continue
		}

		resourceID, err := formatID(value)
		if err != nil {
			return nil, errors.Wrap(err, "unable to serialize resource id")
		}
		resource.SetID(resourceID)
		resourceType, err := tag.GetTypeIdentifier()
		if err != nil {
			return nil, errors.Wrap(err, "unable to serialize resource type")
		}
		resource.SetType(resourceType)
	}

	if ok, existing := cache.Add(resource); !ok {
		return existing, nil
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		tag, err := ParseTag(field)
		if err != nil {
			return nil, errors.Wrap(err, "invalid jsonapi struct tag")
		}

		if tag.IsPrimary() {
			continue
		} else if tag.IsAttribute() {
			attrName, err := tag.GetAttributeName()
			if err != nil {
				return nil, errors.Wrap(err, "unable to serialize attribute key")
			}
			// Using value.Interface() so the json encoder can handle the
			// serialization of attributes
			resource.GetAttributes().Append(attrName, value.Interface())
		} else if tag.IsRelation() {
			relationName, err := tag.GetRelationName()
			if err != nil {
				return nil, errors.Wrap(err, "unable to serialize relation key")
			}
			if value.Kind() == reflect.Ptr {
				value = reflect.Indirect(value)
			}

			var relManager Relationship
			if value.Kind() == reflect.Slice {
				relManager = newSliceRelationship()
				for x := 0; x < value.Len(); x++ {
					relationResource, err := buildResource(cache, value.Index(x), includer, true)
					if err != nil {
						return nil, errors.Wrap(err, "unable to serialize relation resource")
					}
					relManager.Append(relationResource)
				}
			} else {
				relManager = newSingleRelationship()
				relationResource, err := buildResource(cache, value, includer, true)
				if err != nil {
					return nil, errors.Wrap(err, "unable to serialize relation resource")
				}
				relManager.Append(relationResource)
			}
			resource.GetRelationships().Append(relationName, relManager)
		}
	}

	if addIncluded {
		includer.AppendIncluded(resource)
	}

	return resource, nil
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
