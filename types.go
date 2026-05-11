// types.go is part of the Go models package.
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
	"encoding/json"
	"log"
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	// 3rd Party packages
	"github.com/google/uuid"
	"github.com/nyaruka/phonenumbers"
)

const (
	OrcidPattern = `[0-9]{4}-[0-9]{4}-[0-9]{4}-[0-9]{3}[0-9A-Z]`
	RORPattern   = `^0[a-hj-km-np-tv-z|0-9]{6}[0-9]{2}$`
	ISNIPattern  = `[0-9]{4} [0-9]{4} [0-9]{4} [0-9X]{4}|[0-9]{4}-[0-9]{4}-[0-9]{4}-[0-9X]{4}`
	// ISBN patterns - use simple placeholder patterns; real validation done in ValidateISBN
	ISBN10Pattern = `[0-9\- ]{10,17}`
	ISBN13Pattern = `[0-9\- ]{13,26}`
	ISBNPattern   = `[0-9\- ]{10,26}`
	// ISSN pattern
	ISSNPattern = `[0-9\- ]{8,10}`
	// DOI pattern - placeholder, real validation in ValidateDOI function
	DOIPattern = `10\.[0-9/\-_.:]+`
	// ArXiv patterns
	ARXIVPattern = `arxiv:[0-9a-z\./\-]+`
	// EAN pattern (same as ISBN-13)
	EANPattern = ISBN13Pattern
	// PMID pattern
	PMIDPattern = `^[0-9]+$`
	// PMCID pattern
	PMCIDPattern = `^PMC[0-9]+$`
	// FundRef pattern
	FundRefPattern = `^10\.[0-9]{4,9}/[-._;()/:A-Z0-9]+$`
	// LCNAF pattern
	LCNAFPattern = `^[a-zA-Z0-9]+$`
	// VIAF pattern
	VIAFPattern = `^[0-9]+$`
	// SNAC pattern
	SNACPattern = `^[0-9]+$`
	// ARK pattern: optional ark:/ prefix, 5-digit NAAN, slash, name
	ARKPattern = `^ark:/?[0-9]{5}/[-.a-zA-Z0-9_~]+`
	// Wikidata QID pattern: Q followed by one or more digits
	WikidataPattern = `^Q[0-9]+$`
)

var (
	ReORCID    *regexp.Regexp
	ReROR      *regexp.Regexp
	ReISNI     *regexp.Regexp
	ReISBN     *regexp.Regexp
	ReISSN     *regexp.Regexp
	ReDOI      *regexp.Regexp
	ReARXIV    *regexp.Regexp
	ReEAN      *regexp.Regexp
	RePMID     *regexp.Regexp
	RePMCID    *regexp.Regexp
	ReFundRef  *regexp.Regexp
	ReLCNAF    *regexp.Regexp
	ReVIAF     *regexp.Regexp
	ReSNAC     *regexp.Regexp
	ReARK      *regexp.Regexp
	ReWikidata *regexp.Regexp
)

// GenerateROR setups up for an HTML ROR type input element
func GenerateROR() *Element {
	return &Element{
		Type:    "text",
		Pattern: RORPattern,
		Attributes: map[string]string{
			"placeholder": "enter a ROR",
		},
	}
}

// Validate ROR form element
func ValidateROR(elem *Element, formValue string) bool {
	if Debug {
		log.Printf("DEBUG validating elem.Id %q, elem.Type %q, value %q \n", elem.Id, elem.Type, formValue)
	}
	if formValue == "" {
		return true
	}
	if strings.HasPrefix(formValue, "https://ror.org/") {
		formValue = strings.TrimPrefix(formValue, "https://ror.org/")
	}
	if !ReROR.MatchString(formValue) {
		if Debug {
			log.Printf("DEBUG failed to validate pattern elem.Id %q, elem.Type %q, value %q \n", elem.Id, elem.Type, formValue)
		}
		return false
	}
	/*FIXME: Need to figure out how to validate (or unencode) the Crockford base 32 value */
	if Debug {
		log.Printf("DEBUG OK, elem.Id %q, elem.Type %q, value %q\n", elem.Id, elem.Type, formValue)
	}
	return true
}

// GenerateUUID setups up for an HTML uuid type input element
func GenerateUUID() *Element {
	return &Element{
		Type: "text",
		Attributes: map[string]string{
			"placeholder": "enter a uuid",
		},
	}
}

func ValidateUUID(elem *Element, formValue string) bool {
	if Debug {
		log.Printf("DEBUG validating elem.Id %q, elem.Type %q, value %q \n", elem.Id, elem.Type, formValue)
	}
	if formValue == "" {
		return true
	}
	if _, err := uuid.Parse(formValue); err != nil {
		if Debug {
			log.Printf("DEBUG failed to validate elem.Id %q, elem.Type %q, value %q \n", elem.Id, elem.Type, formValue)
		}
		return false
	}
	if Debug {
		log.Printf("DEBUG OK, elem.Id %q, elem.Type %q, value %q\n", elem.Id, elem.Type, formValue)
	}
	return true
}

// GenerateDate setups up for HTML date input element
func GenerateDate() *Element {
	return &Element{
		Type: "date",
		Attributes: map[string]string{
			"placeholder": "enter a date",
		},
	}
}

