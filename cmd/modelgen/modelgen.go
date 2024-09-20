package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/models"

	// 3rd Party Packges
	"gopkg.in/yaml.v3"
)

const (
	helpText = `%{app_name}(1) user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

{app_name} 

# SYNOPSIS

{app_name} [OPTIONS] html|sqlite3 [MODEL_NAME] [OUT_NAME]

# DESCRIPTION

{app_name} is a demonstration of the models package for Go.  It can read
a model expressed as YAML and transform the result into an HTML web form
or SQLite3 database schema.

MODEL_NAME is the name of the YAML file to read. If no filename is provided
then the model is read from standard input.

OUT_NAME is the name of the file to write. If it is loft off then
then standard out is used.

# OPTIONS

-help
: Display help

-version
: Display {app_name} version.

-license
: Display {app_name} license.

# EXAMPLE

~~~
{app_name} html guestbook.yaml guestbook.html
{app_name} sqlite3 guestbook.yaml guestbook.sql
~~~

`
)

var (
	// Standard Options
	showHelp    bool
	showLicense bool
	showVersion bool
)

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
	if ! model.Check(eout) {
		fmt.Fprintf(eout, "ERROR: problem with model")
		os.Exit(1)
	}
    // Now transform the model.	
	verb := strings.ToLower(args[0])
	switch verb {
		case "html":
			if err := model.ToHTML(out); err != nil {
				fmt.Fprintf(eout, "ERROR: %s\n", err)
				os.Exit(1)
			}
		case "sqlite3":
			if err := model.ToSQLiteScheme(out); err != nil {
				fmt.Fprintf(eout, "ERROR: %s\n", err)
				os.Exit(1)
			}
		default:
			fmt.Fprintf(eout, "%q output format not supported\n", verb)
			os.Exit(1)
	}
}

