// model_test.go is part of the Go models package.
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2024, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided
// that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and 
//    the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions
//    and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or
//    promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, 
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE
// USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
package models

import (
	"bytes"
	"testing"

	// 3rd Party packages
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

func inList(l []string, s string) bool {
	for _, val := range l {
		if val == s {
			return true
		}
	}
	return false
}

// TestModel test a model's methods
func TestModel(t *testing.T) {
	m := new(Model)
	if m.HasChanges() {
		t.Errorf("A new empty model should not have changed yet")
	}
	if m.HasElement("id") {
		t.Errorf("A new empty model should not have an id yet")
	}
	if elem, ok := m.GetModelIdentifier(); ok || elem != nil {
		t.Errorf("A new model should not have a identifier assigned yet, got %+v, %t", elem, ok)
	}
	if attrIds := m.GetAttributeIds(); len(attrIds) > 0 {
		t.Errorf("A new model should not have attributes yet, got %+v", attrIds)
	}
	if elemIds := m.GetElementIds(); len(elemIds) > 0 {
		t.Errorf("A new model should not have element ids yet, got %+v", elemIds)
	}
	if elem, ok := m.GetElementById("name"); ok || elem != nil {
		t.Errorf("A new model should not have an element called 'name', got %+v, %t", elem, ok)
	}
	txt := `id: test_model
attributes:
  method: GET
  action: ./
elements:
  - id: id
    type: text
    attributes:
      required: true
      name: id
    is_primary_id: true
  - id: name
    type: text
    attributes:
      name: name
      required: "true"
  - id: msg
    type: textarea
    attributes:
      name: msg
  - id: updated
    type: text
    attributes:
      name: updated
    generator: current_timestamp
  - id: created
    type: text
    atteributes:
      name: created
    generator: created_timestamp
`
	if err := yaml.Unmarshal([]byte(txt), m); err != nil {
		t.Errorf("expected to be able to unmarshal yaml into model, %s", err)
		t.FailNow()
	}
	buf := bytes.NewBuffer([]byte{})
	if !m.Check(buf) {
		t.Errorf("expected valid model, got %s", buf.Bytes())
		t.FailNow()
	}
	expectedAttr := []string{"method", "action", "elements"}
	for _, attr := range m.GetAttributeIds() {
		if !inList(expectedAttr, attr) {
			t.Errorf("expected %q to be in attribute list %+v", attr, expectedAttr)
		}
	}
	expectedElemIds := []string{"id", "name", "msg", "updated"}
	elemIds := m.GetElementIds()
	for _, elemId := range expectedElemIds {
		if !inList(elemIds, elemId) {
			t.Errorf("expected element id %q to be in list %+v", elemId, elemIds)
		}
	}
	primaryId := m.GetPrimaryId()
	if primaryId == "" {
		t.Errorf("expected %q, got %q", "id", primaryId)
	}

	generatedTypes := m.GetGeneratedTypes()
	if len(generatedTypes) != 2 {
		t.Errorf("expected 2 generator type elements, got %d", len(generatedTypes))
	}
	if val, ok := generatedTypes["updated"]; !ok {
		t.Errorf("expected updated to be %t, got %t in generator type", true, ok)
	} else if val != "current_timestamp" {
		t.Errorf("expected %q, got %q", "current_timestamp", val)
	}
	if val, ok := generatedTypes["created"]; !ok {
		t.Errorf("expected created to be %t, got %t in generator type", true, ok)
	} else if val != "created_timestamp" {
		t.Errorf("expected created %q, got %q", "created_timestamp", val)
	}
}

// TestModelBuilding tests creating a newmodel programatticly
func TestModelBuilding(t *testing.T) {
	modelId := "test_model"
	m, err := NewModel(modelId)
	if err != nil {
		t.Errorf("failed to create new model %q, %s", modelId, err)
	}
	m.isChanged = false
	if m.HasChanges() {
		t.Errorf("%s should not have changes yet", modelId)
	}
	// Example YAML expression of a model
	buf := bytes.NewBuffer([]byte{})
	if !m.Check(buf) {
		t.Errorf("expected a valid model, got %s", buf.Bytes())
		t.FailNow()
	}
	/*
	   func (e *Element) Check(buf io.Writer) bool {
	   func IsValidVarname(s string) bool {
	   func NewElement(elementId string) (*Element, error) {
	   func (model *Model) InsertElement(pos int, element *Element) error {
	   func (model *Model) UpdateElement(elementId string, element *Element) error {
	   func (model *Model) RemoveElement(elementId string) error {
	*/
}

// TestHelperFuncs test the funcs from util.go
func TestHelperFuncs(t *testing.T) {
	m := map[string]string{
		"one":   "1",
		"two":   "2",
		"three": "3",
	}
	attrNames := []string{"one", "two", "three"}
	got := getAttributeIds(m)
	if len(got) != 3 {
		t.Errorf("expected 3 attribute ids, got %d %+v", len(got), got)
		t.FailNow()
	}
	for _, expected := range attrNames {
		if !inList(got, expected) {
			t.Errorf("expected %q in %+v, missing", expected, got)
		}
	}
}

// TestValidateModel tests the model validation based on YAML input.
func TestValidateModel(t *testing.T) {
	src := []byte(`id: test_validator
description: This is a test of the validation code
elements:
  - id: pid
    type: text
    attributes:
      name: pid
      required: true
    is_primary_id: true
    label: Personal Identifier
  - id: lived
    type: text
    attributes:
      name: lived
      required: true
    label: Lived Name
  - id: family
    type: text
    attributes:
      name: family
      required: true
    label: Family Name
  - id: orcid
    type: text
    pattern: "[0-9]{4}-[0-9]{4}-[0-9]{4}-[0-9]{3}[0-9A-Z]"
    attributes:
      name: orcid
      required: true
    label: ORCID
`)
	model, err := NewModel("test_model")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := yaml.Unmarshal(src, &model); err != nil {
		t.Error(err)
		t.FailNow()
	}
	SetDefaultTypes(model)

	formData := map[string]string{
		"pid":    "jane-doe",
		"lived":  "Jane",
		"family": "Doe",
		"orcid":  "0000-1111-2222-3333",
	}
	if ok := model.Validate(formData); !ok {
		t.Errorf("%+v failed to validate", formData)
	}
}

// TestValidateMapInterface tests the YAML model mapping
// for decoding and encoding models.
func TestValidateMapInterface(t *testing.T) {
	src := []byte(`id: test_validate_map_inteface
description: This is a test of the validation code
elements:
  - id: pid
    type: text
    attributes:
      name: pid
      required: true
    is_primary_id: true
    label: Personal Identifier
    generator: uuid
  - id: lived
    type: text
    attributes:
      name: lived
      required: true
    label: Lived Name
  - id: family
    type: text
    attributes:
      name: family
      required: true
    label: Family Name
  - id: orcid
    type: text
    pattern: "[0-9]{4}-[0-9]{4}-[0-9]{4}-[0-9]{3}[0-9A-Z]"
    attributes:
      name: orcid
      required: true
    label: ORCID
  - id: created
    type: datetime-local
    attributes:
      required: true
    label: created
    generator: created_timestmap
  - id: updated
    type: datetime-local
    attributes:
      required: true
    generator: current_timestamp
`)
	model, err := NewModel("test_model")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := yaml.Unmarshal(src, &model); err != nil {
		t.Error(err)
		t.FailNow()
	}
	//Debug = true
	SetDefaultTypes(model)

	pid := uuid.New()
	formData := map[string]interface{}{
		"pid":     pid,
		"lived":   "Jane",
		"family":  "Doe",
		"orcid":   "0000-1111-2222-3333",
		"created": "2024-10-03T12:40:00",
		"updated": "2024-10-03 12:41:32",
	}
	if ok := model.ValidateMapInterface(formData); !ok {
		t.Errorf("%+v failed to validate", formData)
	}

	formData = map[string]interface{}{
		"created": "2024-10-03T13:25:24-07:00",
		"family":  "Jetson",
		"lived":   "George",
		"orcid":   "1234-4321-1234-4321",
		"pid":     "0192540f-0806-7631-b08f-4ae5c4d37cca",
		"updated": "2024-10-03T13:25:24-07:00",
	}
	if ok := model.ValidateMapInterface(formData); !ok {
		t.Errorf("%+v failed to validate", formData)
	}
}

// TestModelElements tests the GetGeneratedTypes func for Models.
func TestModelElements(t *testing.T) {
	m := new(Model)
	modelTypes := m.GetGeneratedTypes()
	if len(modelTypes) != 0 {
		t.Errorf("expected zero model types, got %+v", modelTypes)
	}
}

// TestNestedValidation tests validation of nested structures (objects and lists)
func TestNestedValidation(t *testing.T) {
	// Create a model with nested author object
	txt := `id: article
Description: Test model with nested structures
elements:
  - id: doi
    type: doi
    attributes:
      required: true
  - id: title
    type: text
    attributes:
      required: true
  - id: author
    type: text
    is_list: true
    elements:
      - id: given
        type: text
        attributes:
          required: true
      - id: family
        type: text
        attributes:
          required: true
      - id: orcid
        type: orcid
        attributes:
          required: false
  - id: affiliation
    type: text
    is_object: true
    elements:
      - id: name
        type: text
      - id: ror
        type: ror
`
	model := new(Model)
	if err := yaml.Unmarshal([]byte(txt), model); err != nil {
		t.Fatalf("Failed to parse model YAML: %v", err)
	}
	SetDefaultTypes(model)

	// Test valid nested data
	validData := map[string]interface{}{
		"doi": "10.1234/test",
		"title": "Test Article",
		"author": []interface{}{
			map[string]interface{}{
				"given":  "Jane",
				"family": "Doe",
				"orcid":  "0000-0002-1825-0097",
			},
			map[string]interface{}{
				"given":  "John",
				"family": "Smith",
			},
		},
		"affiliation": map[string]interface{}{
			"name": "Caltech",
			"ror":  "05dxps055",
		},
	}

	if !model.ValidateInterface(validData) {
		t.Errorf("Valid nested data failed validation")
	}

	// Test invalid - missing required field in nested author
	invalidAuthor := map[string]interface{}{
		"doi": "10.1234/test",
		"title": "Test Article",
		"author": []interface{}{
			map[string]interface{}{
				"given": "Jane",
				// missing family
			},
		},
	}
	if model.ValidateInterface(invalidAuthor) {
		t.Errorf("Invalid nested data (missing family) passed validation")
	}

	// Test invalid - wrong type for author (should be list)
	invalidAuthorType := map[string]interface{}{
		"doi":    "10.1234/test",
		"title":  "Test Article",
		"author": "Jane Doe", // string instead of list
	}
	if model.ValidateInterface(invalidAuthorType) {
		t.Errorf("Invalid author type (string instead of list) passed validation")
	}

	// Test invalid DOI
	invalidDOI := map[string]interface{}{
		"doi": "not-a-doi",
		"title": "Test Article",
		"author": []interface{}{
			map[string]interface{}{
				"given":  "Jane",
				"family": "Doe",
			},
		},
	}
	if model.ValidateInterface(invalidDOI) {
		t.Errorf("Invalid DOI passed validation")
	}
}

// TestGetNestedElement tests the GetNestedElement method
func TestGetNestedElement(t *testing.T) {
	// Create a model with nested structure
	txt := `id: crossref
Description: CrossRef-like structure
elements:
  - id: doi
    type: doi
  - id: author
    type: text
    is_list: true
    elements:
      - id: given
        type: text
      - id: family
        type: text
      - id: affiliation
        type: text
        is_object: true
        elements:
          - id: name
            type: text
          - id: ror
            type: ror
`
	model := new(Model)
	if err := yaml.Unmarshal([]byte(txt), model); err != nil {
		t.Fatalf("Failed to parse model YAML: %v", err)
	}

	// Test getting top-level element
	if _, ok := model.GetNestedElement("doi"); !ok {
		t.Errorf("Failed to get top-level element 'doi'")
	}

	// Test getting nested element
	if _, ok := model.GetNestedElement("author"); !ok {
		t.Errorf("Failed to get nested element 'author'")
	}

	// Test getting deeply nested element (requires path notation)
	// For now, GetNestedElement only handles direct children
	// This test documents current limitation
	if elem, ok := model.GetNestedElement("author.given"); !ok {
		t.Logf("Note: GetNestedElement does not currently support dotted paths (author.given)")
	} else {
		t.Logf("Got element: %+v", elem)
	}
}

// TestElementCheck tests the Check method for elements with nested structures
func TestElementCheck(t *testing.T) {
	// Create a valid element with nested elements
	elem, err := NewElement("author")
	if err != nil {
		t.Fatalf("Failed to create element: %v", err)
	}
	elem.IsList = true
	
	// Add nested elements
	given, _ := NewElement("given")
	given.Type = "text"
	family, _ := NewElement("family")
	family.Type = "text"
	
	elem.Elements = append(elem.Elements, given)
	elem.Elements = append(elem.Elements, family)

	var buf bytes.Buffer
	if !elem.Check(&buf) {
		t.Errorf("Valid nested element failed Check: %s", buf.String())
	}

	// Create an invalid element (list without nested elements)
	invalidElem, _ := NewElement("invalid_list")
	invalidElem.IsList = true
	// No nested elements added
	
	var buf2 bytes.Buffer
	if invalidElem.Check(&buf2) {
		t.Errorf("Invalid element (list without nested elements) passed Check")
	}
}
