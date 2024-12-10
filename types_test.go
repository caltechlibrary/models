package models

import (
	"regexp"
	"testing"

	// 3rd Party Packages
	"gopkg.in/yaml.v3"
)

// TestVAlidateElementText, checks the validation code
func TestValidateElementText(t *testing.T) {
	// Debug = true
	elem := new(Element)
	elem.Id = "orcid"
	elem.Type = "text"
	elem.Pattern = OrcidPattern
	//val := `0000-1111-2222-333X`
	val := `0000-0003-0900-6903`
	if !ValidateText(elem, val) {
		t.Errorf("Expected true, got false for elem %+v -> val %q", elem, val)
	}
}

// TestORCIDRegExp tests the ORCID validation code
func TestORCIDRegExp(t *testing.T) {
	// Debug = true
	pattern := `[0-9]{4}-[0-9]{4}-[0-9]{4}-[0-9]{3}[0-9A-Z]`
	re := regexp.MustCompilePOSIX(pattern)
	orcid := `0000-0003-0900-6903`
	if !re.MatchString(orcid) {
		t.Errorf("expected true, got false for pattern %q and value %q", pattern, orcid)
	}
	elem := new(Element)
	elem.Id = "orcid"
	elem.Type = "orcid"
	elem.Generator = ""
	//SetDebug(true)
	if !ValidateORCID(elem, orcid) {
		t.Errorf("expected orcid to validate true, got false")
	}
	orcid = `0000-0001-9689-9628`
	if !ValidateORCID(elem, orcid) {
		t.Errorf("expected orcid to validate true, got false")
	}
	// This isn't an valid ORCID throug it could be an INSI
	orcid = `2345-5432-1234-4326`
	if ValidateORCID(elem, orcid) {
		t.Errorf("expected orcid to validate false, got true")
	}
	//SetDebug(false)
}

// TestDatetimeLocal tests the "datetime-local" structure
func TestDatetimeLocal(t *testing.T) {
	// Debug = true
	elem := new(Element)
	elem.Id = "created"
	elem.Type = "datetime-local"
	elem.Generator = "created_timestamp"

	val := "2024-10-03T12:51:01"
	if !ValidateDateTimeLocal(elem, val) {
		t.Errorf("expected true, got false for value %q", val)
	}
	/*
			  "created": "2024-10-03T13:30:28-07:00",
		  "family": "Jetson",
		  "lived": "George",
		  "orcid": "1234-4321-1234-4321",
		  "pid": "01925413-abc0-75c8-aa75-bfc062cd2949",
		  "updated": "2024-10-03T13:30:28-07:00"

	*/
	val = "2024-10-03T13:30:28-07:00"
	if !ValidateDateTimeLocal(elem, val) {
		t.Errorf("expected true, got false for value %q", val)
	}
}

// TestUUID tests UUID generation
func TestUUID(t *testing.T) {
	// Debug = true
	elem := new(Element)
	elem.Id = "pid"
	elem.Type = "uuid"
	elem.Generator = "uuid"
	val := "01925416-3e1a-77a5-9cf5-7452554913c8"
	if !ValidateUUID(elem, val) {
		t.Errorf("expected true, got false for value %q", val)
	}
	val = "01925413-abc0-75c8-aa75-bfc062cd2949"
	if !ValidateUUID(elem, val) {
		t.Errorf("expected true, got false for value %q", val)
	}
}

// TestROR tests the ROR validation func for an element.
func TestROR(t *testing.T) {
	elem := new(Element)
	elem.Id = "ror"
	elem.Type = "ror"
	//SetDebug(true);
	val := `https://ror.org/05dxps055`
	if !ValidateROR(elem, val) {
		t.Errorf("expected ValidateROR(elem, %q) to return true, return false", val)
	}
	val = `05dxps055`
	if !ValidateROR(elem, val) {
		t.Errorf("expected ValidateROR(elem, %q) to return true, return false", val)
	}
	//SetDebug(false);
}

