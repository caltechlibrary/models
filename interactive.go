// interactive.go is part of the Go models package.
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
	"os"
	"regexp"
	"strconv"
	"strings"

	// 3rd Party Packages
	"gopkg.in/yaml.v3"
)

// ModelToYAML renders a Model struct as YAML
func ModelToYAML(out io.Writer, model *Model) error {
	encoder := yaml.NewEncoder(out)
	encoder.SetIndent(2)
	if err := encoder.Encode(model); err != nil {
		return err
	}
	return nil
}

// removeElement removes an element from a model
func removeElement(model *Model, in io.Reader, out io.Writer, eout io.Writer, elementId string) error {
	elemFound := false
	for i, elem := range model.Elements {
		if elem.Id == elementId {
			model.Elements = append(model.Elements[:i], model.Elements[(i+1):]...)
			model.Changed(true)
			elemFound = true
		}
	}
	if !elemFound {
		return fmt.Errorf("failed to find %s.%s", model.Id, elementId)
	}
	return nil
}

// normalizeInputType takes a string and returns a normlized input type
// (e.g. lowercased and spaces trim) and if the type is an alias
// then a pattern is returned along with the normalized type.
// e.g. DATE -> date,""
//
//	eMail -> email, ""
//	ORCID -> orcid, "[0-9]{4}-[0-9]{4}-[0-9]{4}-[0-9]{3}[0-9A-Z]"
func normalizeInputType(inputType string) (string, string) {
	patternMap := map[string]string{
		// This is where we put the aliases to common regexp validated patterns
		"orcid": "[0-9]{4}-[0-9]{4}-[0-9]{4}-[0-9]{3}[0-9A-Z]",
	}
	val := strings.ToLower(inputType)
	if pattern, ok := patternMap[val]; ok {
		return val, pattern
	}
	return val, ""
}

// addElementStub adds an empty element to model.Elements list. Returns a new model list and error value
func addElementStub(model *Model, elementId string) ([]string, error) {
	elementList := model.GetElementIds()
	if !IsValidVarname(elementId) {
		return elementList, fmt.Errorf("%q is not a valid elmeent id", elementId)
	}
	elem, err := NewElement(elementId)
	if err != nil {
		return elementList, err
	}
	model.Elements = append(model.Elements, elem)
	elementList = model.GetElementIds()
	return elementList, nil
}

// modifyModelAttributesTUI provides a text UI for managing a model's attributes
func modifyModelAttributesTUI(model *Model, in io.Reader, out io.Writer, eout io.Writer) error {
	prompt := NewPrompt(in, out, eout)
	if model.Attributes == nil {
		model.Attributes = map[string]string{}
	}
	for quit := false; !quit; {
		attributeList := []string{}
		for k, v := range model.Attributes {
			attributeList = append(attributeList, fmt.Sprintf("%s -> %q", k, v))
		}
		menu, opt := prompt.SelectMenu(
			fmt.Sprintf("Manage %s attributes (none required)", model.Id),
			"Choices [a]dd, [m]odify, [r]emove or press enter when done",
			attributeList, "", "", true)
		if len(menu) > 0 {
			menu = menu[0:1]
		}
		var ok bool
		opt, ok = getIdFromList(attributeList, opt)
		switch menu {
		case "a":
			if opt == "" {
				fmt.Fprintf(out, "Enter attribute name: ")
				opt = prompt.GetAnswer("", true)
			}
			if !IsValidVarname(opt) {
				fmt.Fprintf(eout, "%q is not a valid attribute name\n", opt)
			} else {
				model.Attributes[opt] = ""
				model.Changed(true)
			}
		case "m":
			if opt == "" {
				fmt.Fprintf(out, "Enter attribute name: ")
				opt = prompt.GetAnswer("", true)
			}
			fmt.Fprintf(out, "Enter %s's value: ", opt)
			val := prompt.GetAnswer("", false)
			if val != "" {
				model.Attributes[opt] = val
				model.Changed(true)
			}
		case "r":
			if opt == "" {
				fmt.Fprintf(out, "Enter attribute name to remove: ")
				opt = prompt.GetAnswer("", true)
				opt, ok = getIdFromList(attributeList, opt)
			}
			if ok {
				if _, ok := model.Attributes[opt]; ok {
					delete(model.Attributes, opt)
					model.Changed(true)
				}
			}
			if !ok {
				fmt.Fprintf(eout, "failed to find %q in attributes\n", opt)
			}
		case "q":
			quit = true
		case "":
			quit = true
		default:
			fmt.Fprintf(eout, "failed to underand %q\n", opt)
		}
	}
	return nil
}

