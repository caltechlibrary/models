package models

import (
	"testing"
	"regexp"
)

func TestValidateElementText(t *testing.T) {
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
	pattern := `[0-9]{4}-[0-9]{4}-[0-9]{4}-[0-9]{3}[0-9A-Z]`
	re := regexp.MustCompilePOSIX(pattern)
	orcid := `0000-0003-0900-6903`
	if ! re.MatchString(orcid) {
		t.Errorf("expected true, got false for pattern %q and value %q", pattern, orcid)
	}
}

func TestDatetimeLocal(t *testing.T) {
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
