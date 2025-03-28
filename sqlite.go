// sqlite.go is part of the Go models package.
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