// modifyElementAttributesTUI provides a text UI for managing a model's element's attributes
func modifyElementAttributesTUI(model *Model, in io.Reader, out io.Writer, eout io.Writer, elementId string) error {
	prompt := NewPrompt(in, out, eout)
	elem, _ := model.GetElementById(elementId)
	for quit := false; !quit; {
		var ok bool
		attributeList := []string{}
		for k, v := range elem.Attributes {
			attributeList = append(attributeList, fmt.Sprintf("%s -> %q", k, v))
		}
		val := ""
		menu, opt := prompt.SelectMenu(
			fmt.Sprintf("Modify element %s.%s attributes", model.Id, elementId),
			"Choices [a]dd, [m]odify, [r]emove or press enter when done",
			attributeList, "", "", true)
		if len(menu) > 0 {
			menu = menu[0:1]
		}
		if i := strings.Index(opt, " "); i >= 0 {
			opt, val = opt[0:i], opt[i+1:]
		}
		opt, ok = getIdFromList(attributeList, opt)
		switch menu {
		case "a":
			if opt == "" {
				fmt.Fprintf(out, "Enter attribute name: ")
				opt = prompt.GetAnswer("", true)
			}
			if !IsValidVarname(opt) {
				fmt.Fprintf(eout, "%q is not a valid attribute name\n", opt)
			} else {
				elem.Attributes[opt] = val
				elem.Changed(true)
			}
		case "m":
			if opt == "" {
				fmt.Fprintf(out, "Enter attribute name: ")
				opt = prompt.GetAnswer("", true)
			}
			if val == "" {
				fmt.Fprintf(out, "Enter %s's value: ", opt)
				val = prompt.GetAnswer("", false)
			}
			if val != "" {
				elem.Attributes[opt] = val
				elem.Changed(true)
			}
		case "r":
			if opt == "" {
				fmt.Fprintf(out, "Enter attribute name to remove: ")
				opt = prompt.GetAnswer("", true)
				opt, ok = getIdFromList(attributeList, opt)
			}
			if ok {
				if _, ok = elem.Attributes[opt]; ok {
					delete(elem.Attributes, opt)
					elem.Changed(true)
				}
			}
			if !ok {
				fmt.Fprintf(eout, "failed to find %q in attributes\n", opt)
			}
		case "q":
			quit = true
		case "":
			quit = true
		default:
			fmt.Fprintf(eout, "failed to underand %q\n", opt)
		}
	}
	return nil
}

