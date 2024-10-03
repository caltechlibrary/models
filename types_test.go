package models

import (
	"testing"
	"regexp"
)

func TestValidateElementText(t *testing.T) {
	// Debug = true
	elem := new(Element)
	elem.Id = "orcid"
	elem.Type = "text"
	elem.Pattern = OrcidPattern
    //val := `0000-1111-2222-333X`
	val := `0000-0003-0900-6903`
	if ! ValidateText(elem, val) {
		t.Errorf("Expected true, got false for elem %+v -> val %q", elem, val)
	}
}

func TestORCIDRegExp(t *testing.T) {
	// Debug = true
	pattern := `[0-9]{4}-[0-9]{4}-[0-9]{4}-[0-9]{3}[0-9A-Z]`
	re := regexp.MustCompilePOSIX(pattern)
	orcid := `0000-0003-0900-6903`
	if ! re.MatchString(orcid) {
		t.Errorf("expected true, got false for pattern %q and value %q", pattern, orcid)
	}
	elem := new (Element)
	elem.Id = "orcid"
	elem.Type = "orcid"
	elem.Generator = ""
	//SetDebug(true)
	if ! ValidateORCID(elem, orcid) {
		t.Errorf("expected orcid to validate true, got false")
	}
	// This isn't an valid ORCID throug it could be an INSI
	orcid = `2345-5432-1234-4326`
	if ValidateORCID(elem, orcid) {
		t.Errorf("expected orcid to validate false, got true")
	}
	//SetDebug(false)
}

func TestDatetimeLocal(t *testing.T) {
	// Debug = true
	elem := new(Element)
	elem.Id = "created"
	elem.Type = "datetime-local"
	elem.Generator = "created_timestamp"


	val := "2024-10-03T12:51:01"
	if ! ValidateDateTimeLocal(elem, val) {
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
	if ! ValidateDateTimeLocal(elem, val) {
		t.Errorf("expected true, got false for value %q", val)
	}
}

func TestUUID(t *testing.T) {
	// Debug = true
	elem := new(Element)
	elem.Id = "pid"
	elem.Type = "uuid"
	elem.Generator = "uuid"
	val := "01925416-3e1a-77a5-9cf5-7452554913c8"
	if ! ValidateUUID(elem, val) {
		t.Errorf("expected true, got false for value %q",val)
	}
  	val = "01925413-abc0-75c8-aa75-bfc062cd2949"
	if ! ValidateUUID(elem, val) {
		t.Errorf("expected true, got false for value %q",val)
	}
}