// ValidateDate makes sure the date string conforms to YYYY-MM-DD
func ValidateDate(elem *Element, formValue string) bool {
	if formValue == "" {
		return true
	}
	// FIXME: Need to check against min, max and step values
	if Debug {
		log.Printf("DEBUG validating elem.Id %q, elem.Type %q, value %q \n", elem.Id, elem.Type, formValue)
	}
	if _, err := time.Parse("2006-01-02", formValue); err != nil {
		if Debug {
			log.Printf("DEBUG failed to validate elem.Id %q, elem.Type %q, value %q: %s \n", elem.Id, elem.Type, formValue, err)
		}
		return false
	}
	if Debug {
		log.Printf("DEBUG OK, elem.Id %q, elem.Type %q, value %q\n", elem.Id, elem.Type, formValue)
	}
	return true
}

// GenerateDateTimeLocal sets up for HTML input type "datetime-local"
func GenerateDateTimeLocal() *Element {
	return &Element{
		Type: "datetime-local",
		Attributes: map[string]string{
			"placeholder": "enter a local timestamp",
		},
	}
}

// ValidateDateTimeLocal makes sure the datetime string conforms to
// Spec: https://html.spec.whatwg.org/multipage/common-microsyntaxes.html#valid-local-date-and-time-string
func ValidateDateTimeLocal(elem *Element, formValue string) bool {
	if formValue == "" {
		return true
	}
	// FIXME: Need to check against min, max and step values
	// See https://html.spec.whatwg.org/multipage/common-microsyntaxes.html#valid-local-date-and-time-string for validation steps
	if formValue == "" {
		return true
	}
	if Debug {
		log.Printf("DEBUG validating elem.Id %q, elem.Type %q, value %q \n", elem.Id, elem.Type, formValue)
	}
	// If we have timezone info so handle as RFC3339 validation
	if len(formValue) >= 20 {
		if _, err := time.Parse(time.RFC3339, formValue); err != nil {
			if Debug {
				log.Printf("DEBUG failed to validate elem.Id %q, elem.Type %q, value %q: %s \n", elem.Id, elem.Type, formValue, err)
			}
			return false
		}
		if Debug {
			log.Printf("DEBUG OK, elem.Id %q, elem.Type %q, value %q\n", elem.Id, elem.Type, formValue)
		}
		return true
	}
	// Parse date component first
	if _, err := time.Parse("2006-01-02", formValue[0:10]); err != nil {
		if Debug {
			log.Printf("DEBUG failed to validate elem.Id %q, elem.Type %q, value %q: %s \n", elem.Id, elem.Type, formValue, err)
		}
		return false
	}
	// String doesn't include the time so fail
	if len(formValue) <= 10 {
		if Debug {
			log.Printf("DEBUG failed to validate elem.Id %q, elem.Type %q, value %q: %s \n", elem.Id, elem.Type, formValue, "formValue length less than 10")
		}
		return false
	}
	if formValue[10:11] != "T" && formValue[10:11] != " " {
		if Debug {
			log.Printf("DEBUG failed to validate elem.Id %q, elem.Type %q, value %q: %s \n", elem.Id, elem.Type, formValue, "missing T or space between date and time")
		}
		return false
	}
	if _, err := time.Parse("15:04:05", formValue[11:19]); err != nil {
		if Debug {
			log.Printf("DEBUG failed to validate elem.Id %q, elem.Type %q, value %q: %s \n", elem.Id, elem.Type, formValue, err)
		}
		return false
	}
	if Debug {
		log.Printf("DEBUG OK, elem.Id %q, elem.Type %q, value %q\n", elem.Id, elem.Type, formValue)
	}
	return true
}

// GenerateMonth sets up for HTML input type "month"
func GenerateMonth() *Element {
	return &Element{
		Type: "month",
		Attributes: map[string]string{
			"placeholder": "Enter year dash month, example 2006-01",
		},
	}
}

// ValidateMonth parses the string for a year and month value, i.e. YYYY-MM style date string
func ValidateMonth(elem *Element, formValue string) bool {
	// FIXME: Need to check against min, max and step values
	if _, err := time.Parse("2006-02", formValue); err != nil {
		return false
	}
	return true
}

// GenerateColor sets up for HTML input type "color"
func GenerateColor() *Element {
	return &Element{
		Type: "color",
		Attributes: map[string]string{
			"value":       "#000000",
			"placeholder": "enter a color in hexidecimal format, e.g. green is #00FF00",
		},
	}
}

// ValidateColor checks to see if the value is expressed using Hexidecimal notation
func ValidateColor(elem *Element, formValue string) bool {
	// color should return a hexidecimal value
	_, err := strconv.ParseUint(formValue, 16, 64)
	if err != nil {
		// formValue is not a valid
		return false
	}
	return true
}

// GenerateEmail sets up for HTML input type "email"
func GenerateEmail() *Element {
	return &Element{
		Type: "email",
		Attributes: map[string]string{
			"placecholder": "E.g. jane.doe@example.org",
		},
	}
}

