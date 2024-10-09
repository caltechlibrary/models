package models

import (
	"fmt"
	"io"
	"strings"
)

//
// This package renders Python classes from a Model
//

func ModelToPythonClass(out io.Writer, model *Model) error {
	// Include model.Id and model.Description as an opening comment.
	fmt.Fprintf(out, `#
# Model: %s
#
# %s
#

`, model.Id, model.Description)

	className := model.Id
	if len(className) > 1 {
		className = strings.ToUpper(className[0:1])  + className[1:]
	} else {
		className = strings.ToUpper(className)
	}
	fmt.Fprintf(out, `# %s model's definition
class %s:
`, className, className)
	for _, elem := range model.Elements {
		varName := elem.Id
		//varType := mapTypeToTypeScript(elem)
		fmt.Fprintf(out, "    %s\n", varName)
	}
	fmt.Fprintln(out, "\n    def __init__(self):")
	for _, elem := range model.Elements {
		varName := elem.Id
		varDefault := mapTypeToPythonDefault(elem)
		fmt.Fprintf(out, "        self.%s = %s\n", varName, varDefault)
	}
	return nil
}

func mapTypeToPythonDefault(elem *Element) string {
	dTypes := map[string]string{
		"date":           "",
		"datetime-local": "",
		"month":          "",
		"color":          "",
		"email":          "",
		"number":         "0",
		"range":          "[]",
		"text":           "",
		"tel":            "",
		"time":           "",
		"url":            "",
		"checkbox":       "",
		"password":       "",
		"radio":          "",
		"textarea":       "",
		"orcid":          "",
		"isni":           "",
		"uuid":           "",
		"ror":            "",
	}
	if val, ok := dTypes[elem.Type]; ok {
		if val == "" {
			return `""`
		}
		return val
	}
	return ""
}
