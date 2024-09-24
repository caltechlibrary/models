package models

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/mail"
	"net/url"
	"strconv"
	"strings"
	"time"
	"regexp"

	// 3rd Party packages
	"github.com/nyaruka/phonenumbers"
)

func ModelToHTML(out io.Writer, model *Model) error {
	// FIXME: Handle title if it exists
	// Write opening form element
	if model.Id != "" {
		fmt.Fprintf(out, "<form id=%q", model.Id)
	} else {
		fmt.Fprintf(out, "<form")
	}
	for k, v := range model.Attributes {
		switch k {
		case "checked":
			fmt.Fprintf(out, " checked")
		case "required":
			fmt.Fprintf(out, " required")
		default:
			fmt.Fprintf(out, " %s=%q", k, v)
		}
	}
	cssBaseClass := strings.ReplaceAll(strings.ToLower(model.Id), " ", "_")
	fmt.Fprintf(out, ">\n")
	for _, elem := range model.Elements {
		ElementToHTML(out, cssBaseClass, elem)
	}
	if !model.HasElementType("submit") {
		cssName := fmt.Sprintf("%s-submit", cssBaseClass)
		fmt.Fprintf(out, `  <div class=%q><input class=%q type="submit" value="submit"> <input class=%q type="reset" value="cancel"></div>`,
			cssName, cssName, cssName)
	}

	// Write closing form element
	fmt.Fprintf(out, "\n</form>\n")
	return nil
}

func ElementToHTML(out io.Writer, cssBaseClass string, elem *Element) error {
	cssClass := fmt.Sprintf("%s-%s", cssBaseClass, strings.ToLower(elem.Id))
	fmt.Fprintf(out, "  <div class=%q>", cssClass)
	switch strings.ToLower(elem.Type) {
	case "textarea":
		if elem.Label != "" {
			if name, ok := elem.Attributes["name"]; ok {
				fmt.Fprintf(out, "<label class=%q set=%q>%s</label> <textarea class=%q", cssClass, name, elem.Label, cssClass)
			} else {
				fmt.Fprintf(out, "<label class=%q set=%q>%s</label> <textarea class=%q name=%q", cssClass, elem.Id, elem.Label, cssClass, elem.Id)
			}
		} else {
			fmt.Fprintf(out, "<textarea class=%q", cssClass)
		}
	case "button":
		fmt.Fprintf(out, "<button class=%q", cssClass)
	default:
		if elem.Label != "" {
			if name, ok := elem.Attributes["name"]; ok {
				fmt.Fprintf(out, "<label class=%q set=%q>%s</label> <input class=%q type=%q", cssClass, name, elem.Label, cssClass, elem.Type)
			} else {
				fmt.Fprintf(out, "<label class=%q set=%q>%s</label> <input class=%q name=%q type=%q", cssClass, elem.Id, elem.Label, cssClass, elem.Id, elem.Type)
			}
		} else {
			fmt.Fprintf(out, "<input class=%q type=%q", cssClass, elem.Type)
		}
	}
	if elem.Id != "" {
		fmt.Fprintf(out, " id=%q", elem.Id)
	}
	for k, v := range elem.Attributes {
		switch k {
		case "checked":
			fmt.Fprintf(out, " checked")
		case "required":
			fmt.Fprintf(out, " required")
		default:
			fmt.Fprintf(out, " %s=%q", k, v)
		}
	}
	switch strings.ToLower(elem.Type) {
	case "button":
		fmt.Fprintf(out, " >%s</button>", elem.Label)
	case "textarea":
		fmt.Fprintf(out, " ></textarea>")
	default:
		fmt.Fprintf(out, " >")
	}
	fmt.Fprintf(out, "</div>\n")
	return nil
}

// ValidateDate makes sure the date string conforms to YYYY-MM-DD
func ValidateDate(elem *Element, formValue string) bool {
	// FIXME: Need to check against min, max and step values
	if _, err := time.Parse("2006-01-02", formValue); err != nil {
		return false
	}
	return true
}

// ValidateDateTimeLocal makes sure the datetime string conforms to
// Spec: https://html.spec.whatwg.org/multipage/common-microsyntaxes.html#valid-local-date-and-time-string
func ValidateDateTimeLocal(elem *Element, formValue string) bool {
	// FIXME: Need to check against min, max and step values
	// See https://html.spec.whatwg.org/multipage/common-microsyntaxes.html#valid-local-date-and-time-string for validation steps
	if formValue == "" {
		return true
	}
	// Parse date component first
	if _, err := time.Parse("2006-01-02", formValue[0:10]); err != nil {
		return false
	}
	// String doesn't include the time so fail
	if len(formValue) <= 10 {
		return false
	}
	if formValue[11:12] != "T" && formValue[11:12] != " " {
		return false
	}
	if _, err := time.Parse("19:54:00", formValue[12:19]); err != nil {
		return false
	}
	// If we have timezone info
	if len(formValue) >= 19 {
		if _, err := time.Parse(time.RFC3339, formValue); err != nil {
			return false
		}
	}
	return true
}

