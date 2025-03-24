// transcript.go is part of the Go models package.
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

//
// This package renders TypeScript classes from a Model
//

// ModelToTypeScriptClass render a model as TypeScript interface and class
// @param out : io.Writer, where the generated code is written
// @param model : *Model, the model to render as TypeScript.
// @return error or nil (on success)
func ModelToTypeScriptClass(out io.Writer, model *Model) error {
	// Include model.Id and model.Description as an opening comment.
	fmt.Fprintf(out, `/*
Model: %s

%s
*/

`, model.Id, model.Description)

	className := model.Id
	if len(className) > 1 {
		className = strings.ToUpper(className[0:1]) + className[1:]
	} else {
		className = strings.ToUpper(className)
	}
	interfaceName := className + "Interface"
	fmt.Fprintf(out, `// %s model's %s interface
export interface %s {
`, interfaceName, className, interfaceName)
	for _, elem := range model.Elements {
		varName := elem.Id
		varType := mapTypeToTypeScript(elem)
		fmt.Fprintf(out, "\t%s: %s;\n", varName, varType)
	}
	fmt.Fprintln(out, `}

`)

	fmt.Fprintf(out, `// %s's class definition
export class %s implements %s {
`, className, className, interfaceName)
	for _, elem := range model.Elements {
		varName := elem.Id
		varType := mapTypeToTypeScript(elem)
		switch varType {
		case "string":
			varType = "string = \"\""
		case "number":
			varType = "number = 0.0"
		case "boolean":
			varType = "boolean = false"
		}
		fmt.Fprintf(out, "\t%s: %s;\n", varName, varType)
	}
	fmt.Fprintln(out, `}

`)
	return nil
}

func mapTypeToTypeScript(elem *Element) string {
	dTypes := map[string]string{
		"date":           "string",
		"datetime-local": "string",
		"month":          "string",
		"color":          "string",
		"email":          "string",
		"number":         "number",
		"range":          "number[]",
		"text":           "string",
		"tel":            "string",
		"time":           "string",
		"url":            "string",
		"checkbox":       "boolean",
		"password":       "string",
		"radio":          "string",
		"textarea":       "string",
		"orcid":          "string",
		"isni":           "string",
		"uuid":           "string",
		"ror":            "string",
	}
	if val, ok := dTypes[elem.Type]; ok {
		return val
	}
	return "string"
}