// ValidateEmail parses email address to confirm it is valid
func ValidateEmail(elem *Element, formValue string) bool {
	if _, err := mail.ParseAddress(formValue); err != nil {
		return false
	}
	return true
}

// GenerateText generates an Element setup to hold an HTML text input elements
func GenerateText() *Element {
	return &Element{
		Type: "text",
	}
}

// ValidateText will check to see if pattern is set, if so it will
// evaluate the formValue against the RegExp given in Pattern.
func ValidateText(elem *Element, formValue string) bool {
	if elem.Pattern == "" {
		return true
	}
	/*
		re, err := regexp.CompilePOSIX(elem.Pattern)
		if err != nil {
			return false
		}
		return re.MatchString(formValue)
	*/
	re := regexp.MustCompilePOSIX(elem.Pattern)
	return re.MatchString(formValue)
}

func jsonDecodeNumber(value string) (float64, error) {
	var number float64
	dec := json.NewDecoder(strings.NewReader(value))
	if err := dec.Decode(&number); err != nil {
		return 0, err
	}
	return number, nil
}

// GenerateNumber sets up for an HTML input type "number"
func GenerateNumber() *Element {
	return &Element{
		Type: "number",
		Attributes: map[string]string{
			"value": "0",
		},
	}
}

// ValidateNumber implements a number validation using the json package.
func ValidateNumber(elem *Element, formValue string) bool {
	if _, err := jsonDecodeNumber(formValue); err != nil {
		return false
	}
	return true
}

// GenerateRange sets up for an HTML input "range" (defauting is min 0 to max 100, step 1)
func GenerateRange() *Element {
	return &Element{
		Type: "range",
		Attributes: map[string]string{
			"value": "0",
			"min":   "0",
			"max":   "100",
			"step":  "1",
		},
	}
}

// ValidateRange retrieves the form's value as a float64 then checks if it is in range.  Min and max must befined in
// the attributes of the element since they are required to make the comparison. NOTE: ValidateRange isn't currently checking
// the step value as I don't know if the value of the input element is supposed to be an integer for real number.
func ValidateRange(elem *Element, formValue string) bool {
	var (
		minNumber   float64
		maxNumber   float64
		numberValue float64
		err         error
	)
	// First make sure elem has a minimum and maximum defined in it's attributes
	if val, ok := elem.Attributes["min"]; ok {
		minNumber, err = jsonDecodeNumber(val)
		if err != nil {
			return false
		}
	} else {
		return false
	}
	if val, ok := elem.Attributes["max"]; ok {
		maxNumber, err = jsonDecodeNumber(val)
		if err != nil {
			return false
		}
	} else {
		return false
	}

	numberValue, err = jsonDecodeNumber(formValue)
	if err != nil {
		return false
	}
	if numberValue >= minNumber && numberValue <= maxNumber {
		return true
	}
	return false
}

// GenerateTel sets up for an HTML input type "tel" (i.e. telephone number).
func GenerateTel() *Element {
	return &Element{
		Type: "tel",
		Attributes: map[string]string{
			"placeholder": "e.g. phone like 123-456-7890",
			"pattern":     "[0-9]{3}-[0-9]{3}-[0-9]{4}",
		},
	}
}

// ValidateTel validates formValue conforms to a phone number.
func ValidateTel(elem *Element, formValue string) bool {
	// NOTE: I am defaulting to US numbers because Caltech Library is in the US
	if _, err := phonenumbers.Parse(formValue, "US"); err != nil {
		return false
	}
	return true
}

// GenerateTime sets up for an HTML input type "time"
func GenerateTime() *Element {
	return &Element{
		Type: "time",
		Attributes: map[string]string{
			"placeholder": "E.g. 13:44 would be 1:44pm",
			"pattern":     "[0-2][0-9]:[0-5][0-9]|[0-2][0-9]:[0-5][0-9]:[0-5][0-9]",
		},
	}
}

// ValidateTime validates the formValue is a time format
func ValidateTime(elem *Element, formValue string) bool {
	//FIXME: Need to check against min and max values
	if _, err := time.Parse("15:04:05", formValue); err != nil {
		return false
	}
	return true
}

// GenerateURL sets up for an HTML input type "url"
func GenerateURL() *Element {
	return &Element{
		Type: "url",
		Attributes: map[string]string{
			"placeholder": "https://example.edu",
			"pattern":     "https://.*",
		},
	}
}

// ValidateURL validates a formValue is a URL
func ValidateURL(elem *Element, formValue string) bool {
	if _, err := url.Parse(formValue); err != nil {
		return false
	}
	return true
}

// GenerateWeek generates an Element setup to hold an HTML week input
func GenerateWeek() *Element {
	return &Element{
		Type: "week",
		Attributes: map[string]string{
			"placeholder": "Input as YYYY-WW where WW is week nuber, e.g. 2024-51",
			"pattern":     "[0-9]{4}-[0-5][0-9]",
		},
	}
}

// ValidateWeek attempts to validate a week number with year, string is WW-YYYY formatted
// NOTE: this is a crude validation since some years have 52 weeks other 53 depending on how
// the days of the week line up against the year.
//
// Also noted is this input element isn't widely support by browser so I might drop in the future.
func ValidateWeek(elem *Element, formValue string) bool {
	if strings.Index(formValue, "-") == 2 {
		weekNum, err := strconv.Atoi(formValue[0:2])
		if err != nil {
			return false
		}
		_, err = strconv.Atoi(formValue[3:])
		if err != nil {
			return false
		}
		if weekNum > 0 && weekNum <= 53 {
			return true
		}
	}
	return false
}

