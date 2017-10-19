package jsonapi

import (
	"fmt"
	"strings"
)

const (
	TagName      = "jsonapi"
	tagAttribute = "attr"
	tagPrimary   = "primary"
	tagRelation  = "relation"
)

type Tag interface {
	IsAttribute() bool
	GetAttributeName() (string, error)

	IsPrimary() bool
	GetTypeIdentifier() (string, error)

	IsRelation() bool
	GetRelationName() (string, error)
}

type coreTag struct {
	hasAttribute bool
	hasPrimary   bool
	hasRelation  bool
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

func ParseTag(tagValue string) Tag {
	tag := &coreTag{}
	parts := strings.Split(tagValue, ",")
	tag.hasAttribute = parts[0] == tagAttribute
	tag.hasPrimary = parts[0] == tagPrimary
	tag.hasRelation = parts[0] == tagRelation
	tag.value = parts[1]
	return tag
}
