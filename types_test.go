// types_test.go is part of the Go models package.
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

// TestDOI tests DOI validation
func TestDOI(t *testing.T) {
	elem := &Element{Id: "doi", Type: "doi"}
	valid := []string{
		"10.1234/test",
		"10.22002/bv2pv-2b295",
		"https://doi.org/10.1038/nature12373",
		"doi:10.1000/xyz123",
	}
	for _, v := range valid {
		if !ValidateDOI(elem, v) {
			t.Errorf("expected ValidateDOI to accept %q", v)
		}
	}
	invalid := []string{
		"not-a-doi",
		"10./missing-registrant",
		"10.123/too-short-prefix",
		"11.1234/wrong-prefix",
	}
	for _, v := range invalid {
		if ValidateDOI(elem, v) {
			t.Errorf("expected ValidateDOI to reject %q", v)
		}
	}
}

// TestISBN tests ISBN-10 and ISBN-13 checksum validation
func TestISBN(t *testing.T) {
	elem := &Element{Id: "isbn", Type: "isbn"}
	valid := []string{
		"0-306-40615-2",   // ISBN-10
		"978-3-16-148410-0", // ISBN-13
		"0306406152",      // ISBN-10 no dashes
		"9783161484100",   // ISBN-13 no dashes
	}
	for _, v := range valid {
		if !ValidateISBN(elem, v) {
			t.Errorf("expected ValidateISBN to accept %q", v)
		}
	}
	invalid := []string{
		"0-306-40615-3",   // bad check digit
		"978-3-16-148410-1", // bad check digit
		"12345",
	}
	for _, v := range invalid {
		if ValidateISBN(elem, v) {
			t.Errorf("expected ValidateISBN to reject %q", v)
		}
	}
}

// TestISSN tests ISSN checksum validation
func TestISSN(t *testing.T) {
	elem := &Element{Id: "issn", Type: "issn"}
	valid := []string{
		"0317-8471",
		"1050-124X",
		"ISSN 0317-8471",
	}
	for _, v := range valid {
		if !ValidateISSN(elem, v) {
			t.Errorf("expected ValidateISSN to accept %q", v)
		}
	}
	invalid := []string{
		"0317-8472", // bad check digit
		"1234-5678", // bad check digit
		"notanissn",
	}
	for _, v := range invalid {
		if ValidateISSN(elem, v) {
			t.Errorf("expected ValidateISSN to reject %q", v)
		}
	}
}

// TestPMCID tests PubMed Central ID validation
func TestPMCID(t *testing.T) {
	elem := &Element{Id: "pmcid", Type: "pmcid"}
	valid := []string{"PMC1234567", "PMC9999999", "1234567"}
	for _, v := range valid {
		if !ValidatePMCID(elem, v) {
			t.Errorf("expected ValidatePMCID to accept %q", v)
		}
	}
	invalid := []string{"notapmcid", "PM1234567", "PMC"}
	for _, v := range invalid {
		if ValidatePMCID(elem, v) {
			t.Errorf("expected ValidatePMCID to reject %q", v)
		}
	}
}

// TestARK tests ARK identifier validation
func TestARK(t *testing.T) {
	elem := &Element{Id: "ark", Type: "ark"}
	valid := []string{
		"ark:/99999/fk4cz3dh0",
		"ark:13960/t6m042c11",
		"ark:/12025/654xz321",
	}
	for _, v := range valid {
		if !ValidateARK(elem, v) {
			t.Errorf("expected ValidateARK to accept %q", v)
		}
	}
	invalid := []string{
		"not-an-ark",
		"ark:/999/too-short-naan",
		"http://example.org/ark",
		"ark:",
	}
	for _, v := range invalid {
		if ValidateARK(elem, v) {
			t.Errorf("expected ValidateARK to reject %q", v)
		}
	}
}

// TestWikidata tests Wikidata QID validation
func TestWikidata(t *testing.T) {
	elem := &Element{Id: "wikidata", Type: "wikidata"}
	valid := []string{
		"Q42",
		"Q1234567",
		"https://www.wikidata.org/entity/Q42",
	}
	for _, v := range valid {
		if !ValidateWikidata(elem, v) {
			t.Errorf("expected ValidateWikidata to accept %q", v)
		}
	}
	invalid := []string{
		"q42",            // lowercase Q not accepted
		"42",             // missing Q prefix
		"QQ42",           // double Q
		"not-a-qid",
	}
	for _, v := range invalid {
		if ValidateWikidata(elem, v) {
			t.Errorf("expected ValidateWikidata to reject %q", v)
		}
	}
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
