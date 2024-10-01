package models

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
)

// RenderFunc is a function thation takes an io.Writer and Model then
// renders the model into the io.Writer. It is used to extend the Model to
// support various output formats.
type RenderFunc func(io.Writer, *Model) error

// GenElementFunc is a function which will generate an Element configured represent a model's supported "types"
type GenElementFunc func() *Element

// ValidateFunc is a function that validates form assocaited with the Element and the string value
// received in the web form (value before converting to Go type).
type ValidateFunc func(*Element, string) bool

// Model implements a data structure description inspired by GitHub YAML issue template syntax.
// See <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-issue-forms>
//
// The Model structure describes the HTML elements used to form a record. It can be used in code generation and in validating
// POST and PUT requests in datasetd.
type Model struct {
	// Id is a required field for model, it maps to the HTML element id and name
	Id string `json:"id,required" yaml:"id"`

	// This is a Newt specifc set of attributes to place in the form element of HTML. I.e. it could
	// be form "class", "method", "action", "encoding". It is not defined in the GitHub YAML issue template syntax
	// (optional)
	Attributes map[string]string `json:"attributes,omitempty" yaml:"attributes,omitempty"`

	// Description, A description for the issue form template, which appears in the template chooser interface.
	// (required)
	Description string `json:"description,required" yaml:"description,omitempty"`

	// Elements, Definition of the input types in the form.
	// (required)
	Elements []*Element `json:"elements,required" yaml:"elements,omitempty"`

	// Title, A default title that will be pre-populated in the issue submission form.
	// (optional) only there for compatibility with GitHub YAML Issue Templates
	//Title string `json:"title,omitempty" yaml:"title,omitempty"`

	// isChanged is an internal state used by the modeler to know when a model has changed
	isChanged bool `json:"-" yaml:"-"`

	// renderer is a map of names to RenderFunc functions. A RenderFunc is that take a io.Writer and the model object as parameters then
	// return an error type.  This allows for many renderers to be used with Model by
	// registering the function then envoking render with the name registered.
	renderer map[string]RenderFunc `json:"-" yaml:"-"`

	// genElements holds a map to the "type" pointing to an element generator
	genElements map[string]GenElementFunc `json:"-" yaml:"-"`

	// validators holds a list of validate function associated with types. Key is type name.
	validators map[string]ValidateFunc `json:"-" yaml:"-"`
}

// GenElementType takes an element type and returns an Element struct populated for that type and true or nil and false if type is not supported.
func (model *Model) GenElementType(typeName string) (*Element, bool) {
	if fn, ok := model.genElements[typeName]; ok {
		return fn(), true
	}
	return nil, false
}