// TestValidateModelTypes test model element types from YAML input
func TestValidateModelTypes(t *testing.T) {
	src := []byte(`id: people_model
description: CaltechPEOPLE
elements:
  - type: text
    id: clpid
    attributes:
      name: clpid
      required: "true"
    is_primary_id: true
    label: CL Person Id
  - type: text
    id: display_name
    attributes:
      name: display_name
    label: Display_name
  - type: text
    id: family_name
    attributes:
      name: family_name
      required: "true"
    label: Family_name
  - type: text
    id: given_name
    attributes:
      name: given_name
    label: Given_name
  - type: textarea
    id: bio
    attributes:
      name: bio
    label: Bio
  - type: textarea
    id: education
    attributes:
      name: education
    label: Education
  - type: email
    id: email
    attributes:
      name: email
    label: Email
  - type: text
    id: directory_person_type
    attributes:
      name: directory_person_type
    label: Directory_person_type
  - type: text
    id: directory_user_id
    attributes:
      name: directory_user_id
    label: Directory_user_id
  - type: text
    id: division
    attributes:
      name: division
    label: Division
  - type: orcid
    id: orcid
    attributes:
      name: orcid
    pattern: '[0-9]{4}-[0-9]{4}-[0-9]{4}-[0-9]{3}[0-9A-Z]'
    label: ORCID
  - type: text
    id: ror
    attributes:
      name: ror
    label: Ror
  - type: isni
    id: isni
    attributes:
      name: isni
    label: ISNI
  - type: text
    id: lcnaf
    attributes:
      name: lcnaf
    label: Lcnaf
  - type: text
    id: viaf
    attributes:
      name: viaf
    label: Viaf
  - type: text
    id: wikidata
    attributes:
      name: wikidata
    label: Wikidata
  - type: text
    id: snac
    attributes:
      name: snac
    label: Snac
  - type: text
    id: archivesspace_id
    attributes:
      name: archivesspace_id
    label: Archivesspace_id
  - type: text
    id: authors_id
    attributes:
      name: authors_id
    label: Authors_id
  - type: text
    id: thesis_id
    attributes:
      name: thesis_id
    label: Thesis_id
  - type: text
    id: advisors_id
    attributes:
      name: advisors_id
    label: Advisors_id
  - type: checkbox
    id: caltech
    attributes:
      name: caltech
    label: Caltech
  - type: checkbox
    id: jpl
    attributes:
      name: jpl
    label: JPL
  - type: checkbox
    id: include_in_feeds
    attributes:
      name: include_in_feeds
    label: Include in Feeds
  - type: datetime-local
    id: updated
    attributes:
      name: updated
    label: Updated
`)
	model := new(Model)
	if err := yaml.Unmarshal(src, &model); err != nil {
		t.Error(err)
	}
	if model == nil {
		t.Errorf("Expecte model to be non-nil")
		t.FailNow()
	}
	SetDefaultTypes(model)

	formData := map[string]interface{}{
		"clpid":                 "Doiel-M-S",
		"display_name":          "Doiel, Mark",
		"family_name":           "Doiel",
		"given_name":            "Mark",
		"bio":                   "Emeritus Music Professor",
		"education":             "COC\r\nCal Arts\r\n",
		"email":                 "mdoiel@music.example.edu",
		"directory_person_type": "",
		"directory_user_id":     "",
		"division":              "",
		"orcid":                 "",
		"ror":                   "",
		"isni":                  "",
		"lcnaf":                 "",
		"viaf":                  "",
		"wikidata":              "",
		"snac":                  "",
		"archivesspace_id":      "",
		"authors_id":            "",
		"thesis_id":             "",
		"advisors_id":           "",
		"include_in_feeds":      "true",
		"caltech":               "false",
		"jpl":                   "false",
		"updated":               "2024-10-08T10:11:12",
	}

	//SetDebug(true)
	if !model.ValidateMapInterface(formData) {
		t.Error("Model failed to validate")
		for _, elem := range model.Elements {
			if val, ok := formData[elem.Id]; ok {
				if validator, ok := model.validators[elem.Type]; ok {
					if !validator(elem, val.(string)) {
						t.Errorf("elem.Id %q, elem.Type %q failed to validate value %q", elem.Id, elem.Type, val)
					}
				}
			} else {
				t.Errorf("%q missing from form data %+v\n", elem.Id, formData)
			}
		}
	}
	SetDebug(false)
}