func modifySelectElementTUI(elem *Element, in io.Reader, out io.Writer, eout io.Writer, modelId string) error {
	prompt := NewPrompt(in, out, eout)
	if elem.Options == nil {
		elem.Options = []map[string]string{}
	}
	for quit := false; !quit; {
		optionsList := getValueLabelList(elem.Options)
		menu, opt := prompt.SelectMenu(
			fmt.Sprintf("Manage %s.%s options", modelId, elem.Id),
			"Choices [a]dd, [m]odify, [r]emove or press enter when done",
			optionsList,
			"", "", true)
		if len(menu) > 0 {
			menu = menu[0:1]
		}
		var (
			val   string
			label string
		)
		switch menu {
		case "a":
			if opt == "" {
				fmt.Fprintf(out, "Enter an option value: ")
				val = prompt.GetAnswer("", true)
			} else {
				val = opt
			}
			if val == "" {
				fmt.Fprintf(eout, "Error, an option value required\n")
			} else {
				fmt.Fprintf(out, "Enter an option label: ")
				label = prompt.GetAnswer("", false)
				if label == "" {
					label = val
				}
				option := map[string]string{
					val: label,
				}
				elem.Options = append(elem.Options, option)
				elem.Changed(true)
			}
		case "m":
			pos, ok := -1, false
			if i, err := strconv.Atoi(opt); err == nil {
				// Adjust to a zero based array
				if i > 0 && i <= len(optionsList) {
					pos, ok = (i - 1), true
				}
			}
			if !ok {
				fmt.Fprintf(out, "Enter option number: ")
				pos, ok = prompt.GetDigit(optionsList)

			}
			if ok {
				option := elem.Options[pos]
				val, label, _ := getValAndLabel(option)
				fmt.Fprintf(out, "Enter an option label (for %q): ", val)
				answer := prompt.GetAnswer("", false)
				if answer != "" {
					label = answer
				}
				if label == "" {
					label = val
				}
				option = map[string]string{
					val: label,
				}
				// Replace the option in the options list.
				elem.Options[pos] = option
				elem.Changed(true)
			} else {
				fmt.Fprintf(eout, "Failed to identify option %q\n", opt)
			}
		case "r":
			pos, ok := -1, false
			if i, err := strconv.Atoi(opt); err == nil {
				// Adjust to a zero based array
				if i > 0 && i <= len(optionsList) {
					pos, ok = (i - 1), true
				}
			} else {
				fmt.Fprintf(out, "Enter option number: ")
				pos, ok = prompt.GetDigit(optionsList)
			}
			if ok {
				elem.Options = append(elem.Options[0:pos], elem.Options[(pos+1):]...)
				elem.Changed(true)
			} else {
				fmt.Fprintf(eout, "failed to remove option number (%d) %q\n", pos, opt)
			}
		case "q":
			quit = true
		case "":
			quit = true
		default:
			fmt.Fprintf(eout, "did not understand, %s %s\n", menu, opt)
		}
	}
	return nil
}

