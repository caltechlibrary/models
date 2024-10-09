package models

import (
	"fmt"
	"io"
	"strings"
)

//
// This package renders TypeScript classes from a Model
//

func ModelToTypeScriptClass(out io.Writer, model *Model) error {
	// Include model.Id and model.Description as an opening comment.
	fmt.Fprintf(out, `/*
Model: %s

%s
*/

`, model.Id, model.Description)

	className := strings.ToTitle(model.Id)
	if len(className) > 1 {
		className = strings.ToUpper(className[0:1])  + strings.ToUpper(className[1:])
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
		fmt.Fprintf(out, "\t%s: %s;\n", varName, varType)
	}
	fmt.Fprintln(out, `}

`)
	return nil
}

func mapTypeToTypeScript(elem *Element) string {
	dTypes := map[string]string{
		"date":           "Date",
		"datetime-local": "Date",
		"month":          "Date",
		"color":          "string",
		"email":          "string",
		"number":         "number",
		"range":          "[]",
		"text":           "string",
		"tel":            "string",
		"time":           "Date",
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
