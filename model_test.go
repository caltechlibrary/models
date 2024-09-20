package dataset

import (
	"bytes"
	"testing"

	// 3rd Party packages
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
	expectedElemIds := []string{"id", "name", "msg"}
	elemIds := m.GetElementIds()
	for _, elemId := range expectedElemIds {
		if !inList(elemIds, elemId) {
			t.Errorf("expected element id %q to be in list %+v", elemId, elemIds)
		}
	}
}

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
	   func isValidVarname(s string) bool {
	   func NewElement(elementId string) (*Element, error) {
	   func (model *Model) InsertElement(pos int, element *Element) error {
	   func (model *Model) UpdateElement(elementId string, element *Element) error {
	   func (model *Model) RemoveElement(elementId string) error {
	*/
}

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
