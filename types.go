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
	RORPattern = `^0[a-hj-km-np-tv-z|0-9]{6}[0-9]{2}$`
	ISNIPattern = `[0-9]{4} [0-9]{4} [0-9]{4] [0-9X]{4}|[0-9]{4}-[0-9]{4}-[0-9]{4]-[0-9X]{4}`
)

var (
	ReORCID *regexp.Regexp
	ReROR *regexp.Regexp
	ReISNI *regexp.Regexp
)

// GenerateROR setups up for an HTML ROR type input element
func GenerateROR() *Element {
	return &Element{
		Type: "text",
		Pattern: RORPattern,
		Attributes:map[string]string{
			"placeholder": "enter a uuid",
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
	if ! ReROR.MatchString(formValue) {
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
		Attributes:map[string]string{
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
			"pattern":   ISNIPattern,
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
	ck = (12 - r % 11) % 11
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
		if Debug  {
			log.Printf("DEBUG testing orcid ranges: %t, elem.Id %q, elem.Type %q, value %q", (val >= 15000000) && (val <= 35000000), elem.Id, elem.Type, formValue)
		}
		return (val >= 15000000) && (val <= 35000000)
	}
	if Debug {
		log.Printf("DEBUG failed to validate elem.Id %q, elem.Type %q, value %q: %s \n", elem.Id, elem.Type, formValue, "does not confirm to isni")
	}
	return false
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

	// NOTE: The following are not in the default but their usefulness
	// in the context of persisting data is not clear.
	//
	//model.Define("search", DefeaultSearch, ValidateSearch)
	//model.Define("reset", ValidateReset)
	//model.Define("submit", ValidateSubmit)
	//model.Define("button", ValidateButton)
	//model.Define("week", ValidateWeek)
	//model.Define("image", ValidateImage)
}

func init() {
	ReORCID = regexp.MustCompilePOSIX(OrcidPattern)
	ReROR = regexp.MustCompilePOSIX(RORPattern)
	ReISNI = regexp.MustCompilePOSIX(ISNIPattern)
}