// GenerateCheckbox sets up for an HTML input type "checkbox"
func GenerateCheckbox() *Element {
	return &Element{
		Type:       "checkbox",
		Attributes: map[string]string{},
	}
}

// ValidateCheckbox checks is the form value was provided, returns false if empty string recieved for value.
func ValidateCheckbox(elem *Element, formValue string) bool {
	// Checkbox return their string value if checked.
	return strings.TrimSpace(formValue) != ""
}

// GenerateImage sets up for an HTML input type "image"
func GenerateImage() *Element {
	return &Element{
		Type:       "image",
		Attributes: map[string]string{},
	}
}

// ValidateImage, if value is empty string this returns true.
// NOTE: this func maybe depreciated as this is not a common form element
func ValidateImage(elem *Element, formValue string) bool {
	// The element value should be none, see https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/image#technical_summary
	return formValue == ""
}

// GeneratePassword sets up for an HTML input type "password"
func GeneratePassword() *Element {
	return &Element{
		Type:       "password",
		Attributes: map[string]string{},
	}
}

// ValidatePassword makes sure an password input element holds a single string
func ValidatePassword(elem *Element, formValue string) bool {
	// Passwords must be a single line of text, see https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/password
	if strings.Index(formValue, "\r") > -1 || strings.Index(formValue, "\n") > -1 {
		return false
	}
	return ValidateText(elem, formValue)
}

// GenerateRadio sets up for an HTML input type "radio"
func GenerateRadio() *Element {
	return &Element{
		Type:       "radio",
		Attributes: map[string]string{},
	}
}

func ValidateRadio(elem *Element, formValue string) bool {
	// Checkbox return their string value if checked.
	return strings.TrimSpace(formValue) != ""
}

// GenerateButton sets up for an HTML input type "button"
func GenerateButton() *Element {
	return &Element{
		Type:       "button",
		Attributes: map[string]string{},
	}
}

func ValidateButton(elem *Element, formValue string) bool {
	return true
}

// GenerateReset sets up for an HTML input type "reset"
func GenerateReset() *Element {
	return &Element{
		Type: "reset",
		Attributes: map[string]string{
			"value": "reset",
		},
	}
}

func ValidateReset(elem *Element, formValue string) bool {
	return true
}

// GenerateSubmit sets up for an HTML input type "submit"
func GenerateSubmit() *Element {
	return &Element{
		Type: "submit",
		Attributes: map[string]string{
			"value": "submit",
		},
	}
}

func ValidateSubmit(elem *Element, formValue string) bool {
	return true
}

// GenerateSearch sets up HTML input type "search"
func GenerateSearch() *Element {
	return &Element{
		Type:       "search",
		Attributes: map[string]string{},
	}
}

func ValidateSearch(elem *Element, formValue string) bool {
	return ValidateText(elem, formValue)
}

// GenerateTextarea sets up for HTML textarea input
func GenerateTextarea() *Element {
	return &Element{
		Type:       "textarea",
		Attributes: map[string]string{},
	}
}

func ValidateTextarea(elem *Element, formValue string) bool {
	return ValidateText(elem, formValue)
}

// GenerateISNI sets up for an HTML input type "text" with a pattern for INSI input
func GenerateISNI() *Element {
	return &Element{
		Type: "text",
		Attributes: map[string]string{
			"pattern":     ISNIPattern,
			"placeholder": "e.g. 1111 2222 3333 444X",
		},
	}
}

func ValidateISNI(elem *Element, formValue string) bool {
	if Debug {
		log.Printf("DEBUG validating isni elem.Id %q, elem.Type %q, value %q\n", elem.Id, elem.Type, formValue)
	}
	if formValue == "" {
		return true
	}
	formValue = strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(formValue, "-", ""), " ", ""))
	if len(formValue) != 16 {
		if Debug {
			log.Printf("DEBUG validating isni elem.Id %q, elem.Type %q, value %q: %s\n", elem.Id, elem.Type, formValue, "length != 16")
		}
		return false
	}
	r := 0
	ck := 0
	for pos := 0; pos < 15; pos++ {
		x, err := strconv.Atoi(formValue[pos : pos+1])
		if err != nil {
			if Debug {
				log.Printf("DEBUG validating isni elem.Id %q, elem.Type %q, value %q: %s\n", elem.Id, elem.Type, formValue, err)
			}
			return false
		}
		r = (r + x) * 2
	}
	lastDigit, err := strconv.Atoi(formValue[len(formValue)-1:])
	if err != nil {
		if Debug {
			log.Printf("DEBUG validating isni elem.Id %q, elem.Type %q, value %q: %s\n", elem.Id, elem.Type, formValue, err)
		}
		return false
	}
	ck = (12 - r%11) % 11
	if Debug {
		log.Printf("DEBUG validating isni elem.Id %q, elem.Type %q, value %q\n, result: %t", elem.Id, elem.Type, formValue, (ck == lastDigit))
	}
	return ck == lastDigit
}