// ValidateMonth parses the string for a year and month value, i.e. YYYY-MM style date string
func ValidateMonth(elem *Element, formValue string) bool {
	// FIXME: Need to check against min, max and step values
	if _, err := time.Parse("2006-02", formValue); err != nil {
		return false
	}
	return true
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

// ValidateEmailAddress parses email address to confirm it is valid
func ValidateEmailAddress(elem *Element, formValue string) bool {
	if _, err := mail.ParseAddress(formValue); err != nil {
		return false
	}
	return true
}

// ValidateTextElement will check to see if pattern is set, if so it will
// evaluate the formValue against the RegExp given in Pattern.
func ValidateTextElement(elem *Element, formValue string) bool {
	if elem.Pattern == "" {
		return true
	}
	re, err := regexp.CompilePOSIX(elem.Pattern)
	if err != nil {
		return false
	}
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

// ValidateNumber implements a number validation using the json package.
func ValidateNumber(elem *Element, formValue string) bool {
	if _, err := jsonDecodeNumber(formValue); err != nil {
		return false
	}
	return true
}

// ValidateRange retrieves the form's value as a float64 then checks if it is in range.  Min and max must befined in
// the attributes of the element since they are required to make the comparison. NOTE: ValidateRange isn't currently checking
// the step value as I don't know if the value of the input element is supposed to be an integer for real number. 
func ValidateRange(elem *Element, formValue string) bool {
	var (
		minNumber float64
		maxNumber float64
		numberValue float64
		err error
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

// ValidateTel validates formValue conforms to a phone number.
func ValidateTel(elem *Element, formValue string) bool {
	// NOTE: I am defaulting to US numbers because Caltech Library is in the US
	if _, err := phonenumbers.Parse(formValue, "US"); err != nil {
		return false
	}
	return true
}

// ValidateTime validates the formValue is a time format
func ValidateTime(elem *Element, formValue string) bool {
	//FIXME: Need to check against min and max values
	if _, err := time.Parse("15:04:05", formValue); err != nil {
		return false
	}
	return true
}

// ValidateURL validates a formValue is a URL
func ValidateURL(elem *Element, formValue string) bool {
	if _, err := url.Parse(formValue); err != nil {
		return false
	}
	return true
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


// DefaultValidator implement basic HTML5 form element validation
func DefaultValidator(elem *Element, formValue string) bool {
	switch strings.ToLower(elem.Type) {
	case "button":
		// A button value can be empty or a string, so always returns true.
		return true;
	case "checkbox":
		// Checkbox return their string value if checked.
		return strings.TrimSpace(formValue) != ""
	case "color":
		return ValidateColor(elem, formValue)
	case "date":
		return ValidateDate(elem, formValue)
	case "datetime-local":
		return ValidateDateTimeLocal(elem, formValue)
	case "email":
		return ValidateEmailAddress(elem, formValue)
	case "hidden":
		return ValidateTextElement(elem, formValue)
	case "image":
		// The element value should be none, see https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/image#technical_summary
		return formValue == ""
	case "month":
		return ValidateMonth(elem, formValue)
	case "number":
		return ValidateNumber(elem, formValue)
	case "password":
		// Passwords must be a single line of text, see https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/password
		if strings.Index(formValue, "\r") > -1 || strings.Index(formValue, "\n") > -1 {
			return false
		}
		return ValidateTextElement(elem, formValue)
	case "radio":
		// Checkbox return their string value if checked.
		return strings.TrimSpace(formValue) != ""
	case "range":
		return ValidateRange(elem, formValue)
	case "reset":
		return true;
	case "search":
		return ValidateTextElement(elem, formValue)
	case "submit":
		return true;
	case "tel":
		return ValidateTel(elem, formValue)	
	case "text":
		return ValidateTextElement(elem, formValue)
	case "textarea":
		return ValidateTextElement(elem, formValue)
	case "time":
		return ValidateTime(elem, formValue)
	case "url":
		return ValidateURL(elem, formValue)
	case "week":
		return ValidateWeek(elem, formValue)
	}
	log.Printf("DefaultValidator not implemented for %q with %q", elem.Type, formValue)
	return false
}
