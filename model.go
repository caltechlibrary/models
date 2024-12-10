package models

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
)

// RenderFunc is a function thation takes an io.Writer and Model then
// renders the model into the io.Writer. It is used to extend the Model to
// support various output formats.
type RenderFunc func(io.Writer, *Model) error

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
	ids := model.GetElementIds()
	if len(ids) != len(formData) {
		return false
	}
	for k, v := range formData {
		if elem, ok := model.GetElementById(k); ok {
			if validator, ok := model.validators[elem.Type]; ok {
				if !validator(elem, v) {
					if Debug {
						log.Printf("DEBUG failed to validate elem.Id %q, elem.Type %q, value %q", elem.Id, elem.Type, v)
					}
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

// ValidateMapInterface normalizes the map inteface values before calling
// the element's validator function.
func (model *Model) ValidateMapInterface(data map[string]interface{}) bool {
	if model == nil {
		if Debug {
			log.Printf("model is nil, can't validate")
		}
		return false
	}
	ids := model.GetElementIds()
	if len(ids) != len(data) {
		if Debug {
			log.Printf("DEBUG expected len(ids) %d, got len(data) %d", len(ids), len(data))
		}
		return false
	}
	for k, v := range data {
		var val string
		switch v.(type) {
		case string:
			val = v.(string)
		case int:
			val = fmt.Sprintf("%d", v)
		case float64:
			val = fmt.Sprintf("%f", v)
		case json.Number:
			val = fmt.Sprintf("%s", v)
		case bool:
			val = fmt.Sprintf("%t", v)
		default:
			val = fmt.Sprintf("%+v", v)
		}
		if elem, ok := model.GetElementById(k); ok {
			if validator, ok := model.validators[elem.Type]; ok {
				if !validator(elem, val) {
					if Debug {
						log.Printf("DEBUG failed to validate elem.Id %q, value %q", elem.Id, val)
					}
					return false
				}
			} else {
				if Debug {
					log.Printf("DEBUG failed to validate elem.Id %q, value %q, missing validator", elem.Id, val)
				}
				return false
			}
		} else {
			if Debug {
				log.Printf("DEBUG failed to validate elem.Id %q, value %q, missing elemnt %q", elem.Id, val, k)
			}
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