// GenerateORCID sets up for an HTML input type text using a pattern for ORCID
func GenerateORCID() *Element {
	return &Element{
		Type: "text",
		Attributes: map[string]string{
			"pattern": OrcidPattern,
		},
	}
}

func ValidateORCID(elem *Element, formValue string) bool {
	if Debug {
		log.Printf("DEBUG validating elem.Id %q, elem.Type %q, value %q \n", elem.Id, elem.Type, formValue)
	}
	if formValue == "" {
		return true
	}
	/* Based on https://idutils.readthedocs.io/en/latest/_modules/idutils.html#is_orcid */
	if strings.HasPrefix(formValue, "https://orcid.org/") {
		formValue = strings.TrimPrefix(formValue, "https://orcid.org/")
	}
	formValue = strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(formValue, "-", ""), " ", ""))
	if ValidateISNI(elem, formValue) {
		// Remove tailing check digit, then convert to integer
		//log.Printf("DEBUG formValue: %q -> formValue[0:len(formValue) -1]: %q", formValue, formValue[0:len(formValue)-1])
		val, err := strconv.Atoi(formValue[0 : len(formValue)-1])
		if err != nil {
			if Debug {
				log.Printf("DEBUG failed to validate elem.Id %q, elem.Type %q, value %q: %s \n", elem.Id, elem.Type, formValue, err)
			}
			return false
		}
		if Debug {
			log.Printf("DEBUG testing orcid ranges: %t, elem.Id %q, elem.Type %q, value %q", (val >= 15000000) && (val <= 35000000), elem.Id, elem.Type, formValue)
		}
		return (val >= 15000000) && (val <= 35000000)
	}
	if Debug {
		log.Printf("DEBUG failed to validate elem.Id %q, elem.Type %q, value %q: %s \n", elem.Id, elem.Type, formValue, "does not confirm to isni")
	}
	return false
}

// GenerateISBN sets up for an HTML input type text using a pattern for ISBN (10 or 13)
func GenerateISBN() *Element {
	return &Element{
		Type: "text",
		Attributes: map[string]string{
			"pattern":     ISBNPattern,
			"placeholder": "e.g. 978-3-16-148410-0 or 0-306-40615-2",
		},
	}
}

// ValidateISBN validates ISBN-10 or ISBN-13 checksum
func ValidateISBN(elem *Element, formValue string) bool {
	if formValue == "" {
		return true
	}
	cleanISBN := strings.ReplaceAll(strings.ReplaceAll(strings.ToUpper(formValue), "-", ""), " ", "")
	if len(cleanISBN) == 10 {
		return validateISBN10(cleanISBN)
	}
	if len(cleanISBN) == 13 {
		return validateISBN13(cleanISBN)
	}
	return false
}

// validateISBN10 validates ISBN-10 checksum (internal helper)
func validateISBN10(isbn string) bool {
	if len(isbn) != 10 {
		return false
	}
	checksum := 0
	for i := 0; i < 9; i++ {
		digit, err := strconv.Atoi(isbn[i : i+1])
		if err != nil {
			return false
		}
		checksum += digit * (10 - i)
	}
	checkDigit := isbn[9]
	if checkDigit == 'X' || checkDigit == 'x' {
		checksum += 10
	} else {
		digit, err := strconv.Atoi(string(checkDigit))
		if err != nil {
			return false
		}
		checksum += digit
	}
	return checksum%11 == 0
}

// validateISBN13 validates ISBN-13 checksum (internal helper)
func validateISBN13(isbn string) bool {
	if len(isbn) != 13 {
		return false
	}
	checksum := 0
	for i := 0; i < 12; i++ {
		digit, err := strconv.Atoi(isbn[i : i+1])
		if err != nil {
			return false
		}
		// Weight is 1 for even positions, 3 for odd positions (0-indexed)
		weight := 1
		if i%2 == 1 {
			weight = 3
		}
		checksum += digit * weight
	}
	checkDigit, err := strconv.Atoi(isbn[12:13])
	if err != nil {
		return false
	}
	return (checksum+checkDigit)%10 == 0
}

// GenerateISSN sets up for an HTML input type text using a pattern for ISSN
func GenerateISSN() *Element {
	return &Element{
		Type: "text",
		Attributes: map[string]string{
			"pattern":     ISSNPattern,
			"placeholder": "e.g. ISSN 1234-5678",
		},
	}
}

// ValidateISSN validates ISSN checksum
func ValidateISSN(elem *Element, formValue string) bool {
	if formValue == "" {
		return true
	}
	// Strip prefix and normalize
	bareISSN := strings.ToUpper(strings.TrimSpace(formValue))
	bareISSN = strings.TrimPrefix(bareISSN, "ISSN")
	bareISSN = strings.ReplaceAll(bareISSN, "-", "")
	bareISSN = strings.ReplaceAll(bareISSN, " ", "")
	if len(bareISSN) != 8 {
		return false
	}
	return validateISSNChecksum(bareISSN)
}