// Validate form data expressed as map[string]string.
func (model *Model) Validate(formData map[string]string) bool {
	for k, v := range formData {
		if elem, ok := model.GetElementById(k); ok {
			if validator, ok := model.validators[elem.Type]; ok {
				if !validator(elem, v) {
					return false
				}
			} else {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

// HasChanges checks if the model's elements have changed
func (model *Model) HasChanges() bool {
	if model.isChanged {
		return true
	}
	for _, e := range model.Elements {
		if e.isChanged {
			return true
		}
	}
	return false
}

// Changed sets the change state
func (model *Model) Changed(state bool) {
	model.isChanged = state
}

// HasElement checks if the model has a given element id
func (model *Model) HasElement(elementId string) bool {
	for _, e := range model.Elements {
		if e.Id == elementId {
			return true
		}
	}
	return false
}

// HasElementType checks if an element type matches given type.
func (model *Model) HasElementType(elementType string) bool {
	for _, e := range model.Elements {
		if strings.ToLower(e.Type) == strings.ToLower(elementType) {
			return true
		}
	}
	return false
}

// GetModelIdentifier() returns the element which describes the model identifier.
// Returns the element and a boolean set to true if found.
func (m *Model) GetModelIdentifier() (*Element, bool) {
	for _, e := range m.Elements {
		if e.IsObjectId {
			return e, true
		}
	}
	return nil, false
}

// GetAttributeIds returns a slice of attribute ids found in the model's .Elements
func (m *Model) GetAttributeIds() []string {
	return getAttributeIds(m.Attributes)
}

// GetElementIds returns a slice of element ids found in the model's .Elements
func (m *Model) GetElementIds() []string {
	ids := []string{}
	for _, elem := range m.Elements {
		if elem.Id != "" {
			ids = append(ids, elem.Id)
		}
	}
	return ids
}

// GetPrimaryId returns the primary id
func (m *Model) GetPrimaryId() string {
	for _, elem := range m.Elements {
		if elem.IsObjectId {
			return elem.Id
		}
	}
	return ""
}

// GetGeneratedTypes returns a map of elemend id and value held by .Generator
func (m *Model) GetGeneratedTypes() map[string]string {
	gt := map[string]string{}
	for _, elem := range m.Elements {
		if elem.Generator != "" {
			gt[elem.Id] = elem.Generator
		}
	}
	return gt
}

// GetElementById returns a Element from the model's .Elements.
func (m *Model) GetElementById(id string) (*Element, bool) {
	for _, elem := range m.Elements {
		if elem.Id == id {
			return elem, true
		}
	}
	return nil, false
}

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

// IsValidVarname tests a sting confirms it conforms to Model's naming rule.
func IsValidVarname(s string) bool {
	if len(s) == 0 {
		return false
	}
	// NOTE: variable names must start with a latter and maybe followed by
	// one or more letters, digits and underscore.
	vRe := regexp.MustCompile(`^([a-zA-Z]|[a-zA-Z][0-9a-zA-Z\_]+)$`)
	return vRe.Match([]byte(s))
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

// NewModel, makes sure model id is valid, populates a Model with the identifier element providing
// returns a *Model and error value.
func NewModel(modelId string) (*Model, error) {
	if !IsValidVarname(modelId) {
		return nil, fmt.Errorf("invalid model id, %q", modelId)
	}
	model := new(Model)
	model.Id = modelId
	model.Description = fmt.Sprintf("... description of %q goes here ...", modelId)
	model.Attributes = map[string]string{}
	model.Elements = []*Element{}
	// Make the required element ...
	element := new(Element)
	element.Id = "id"
	element.IsObjectId = true
	element.Type = "text"
	element.Attributes = map[string]string{"required": "true"}
	if err := model.InsertElement(0, element); err != nil {
		return nil, err
	}
	return model, nil
}

// Check analyze the model and make sure at least one element exists and the
// model has a single identifier (e.g. "identifier")
func (model *Model) Check(buf io.Writer) bool {
	if model == nil {
		fmt.Fprintf(buf, "model is nil\n")
		return false
	}
	if model.Elements == nil {
		fmt.Fprintf(buf, "missing %s.body\n", model.Id)
		return false
	}
	// Check to see if we have at least one element in Elements
	if len(model.Elements) > 0 {
		ok := true
		hasModelId := false
		for i, e := range model.Elements {
			// Check to make sure each element is valid
			if !e.Check(buf) {
				fmt.Fprintf(buf, "error for %s.%s\n", model.Id, e.Id)
				ok = false
			}
			if e.IsObjectId {
				if hasModelId == true {
					fmt.Fprintf(buf, "duplicate model identifier element (%d) %s.%s\n", i, model.Id, e.Id)
					ok = false
				}
				hasModelId = true
			}
		}
		if !hasModelId {
			fmt.Fprintf(buf, "missing required object identifier for model %s\n", model.Id)
			ok = false
		}
		return ok
	}
	fmt.Fprintf(buf, "Missing elements for model %q\n", model.Id)
	return false
}

// InsertElement will add a new element to model.Elements in the position indicated,
// It will also set isChanged to true on additional.
func (model *Model) InsertElement(pos int, element *Element) error {
	if model.Elements == nil {
		model.Elements = []*Element{}
	}
	if !IsValidVarname(element.Id) {
		return fmt.Errorf("element id is not value")
	}
	if model.HasElement(element.Id) {
		return fmt.Errorf("duplicate element id, %q", element.Id)
	}
	if pos < 0 {
		pos = 0
	}
	if pos > len(model.Elements) {
		model.Elements = append(model.Elements, element)
		model.isChanged = true
		return nil
	}
	if pos < len(model.Elements) {
		elements := append(model.Elements[:pos], element)
		model.Elements = append(elements, model.Elements[(pos+1):]...)
	} else {
		model.Elements = append(model.Elements, element)
	}
	model.isChanged = true
	return nil
}

// UpdateElement will update an existing element with element id will the new element.
func (model *Model) UpdateElement(elementId string, element *Element) error {
	if !model.HasElement(elementId) {
		return fmt.Errorf("%q element id not found", elementId)
	}
	for i, e := range model.Elements {
		if e.Id == elementId {
			model.Elements[i] = element
			model.isChanged = true
			return nil
		}
	}
	return fmt.Errorf("failed to find %q to update", elementId)
}

// RemoveElement removes an element by id from the model.Elements
func (model *Model) RemoveElement(elementId string) error {
	if !model.HasElement(elementId) {
		return fmt.Errorf("%q element id not found", elementId)
	}
	for i, e := range model.Elements {
		if e.Id == elementId {
			model.Elements = append(model.Elements[:i], model.Elements[(i+1):]...)
			model.isChanged = true
			return nil
		}
	}
	return fmt.Errorf("%q element id is missing", elementId)
}

// ToSQLiteScheme takes a model and trys to render a SQLite3 SQL create statement.
func (model *Model) ToSQLiteScheme(out io.Writer) error {
	return ModelToSQLiteScheme(out, model)
}

// ToHTML takes a model and trys to render an HTML web form
func (model *Model) ToHTML(out io.Writer) error {
	return ModelToHTML(out, model)
}

// ModelInteractively takes a model and interactively prompts to create
// a YAML model file.
func (model *Model) ModelToYAML(out io.Writer) error {
	return ModelToYAML(out, model)
}

// Register takes a name (string) and a RenderFunc and registers it with the model.
// Registered names then can be invoke by the register name.
func (model *Model) Register(name string, fn RenderFunc) {
	if model.renderer == nil {
		model.renderer = map[string]RenderFunc{}
	}
	model.renderer[name] = fn
}

// Render takes a register render io.Writer and register name envoking the function
// with the model.
func (model *Model) Render(out io.Writer, name string) error {
	if fn, ok := model.renderer[name]; ok {
		return fn(out, model)
	}
	return fmt.Errorf("%s is not a registered rendering function", name)
}

// IsSupportedElementType checks if the element type is supported by Newt, returns true if OK false is it is not
func (model *Model) IsSupportedElementType(eType string) bool {
	for sType, _ := range model.genElements {
		if eType == sType {
			return true
		}
	}
	return false
}

// getAttributeIds returns a list of attribue keys in a maps[string]interface{} structure
func getAttributeIds(m map[string]string) []string {
	ids := []string{}
	for k, _ := range m {
		if k != "" {
			ids = append(ids, k)
		}
	}
	if len(ids) > 0 {
		sort.Strings(ids)
	}
	return ids
}

// Define takes a model and attaches a type definition (an element generator) and validator for the named type
func (model *Model) Define(typeName string, genElementFn GenElementFunc, validateFn ValidateFunc) {
	if model.genElements == nil {
		model.genElements = map[string]GenElementFunc{}
	}
	model.genElements[typeName] = genElementFn
	if model.validators == nil {
		model.validators = map[string]ValidateFunc{}
	}
	model.validators[typeName] = validateFn
}
