package models

import (
	"fmt"
	"io"
	"strings"
)

// ModelToSQLiteScheme takess a model and renders the SQLite DB Schema to out.
// @param out: io.Writer, the target to render the text into
// @param model: *Model, the model to be rendered.
func ModelToSQLiteScheme(out io.Writer, model *Model) error {
	if !IsValidVarname(model.Id) {
		return fmt.Errorf("model id that can't be used for table name, %q", model.Id)
	}
	if model.Description != "" {
		fmt.Fprintf(out, "-- %s\n", strings.ReplaceAll(model.Description, "\n", "\n-- "))
	}
	fmt.Fprintf(out, "create table %s if not exists (\n", model.Id)
	addNL := false
	for i, elem := range model.Elements {
		if !IsValidVarname(elem.Id) {
			return fmt.Errorf("element id can't be used for column name, %q", elem.Id)
		}
		if i > 0 {
			fmt.Fprintf(out, ",\n")
			addNL = true
		}
		//NOTE: Map HTML5 types to SQLite3 type
		var columnType string
		switch strings.ToLower(elem.Type) {
		case "int":
			columnType = "int"
		case "integer":
			columnType = "int"
		case "float":
			columnType = "real"
		case "real":
			columnType = "real"
		case "numeric":
			columnType = "num"
		case "number":
			columnType = "num"
		case "date":
			columnType = "text"
		case "datetime-local":
			columnType = "text"
		case "checkbox":
			columnType = "boolean"
		default:
			columnType = "text"
		}
		if elem.Generator != "" {
			switch elem.Generator {
			case "autoincrement":
				columnType = fmt.Sprintf("%s autoincrement", columnType)
			case "date":
				columnType = fmt.Sprintf("%s default current_date not null", columnType)
			case "created_date":
				columnType = fmt.Sprintf("%s default current_date not null", columnType)
			case "current_date":
				columnType = fmt.Sprintf("%s default current_date not null", columnType)
			case "timestamp":
				columnType = fmt.Sprintf("%s default current_timestamp not null", columnType)
			case "created_timestamp":
				columnType = fmt.Sprintf("%s default current_timestamp not null", columnType)
			case "current_timestamp":
				columnType = fmt.Sprintf("%s default current_timestamp not null", columnType)
			}
		}
		if elem.IsObjectId {
			columnType = fmt.Sprintf(" %s primary key", columnType)
		}
		fmt.Fprintf(out, "  %s %s", elem.Id, columnType)
	}
	if addNL {
		fmt.Fprintf(out, "\n")
	}
	fmt.Fprintf(out, ");\n")
	return nil
}