// validateISSNChecksum validates ISSN checksum using the algorithm (internal helper)
func validateISSNChecksum(issn string) bool {
	if len(issn) != 8 {
		return false
	}
	digits := issn[0:7]
	checkDigit := strings.ToUpper(string(issn[7]))
	
	checksum := 0
	for i := 0; i < 7; i++ {
		digit, err := strconv.Atoi(digits[i : i+1])
		if err != nil {
			return false
		}
		checksum += digit * (8 - i)
	}
	remainder := checksum % 11
	var expectedCheckDigit string
	if remainder == 0 {
		expectedCheckDigit = "0"
	} else if remainder == 1 {
		expectedCheckDigit = "X"
	} else {
		expectedCheckDigit = strconv.Itoa(11 - remainder)
	}
	return checkDigit == expectedCheckDigit
}

// GenerateDOI sets up for an HTML input type text using a pattern for DOI
func GenerateDOI() *Element {
	return &Element{
		Type: "text",
		Attributes: map[string]string{
			"pattern":     DOIPattern,
			"placeholder": "e.g. 10.22002/bv2pv-2b295",
		},
	}
}

// ValidateDOI validates DOI format
func ValidateDOI(elem *Element, formValue string) bool {
	if formValue == "" {
		return true
	}
	// Strip URL prefix if present
	val := strings.TrimSpace(formValue)
	if strings.HasPrefix(val, "https://doi.org/") {
		val = strings.TrimPrefix(val, "https://doi.org/")
	} else if strings.HasPrefix(val, "http://doi.org/") {
		val = strings.TrimPrefix(val, "http://doi.org/")
	} else if strings.HasPrefix(val, "doi:") {
		val = strings.TrimPrefix(val, "doi:")
	}
	// DOI format: 10.NNNN/suffix where NNNN is 4+ digits
	if !strings.HasPrefix(val, "10.") {
		return false
	}
	// Remove prefix
	val = strings.TrimPrefix(val, "10.")
	if val == "" {
		return false
	}
	// Find the slash
	slashIndex := strings.Index(val, "/")
	if slashIndex < 4 { // At least 4 digits before slash
		return false
	}
	// Check digits before slash
	prefix := val[:slashIndex]
	if !isAllDigits(prefix) || len(prefix) < 4 {
		return false
	}
	// Check suffix exists
	if slashIndex >= len(val)-1 {
		return false
	}
	return true
}

// GenerateArXiv sets up for an HTML input type text using a pattern for ArXiv
func GenerateArXiv() *Element {
	return &Element{
		Type: "text",
		Attributes: map[string]string{
			"pattern":     ARXIVPattern,
			"placeholder": "e.g. arxiv:2412.03631 or arxiv:hep-th/9901001",
		},
	}
}

// ValidateArXiv validates ArXiv identifier format
func ValidateArXiv(elem *Element, formValue string) bool {
	if formValue == "" {
		return true
	}
	val := strings.ToLower(strings.TrimSpace(formValue))
	// Strip prefix if present
	val = strings.TrimPrefix(val, "arxiv:")
	if val == "" {
		return false
	}
	// Check new format: YYYY.NNNNN or YYYY.NNNNNvN
	if strings.Contains(val, ".") {
		parts := strings.Split(val, ".")
		if len(parts) != 2 {
			return false
		}
		// Year part: 4 digits
		if len(parts[0]) != 4 || !isAllDigits(parts[0]) {
			return false
		}
		// ID part: at least 4 digits, optional vN suffix
		idPart := parts[1]
		if len(idPart) < 4 {
			return false
		}
		// Check for version suffix
		if strings.HasPrefix(idPart, "v") {
			idPart = strings.TrimPrefix(idPart, "v")
			if len(idPart) < 1 || !isAllDigits(idPart) {
				return false
			}
		} else {
			if !isAllDigits(idPart) {
				return false
			}
		}
		return true
	}
	// Check old format: archive/NNNNNNN or archive/NNNNNNNvN
	if strings.Contains(val, "/") {
		parts := strings.Split(val, "/")
		if len(parts) != 2 {
			return false
		}
		// Archive part: alphanumeric with hyphens
		if len(parts[0]) < 1 {
			return false
		}
		// ID part: exactly 7 digits (or more with version)
		idPart := parts[1]
		if len(idPart) < 7 {
			return false
		}
		// Check for version suffix
		if strings.HasPrefix(idPart, "v") {
			idPart = strings.TrimPrefix(idPart, "v")
		}
		if len(idPart) != 7 || !isAllDigits(idPart) {
			return false
		}
		return true
	}
	return false
}

// isAllDigits checks if a string contains only digits
func isAllDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return s != ""
}

// GenerateEAN sets up for an HTML input type text using a pattern for EAN (same as ISBN-13)
func GenerateEAN() *Element {
	return &Element{
		Type: "text",
		Attributes: map[string]string{
			"pattern":     EANPattern,
			"placeholder": "e.g. 9780306406157",
		},
	}
}

// ValidateEAN validates EAN-13 (same as ISBN-13 validation)
func ValidateEAN(elem *Element, formValue string) bool {
	if formValue == "" {
		return true
	}
	cleanEAN := strings.ReplaceAll(strings.ReplaceAll(formValue, "-", ""), " ", "")
	if len(cleanEAN) == 13 {
		return validateISBN13(cleanEAN)
	}
	return false
}

