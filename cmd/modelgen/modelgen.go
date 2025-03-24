// modelgen.go is part of the Go models package.
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
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/models"

	// 3rd Party Packges
	"github.com/pkg/fileutils"
	"gopkg.in/yaml.v3"
)

const (
	helpText = `%{app_name}(1) user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

{app_name} 

# SYNOPSIS

{app_name} [OPTIONS] ACTION [MODEL_NAME] [OUT_NAME]

# DESCRIPTION

{app_name} is a demonstration of the models package for Go.  It can read
a model expressed as YAML and transform the result into an HTML web form
or SQLite3 database schema.

MODEL_NAME is the name of the YAML file to read. If no filename is provided
then the model is read from standard input.

OUT_NAME is the name of the file to write. If it is loft off then
then standard out is used.

# ACTION

An action can be "model", "html" or "sqlite". Actions result in a file or
content generation rendering a model.

model MODEL_NAME
: This action is an interactive modeler. It generates YAML file holding
the model.  MODEL_NAME is required is used as the filename for the model.

html
: This action will render a YAML model as HTML. If no MODEL_NAME is provided
then the YAML is read from standard input.

sqlite
: This action will render a SQL file suitable for use with SQLite 3.

typescript
: This action will render a TypeScript class definition

python
: This action with render a Python class definition

# OPTIONS

-help
: Display help

-version
: Display {app_name} version.

-license
: Display {app_name} license.

# EXAMPLE

In this example we create a new model YAML file interactively using
the "model" action. Then create an HTML page followed by SQL file
holding the SQL schema for SQLite 3.

~~~
{app_name} model guestbook.yaml
{app_name} html guestbook.yaml guestbook.html
{app_name} sqlite guestbook.yaml guestbook.sql
~~~

`
)

var (
	// Standard Options
	showHelp    bool
	showLicense bool
	showVersion bool
)

// getAnswer get a Y/N response from buffer
func getAnswer(buf *bufio.Reader, defaultAnswer string, lower bool) string {
	answer, err := buf.ReadString('\n')
	if err != nil {
		return ""
	}
	answer = strings.TrimSpace(answer)
	if answer == "" {
		answer = defaultAnswer
	}
	if lower {
		return strings.ToLower(answer)
	}
	return answer
}


func main() {
	appName := path.Base(os.Args[0])

	// Standard Options
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "version", false, "display version")

	// We're ready to process args
	flag.Parse()
	args := flag.Args()

	var err error
	in := os.Stdin
	out := os.Stdout
	eout := os.Stderr

	if showHelp {
		fmt.Fprintf(out, "%s\n", models.FmtHelp(helpText, appName, models.Version, models.ReleaseDate, models.ReleaseHash))
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprintf(out, "%s\n", models.LicenseText)
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintf(out, "%s %s %s\n", appName, models.Version, models.ReleaseHash)
		os.Exit(0)
	}
	if len(args) == 0 {
		fmt.Fprintf(eout, "%s\n", models.FmtHelp(helpText, appName, models.Version, models.ReleaseDate, models.ReleaseHash))
		os.Exit(1)
	}
	// Now transform the model.
	verb := strings.ToLower(args[0])
	if verb == "model" {
		if len(args) < 2 {
			fmt.Fprintf(eout, "ERROR: must provide a name for the YAML model file")
			os.Exit(1)
		}
		fName := args[1]
		modelId := path.Base(fName)
		modelId = strings.ToLower(strings.TrimSuffix(modelId, ".yaml"))
		model, err := models.NewModel(modelId)
		if err != nil {
			fmt.Fprintf(eout, "ERROR: %s\n", err)
			os.Exit(5)
		}
		// If file exists, make backup then load the contents into memory
		backupFile := false
		if _, err := os.Stat(fName); err == nil {
			backupFile = true
			src, err := os.ReadFile(fName)
			if err != nil {
				fmt.Fprintf(eout, "ERROR: %s\n", err)
				os.Exit(2)
			}
			if err := yaml.Unmarshal(src, &model); err != nil {
				fmt.Fprintf(eout, "ERROR: %s\n", err)
				os.Exit(3)
			}
		}
		// Decide if I'm going to create or open an existing YAML file.
		fout, err := os.OpenFile(args[1], os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			fmt.Fprintf(eout, "ERROR: %s\n", err)
			os.Exit(4)
		}
		defer fout.Close()

		models.SetDefaultTypes(model)
		model.Register("yaml", models.ModelToYAML)
		if err := models.ModelInteractively(model); err != nil {
			fmt.Fprintf(eout, "ERROR: %s\n", err)
			os.Exit(6)
		}
		if model.HasChanges() {
			buf := bufio.NewReader(in)
			fmt.Fprintf(out, "Save %s (Y/n)?", fName)
			answer := getAnswer(buf, "y", true)
			if answer == "y" {
				// backup file if needed
				if backupFile {
					if err := fileutils.CopyFile(fName + ".bak", fName) ; err != nil {
						fmt.Fprintf(eout, "ERROR: %s\n", err)
						os.Exit(7)
					}
				}
				encoder := yaml.NewEncoder(fout)
				encoder.SetIndent(2)
				if err := encoder.Encode(model); err != nil {
					fmt.Fprintf(eout, "ERROR: %s\n", err)
					os.Exit(8)
				}
			}
		}
		os.Exit(0)
	}

	if len(args) > 1 {
		in, err = os.Open(args[1])
		if err != nil {
			fmt.Fprintf(eout, "ERROR: %s\n", err)
			os.Exit(1)
		}
		defer in.Close()
	}

	if len(args) > 2 {
		out, err = os.Create(args[2])
		if err != nil {
			fmt.Fprintf(eout, "ERROR: %s\n", err)
			os.Exit(1)
		}
		defer out.Close()
	}
	src, err := io.ReadAll(in)
	if err != nil {
		fmt.Fprintf(eout, "ERROR: %s\n", err)
		os.Exit(1)
	}
	model := new(models.Model)
	if err := yaml.Unmarshal(src, model); err != nil {
		fmt.Fprintf(eout, "ERROR: %s\n", err)
		os.Exit(1)
	}
	if !model.Check(eout) {
		fmt.Fprintf(eout, "ERROR: problem with model")
		os.Exit(1)
	}
	model.Register("html", models.ModelToHTML)
	model.Register("sqlite", models.ModelToSQLiteScheme)
	model.Register("sqlite3", models.ModelToSQLiteScheme)
	model.Register("typescript", models.ModelToTypeScriptClass)
	model.Register("python", models.ModelToPythonClass)
	if err := model.Render(out, verb); err != nil {
		fmt.Fprintf(eout, "ERROR: %s\n", err)
		os.Exit(1)
	}
}
