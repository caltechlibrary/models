package models

import (
	"fmt"
	"io"
	"strings"
)

// GenElementFunc is a function which will generate an Element configured represent a model's supported "types"
type GenElementFunc func() *Element

// ValidateFunc is a function that validates form assocaited with the Element and the string value
// received in the web form (value before converting to Go type).
type ValidateFunc func(*Element, string) bool


// Element implementes the GitHub YAML issue template syntax for an input element.
// The input element YAML is described at <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-githubs-form-schema>
//
// While the syntax most closely express how to setup an HTML representation it is equally
// suitable to expressing, through inference, SQL column type definitions. E.g. a bare `input` type is a `varchar`,
// a `textarea` is a `text` column type, an `input[type=date]` is a date column type.
type Element struct {
	// Type, The type of element that you want to input. It is required. Valid values are
	// checkboxes, dropdown, input, markdown and text area.
	//
	// The input type corresponds to the native input types defined for HTML 5. E.g. text, textarea,
	// email, phone, date, url, checkbox, radio, button, submit, cancel, select
	// See MDN developer docs for input, <https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input>
	Type string `json:"type,required" yaml:"type,omitempty"`

	// Id for the element, except when type is set to markdown. Can only use alpha-numeric characters,
	//  -, and _. Must be unique in the form definition. If provided, the id is the canonical identifier
	//  for the field in URL query parameter prefills.
	Id string `json:"id,omitempty" yaml:"id,omitempty"`

	// Attributes, a set of key-value pairs that define the properties of the element.
	// This is a required element as it holds the "value" attribute when expressing
	// HTML content. Other commonly use attributes
	Attributes map[string]string `json:"attributes,omitempty" yaml:"attributes,omitempty"`

	// Pattern holds a validation pattern. When combined with an input type (or input type alias, e.g. orcid)
	// produces a form element that sports a specific client side validation exceptation.  This intern can be used
	// to generate appropriate validation code server side.
	Pattern string `json:"pattern,omitempty" yaml:"pattern,omitempty"`

	// Options holds a list of values and their labels used for HTML select elements in rendering their option child elements
	Options []map[string]string `json:"optoins,omitempty" yaml:"options,omitempty"`

	// IsObjectId (i.e. is the identifier of the object) used by for the modeled data.
	// It is used in calculating routes and templates where the object identifier is required.
	IsObjectId bool `json:"is_primary_id,omitempty" yaml:"is_primary_id,omitempty"`

	// Generator indicates the type of automatic population of a field. It is used to
	// indicate autoincrement and uuids for primary keys and timestamps for datetime oriented fields.
	Generator string `json:"generator,omitempty" yaml:"generator,omitempty"`

	// Label is used when rendering an HTML form as a label element tied to the input element via the set attribute and
	// the element's id.
	Label string `json:"label,omitempty" yaml:"label,omitempty"`

	//
	// These fields are used by the modeler to manage the models and their elements
	//
	isChanged bool `json:"-" yaml:"-"`
}

// NewElement, makes sure element id is valid, populates an element as a basic input type.
// The new element has the attribute "name" and label set to default values.
func NewElement(elementId string) (*Element, error) {
	if !IsValidVarname(elementId) {
		return nil, fmt.Errorf("invalid element id, %q", elementId)
	}
	element := new(Element)
	element.Id = elementId
	element.Attributes = map[string]string{"name": elementId}
	element.Type = "text"
	element.Label = strings.ToUpper(elementId[0:1]) + elementId[1:]
	element.IsObjectId = false
	element.isChanged = true
	return element, nil
}

// HasChanged checks to see if the Element has been changed.
func (e *Element) HasChanged() bool {
	return e.isChanged
}

// Changed sets the change state on element
func (e *Element) Changed(state bool) {
	e.isChanged = state
}

// Check reviews an Element to make sure if is value.
func (e *Element) Check(buf io.Writer) bool {
	ok := true
	if e == nil {
		fmt.Fprintf(buf, "element is nil\n")
		ok = false
	}
	if e.Id == "" {
		fmt.Fprintf(buf, "element missing id\n")
		ok = false
	}
	if e.Type == "" {
		fmt.Fprintf(buf, "element, %q, missing type\n", e.Id)
		ok = false
	}
	return ok
}