// GeneratePMID sets up for an HTML input type text using a pattern for PubMed ID
func GeneratePMID() *Element {
	return &Element{
		Type: "text",
		Attributes: map[string]string{
			"pattern":     PMIDPattern,
			"placeholder": "e.g. 1234567",
		},
	}
}

// ValidatePMID validates PubMed ID format
func ValidatePMID(elem *Element, formValue string) bool {
	if formValue == "" {
		return true
	}
	val := strings.TrimSpace(formValue)
	// Remove common prefixes
	val = strings.TrimPrefix(val, "PMID:")
	val = strings.TrimPrefix(val, "pmid:")
	val = strings.ReplaceAll(val, " ", "")
	return RePMID.MatchString(val)
}

// GeneratePMCID sets up for an HTML input type text using a pattern for PubMed Central ID
func GeneratePMCID() *Element {
	return &Element{
		Type: "text",
		Attributes: map[string]string{
			"pattern":     PMCIDPattern,
			"placeholder": "e.g. PMC1234567",
		},
	}
}

// ValidatePMCID validates PubMed Central ID format
func ValidatePMCID(elem *Element, formValue string) bool {
	if formValue == "" {
		return true
	}
	val := strings.ToUpper(strings.TrimSpace(formValue))
	if !strings.HasPrefix(val, "PMC") {
		val = "PMC" + val
	}
	return RePMCID.MatchString(val)
}

// GenerateFundRef sets up for an HTML input type text using a pattern for FundRef DOI
func GenerateFundRef() *Element {
	return &Element{
		Type: "text",
		Attributes: map[string]string{
			"pattern":     FundRefPattern,
			"placeholder": "e.g. 10.13039/100000001",
		},
	}
}

// ValidateFundRef validates FundRef DOI format
func ValidateFundRef(elem *Element, formValue string) bool {
	if formValue == "" {
		return true
	}
	val := strings.ToLower(strings.TrimSpace(formValue))
	return ReFundRef.MatchString(val)
}

// GenerateLCNAF sets up for an HTML input type text using a pattern for LCNAF
func GenerateLCNAF() *Element {
	return &Element{
		Type: "text",
		Attributes: map[string]string{
			"pattern":     LCNAFPattern,
			"placeholder": "LCNAF identifier",
		},
	}
}

// ValidateLCNAF validates LCNAF format
func ValidateLCNAF(elem *Element, formValue string) bool {
	if formValue == "" {
		return true
	}
	val := strings.TrimSpace(formValue)
	return ReLCNAF.MatchString(val)
}

// GenerateVIAF sets up for an HTML input type text using a pattern for VIAF
func GenerateVIAF() *Element {
	return &Element{
		Type: "text",
		Attributes: map[string]string{
			"pattern":     VIAFPattern,
			"placeholder": "VIAF identifier",
		},
	}
}

// ValidateVIAF validates VIAF format
func ValidateVIAF(elem *Element, formValue string) bool {
	if formValue == "" {
		return true
	}
	val := strings.TrimSpace(formValue)
	return ReVIAF.MatchString(val)
}

// GenerateSNAC sets up for an HTML input type text using a pattern for SNAC
func GenerateSNAC() *Element {
	return &Element{
		Type: "text",
		Attributes: map[string]string{
			"pattern":     SNACPattern,
			"placeholder": "SNAC identifier",
		},
	}
}

// ValidateSNAC validates SNAC format
func ValidateSNAC(elem *Element, formValue string) bool {
	if formValue == "" {
		return true
	}
	val := strings.TrimSpace(formValue)
	return ReSNAC.MatchString(val)
}

// GenerateList sets up a list element (a repeatable sequence of sub-elements).
// At the web form level this is typically rendered as a JSON-encoded textarea.
func GenerateList() *Element {
	return &Element{
		Type: "list",
	}
}

// ValidateList is a no-op at the scalar level; list contents are validated
// recursively by the model's validateListErrors method.
func ValidateList(elem *Element, formValue string) bool {
	return true
}

// GenerateObject sets up an object element (a named set of sub-elements).
// At the web form level this is typically rendered as a JSON-encoded textarea.
func GenerateObject() *Element {
	return &Element{
		Type: "object",
	}
}

// ValidateObject is a no-op at the scalar level; object contents are validated
// recursively by the model's validateObjectErrors method.
func ValidateObject(elem *Element, formValue string) bool {
	return true
}

// GenerateARK sets up for an HTML input type text using a pattern for ARK identifiers
func GenerateARK() *Element {
	return &Element{
		Type: "text",
		Attributes: map[string]string{
			"pattern":     ARKPattern,
			"placeholder": "e.g. ark:/99999/fk4cz3dh0",
		},
	}
}

// ValidateARK validates an ARK (Archival Resource Key) identifier.
// Accepts bare NAANs (ark:/NAAN/name) or fully qualified ARKs with qualifiers.
// The NAAN (Name Assigning Authority Number) must be exactly 5 digits.
func ValidateARK(elem *Element, formValue string) bool {
	if formValue == "" {
		return true
	}
	val := strings.TrimSpace(formValue)
	// Normalize: ensure ark: prefix is present
	if !strings.HasPrefix(strings.ToLower(val), "ark:") {
		return false
	}
	return ReARK.MatchString(val)
}

