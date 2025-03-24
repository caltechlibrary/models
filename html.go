// html.go is part of the Go models package.
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
	"fmt"
	"io"
	"strings"
)

// ModelToHTML takes a model and renders an input form. The form is not
// populated with values through that could be done easily via JavaScript and DOM calls.
func ModelToHTML(out io.Writer, model *Model) error {
	// Include the description as an HTML comment.
	// Write opening form element
	if model.Id != "" {
		fmt.Fprintf(out, "<!-- %s: %s -->\n", model.Id, model.Description)
		fmt.Fprintf(out, "<form id=%q", model.Id)
	} else {
		fmt.Fprintf(out, "<!-- %s -->\n", model.Description)
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

// ElementToHTML renders an individual element as HTML, includes label as well as input element.
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
