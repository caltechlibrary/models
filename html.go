package models

import (
	"fmt"
	"io"
)

func ModelToHTML(out io.Writer, model *Model) error {
	// Write opening form element
	fmt.Fprintf(out, "<form id=%q", model.id)
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
	fmt.Fprintf(out, ">\n")
	for _, elem := range model.Elements {
		ElementToHTML(out, elem)
	}

	// Write closing form element
	fmt.Fprintf(out, "</form>\n")
	return nil
}

func ElementToHTML(out io.Writer, elem *Element) error {
	return fmt.Errorf("ElementToHTML() not implemented")
}