// GenerateWikidata sets up for an HTML input type text for Wikidata QIDs
func GenerateWikidata() *Element {
	return &Element{
		Type: "text",
		Attributes: map[string]string{
			"pattern":     WikidataPattern,
			"placeholder": "e.g. Q42",
		},
	}
}

// ValidateWikidata validates a Wikidata QID (e.g. Q42, Q1234567).
// Strips the https://www.wikidata.org/entity/ prefix if present.
func ValidateWikidata(elem *Element, formValue string) bool {
	if formValue == "" {
		return true
	}
	val := strings.TrimSpace(formValue)
	if strings.HasPrefix(val, "https://www.wikidata.org/entity/") {
		val = strings.TrimPrefix(val, "https://www.wikidata.org/entity/")
	} else if strings.HasPrefix(val, "http://www.wikidata.org/entity/") {
		val = strings.TrimPrefix(val, "http://www.wikidata.org/entity/")
	}
	return ReWikidata.MatchString(val)
}

func SetDefaultTypes(model *Model) {
	model.Define("date", GenerateDate, ValidateDate)
	model.Define("datetime-local", GenerateDateTimeLocal, ValidateDateTimeLocal)
	model.Define("month", GenerateMonth, ValidateMonth)
	model.Define("color", GenerateColor, ValidateColor)
	model.Define("email", GenerateEmail, ValidateEmail)
	model.Define("text", GenerateText, ValidateText)
	model.Define("number", GenerateNumber, ValidateNumber)
	model.Define("range", GenerateRange, ValidateRange)
	model.Define("tel", GenerateTel, ValidateTel)
	model.Define("time", GenerateTime, ValidateTime)
	model.Define("url", GenerateURL, ValidateURL)
	model.Define("checkbox", GenerateCheckbox, ValidateCheckbox)
	model.Define("password", GeneratePassword, ValidatePassword)
	model.Define("radio", GenerateRadio, ValidateRadio)
	model.Define("textarea", GenerateTextarea, ValidateTextarea)
	model.Define("orcid", GenerateORCID, ValidateORCID)
	model.Define("isni", GenerateISNI, ValidateISNI)
	model.Define("uuid", GenerateUUID, ValidateUUID)
	model.Define("ror", GenerateROR, ValidateROR)
	// Identifier types
	model.Define("isbn", GenerateISBN, ValidateISBN)
	model.Define("issn", GenerateISSN, ValidateISSN)
	model.Define("doi", GenerateDOI, ValidateDOI)
	model.Define("arxiv", GenerateArXiv, ValidateArXiv)
	model.Define("ean", GenerateEAN, ValidateEAN)
	model.Define("pmid", GeneratePMID, ValidatePMID)
	model.Define("pmcid", GeneratePMCID, ValidatePMCID)
	model.Define("fundref", GenerateFundRef, ValidateFundRef)
	model.Define("lcnaf", GenerateLCNAF, ValidateLCNAF)
	model.Define("viaf", GenerateVIAF, ValidateVIAF)
	model.Define("snac", GenerateSNAC, ValidateSNAC)
	model.Define("ark", GenerateARK, ValidateARK)
	model.Define("wikidata", GenerateWikidata, ValidateWikidata)
	model.Define("list", GenerateList, ValidateList)
	model.Define("object", GenerateObject, ValidateObject)

	// NOTE: The following are not in the default but their usefulness
	// in the context of persisting data is not clear.
	//
	//model.Define("search", GenerateSearch, ValidateSearch)
	//model.Define("reset", GenerateReset, ValidateReset)
	//model.Define("submit", GenerateSubmit, ValidateSubmit)
	//model.Define("button", GenerateButton, ValidateButton)
	//model.Define("week", GenerateWeek, ValidateWeek)
	//model.Define("image", GenerateImage, ValidateImage)
}

func init() {
	ReORCID = regexp.MustCompilePOSIX(OrcidPattern)
	ReROR = regexp.MustCompilePOSIX(RORPattern)
	ReISNI = regexp.MustCompilePOSIX(ISNIPattern)
	ReISBN = regexp.MustCompilePOSIX(ISBNPattern)
	ReISSN = regexp.MustCompilePOSIX(ISSNPattern)
	ReDOI = regexp.MustCompilePOSIX(DOIPattern)
	ReARXIV = regexp.MustCompilePOSIX(ARXIVPattern)
	ReEAN = regexp.MustCompilePOSIX(EANPattern)
	RePMID = regexp.MustCompilePOSIX(PMIDPattern)
	RePMCID = regexp.MustCompilePOSIX(PMCIDPattern)
	ReFundRef = regexp.MustCompilePOSIX(FundRefPattern)
	ReLCNAF = regexp.MustCompilePOSIX(LCNAFPattern)
	ReVIAF = regexp.MustCompilePOSIX(VIAFPattern)
	ReSNAC = regexp.MustCompilePOSIX(SNACPattern)
	ReARK = regexp.MustCompile(ARKPattern)
	ReWikidata = regexp.MustCompilePOSIX(WikidataPattern)
}