func modifyElementTUI(model *Model, in io.Reader, out io.Writer, eout io.Writer, elementId string) error {
	prompt := NewPrompt(in, out, eout)
	elem, ok := model.GetElementById(elementId)
	if !ok {
		return fmt.Errorf("could not find %q element", elementId)
	}
	var (
		menu string
		opt  string
	)
	for quit := false; !quit; {
		attributeList := getAttributeIds(elem.Attributes)
		switch elem.Type {
		case "select":
			optionsList := getValueLabelList(elem.Options)
			menu, opt = prompt.SelectMenu(
				fmt.Sprintf("Modify element %s.%s", model.Id, elementId),
				"Choices [t]ype, [l]abel, [a]ttributes, [o]ptions or press enter when done",
				[]string{
					fmt.Sprintf("id %s", elementId),
					fmt.Sprintf("type %s", elem.Type),
					fmt.Sprintf("label %s", elem.Label),
					fmt.Sprintf("attributes:\n\t\t%s", strings.Join(attributeList, ",\n\t\t")),
					fmt.Sprintf("options:\n\t\t%s", strings.Join(optionsList, ",\n\t\t")),
				},
				"", "", true)
		case "textarea":
			menu, opt = prompt.SelectMenu(
				fmt.Sprintf("Modify element %s.%s", model.Id, elementId),
				"Choices [t]ype, [l]abel, [a]ttributes or press enter when done",
				[]string{
					fmt.Sprintf("id %s", elementId),
					fmt.Sprintf("type %s", elem.Type),
					fmt.Sprintf("label %s", elem.Label),
					fmt.Sprintf("attributes:\n\t\t%s", strings.Join(attributeList, ",\n\t\t")),
				},
				"", "", true)
		default:
			menu, opt = prompt.SelectMenu(
				fmt.Sprintf("Manage %s.%s element", model.Id, elementId),
				"Choices [t]ype, [l]abel, [o]bject identifier, [p]attern, [a]ttributes, [g]enerator or press enter when done",
				[]string{
					fmt.Sprintf("id %s", elementId),
					fmt.Sprintf("type %s", elem.Type),
					fmt.Sprintf("label %s", elem.Label),
					fmt.Sprintf("pattern %s", elem.Pattern),
					fmt.Sprintf("attributes:\n\t\t%s", strings.Join(attributeList, ",\n\t\t")),
					fmt.Sprintf("object identifier? %t", elem.IsObjectId),
					fmt.Sprintf("generator %s", elem.Generator),
				},
				"", "", true)
		}
		if len(menu) > 0 {
			menu = menu[0:1]
		}
		switch menu {
		case "t":
			if opt == "" {
				fmt.Fprintf(out, `Enter type string (e.g. text, email, date, textarea, select, orcid): `)
				opt = prompt.GetAnswer("", false)
			}
			if opt != "" {
				eType := opt
				if ok := model.IsSupportedElementType(eType); !ok {
					fmt.Fprintf(eout, "%q is not a supported element type", opt)
				} else {
					// FIXME: If elem.Type is select I need to provide an choices TUI for values and labels
					elem.Type, elem.Pattern = normalizeInputType(eType)
					elem.Changed(true)
					if elem.Type == "select" {
						if err := modifySelectElementTUI(elem, in, out, eout, model.Id); err != nil {
							fmt.Fprintf(eout, "ERROR (%q): %s\n", elementId, err)
						}
					}
				}
			}
		case "p":
			if opt == "" {
				fmt.Fprintf(out, "Enter the regexp to use: ")
				opt = prompt.GetAnswer("", false)
			}
			// FIXME: how do I clear a pattern if I nolonger want to use one?
			if opt != "" {
				// NOTE: remove a pattern by having it match everything, i.e. astrix
				if opt == "*" {
					elem.Pattern = ""
				} else {
					elem.Pattern = opt
				}
				model.Changed(true)
			}
		case "l":
			if opt == "" {
				fmt.Fprintf(out, "Enter a label: ")
				opt = prompt.GetAnswer("", false)
			}
			if opt != "" {
				// NOTE: remove a pattern by having it match everything, i.e. astrix
				if opt != elem.Label {
					elem.Label = opt
				}
				model.Changed(true)
			}
		case "a":
			if err := modifyElementAttributesTUI(model, in, out, eout, elementId); err != nil {
				fmt.Fprintf(eout, "%s\n", err)
			}
		case "o":
			if elem.Type == "select" {
				if err := modifySelectElementTUI(elem, in, out, eout, model.Id); err != nil {
					fmt.Fprintf(eout, "ERROR (%q): %s\n", elem.Id, err)
				}
			} else {
				elem.IsObjectId = !elem.IsObjectId
				elem.Changed(true)
			}
		case "g":
			if opt == "" {
				fmt.Fprintf(out, "Enter generator (e.g. autoincrement, uuid, current_timestamp, created_timestamp, current_date, created_date) ")
				opt = prompt.GetAnswer("", true)
			}
			fmt.Fprintf(out, "DEBUG opt -> %q, elem.Generator -> %q\n", opt, elem.Generator)
			if opt != "" {
				if strings.HasPrefix(opt, "\"") || strings.HasPrefix(opt, "'") {
					elem.Generator = ""
					fmt.Fprintf(out, "DEBUG opt -> %q, elem.Generator -> %q\n", opt, elem.Generator)
					elem.Changed(true)
				} else if opt != elem.Generator {
					elem.Generator = opt
					elem.Changed(true)
				}
			}
		case "q":
			quit = true
		case "":
			quit = true
		default:
			fmt.Fprintf(eout, "did not understand %q\n", menu)
		}
	}
	return nil
}

func removeElementFromModel(model *Model, elementId string) error {
	for i, elem := range model.Elements {
		if elem.Id == elementId {
			model.Elements = append(model.Elements[0:i], model.Elements[(i+1):]...)
			model.Changed(true)
			return nil
		}
	}
	return fmt.Errorf("%s not found", elementId)
}

