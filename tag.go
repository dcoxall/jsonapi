package jsonapi

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	tagName      = "jsonapi"
	tagAttribute = "attr"
	tagPrimary   = "primary"
	tagRelation  = "relation"
	tagIgnore    = "-"
)

type Tag interface {
	IsAttribute() bool
	GetAttributeName() (string, error)

	IsPrimary() bool
	GetTypeIdentifier() (string, error)

	IsRelation() bool
	GetRelationName() (string, error)

	IsIgnore() bool
}

type coreTag struct {
	hasAttribute bool
	hasPrimary   bool
	hasRelation  bool
	hasIgnore    bool
	value        string
}

func (tag *coreTag) IsAttribute() bool {
	return tag.hasAttribute
}

func (tag *coreTag) IsPrimary() bool {
	return tag.hasPrimary
}

func (tag *coreTag) IsRelation() bool {
	return tag.hasRelation
}

func (tag *coreTag) IsIgnore() bool {
	return tag.hasIgnore
}

func (tag *coreTag) GetAttributeName() (string, error) {
	if tag.IsAttribute() {
		return tag.value, nil
	}
	return "", fmt.Errorf("Tag does not represent an attribute")
}

func (tag *coreTag) GetRelationName() (string, error) {
	if tag.IsRelation() {
		return tag.value, nil
	}
	return "", fmt.Errorf("Tag does not represent a relation")
}

func (tag *coreTag) GetTypeIdentifier() (string, error) {
	if tag.IsPrimary() {
		return tag.value, nil
	}
	return "", fmt.Errorf("Tag does not represent the model type")
}

func ParseTag(field reflect.StructField) (Tag, error) {
	tagValue := field.Tag.Get(tagName)
	tag := &coreTag{}
	if tagValue == "" {
		// No tag then we assume an attribute
		tag.hasAttribute = true
		tag.value = field.Name
		return tag, nil
	}
	parts := strings.Split(tagValue, ",")
	switch parts[0] {
	case tagAttribute:
		tag.hasAttribute = true
	case tagPrimary:
		tag.hasPrimary = true
	case tagRelation:
		tag.hasRelation = true
	case tagIgnore:
		tag.hasIgnore = true
		return tag, nil
	default:
		return nil, fmt.Errorf("'%s' is an invalid jsonapi struct tag", tagValue)
	}

	if !tag.IsIgnore() && (len(parts) < 2 || len(parts) > 3) {
		return nil, fmt.Errorf("'%s' is an invalid jsonapi struct tag", tagValue)
	}

	tag.value = parts[1]

	return tag, nil
}
