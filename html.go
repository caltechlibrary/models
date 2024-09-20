package models

import (
	"fmt"
	"io"
	"strings"
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