// modifyElementsTUI modify a specific model's element list.
func modifyElementsTUI(model *Model, in io.Reader, out io.Writer, eout io.Writer) error {
	var (
		err    error
		answer string
	)
	prompt := NewPrompt(in, out, eout)
	// FIXME: Need to support editing model attributes, then allow for modifying model's body to be modified.
	for quit := false; !quit; {
		elementList := model.GetElementIds()
		menu, opt := prompt.SelectMenu(
			fmt.Sprintf("Manage %s elements", model.Id),
			"Choices [a]dd, [m]odify, [r]emove or press enter when done",
			elementList,
			"", "", true)
		if len(menu) > 1 {
			menu = menu[0:1]
		}
		elementId, ok := getIdFromList(elementList, opt)
		switch menu {
		case "a":
			if !ok {
				fmt.Fprintf(out, "Enter element id to add: ")
				opt = prompt.GetAnswer("", false)
				elementId, ok = getIdFromList(elementList, opt)
			}
			if ok {
				elementList, err = addElementStub(model, elementId)
				if err != nil {
					fmt.Fprintf(eout, "WARNING: %s\n", err)
				}
			}
		case "m":
			if elementId == "" {
				fmt.Fprintf(out, "Enter element id to modify: ")
				opt = prompt.GetAnswer("", false)
				elementId, ok = getIdFromList(elementList, opt)
			}
			if err := modifyElementTUI(model, in, out, eout, elementId); err != nil {
				fmt.Fprintf(eout, "ERROR (%q): %s\n", elementId, err)
			}
		case "r":
			if elementId == "" {
				fmt.Fprintf(out, "Enter element id to remove: ")
				opt = prompt.GetAnswer("", false)
				elementId, ok = getIdFromList(elementList, opt)
			}
			if ok {
				if err := removeElementFromModel(model, elementId); err != nil {
					fmt.Fprintf(eout, "ERROR (%q): %s\n", elementId, err)
				}
			}
		case "q":
			quit = true
		case "":
			quit = true
		default:
			fmt.Fprintf(eout, "\n\nERROR: Did not understand %q\n\n", answer)
		}

	}
	return nil
}

func modifyModelMetadataTUI(model *Model, in io.Reader, out io.Writer, eout io.Writer) error {
	prompt := NewPrompt(in, out, eout)
	for quit := false; !quit; {
		menu, opt := prompt.SelectMenu(
			"Manage Model Metadata",
			"Menu [i]d, [d]escription or press enter when done",
			[]string{
				fmt.Sprintf("id: %q", model.Id),
				//fmt.Sprintf("title: %q", model.Title),
				fmt.Sprintf("description: %q", model.Description),
			},
			"", "", true)

		if len(menu) > 0 {
			menu = menu[0:1]
		}
		switch menu {
		case "i":
			if opt == "" {
				fmt.Fprintf(out, `Enter id: `)
				opt = prompt.GetAnswer("", false)
			}
			if opt != "" {
				if opt != model.Id {
					model.Id = opt
					model.Changed(true)
				}
			}
		case "d":
			if opt == "" {
				fmt.Fprintf(out, `Enter Description: `)
				opt = prompt.GetAnswer("", false)
			}
			if opt != "" {
				if opt != model.Description {
					model.Description = opt
					model.Changed(true)
				}
			}
		case "q":
			quit = true
		default:
			quit = true
		}
	}
	return nil
}

func ModelInteractively(model *Model) error {
	in := os.Stdin
	out := os.Stdout
	eout := os.Stderr

	// Manage Model Metadata
	if err := modifyModelMetadataTUI(model, in, out, eout); err != nil {
		return err
	}
	/*
		// Manage Model Attributes
		if err := modifyModelAttributesTUI(model, in, out, eout); err != nil {
			return err
		}
	*/
	// Manage Model Elements
	if err := modifyElementsTUI(model, in, out, eout); err != nil {
		return err
	}
	return nil
}

//
// Misc functions
//

func getIdFromList(list []string, id string) (string, bool) {
	nRe := regexp.MustCompile(`^[0-9]+$`)
	// See if we have been given a model number or a name
	if isDigit := nRe.Match([]byte(id)); isDigit {
		mNo, err := strconv.Atoi(id)
		if err == nil {
			// Adjust provided integer for zero based index.
			if mNo > 0 {
				mNo--
			} else {
				mNo = 0
			}
			if mNo < len(list) {
				if strings.Contains(list[mNo], " ") {
					parts := strings.SplitN(list[mNo], " ", 2)
					return parts[0], true
				}
				return list[mNo], true
			}
		}
	}
	if IsValidVarname(id) {
		return id, true
	}
	return "", false
}

// Get return the first key and value pair
func getValAndLabel(option map[string]string) (string, string, bool) {
	for val, label := range option {
		return val, label, true
	}
	return "", "", false
}

// getValueLabelList takes an array of map[string]string and yours a list of
// strings indicating the value and label
func getValueLabelList(list []map[string]string) []string {
	options := []string{}
	for _, m := range list {
		val, label, ok := getValAndLabel(m)
		if ok {
			options = append(options, fmt.Sprintf("%s %s", val, label))
		}
	}
	return options
}
