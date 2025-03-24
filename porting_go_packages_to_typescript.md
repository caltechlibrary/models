I am porting a project called [Newt](https://caltechlibrary.github.io/newt) from Go to TypeScript running in Deno 2.2. I would like to keep the same file organization as the code is ported. The first source file to port is called `ast.go`. The code follows below.

```go
/**
 * ast.go holds the data structure that defines Newt applications.
 *
 * @author R. S. Doiel
 */
package newt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	// Caltech Library Packages
	"github.com/caltechlibrary/models"

	// 3rd Party
	"github.com/aymerick/raymond"
	"gopkg.in/yaml.v3"
)

// AST holds a configuration for Newt for the data router and code generator.
type AST struct {
	// AppMetadata holds your application's metadata such as needed to render an "about" page in your final app.
	AppMetadata *AppMetadata `json:"app_metadata,omitempty" yaml:"app_metadata,omitempty"`

	// Services holds definitions of the services used to compose your application.
	// and enough metadata to generated appropriate Systemd and Luanchd configurations.
	Services []*Service `json:"services,omitempty" yaml:"services,omitempty"`

	// Models holds a list of data models. It is used by
	// both the data router and code generator.
	Models []*models.Model `json:"models,omitempty" yaml:"models,omitempty"`

	// Routes holds an array of maps of route definitions used by
	// the data router and code generator
	Routes []*Route `json:"routes,omitempty" yaml:"routes,omitempty"`

	// Templates holds an array of maps the request to template to request for
	// Newt (Handlebars) template engine
	Templates []*Template `json:"templates,omitempty" yaml:"templates,omitempty"`

	// isChanged is a convience variable for tracking if the data structure has changed.
	isChanged bool `json:"-" yaml:"-"`
}

// AppMetadata holds metadata about your Newt Service
// This is primarily used in generated Handlbars partials
type AppMetadata struct {
	AppName string `json:"name,omitempty" yaml:"app_name,omitempty"`
	AppTitle string `json:"title,omitempty" yaml:"app_title,omitempty"`
	CopyrightYear string `json:"copyright_year,omitempty" yaml:"copyright_year,omitempty"`
	CopyrightLink string `json:"copyright_link,omitempty" yaml:"copyright_link,omitempty"`
	CopyrightText string `json:"copyright_text,omitempty" yaml:"copyright_text,omitempty"`
	LogoLink string `json:"logo_link,omitempty" yaml:"logo_link,omitempty"`
	LogoText string `json:"logo_text,omitempty" yaml:"logo_text,omitempty"`
	LicenseLink string `json:"license_link,omitempty" yaml:"license_link,omitempty"`
	LicenseText string `json:"license_text,omitempty" yaml:"license_text,omitempty"`
	CSSPath string `json:"css_path,omitempty" yaml:"css_path,omitempty"`
	HeaderLink string `json:"header_link,omitempty" yaml:"header_link,omitempty"`
	HeaderText string `json:"header_text,omitempty" yaml:"header_text,omitempty"`
	ContactAddress string `json:"contact_address,omitempty" yaml:"contact_address,omitempty"`
	ContactPhone string `json:"contact_phone,omitempty" yaml:"contact_phone,omitempty"`
	ContactEMail string `json:"contact_email,omitempty" yaml:"contact_email,omitempty"`
}

/** DEPRECIATED: This is being removed because it causes a rewrite when the optional applications change.
// Services holds the runtime information for newt router, generator,
// template engine.
type Services struct {
	// Newt Router runtime config
	Router *Service `json:"router,omitempty" yaml:"router,omitempty"`

	// TemplateEngine holds Handlebars runtime configuration for Newt template engine
	TemplateEngine *Service `json:"template_engine,omitempty" yaml:"template_engine,omitempty"`

	// Dataset runtime config
	Datasetd *Service `json:"dataset,omitempty" yaml:"dataset,omitempty"`

	// Postgres runtime config, e.g. port number to use for connecting.
	Postgres *Service `json:"postgres,omitempty" yaml:"postgres,omitempty"`

	// PostgREST runtime config
	PostgREST *Service `json:"postgrest,omitempty" yaml:"postgrest,omitempty"`

	// Environment holds a list of OS environment variables that can be made
	// available to the web services.
	Environment []string `json:"environment,omitempty" yaml:"enviroment,omitempty"`

	// Options is a map of name to string values, it is where the
	// the environment variable valuess are stored.
	Options map[string]interface{} `json:"options,omitempty" yaml:"options,omitempty"`
}
*/

// Service implements runtime config for Newt and off the shelf programs used to compose
// your Newt based application.
type Service struct {
	// AppName holds the name of the application, e.g. Postgres, PostgREST
	AppName string `josn:"name,omitempty" yaml:"name,omitempty"`

	// AppPath holds the path to the binary application, e.g. PostgREST
	// This property provides the location of the service to run.
	AppPath string `json:"path,omitempty" yaml:"path,omitempty"`

	// ConfPath holds teh path to the configuration file (e.g. PostgREST configuration file)
	ConfPath string `json:"conf_path,omitempty" yaml:"conf_path,omitempty"`

	// Namespace holds the Postgres Schema name It is used to generate
	// a setup.sql file using the -pg-setup option in newt cli.
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`

	// CName is the name of the dataset collection you wish to use/generate.
	CName string `json:"c_name,omitempty" yaml:"c_name,omitempty"`

	// Port is the name of the localhost port Newt will listen on.
	Port int `json:"port,omitempty" yaml:"port,omitempty"`

	// Timeout is a duration, it is used to set timeouts and the application.
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`

	// Htdocs holds any static files you want to make available through
	// Newt router.
	Htdocs string `json:"htdocs,omitempty" yaml:"htdocs,omitempty"`

	// BaseDir is used by Handlebars, usually holds the "views" directory.
	BaseDir string `json:"base_dir,omitempty" yaml:"base_dir,omitempty"`

	// ExtName is used by Handlebars to set the expected extension (e.g. ".hbs")
	ExtName string `json:"ext_name,omitempty" yaml:"ext_name,omitempty"`

	// PartialsDir is used by Handlebars to find partial templates, usually inside the views directory
	PartialsDir string `json:"partials_dir,omitempty" yaml:"partials_dir,omitempty"`

	// DSN, data ast name is a URI connection string
	DSN string `json:"dsn,omitemity" yaml:"dsn,omitempty"`

	// Environment holds a list of OS environment variables that can be made
	// available to the web services.
	Environment []string `json:"environment,omitempty" yaml:"enviroment,omitempty"`

	// Options is a map of name to string values, it is where the
	// the environment variable valuess are stored.
	Options map[string]interface{} `json:"options,omitempty" yaml:"options,omitempty"`
}

// NewServices generates a default set of applications for your Newt project.
func NewServices() []*Service {
	var applications []*Service
	for _, appName := range []string{"router", "template_engine", "postgres", "postgrest"} {
		//FIXME: Postgres supports specific environment variables, these should be automatically included
		app := &Service{
			AppName: appName,
		}
		applications = append(applications, app)
	}
	return applications
}


// NewAST will create an empty AST with top level attributes
func NewAST() *AST {
	ast := new(AST)
	ast.Services = NewServices()
	return ast
}

// GetService takes a list of applications, `[]*Service`, and returns the application name in the list or nil.
func (ast *AST) GetService(appName string) *Service {
	if ast.Services != nil {
		for _, app := range ast.Services {
			if app.AppName == appName {
				return app
			}
		}
	}
	return nil
}

// RemoveService takes a list of applications, `[]*Service`, and remove the target item.
func (ast *AST) RemoveService(appName string) error {
	if ast.Services != nil {
		// Find the position of the application in list
		for pos, app := range ast.Services {
			if app.AppName == appName {
				ast.Services = append(ast.Services[:pos], ast.Services[pos+1:]...)
				return nil
			}
		}
	}
	return fmt.Errorf("could not remove %q, not found", appName)
}

// UnmarshalAST will read []byte of YAML or JSON,
// populate the provided *AST object and return an error.
//
// ```
// src, _ := os.ReadFile("app.yaml")
// ast := new(AST)
//
//	if err := UnmarshalAST(src, ast); err != nil {
//	    // ... handle error
//	}
//
// ```
func UnmarshalAST(src []byte, ast *AST) error {
	if bytes.HasPrefix(src, []byte("{")) {
		if err := json.Unmarshal(src, &ast); err != nil {
			return err
		}
	} else {
		if err := yaml.Unmarshal(src, &ast); err != nil {
			return err
		}
	}
	if ast.Services == nil {
		ast.Services = NewServices()
	}
	return nil
}

// LoadAST read a YAML file, merges environment variables
// and returns a AST object and error value.
//
// ```
// ast, err := LoadAST("app.yaml")
//
//	if err != nil {
//	    // ... handle error
//	}
//
// ```
func LoadAST(configFName string) (*AST, error) {
	ast := new(AST)
	if configFName != "" {
		src, err := os.ReadFile(configFName)
		if err != nil {
			return nil, fmt.Errorf("failed to read %q, %s", configFName, err)
		}
		if err := UnmarshalAST(src, ast); err != nil {
			return nil, fmt.Errorf("failed to read %q, %s", configFName, err)
		}
	}

	if ast.Services == nil {
		ast.Services = NewServices()
	}
	// Load environment if missing from config file.
	for _, app := range ast.Services {
		for _, envar := range app.Environment {
			// YAML settings take presidence over environment, check for conflicts
			if _, conflict := app.Options[envar]; !conflict {
				app.Options[envar] = os.Getenv(envar)
			}
		}
	}
	ast.isChanged = false
	return ast, nil
}

func (ast *AST) HasChanges() bool {
	if ast.isChanged {
		return true
	}
	for _, m := range ast.Models {
		if m.HasChanges() {
			return true
		}
	}
	return false
}

func (ast *AST) Encode() ([]byte, error) {
	// Now output the YAML
	timeStamp := (time.Now()).Format("2006-01-02")
	userName := os.Getenv("USER")
	comment := []byte(fmt.Sprintf(`#!/usr/bin/env newt check
#
# This was generated by %s on %s with %s version %s %s.
#
`, userName, timeStamp, path.Base(os.Args[0]), Version, ReleaseHash))
	data := bytes.NewBuffer(comment)
	encoder := yaml.NewEncoder(data)
	encoder.SetIndent(2)
	if err := encoder.Encode(ast); err != nil {
		return nil, fmt.Errorf("failed to generate configuration, %s\n", err)
	}
	return data.Bytes(), nil
}

// SaveAs writes the *AST to a YAML file.
func (ast *AST) SaveAs(configName string) error {
	if _, err := os.Stat(configName); err == nil {
		if err := backupFile(configName); err != nil {
			return err
		}
	}
	fp, err := os.Create(configName)
	if err != nil {
		return err
	}
	defer fp.Close()
	src, err := ast.Encode()
	if err != nil {
		return err
	}
		fmt.Fprintf(fp, "%s", src)
	for _, model := range ast.Models {
		for _, element := range model.Elements {
			element.Changed(false)
		}
		model.Changed(false)
	}
	ast.isChanged = false
	return nil
}

// GetModelIds returns a list of model ids
func (ast *AST) GetModelIds() []string {
	if ast.Models == nil {
		ast.Models = []*models.Model{}
	}
	ids := []string{}
	for _, m := range ast.Models {
		if m.Id != "" {
			ids = append(ids, m.Id)
		}
	}
	return ids
}

// GetModelById return a specific model by it's id
func (ast *AST) GetModelById(id string) (*models.Model, bool) {
	for _, m := range ast.Models {
		if m.Id == id {
			return m, true
		}
	}
	return nil, false
}

// AddModel takes a new Model, checks if the model exists in the list (i.e.
// has an existing model id that matches the new model and if not appends
// it so the list.
func (ast *AST) AddModel(model *models.Model) error {
	// Make sure we have a Models lists to work with.
	if ast.Models == nil {
		ast.Models = []*models.Model{}
	}
	// Check to see if this is a duplicate, return error if it is
	for i, m := range ast.Models {
		if m.Id == model.Id {
			return fmt.Errorf("failed, model %d is a duplicate model id, %q", i, m.Id)
		}
	}
	ast.Models = append(ast.Models, model)
	ast.isChanged = true
	return nil
}

// UpdateModel takes a model id and new model struct replacing the
// existing one.
func (ast *AST) UpdateModel(id string, model *models.Model) error {
	// Make sure we have a Models lists to work with.
	if ast.Models == nil {
		return fmt.Errorf("no models defined")
	}
	for i, m := range ast.Models {
		if m.Id == id {
			ast.Models[i] = model
			ast.isChanged = true
			return nil
		}
	}
	return fmt.Errorf("failed to find model %q", id)
}

// RemoveModelById find the model with the model id and remove it
func (ast *AST) RemoveModelById(id string) error {
	// Make sure we have a Models lists to work with.
	if ast.Models == nil {
		return fmt.Errorf("no models defined")
	}
	for i, m := range ast.Models {
		if m.Id == id {
			ast.Models = append(ast.Models[:i], ast.Models[(i+1):]...)
			ast.isChanged = true
			return nil
		}
	}
	return fmt.Errorf("failed to find model %q", id)
}

// RemoveRouteById find the route with route id and remove it
func (ast *AST) RemoveRouteById(id string) error {
	routeFound := false
	for i, r := range ast.Routes {
		// NOTE: A route id ties one or more requests together, e.g. retrieve a web form (GET), then handle it (POST)
		if r.Id == id {
			ast.Routes = append(ast.Routes[:i], ast.Routes[(i+1):]...)
			ast.isChanged = true
			routeFound = true
		}
	}
	if !routeFound {
		return fmt.Errorf("failed to find route %s", id)
	}
	return nil
}

// RemoveTemplateById() find the template id and remove it from the .Templates structure
func (ast *AST) RemoveTemplateById(id string) error {
	templateFound := false
	for i, t := range ast.Templates {
		if t.Id == id {
			ast.Templates = append(ast.Templates[:i], ast.Templates[(i+1):]...)
			ast.isChanged = true
			templateFound = true
		}
	}
	if !templateFound {
		return fmt.Errorf("failed to find template %s", id)
	}
	return nil
}

// GetRouteIds returns a list of Router ids found in ast.Routes
func (ast *AST) GetRouteIds() []string {
	rIds := []string{}
	for _, r := range ast.Routes {
		if r.Id != "" {
			rIds = append(rIds, r.Id)
		}
	}
	return rIds
}

// GetTemplateIds return a list of template ids.
func (ast *AST) GetTemplateIds() []string {
	tIds := []string{}
	for _, t := range ast.Templates {
		if t.Id != "" {
			tIds = append(tIds, t.Id)
		}
	}
	return tIds
}

// GetPrimaryTemplates return a list of primary template filenames
func (ast *AST) GetPrimaryTemplates() []string {
	fNames := []string{}
	for _, t := range ast.Templates {
		if t.Template != "" {
			fNames = append(fNames, t.Template)
		}
	}
	return fNames
}

// GetAllTemplates returns a list of templates, including partials defined
// in the .Templates property. Part template names are indented with a "\t"
func (ast *AST) GetAllTemplates() []string {
	fNames := []string{}
	for _, t := range ast.Templates {
		if t.Template != "" {
			fNames = append(fNames, t.Template)
		}
	}
	return fNames
}

// GetTemplateByPrimary returns the template entry using primary template filename
func (ast *AST) GetTemplateByPrimary(fName string) (*Template, bool) {
	if ast.Templates != nil {
		for _, t := range ast.Templates {
			if t.Template == fName {
				return t, true
			}
		}
	}
	return nil, false
}

// Check reviews the ast *AST and reports and issues, return true
// if no errors found and false otherwise.  The "buf" will hold the error output.
func (ast *AST) Check(buf io.Writer) bool {
	ok := true
	if ast.Services == nil {
		fmt.Fprintf(buf, "no applications defined\n")
		ok = false
	}
	postgres := ast.GetService("postgres")
	datasetd := ast.GetService("datasetd")
	router := ast.GetService("router")
	templateEngine := ast.GetService("template_engine")
	if postgres != nil || datasetd != nil {
		if ast.Models == nil || len(ast.Models) == 0 {
			fmt.Fprintf(buf, "no models defined for applications\n")
			ok = false
		} else {
			for i, m := range ast.Models {
				if !m.Check(buf) {
					fmt.Fprintf(buf, "model #%d is invalid\n", i)
					ok = false
				}
			}
		}
	}

	if router != nil {
		if ast.Routes == nil || len(ast.Routes) == 0 {
			fmt.Fprintf(buf, "no routes defined for Newt Router\n")
			ok = false
		}
		if router.Port == 0 {
			fmt.Fprintf(buf, "application.router.port not set\n")
			ok = false
		}
		for i, r := range ast.Routes {
			if !r.Check(buf) {
				fmt.Fprintf(buf, "route (#%d) errors\n", i)
				ok = false
			}
		}
	}

	if templateEngine != nil {
		if ast.Templates == nil || len(ast.Templates) == 0 {
			fmt.Fprintf(buf, "template engine is defined but not templates are configured\n")
			ok = false
		} else {
			t, err := NewTemplateEngine(ast)
			if err != nil {
				fmt.Fprintf(buf, fmt.Sprintf("application.template_engine not configured, %s\n", err))
				ok = false
			} else if !t.Check(buf) {
				ok = false
			}
		}
	}
	return ok
}

// TemplateEngine defines the `nte` application YAML file. It joins some of the Service struct
// with an array of templates so that "check" can validate the YAML.
type TemplateEngine struct {
	// Port is the name of the localhost port Newt will listen on.
	Port int `json:"port,omitempty" yaml:"port,omitempty"`

	// BaseDir is holds the "views" for that are formed from the templates.
	BaseDir string `json:"base_dir,omitempty" yaml:"base_dir,omitempty"`

	// ExtName is used to set the expected extension (e.g. ".hbs")
	ExtName string `json:"ext_name,omitempty" yaml:"ext_name,omitempty"`

	// PartialsDir is used to find partial templates, usually inside the views directory
	PartialsDir string `json:"partials_dir,omitempty" yaml:"partials_dir,omitempty"`

	// Timeout is a duration, it is used to set timeouts and the application.
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`

	// Templates defined for the service
	Templates []*Template `json:"templates,omitempty" yaml:"templates,omitempty"`
}

// Template hold the request to template mapping for in the TemplateEngine
type Template struct {
	// Id ties a set of one or more template together, e.g. a web form and its response
	Id string `json:"id,required" yaml:"id,omitempty"`

	// Description describes the purpose of the tempalte mapping. It is used to debug Newt YAML files.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Pattern holds a request path, e.g. `/blog_post`. NOTE: the method is ignored. A POST
	// is presumed to hold data that will be processed by the template engine. A GET retrieves the
	// unresolved template.
	Pattern string `json:"request,required" yaml:"request,omitempty"`

	// Template holds a path to the primary template (aka view) file for this route. Path can be relative
	// to the current working directory.
	Template string `json:"template,required" yaml:"template,omitempty"`

	// Debug logs more verbosely if true
	Debug bool `json:"debug,omitempty" yaml:"debug,omitempty"`

	// Document hold the a map of values passed into it from the Newt YAML file in the applications
	// property. These are a way to map in environment or application wide values. These are exposed in
	// the Newt template engine `options`.
	Document map[string]interface{} `json:"document,omitempty" yaml:"document,omitempty"`

	// Vars holds the names of any variables expressed in the pattern, these an be used to replace elements of
	// the output object.
	Vars []string `json:"-" yaml:"-"`

	// Body holds a map of data to process with the template
	Body map[string]interface{} `json:"-" yaml:"-"`

	// The follow are used to simplify individual template invocation.
	// They are populated from the TemplateEngine object.

	/*FIXME: I want to support both Mustache and Handlebars templates.

			 I need to review both mustache and handlebars implementations so figure out an appropriate
			 or wrapper then Tmpl should point to that interface. Gofiber does this with "template views"
			 but I don't want to have to pull in Gofiber's template engine, it's big and provides too many
			 choices to implement smoothly.

	         I need to decide how to specify the template language and if that is per template or engine wide.
			 One approach would be to pick the template engine based on the file extension. Another approoach would
			 be to make it a property of the engine that inherits like the BaseDir, ExtName, etc.
	*/

	// Tmpl points to the compied template
	Tmpl *raymond.Template `json:"-" yaml:"-"`

	// BaseDir is used by holds the "views" directory.
	BaseDir string `json:"-" yaml:"-"`

	// ExtName is used by set the expected extension (e.g. ".hbs")
	ExtName string `json:"-" yaml:"-"`

	// Partials holds partials directory
	PartialsDir string `json:"-" yaml:"-"`
}

// NewTemplateEngine create a new TemplateEngine struct. If a filename
// is provided it reads the file and sets things up accordingly.
func NewTemplateEngine(ast *AST) (*TemplateEngine, error) {
	templateEngine := ast.GetService("template_engine")
	if templateEngine == nil {
		return nil, fmt.Errorf("template engine is nil")
	}

	// Copy our options so we can expose them in the template's .document
	docvars := map[string]interface{}{}
	// Copy in options to envars
	if templateEngine.Options != nil && len(templateEngine.Options) > 0 {
		for k, v := range templateEngine.Options {
			docvars[k] = v
		}
	}
	te := &TemplateEngine{
		Port:        TEMPLATE_ENGINE_PORT,
		ExtName:     TEMPLATE_ENGINE_EXT_NAME,
		BaseDir:     TEMPLATE_ENGINE_BASE_DIR,
		PartialsDir: TEMPLATE_ENGINE_PARTIALS_DIR,
	}
	if templateEngine.Port != 0 {
		te.Port = templateEngine.Port
	}
	if templateEngine.BaseDir != "" {
		te.BaseDir = templateEngine.BaseDir
	}
	if templateEngine.ExtName != "" {
		te.ExtName = templateEngine.ExtName
	}
	if templateEngine.PartialsDir != "" {
		te.PartialsDir = templateEngine.PartialsDir
	}
	// FIXME: Need to copy in environment variables from ast.Options and set t.Document content.
	errMsgs := []string{}
	if ast.Templates != nil && len(ast.Templates) > 0 {
		// Map in the BaseDir, PartialsDir, and ExtName for the templates.
		for _, t := range ast.Templates {
			t.ExtName = te.ExtName
			t.BaseDir = te.BaseDir
			t.PartialsDir = te.PartialsDir
			if t.Document == nil {
				t.Document = map[string]interface{}{}
			}
			// Copy in options to .document
			if len(docvars) > 0 {
				for k, v := range docvars {
					t.Document[k] = v
				}
			}
		}
		// Add the resulting templates into struct.
		te.Templates = append([]*Template{}, ast.Templates...)
	}
	if len(errMsgs) > 0 {
		return te, fmt.Errorf("%s", strings.Join(errMsgs, "\n"))
	}
	return te, nil
}

// Check makes sure the TemplateEngine struct is populated
func (tEng *TemplateEngine) Check(buf io.Writer) bool {
	if tEng == nil {
		fmt.Fprintf(buf, "template engine not defined\n")
		return false
	}
	errMsgs := []string{}
	ok := true
	if tEng.Port == 0 {
		errMsgs = append(errMsgs, "template engine port not set")
		ok = false
	} else {
		errMsgs = append(errMsgs, fmt.Sprintf("template engine will listen on port %d", tEng.Port))
	}
	if tEng.BaseDir == "" {
		errMsgs = append(errMsgs, "base directory not set for templates")
		ok = false
	}
	if tEng.ExtName == "" {
		errMsgs = append(errMsgs, "template extension is not set")
		ok = false
	}
	if tEng.Templates == nil || len(tEng.Templates) == 0 {
		errMsgs = append(errMsgs, "no templates found")
		ok = false
	} else {
		errMsgs = append(errMsgs, fmt.Sprintf("templates are located in %q", tEng.BaseDir))
		if tEng.PartialsDir != "" {
			errMsgs = append(errMsgs, fmt.Sprintf("partials are located in %q", path.Join(tEng.BaseDir, tEng.PartialsDir)))
		}
		if tEng.ExtName != "" {
			errMsgs = append(errMsgs, fmt.Sprintf("template extension is set to %q", tEng.ExtName))
		}
		//FIXME: add check for helpers
		errMsgs = append(errMsgs, fmt.Sprintf("%d template path(s) mapped", len(tEng.Templates)))
		for i, t := range tEng.Templates {
			tBuf := bytes.NewBuffer([]byte{})
			if !t.Check(tBuf) {
				errMsgs = append(errMsgs, fmt.Sprintf("template (#%d) failed check, %s\n", i, tBuf.Bytes()))
				ok = false
			}
		}
	}
	fmt.Fprintf(buf, "%s\n", strings.Join(errMsgs, "\n"))
	return ok
}

// Check evaluates the *Template and outputs finding. Returns true of no error, false if errors found
func (tmpl *Template) Check(buf io.Writer) bool {
	ok := true
	if tmpl == nil {
		fmt.Fprintf(buf, "template is nil\n")
		return false
	}
	if tmpl.Pattern == "" {
		fmt.Fprintf(buf, "template does not have an associated path/pattern\n")
		ok = false
	}
	if tmpl.Template == "" {
		fmt.Fprintf(buf, "missing path to template for %s\n", tmpl.Pattern)
		ok = false
	} else {
		fmt.Fprintf(buf, "template name %s\n", tmpl.Template)
	}
	return ok
}
```

What would this look like in TypeScript?

---

Replace `import * as yaml from "https://deno.land/std/encoding/yaml.ts";` with `import * as yaml from "@std/yaml";`.  Replace `import * as path from "https://deno.land/std/path/mod.ts";` with `import * as path from "@std/path";`.

---

Instead of importing `import { readFileStr, writeFileStr } from "https://deno.land/std/fs/mod.ts";` use `Deno.readTextFile` and `Deno.writeTextFile` instead.

---

`Route` and `Model` interfaces will be imported using

```typescript
import { Route } from "./route.ts";
import { Model } from "./model.ts";
```

---

You would add the following import statement for version and license content.

```
import { version, releaseDate, releaseHash, licenseText } from "./version.ts";
```

---

The line 

```
# This was generated by ${userName} on ${timeStamp} with newt version ${version}.
```

Should be

```
# This was generated by ${userName} on ${timeStamp} with newt version ${version} ${releaseHash}.
```

---

Add the following import statement. 

```
import { Element } from "./element.ts"
```

The defines the `Element` type of `model.element`.

---

Change `throw new Error(`Failed to unmarshal AST: ${err.message}`);` to
`throw new Error(`Failed to unmarshal AST: ${err}`);`.

---

`Template` should be a class.

`AppMetadata` should be a class.


---

`errors` are type `string[]`. They an attribute in each class.  When an error is detected an error string is pushed into the error array.

The `check` method should return a boolean value of `true` if the length of the `errors` array is greater than zero. 

---

In the `Template` class `id: string;` should be `id: string = "";`,
`pattern: string;` should be `pattern: string = "";` and the
`template: string;` should be `template: string = "";`.

---

`Service` should be a class.

---

In `Service` class the `options` attribute should be initialized like this `options: {[key: string]: any } = {};`.

---

The following code looks wrong, you can't map the function `check`.

```typescript
  te.templates = ast.templates.map(t => ({
    ...t,
    baseDir: te.baseDir,
    extName: te.extName,
    partialsDir: te.partialsDir,
    document: { ...templateEngine.options, ...t.document }
  }));
```

---

Now let's port `ast_test.go` to TypeScript running under Deno 2.2. You should use `Deno.test` and import assertions from `@std/assert`.  This is the Go code to port.

```go
package newt

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"testing"
)

// TestNewtAST() tests NewAST() and NewApplications()
func TestNewAST(t *testing.T) {
	ast := NewAST()
	if ast.Applications == nil {
		t.Errorf("ast.Applications should not be nil")
	}
}

// TestUnmarshalAST tests unmarshalling YAML into a Newt AST object
func TestUnmarshalAST(t *testing.T) {
	configFiles := []string{
		path.Join("testdata", "birds.yaml"),
		path.Join("testdata", "blog.yaml"),
		path.Join("testdata", "bundler_test.yaml"),
	}
	for _, fName := range configFiles {
		src, err := os.ReadFile(fName)
		if err != nil {
			t.Errorf("failed to read %q, %s", fName, err)
		} else {
			ast := new(AST)
			if err := UnmarshalAST(src, ast); err != nil {
				t.Errorf("failed tn UnmarshalAST %q, %s", fName, err)
			} else {
				buf := bytes.NewBuffer([]byte{})
				if ok := ast.Check(buf); !ok {
					t.Errorf("UnmarshalAST %q, failed to pass check -> %s", fName, buf.Bytes())
				}
			}
		}
	}

}

// TestLoadAST tests reading on and populating the shared YAML configuration used
// by Newt applications.
func TestLoadAST(t *testing.T) {
	configFiles := []string{
		path.Join("testdata", "birds.yaml"),
		path.Join("testdata", "blog.yaml"),
		path.Join("testdata", "bundler_test.yaml"),
	}
	for _, fName := range configFiles {
		ast, err := LoadAST(fName)
		if err != nil {
			t.Errorf("failed to load %q, %s", fName, err)
		}
		if ast == nil {
			t.Errorf("something went wrong, ast is nil for %q", fName)
		}
		if ast.Applications == nil {
			t.Errorf("ast.Applications is nil (%q), %+v", fName, ast)
		}
		ids := ast.GetModelIds()
		if len(ids) == 0 {
			t.Errorf("expected model ids for %q", fName)
		} else {
			mId := ids[0]
			model, ok := ast.GetModelById(mId)
			if !ok {
				t.Errorf("expected model for %q in %q, %s", mId, fName, err)
			}
			if model == nil {
				t.Errorf("expceted model content for %q in %q, got nil", mId, fName)
			}
		}
	}
}

// TestHasChanges reads in our test YAML files, checks for changes, then modifies them and checkes for changes again.
func TestHasChanges(t *testing.T) {
	configFiles := []string{
		path.Join("testdata", "birds.yaml"),
		path.Join("testdata", "blog.yaml"),
		path.Join("testdata", "bundler_test.yaml"),
	}
	for i, fName := range configFiles {
		ast, err := LoadAST(fName)
		if err != nil {
			t.Errorf("failed to load %q, %s", fName, err)
		}
		if ast.HasChanges() {
			t.Errorf("should not have changes after LoadAST(%q)", fName)
		}
		whatChanged := ""
		switch i {
		case 0:
			ast.Applications = nil
			ast.isChanged = true
			whatChanged = fmt.Sprintf("removed .applications from ast for %q", fName)
		case 1:
			modelList := ast.GetModelIds()
			modelId := modelList[0]
			whatChanged = fmt.Sprintf("removed model %q from %q", modelId, fName)
			ast.RemoveModelById(modelId)
		case 2:
			t := ast.Templates[0]
			ast.Templates = ast.Templates[1:]
			ast.isChanged = true
			whatChanged = fmt.Sprintf("removed template %q -> %q for %q", t.Pattern, t.Template, fName)
		}
		if !ast.HasChanges() {
			t.Errorf("%s, did not detect change for %q", whatChanged, fName)
		}
	}
}
```

---

The classes in `ast.ts` need to be exported so we can test them.

---

In `ast_test.ts` you are trying to set a private attribute `isChanged` in your test. 

---

Let's port `element.go` to `element.ts`. Here is the Go source code.

```go
import (
	"fmt"
	"io"
	"strings"
)

// GenElementFunc is a function which will generate an Element configured represent a model's supported "types"
type GenElementFunc func() *Element

// ValidateFunc is a function that validates form assocaited with the Element and the string value
// received in the web form (value before converting to Go type).
type ValidateFunc func(*Element, string) bool


// Element implementes the GitHub YAML issue template syntax for an input element.
// The input element YAML is described at <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-githubs-form-schema>
//
// While the syntax most closely express how to setup an HTML representation it is equally
// suitable to expressing, through inference, SQL column type definitions. E.g. a bare `input` type is a `varchar`,
// a `textarea` is a `text` column type, an `input[type=date]` is a date column type.
type Element struct {
	// Type, The type of element that you want to input. It is required. Valid values are
	// checkboxes, dropdown, input, markdown and text area.
	//
	// The input type corresponds to the native input types defined for HTML 5. E.g. text, textarea,
	// email, phone, date, url, checkbox, radio, button, submit, cancel, select
	// See MDN developer docs for input, <https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input>
	Type string `json:"type,required" yaml:"type,omitempty"`

	// Id for the element, except when type is set to markdown. Can only use alpha-numeric characters,
	//  -, and _. Must be unique in the form definition. If provided, the id is the canonical identifier
	//  for the field in URL query parameter prefills.
	Id string `json:"id,omitempty" yaml:"id,omitempty"`

	// Attributes, a set of key-value pairs that define the properties of the element.
	// This is a required element as it holds the "value" attribute when expressing
	// HTML content. Other commonly use attributes
	Attributes map[string]string `json:"attributes,omitempty" yaml:"attributes,omitempty"`

	// Pattern holds a validation pattern. When combined with an input type (or input type alias, e.g. orcid)
	// produces a form element that sports a specific client side validation exceptation.  This intern can be used
	// to generate appropriate validation code server side.
	Pattern string `json:"pattern,omitempty" yaml:"pattern,omitempty"`

	// Options holds a list of values and their labels used for HTML select elements in rendering their option child elements
	Options []map[string]string `json:"optoins,omitempty" yaml:"options,omitempty"`

	// IsObjectId (i.e. is the identifier of the object) used by for the modeled data.
	// It is used in calculating routes and templates where the object identifier is required.
	IsObjectId bool `json:"is_primary_id,omitempty" yaml:"is_primary_id,omitempty"`

	// Generator indicates the type of automatic population of a field. It is used to
	// indicate autoincrement and uuids for primary keys and timestamps for datetime oriented fields.
	Generator string `json:"generator,omitempty" yaml:"generator,omitempty"`

	// Label is used when rendering an HTML form as a label element tied to the input element via the set attribute and
	// the element's id.
	Label string `json:"label,omitempty" yaml:"label,omitempty"`

	//
	// These fields are used by the modeler to manage the models and their elements
	//
	isChanged bool `json:"-" yaml:"-"`
}

// NewElement, makes sure element id is valid, populates an element as a basic input type.
// The new element has the attribute "name" and label set to default values.
func NewElement(elementId string) (*Element, error) {
	if !IsValidVarname(elementId) {
		return nil, fmt.Errorf("invalid element id, %q", elementId)
	}
	element := new(Element)
	element.Id = elementId
	element.Attributes = map[string]string{"name": elementId}
	element.Type = "text"
	element.Label = strings.ToUpper(elementId[0:1]) + elementId[1:]
	element.IsObjectId = false
	element.isChanged = true
	return element, nil
}

// HasChanged checks to see if the Element has been changed.
func (e *Element) HasChanged() bool {
	return e.isChanged
}

// Changed sets the change state on element
func (e *Element) Changed(state bool) {
	e.isChanged = state
}

// Check reviews an Element to make sure if is value.
func (e *Element) Check(buf io.Writer) bool {
	ok := true
	if e == nil {
		fmt.Fprintf(buf, "element is nil\n")
		ok = false
	}
	if e.Id == "" {
		fmt.Fprintf(buf, "element missing id\n")
		ok = false
	}
	if e.Type == "" {
		fmt.Fprintf(buf, "element, %q, missing type\n", e.Id)
		ok = false
	}
	return ok
}
```

---

In class `Element` to attribute need initializers.

---

Please create `element_test.ts`. Use `Deno.test` and assertions from `@std/assert`.

---

In `element_test.ts` change when you assert that elements are equal and they fail the assertion should the value they should and the value that failed the assertion.

---

In `element.ts` the `newElement` function should use a default label value should be set to the value of `elementId` where the first letter is capitalized.

---

In `element_test.ts` the line 

```
  assertEquals(element?.label, "Testelement", "Element label should be 'Testelement' but got " + element?.label);
```

should be

```
  assertEquals(element?.label, "TestElement", "Element label should be 'TestElement' but got " + element?.label);
```

---

I want to port `model.go` to `model.ts` written in TypeScript to run in Deno 2.2.  This is the Go source code.

```go
package models

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
)

// RenderFunc is a function thation takes an io.Writer and Model then
// renders the model into the io.Writer. It is used to extend the Model to
// support various output formats.
type RenderFunc func(io.Writer, *Model) error

// Model implements a data structure description inspired by GitHub YAML issue template syntax.
// See <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-issue-forms>
//
// The Model structure describes the HTML elements used to form a record. It can be used in code generation and in validating
// POST and PUT requests in datasetd.
type Model struct {
	// Id is a required field for model, it maps to the HTML element id and name
	Id string `json:"id,required" yaml:"id"`

	// This is a Newt specifc set of attributes to place in the form element of HTML. I.e. it could
	// be form "class", "method", "action", "encoding". It is not defined in the GitHub YAML issue template syntax
	// (optional)
	Attributes map[string]string `json:"attributes,omitempty" yaml:"attributes,omitempty"`

	// Description, A description for the issue form template, which appears in the template chooser interface.
	// (required)
	Description string `json:"description,required" yaml:"description,omitempty"`

	// Elements, Definition of the input types in the form.
	// (required)
	Elements []*Element `json:"elements,required" yaml:"elements,omitempty"`

	// Title, A default title that will be pre-populated in the issue submission form.
	// (optional) only there for compatibility with GitHub YAML Issue Templates
	//Title string `json:"title,omitempty" yaml:"title,omitempty"`

	// isChanged is an internal state used by the modeler to know when a model has changed
	isChanged bool `json:"-" yaml:"-"`

	// renderer is a map of names to RenderFunc functions. A RenderFunc is that take a io.Writer and the model object as parameters then
	// return an error type.  This allows for many renderers to be used with Model by
	// registering the function then envoking render with the name registered.
	renderer map[string]RenderFunc `json:"-" yaml:"-"`

	// genElements holds a map to the "type" pointing to an element generator
	genElements map[string]GenElementFunc `json:"-" yaml:"-"`

	// validators holds a list of validate function associated with types. Key is type name.
	validators map[string]ValidateFunc `json:"-" yaml:"-"`
}

// GenElementType takes an element type and returns an Element struct populated for that type and true or nil and false if type is not supported.
func (model *Model) GenElementType(typeName string) (*Element, bool) {
	if fn, ok := model.genElements[typeName]; ok {
		return fn(), true
	}
	return nil, false
}

// Validate form data expressed as map[string]string.
func (model *Model) Validate(formData map[string]string) bool {
	ids := model.GetElementIds()
	if len(ids) != len(formData) {
		return false
	}
	for k, v := range formData {
		if elem, ok := model.GetElementById(k); ok {
			if validator, ok := model.validators[elem.Type]; ok {
				if !validator(elem, v) {
					if Debug {
						log.Printf("DEBUG failed to validate elem.Id %q, elem.Type %q, value %q", elem.Id, elem.Type, v)
					}
					return false
				}
			} else {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

// ValidateMapInterface normalizes the map inteface values before calling
// the element's validator function.
func (model *Model) ValidateMapInterface(data map[string]interface{}) bool {
	if model == nil {
		if Debug {
			log.Printf("model is nil, can't validate")
		}
		return false
	}
	ids := model.GetElementIds()
	if len(ids) != len(data) {
		if Debug {
			log.Printf("DEBUG expected len(ids) %d, got len(data) %d", len(ids), len(data))
		}
		return false
	}
	for k, v := range data {
		var val string
		switch v.(type) {
		case string:
			val = v.(string)
		case int:
			val = fmt.Sprintf("%d", v)
		case float64:
			val = fmt.Sprintf("%f", v)
		case json.Number:
			val = fmt.Sprintf("%s", v)
		case bool:
			val = fmt.Sprintf("%t", v)
		default:
			val = fmt.Sprintf("%+v", v)
		}
		if elem, ok := model.GetElementById(k); ok {
			if validator, ok := model.validators[elem.Type]; ok {
				if !validator(elem, val) {
					if Debug {
						log.Printf("DEBUG failed to validate elem.Id %q, value %q", elem.Id, val)
					}
					return false
				}
			} else {
				if Debug {
					log.Printf("DEBUG failed to validate elem.Id %q, value %q, missing validator", elem.Id, val)
				}
				return false
			}
		} else {
			if Debug {
				log.Printf("DEBUG failed to validate elem.Id %q, value %q, missing elemnt %q", elem.Id, val, k)
			}
			return false
		}
	}
	return true
}

// HasChanges checks if the model's elements have changed
func (model *Model) HasChanges() bool {
	if model.isChanged {
		return true
	}
	for _, e := range model.Elements {
		if e.isChanged {
			return true
		}
	}
	return false
}

// Changed sets the change state
func (model *Model) Changed(state bool) {
	model.isChanged = state
}

// HasElement checks if the model has a given element id
func (model *Model) HasElement(elementId string) bool {
	for _, e := range model.Elements {
		if e.Id == elementId {
			return true
		}
	}
	return false
}

// HasElementType checks if an element type matches given type.
func (model *Model) HasElementType(elementType string) bool {
	for _, e := range model.Elements {
		if strings.ToLower(e.Type) == strings.ToLower(elementType) {
			return true
		}
	}
	return false
}

// GetModelIdentifier() returns the element which describes the model identifier.
// Returns the element and a boolean set to true if found.
func (m *Model) GetModelIdentifier() (*Element, bool) {
	for _, e := range m.Elements {
		if e.IsObjectId {
			return e, true
		}
	}
	return nil, false
}

// GetAttributeIds returns a slice of attribute ids found in the model's .Elements
func (m *Model) GetAttributeIds() []string {
	return getAttributeIds(m.Attributes)
}

// GetElementIds returns a slice of element ids found in the model's .Elements
func (m *Model) GetElementIds() []string {
	ids := []string{}
	for _, elem := range m.Elements {
		if elem.Id != "" {
			ids = append(ids, elem.Id)
		}
	}
	return ids
}

// GetPrimaryId returns the primary id
func (m *Model) GetPrimaryId() string {
	for _, elem := range m.Elements {
		if elem.IsObjectId {
			return elem.Id
		}
	}
	return ""
}

// GetGeneratedTypes returns a map of elemend id and value held by .Generator
func (m *Model) GetGeneratedTypes() map[string]string {
	gt := map[string]string{}
	for _, elem := range m.Elements {
		if elem.Generator != "" {
			gt[elem.Id] = elem.Generator
		}
	}
	return gt
}

// GetElementById returns a Element from the model's .Elements.
func (m *Model) GetElementById(id string) (*Element, bool) {
	for _, elem := range m.Elements {
		if elem.Id == id {
			return elem, true
		}
	}
	return nil, false
}

// NewModel, makes sure model id is valid, populates a Model with the identifier element providing
// returns a *Model and error value.
func NewModel(modelId string) (*Model, error) {
	if !IsValidVarname(modelId) {
		return nil, fmt.Errorf("invalid model id, %q", modelId)
	}
	model := new(Model)
	model.Id = modelId
	model.Description = fmt.Sprintf("... description of %q goes here ...", modelId)
	model.Attributes = map[string]string{}
	model.Elements = []*Element{}
	// Make the required element ...
	element := new(Element)
	element.Id = "id"
	element.IsObjectId = true
	element.Type = "text"
	element.Attributes = map[string]string{"required": "true"}
	if err := model.InsertElement(0, element); err != nil {
		return nil, err
	}
	return model, nil
}

// Check analyze the model and make sure at least one element exists and the
// model has a single identifier (e.g. "identifier")
func (model *Model) Check(buf io.Writer) bool {
	if model == nil {
		fmt.Fprintf(buf, "model is nil\n")
		return false
	}
	if model.Elements == nil {
		fmt.Fprintf(buf, "missing %s.body\n", model.Id)
		return false
	}
	// Check to see if we have at least one element in Elements
	if len(model.Elements) > 0 {
		ok := true
		hasModelId := false
		for i, e := range model.Elements {
			// Check to make sure each element is valid
			if !e.Check(buf) {
				fmt.Fprintf(buf, "error for %s.%s\n", model.Id, e.Id)
				ok = false
			}
			if e.IsObjectId {
				if hasModelId == true {
					fmt.Fprintf(buf, "duplicate model identifier element (%d) %s.%s\n", i, model.Id, e.Id)
					ok = false
				}
				hasModelId = true
			}
		}
		if !hasModelId {
			fmt.Fprintf(buf, "missing required object identifier for model %s\n", model.Id)
			ok = false
		}
		return ok
	}
	fmt.Fprintf(buf, "Missing elements for model %q\n", model.Id)
	return false
}

// InsertElement will add a new element to model.Elements in the position indicated,
// It will also set isChanged to true on additional.
func (model *Model) InsertElement(pos int, element *Element) error {
	if model.Elements == nil {
		model.Elements = []*Element{}
	}
	if !IsValidVarname(element.Id) {
		return fmt.Errorf("element id is not value")
	}
	if model.HasElement(element.Id) {
		return fmt.Errorf("duplicate element id, %q", element.Id)
	}
	if pos < 0 {
		pos = 0
	}
	if pos > len(model.Elements) {
		model.Elements = append(model.Elements, element)
		model.isChanged = true
		return nil
	}
	if pos < len(model.Elements) {
		elements := append(model.Elements[:pos], element)
		model.Elements = append(elements, model.Elements[(pos+1):]...)
	} else {
		model.Elements = append(model.Elements, element)
	}
	model.isChanged = true
	return nil
}

// UpdateElement will update an existing element with element id will the new element.
func (model *Model) UpdateElement(elementId string, element *Element) error {
	if !model.HasElement(elementId) {
		return fmt.Errorf("%q element id not found", elementId)
	}
	for i, e := range model.Elements {
		if e.Id == elementId {
			model.Elements[i] = element
			model.isChanged = true
			return nil
		}
	}
	return fmt.Errorf("failed to find %q to update", elementId)
}

// RemoveElement removes an element by id from the model.Elements
func (model *Model) RemoveElement(elementId string) error {
	if !model.HasElement(elementId) {
		return fmt.Errorf("%q element id not found", elementId)
	}
	for i, e := range model.Elements {
		if e.Id == elementId {
			model.Elements = append(model.Elements[:i], model.Elements[(i+1):]...)
			model.isChanged = true
			return nil
		}
	}
	return fmt.Errorf("%q element id is missing", elementId)
}

// ToSQLiteScheme takes a model and trys to render a SQLite3 SQL create statement.
func (model *Model) ToSQLiteScheme(out io.Writer) error {
	return ModelToSQLiteScheme(out, model)
}

// ToHTML takes a model and trys to render an HTML web form
func (model *Model) ToHTML(out io.Writer) error {
	return ModelToHTML(out, model)
}

// ModelInteractively takes a model and interactively prompts to create
// a YAML model file.
func (model *Model) ModelToYAML(out io.Writer) error {
	return ModelToYAML(out, model)
}

// Register takes a name (string) and a RenderFunc and registers it with the model.
// Registered names then can be invoke by the register name.
func (model *Model) Register(name string, fn RenderFunc) {
	if model.renderer == nil {
		model.renderer = map[string]RenderFunc{}
	}
	model.renderer[name] = fn
}

// Render takes a register render io.Writer and register name envoking the function
// with the model.
func (model *Model) Render(out io.Writer, name string) error {
	if fn, ok := model.renderer[name]; ok {
		return fn(out, model)
	}
	return fmt.Errorf("%s is not a registered rendering function", name)
}

// IsSupportedElementType checks if the element type is supported by Newt, returns true if OK false is it is not
func (model *Model) IsSupportedElementType(eType string) bool {
	for sType, _ := range model.genElements {
		if eType == sType {
			return true
		}
	}
	return false
}

// Define takes a model and attaches a type definition (an element generator) and validator for the named type
func (model *Model) Define(typeName string, genElementFn GenElementFunc, validateFn ValidateFunc) {
	if model.genElements == nil {
		model.genElements = map[string]GenElementFunc{}
	}
	model.genElements[typeName] = genElementFn
	if model.validators == nil {
		model.validators = map[string]ValidateFunc{}
	}
	model.validators[typeName] = validateFn
}
```

---

In `Model` class the method `getElementById` should return either an `Element` or `null` not `Element | null, boolean`.

---

In `Model` class `genElementByType` should return an `Element | null` not `Element | null, boolean`.

---

The `Model` class should have an `errors: string[] = [];` attribute. In the `check` method errors should be push into the `errors` arrays rather than using a WritableStream.

---

`newModel` should return `Model | null` not `[Model | null, string | null]`;

---

Next I'd like to port `model_test.go` to `model_test.ts`. Use `Deno.test` and assertions from `@std/assert`.
This is the Go source code.

```go
package models

import (
	"bytes"
	"testing"

	// 3rd Party packages
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

func inList(l []string, s string) bool {
	for _, val := range l {
		if val == s {
			return true
		}
	}
	return false
}

// TestModel test a model's methods
func TestModel(t *testing.T) {
	m := new(Model)
	if m.HasChanges() {
		t.Errorf("A new empty model should not have changed yet")
	}
	if m.HasElement("id") {
		t.Errorf("A new empty model should not have an id yet")
	}
	if elem, ok := m.GetModelIdentifier(); ok || elem != nil {
		t.Errorf("A new model should not have a identifier assigned yet, got %+v, %t", elem, ok)
	}
	if attrIds := m.GetAttributeIds(); len(attrIds) > 0 {
		t.Errorf("A new model should not have attributes yet, got %+v", attrIds)
	}
	if elemIds := m.GetElementIds(); len(elemIds) > 0 {
		t.Errorf("A new model should not have element ids yet, got %+v", elemIds)
	}
	if elem, ok := m.GetElementById("name"); ok || elem != nil {
		t.Errorf("A new model should not have an element called 'name', got %+v, %t", elem, ok)
	}
	txt := `id: test_model
attributes:
  method: GET
  action: ./
elements:
  - id: id
    type: text
    attributes:
      required: true
      name: id
    is_primary_id: true
  - id: name
    type: text
    attributes:
      name: name
      required: "true"
  - id: msg
    type: textarea
    attributes:
      name: msg
  - id: updated
    type: text
    attributes:
      name: updated
    generator: current_timestamp
  - id: created
    type: text
    atteributes:
      name: created
    generator: created_timestamp
`
	if err := yaml.Unmarshal([]byte(txt), m); err != nil {
		t.Errorf("expected to be able to unmarshal yaml into model, %s", err)
		t.FailNow()
	}
	buf := bytes.NewBuffer([]byte{})
	if !m.Check(buf) {
		t.Errorf("expected valid model, got %s", buf.Bytes())
		t.FailNow()
	}
	expectedAttr := []string{"method", "action", "elements"}
	for _, attr := range m.GetAttributeIds() {
		if !inList(expectedAttr, attr) {
			t.Errorf("expected %q to be in attribute list %+v", attr, expectedAttr)
		}
	}
	expectedElemIds := []string{"id", "name", "msg", "updated"}
	elemIds := m.GetElementIds()
	for _, elemId := range expectedElemIds {
		if !inList(elemIds, elemId) {
			t.Errorf("expected element id %q to be in list %+v", elemId, elemIds)
		}
	}
	primaryId := m.GetPrimaryId()
	if primaryId == "" {
		t.Errorf("expected %q, got %q", "id", primaryId)
	}

	generatedTypes := m.GetGeneratedTypes()
	if len(generatedTypes) != 2 {
		t.Errorf("expected 2 generator type elements, got %d", len(generatedTypes))
	}
	if val, ok := generatedTypes["updated"]; !ok {
		t.Errorf("expected updated to be %t, got %t in generator type", true, ok)
	} else if val != "current_timestamp" {
		t.Errorf("expected %q, got %q", "current_timestamp", val)
	}
	if val, ok := generatedTypes["created"]; !ok {
		t.Errorf("expected created to be %t, got %t in generator type", true, ok)
	} else if val != "created_timestamp" {
		t.Errorf("expected created %q, got %q", "created_timestamp", val)
	}
}

// TestModelBuilding tests creating a newmodel programatticly
func TestModelBuilding(t *testing.T) {
	modelId := "test_model"
	m, err := NewModel(modelId)
	if err != nil {
		t.Errorf("failed to create new model %q, %s", modelId, err)
	}
	m.isChanged = false
	if m.HasChanges() {
		t.Errorf("%s should not have changes yet", modelId)
	}
	// Example YAML expression of a model
	buf := bytes.NewBuffer([]byte{})
	if !m.Check(buf) {
		t.Errorf("expected a valid model, got %s", buf.Bytes())
		t.FailNow()
	}
	/*
	   func (e *Element) Check(buf io.Writer) bool {
	   func IsValidVarname(s string) bool {
	   func NewElement(elementId string) (*Element, error) {
	   func (model *Model) InsertElement(pos int, element *Element) error {
	   func (model *Model) UpdateElement(elementId string, element *Element) error {
	   func (model *Model) RemoveElement(elementId string) error {
	*/
}

// TestHelperFuncs test the funcs from util.go
func TestHelperFuncs(t *testing.T) {
	m := map[string]string{
		"one":   "1",
		"two":   "2",
		"three": "3",
	}
	attrNames := []string{"one", "two", "three"}
	got := getAttributeIds(m)
	if len(got) != 3 {
		t.Errorf("expected 3 attribute ids, got %d %+v", len(got), got)
		t.FailNow()
	}
	for _, expected := range attrNames {
		if !inList(got, expected) {
			t.Errorf("expected %q in %+v, missing", expected, got)
		}
	}
}

// TestValidateModel tests the model validation based on YAML input.
func TestValidateModel(t *testing.T) {
	src := []byte(`id: test_validator
description: This is a test of the validation code
elements:
  - id: pid
    type: text
    attributes:
      name: pid
      required: true
    is_primary_id: true
    label: Personal Identifier
  - id: lived
    type: text
    attributes:
      name: lived
      required: true
    label: Lived Name
  - id: family
    type: text
    attributes:
      name: family
      required: true
    label: Family Name
  - id: orcid
    type: text
    pattern: "[0-9]{4}-[0-9]{4}-[0-9]{4}-[0-9]{3}[0-9A-Z]"
    attributes:
      name: orcid
      required: true
    label: ORCID
`)
	model, err := NewModel("test_model")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := yaml.Unmarshal(src, &model); err != nil {
		t.Error(err)
		t.FailNow()
	}
	SetDefaultTypes(model)

	formData := map[string]string{
		"pid":    "jane-doe",
		"lived":  "Jane",
		"family": "Doe",
		"orcid":  "0000-1111-2222-3333",
	}
	if ok := model.Validate(formData); !ok {
		t.Errorf("%+v failed to validate", formData)
	}
}

// TestValidateMapInterface tests the YAML model mapping
// for decoding and encoding models.
func TestValidateMapInterface(t *testing.T) {
	src := []byte(`id: test_validate_map_inteface
description: This is a test of the validation code
elements:
  - id: pid
    type: text
    attributes:
      name: pid
      required: true
    is_primary_id: true
    label: Personal Identifier
    generator: uuid
  - id: lived
    type: text
    attributes:
      name: lived
      required: true
    label: Lived Name
  - id: family
    type: text
    attributes:
      name: family
      required: true
    label: Family Name
  - id: orcid
    type: text
    pattern: "[0-9]{4}-[0-9]{4}-[0-9]{4}-[0-9]{3}[0-9A-Z]"
    attributes:
      name: orcid
      required: true
    label: ORCID
  - id: created
    type: datetime-local
    attributes:
      required: true
    label: created
    generator: created_timestmap
  - id: updated
    type: datetime-local
    attributes:
      required: true
    generator: current_timestamp
`)
	model, err := NewModel("test_model")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := yaml.Unmarshal(src, &model); err != nil {
		t.Error(err)
		t.FailNow()
	}
	//Debug = true
	SetDefaultTypes(model)

	pid := uuid.New()
	formData := map[string]interface{}{
		"pid":     pid,
		"lived":   "Jane",
		"family":  "Doe",
		"orcid":   "0000-1111-2222-3333",
		"created": "2024-10-03T12:40:00",
		"updated": "2024-10-03 12:41:32",
	}
	if ok := model.ValidateMapInterface(formData); !ok {
		t.Errorf("%+v failed to validate", formData)
	}

	formData = map[string]interface{}{
		"created": "2024-10-03T13:25:24-07:00",
		"family":  "Jetson",
		"lived":   "George",
		"orcid":   "1234-4321-1234-4321",
		"pid":     "0192540f-0806-7631-b08f-4ae5c4d37cca",
		"updated": "2024-10-03T13:25:24-07:00",
	}
	if ok := model.ValidateMapInterface(formData); !ok {
		t.Errorf("%+v failed to validate", formData)
	}
}

// TestModelElements tests the GetGeneratedTypes func for Models.
func TestModelElements(t *testing.T) {
	m := new(Model)
	modelTypes := m.GetGeneratedTypes()
	if len(modelTypes) != 0 {
		t.Errorf("expected zero model types, got %+v", modelTypes)
	}
}
```

---

`Model` should be imported as `import { Model } from "./model.ts";`.

`Element` should be imported with `import { Element } from "./element.ts";`.

Use the built in `crypto.randomUUID()` to create v4 UUID.

---

`const elemIds = model.getElementIds();` is declated twice in the same block. The first
one should be implemented using "let", `let elemIds = model.getElementIds();` The second
should be an assignment, `elemIds = model.getElementIds();`.

---


In `model.ts` remove `import { assert, assertEquals } from "@std/assert";`.

---

In the `Model` class add a `fromObject` method. It takes a parameter of `{[key: string]:any}`. This object is mapped into the attributes of a `Model` including making sure that the types instantiate properly. Example the `this.elements` array needs to hold `Element` instances.

---

In `model_test.ts` change this block of code

```typescript
  const model = Model.newModel("test_model");
  assert(model !== null, "Failed to create new model");
  const yamlData = await yaml.parse(src);
  Object.assign(model, yamlData);
```

to

```typescript
  const model = Model.newModel("test_model");
  assert(model !== null, "Failed to create new model");
  const yamlData = await yaml.parse(src) as {[key: string]:any};
  model.fromObject(yamlData);
```

---

In `model_test.ts` change line 9's `const model = new Model();` to `let model: Model | null = new Model();` and line 55 change line `const model = Model.newModel("test_model");` to `model = Model.newModel("test_model");`

---

In `model_test.ts` remove the following import `import { crypto } from "https://deno.land/std@0.152.0/crypto/mod.ts";`. It is not necessary since `cypto.randomUUID()` is built-in in Deno 2.2.

---

In `Model` class `getModelIdentifier` should return `Element | null` not `[Element | null, boolean]`.

---

Please up `model_test.ts` to conform to the changes in the `Model` class.

---

In `element.ts` the `Element` should have an attribute `errors: string[] = [];`. The `check()` method should the errors encounter to `this.errors` instead of displaying with with `console.error`.

---

In `Element` add a `fromObject` method. The method takes one parameter of type `{[key: string]: any}`. It should map the attributes into the `Element`'s attributes and ensuring type. If the parameter contains `is_primary_key` or `is_object_id` and it is a boolean value then map the boolean value into `Element`'s `isObjectId` attribute.

---

Update `element_test.ts` to include changes to the `Element` class.

---

In `element.ts` the function `isValidVarname` should be exported.

The function `newElement` should be a static method of `Element`.

---

Update `element_test.ts` to include changes to the `Element` class.

---

On line 10 of `element_test.ts` change `  assertEquals(element?.label, "Testelement", "Element label should be 'Testelement'");`
` to `assertEquals(element?.label, "TestElement", "Element label should be 'TestElement', got " + (element.label || "''"));`

---

In `Model` the `fromObject` change the code block

```typescript
        } else if (typeof elem === "object") {
          return new Element(elem);
        }
```

to 

```
        } else if (typeof elem === "object") {
			const newElem = new Element(elem);
			newElem.fromObject(elem);
          return newElem;
        }
```

---

In `model.ts`, `newModel` change the code block

```
    const element = Element.newElement("id");
    element.isObjectId = true;
    element.type = "text";
    element.attributes = { required: "true" };
    model.insertElement(0, element);
```

to 

```
    const element = Element.newElement("id");
	if (element === null ) {
		return null;
	}
    element.isObjectId = true;
    element.type = "text";
    element.attributes = { required: "true" };
    model.insertElement(0, element);
```

---

In `Model` the method `changed` the value `state` also needs to be applied to the `elements` array.

---

The `Model` where methods have used `console.error` the string should be pushed on to the `errors` array.

---

In `Model` update the `newModel` method to look like this.

```
  static newModel(modelId: string): Model {
    const model = new Model();
    model.id = modelId;
    model.description = `... description of ${modelId} goes here ...`;
    model.attributes = {};
    model.elements = [];
    // Make the required element ...
    const element = Element.newElement("id");
    element.isObjectId = true;
    element.type = "text";
    element.attributes = { required: "true" };
    model.insertElement(0, element);
    model.check();
    return model;
  }
```

---

In `Model` update `check` method to

```
  check(): boolean {
    this.errors = [];
    if (!Element.isValidVarname(this.id)) {
      this.errors.push(`Invalid model id, ${this.id}`);
    }
    if (this.elements.length === 0) {
      this.errors.push("Missing model elements.");
      return false;
    }
    let hasModelId = false;
    for (const [i, e] of this.elements.entries()) {
      // Check to make sure each element is valid
      if (!e.check()) {
        this.errors.push(`Error for ${this.id}.${e.id}`);
      }
      if (e.isObjectId) {
        if (hasModelId) {
          this.errors.push(`Duplicate model identifier element (${i}) ${this.id}.${e.id}`);
        }
        hasModelId = true;
      }
    }
    if (!hasModelId) {
      this.errors.push(`Missing required object identifier for model ${this.id}`);
    }
    return this.errors.length === 0;
  }
```

---

In `Element` update the `newElement` method to 

```
  static newElement(elementId: string): Element {
    const element = new Element();
    element.id = elementId;
    element.attributes = { name: elementId };
    element.type = "text";
    element.label = elementId.charAt(0).toUpperCase() + elementId.slice(1);
    element.isObjectId = false;
    element.changed(true);
    element.check();
    return element;
  }
```


In `Element` update the `check` method to

```
  check(): boolean {
    this.errors = [];
    if (!Element.isValidVarname(this.id)) {
      this.errors.push(`Invalid model id, ${this.id}`);
    }
    if (this.id === "") {
      this.errors.push("element missing id");
    }
    if (this.type === "") {
      this.errors.push(`element, ${this.id}, missing type`);
    }
    return this.errors.length === 0;
  }
```

---

Update `Element` method `check` to

```
  check(): boolean {
    this.errors = [];
    if (this.id === "") {
      this.errors.push("element missing id");
    }
    if (!Element.isValidVarname(this.id)) {
      this.errors.push(`Invalid element id, ${this.id}`);
    }
    if (this.type === "") {
      this.errors.push(`element, ${this.id}, missing type`);
    }
    return this.errors.length === 0;
  }
```

---

Let's move the static method in the class `Element` called `isValidVarname` to an exported  function in the TypeScript file named `util.ts`. 

--

Update `element_test.ts` moving the test labeled "TestIsValidVarname" to the file `util_test.ts`.

--

I have a new Go source file to port called `html.go`. It includes functions that get registered with our model to render HTML. The TypeScript file will be called `html.ts` and run under Deno 2.2.  The Go source code follows.

```go
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
```

---

Running `deno check html.ts` yeilds the following error. Can you recommend a fix?

```
TS2339 [ERROR]: Property 'write' does not exist on type 'WritableStream<any>'.
    await out.write(`<!-- ${model.id}: ${model.description} -->\n`);
              ~~~~~
    at file:///Users/rsdoiel/src/github.com/caltechlibrary/models/html.ts:10:15

```

---

I need a TypeScript test module to confirm this works. It should be called `html_test.ts`. It should use `Deno.test` and the assertions from `@std/assert`.

---

When I run `deno check html_test.ts` I get the follow error, can you suggest a fix?

```
TS2551 [ERROR]: Property 'writer' does not exist on type 'UnderlyingSink<any>'. Did you mean 'write'?
      this.writer = controller.getWriter();
           ~~~~~~
    at file:///Users/rsdoiel/src/github.com/caltechlibrary/models/html_test.ts:31:12

    'write' is declared here.
      write?: UnderlyingSinkWriteCallback<W>;
      ~~~~~
        at asset:///lib.deno.web.d.ts:731:3
```

---

When I run `deno test html_test.ts` I get the following test error, can you suggest a fix in either `html.ts` or `html_test.ts`?

```
running 1 test from ./html_test.ts
TestModelToHTML ... FAILED (2ms)

 ERRORS 

TestModelToHTML => ./html_test.ts:7:6
error: TypeError: The writable stream is locked, therefore cannot be closed.
  await output.close();
               ^
    at WritableStream.close (ext:deno_web/06_streams.js:6442:9)
    at file:///Users/rsdoiel/src/github.com/caltechlibrary/models/html_test.ts:43:16

 FAILURES 

TestModelToHTML => ./html_test.ts:7:6

FAILED | 0 passed | 1 failed (4ms)

error: Test failed
```

---

In `html_test.ts` the test case has a wrong expected value. The line `<form id="testModel" required>` is not correct. A form elememt should not have a "required" attribute. Please correct the test.

---

In `html_test.ts` the test case has the wrong expected value. The line with 

```
     <div class="testmodel-textinput"><label class="testmodel-textinput" for="textInput">Textinput</label> <input class="testmodel-textinput" type="text" id="textInput" name="textInput" required></div>\n
```
should should be 

```
     <div class="testmodel-textinput"><label class="testmodel-textinput" for="textInput">TextInput</label> <input class="testmodel-textinput" type="text" id="textInput" name="textInput" required></div>\n
```

Please correct the test.

---

In `html_test.ts` the test case has the wrong expected value. The line with 

```
     <div class="testmodel-submitbutton"><input class="testmodel-submitbutton" type="submit" id="submitButton" value="Submit"></div>\n
```

should be

```
     <div class="testmodel-submitbutton"><label class="testmodel-submitbutton" for="submitButton">SubmitButton</label> <input class="testmodel-submitbutton" type="submit" id="submitButton" value="Submit"></div>\n
```

---

In `html_test.ts` the test case has the wrong expected value. The lines

```
     <div class="testmodel-id"><label class="testmodel-id" for="id">Id</label> <input class="testmodel-id" type="text" id="id" required></div>\n
   \n
```

are not expected.

---

Let's try a different approach. In `html_test.ts`, "TestModelToHTML" instead of building the build using JavaScript, create a YAML document that represents the same structured form. Import `@std/yaml` using the parse function to return a `{[key:string]: any}` type. Create an instance of the `Model` class and use the method `fromObject` to populate it.

```

---

The `html_test.ts` seems to work correctly. There is an error in the output from `html.ts`'s implementation of `modelToHTML`. Can you fix the error in `html.ts` that is identified in by this test result?

```
TestModelToHTML ... FAILED (6ms)

 ERRORS 

TestModelToHTML => ./html_test.ts:8:6
error: AssertionError: Values are not equal: The rendered HTML should match the expected output


    [Diff] Actual / Expected


    <!-- testModel: ... description of testModel goes here ... -->\n
    <form id="testModel">\n
-     <div class="testmodel-textinput"><input class="testmodel-textinput" type="text" id="textInput" name="textInput" required></div>\n
-     <div class="testmodel-textarea"><textarea class="testmodel-textarea" id="textArea" name="textArea" required></textarea></div>\n
-     <div class="testmodel-submitbutton"><input class="testmodel-submitbutton" type="submit" id="submitButton" value="Submit"></div>\n
-   \n
+     <div class="testmodel-textinput"><label class="testmodel-textinput" for="textInput">TextInput</label> <input class="testmodel-textinput" type="text" id="textInput" name="textInput" required></div>\n
+     <div class="testmodel-textarea"><label class="testmodel-textarea" for="textArea">TextArea</label> <textarea class="testmodel-textarea" id="textArea" name="textArea" required></textarea></div>\n
+     <div class="testmodel-submitbutton"><label class="testmodel-submitbutton" for="submitButton">SubmitButton</label> <input class="testmodel-submitbutton" type="submit" id="submitButton" value="Submit"></div>\n
    </form>\n
    


  throw new AssertionError(message);
        ^
    at assertEquals (https://jsr.io/@std/assert/1.0.11/equals.ts:64:9)
    at file:///Users/rsdoiel/src/github.com/caltechlibrary/models/html_test.ts:65:3

 FAILURES 

TestModelToHTML => ./html_test.ts:8:6

FAILED | 0 passed | 1 failed (8ms)

error: Test failed
```

---

Update `html_test.ts` to have `yamlDocument` set to the following YAML document value.

```yaml
  id: testModel
  description: ... description of testModel goes here ...
  attributes: {}
  elements:
    - id: key
      type: text
      attributes:
        name: key
        required: true
      is_object_id: true
    - id: textInput
      type: text
      attributes:
        name: textInput
        required: "true"
    - id: textArea
      type: textarea
      attributes:
        name: textArea
        required: "true"
    - id: submitButton
      type: submit
      attributes:
        value: Submit
```

Update the "TestModelToHTML" to use the new YAML document.

---

What isn't a default label element being generated by `modelToHTML`?

---

Why is there an extra empty new line before the closing form element in `modelToHTML`?

---

Please add the following test to `model_test.ts`,

```
Deno.test("TestYAMLConversion", () => {
  const model = new Model();
  const yamlDocument = `
  id: testModel
  description: ... description of testModel goes here ...
  attributes: {}
  elements:
    - id: key
      type: text
      attributes:
        name: key
        required: true
      is_object_id: true
    - id: textInput
      type: text
      attributes:
        name: textInput
        required: "true"
    - id: textArea
      type: textarea
      attributes:
        name: textArea
        required: "true"
    - id: submitButton
      type: submit
      attributes:
        value: Submit
  `;
  const obj = yaml.parse(yamlDocument) as { [key: string ]: any};
  model.fromObject(obj);
  assertEquals(model.check(), true, `expected true, got false for "${model.errors.join(', ')}"`);
});
```

---

Please display the complete `model_test.ts` result.

---

You seem to forgotten some tests `model_test.ts`. Please add the following tests to `model_test.ts` and display the complete file.

```
Deno.test("TestModelBuilding", async () => {
  const modelId = "test_model";
  const model = Model.newModel(modelId);
  assert(model !== null, `Failed to create new model ${modelId}`);
  model.changed(false);
  assert(!model.hasChanges(), `${modelId} should not have changes yet`);

  const buf = new WritableStream();
  assert(model.check(), "Expected a valid model");
});

Deno.test("TestHelperFuncs", () => {
  const m = {
    one: "1",
    two: "2",
    three: "3",
  };
  const attrNames = ["one", "two", "three"];
  const got = Object.keys(m);
  assertEquals(got.length, 3, "Expected 3 attribute ids");
  for (const expected of attrNames) {
    assert(got.includes(expected), `Expected ${expected} in ${got}`);
  }
});

Deno.test("TestValidateModel", async () => {
  const src = `
id: test_validator
description: This is a test of the validation code
elements:
  - id: pid
    type: text
    attributes:
      name: pid
      required: true
    is_primary_id: true
    label: Personal Identifier
  - id: lived
    type: text
    attributes:
      name: lived
      required: true
    label: Lived Name
  - id: family
    type: text
    attributes:
      name: family
      required: true
    label: Family Name
  - id: orcid
    type: text
    pattern: "[0-9]{4}-[0-9]{4}-[0-9]{4}-[0-9]{3}[0-9A-Z]"
    attributes:
      name: orcid
      required: true
    label: ORCID
`;

  const model = Model.newModel("test_model");
  assert(model !== null, "Failed to create new model");
  const yamlData = await yaml.parse(src) as { [key: string]: any };
  model.fromObject(yamlData);

  const formData = {
    pid: "jane-doe",
    lived: "Jane",
    family: "Doe",
    orcid: "0000-1111-2222-3333",
  };

  assert(model.validate(formData), "Form data failed to validate");
});

Deno.test("TestValidateMapInterface", async () => {
  const src = `
id: test_validate_map_interface
description: This is a test of the validation code
elements:
  - id: pid
    type: text
    attributes:
      name: pid
      required: true
    is_primary_id: true
    label: Personal Identifier
    generator: uuid
  - id: lived
    type: text
    attributes:
      name: lived
      required: true
    label: Lived Name
  - id: family
    type: text
    attributes:
      name: family
      required: true
    label: Family Name
  - id: orcid
    type: text
    pattern: "[0-9]{4}-[0-9]{4}-[0-9]{4}-[0-9]{3}[0-9A-Z]"
    attributes:
      name: orcid
      required: true
    label: ORCID
  - id: created
    type: datetime-local
    attributes:
      required: true
    label: created
    generator: created_timestamp
  - id: updated
    type: datetime-local
    attributes:
      required: true
    generator: current_timestamp
`;

  const model = Model.newModel("test_model");
  assert(model !== null, "Failed to create new model");
  const yamlData = await yaml.parse(src) as { [key: string]: any };
  model.fromObject(yamlData);

  const pid = crypto.randomUUID();
  const formData = {
    pid,
    lived: "Jane",
    family: "Doe",
    orcid: "0000-1111-2222-3333",
    created: "2024-10-03T12:40:00",
    updated: "2024-10-03 12:41:32",
  };

  assert(model.validateMapInterface(formData), "Form data failed to validate");

  const formData2 = {
    created: "2024-10-03T13:25:24-07:00",
    family: "Jetson",
    lived: "George",
    orcid: "1234-4321-1234-4321",
    pid: "0192540f-0806-7631-b08f-4ae5c4d37cca",
    updated: "2024-10-03T13:25:24-07:00",
  };

  assert(model.validateMapInterface(formData2), "Form data failed to validate");
});

Deno.test("TestModelElements", () => {
  const model = new Model();
  const modelTypes = model.getGeneratedTypes();
  assertEquals(Object.keys(modelTypes).length, 0, "Expected zero model types");
});
```

---

The updated `model_test.ts` has revealed problems in `model.ts`. The following tests are failing, can you fix this?

```
TestModelCheck => ./model_test.ts:18:6
error: AssertionError: Values are not equal: Valid model should pass the check


    [Diff] Actual / Expected


-   false
+   true

  throw new AssertionError(message);
        ^
    at assertEquals (https://jsr.io/@std/assert/1.0.11/equals.ts:64:9)
    at file:///Users/rsdoiel/src/github.com/caltechlibrary/models/model_test.ts:27:3

TestModelInsertElement => ./model_test.ts:143:6
error: AssertionError: Expected error message to include "Inserting duplicate element should throw an error", but got "Duplicate element id, element1".
    throw new AssertionError(msg);
          ^
    at assertIsError (https://jsr.io/@std/assert/1.0.11/is_error.ts:63:11)
    at assertThrows (https://jsr.io/@std/assert/1.0.11/throws.ts:96:7)
    at file:///Users/rsdoiel/src/github.com/caltechlibrary/models/model_test.ts:152:3

TestModelUpdateElement => ./model_test.ts:155:6
error: AssertionError: Expected error message to include "Updating non-existent element should throw an error", but got "Element id element2 not found".
    throw new AssertionError(msg);
          ^
    at assertIsError (https://jsr.io/@std/assert/1.0.11/is_error.ts:63:11)
    at assertThrows (https://jsr.io/@std/assert/1.0.11/throws.ts:96:7)
    at file:///Users/rsdoiel/src/github.com/caltechlibrary/models/model_test.ts:164:3

TestModelRemoveElement => ./model_test.ts:167:6
error: AssertionError: Expected error message to include "Removing non-existent element should throw an error", but got "Element id element1 not found".
    throw new AssertionError(msg);
          ^
    at assertIsError (https://jsr.io/@std/assert/1.0.11/is_error.ts:63:11)
    at assertThrows (https://jsr.io/@std/assert/1.0.11/throws.ts:96:7)
    at file:///Users/rsdoiel/src/github.com/caltechlibrary/models/model_test.ts:175:3
```

---

We have removed `isValidVarname` static method from `Element`. Update `model.ts` to use the `isValidVarname` imported from `./util.ts` instead. Display the updated `model.ts`.

---

In `model_test.ts` change the line 

```typescript
  assertEquals(model.check(), true, "Valid model should pass the check");
```

to 

```typescript
  assertEquals(model.check(), true, `Valid model should pass the check, errors ${model.errors.join(', ')}`);
```
---

Line 23 of `model_test.ts` should be 

```
    new Element({ id: "element1", type: "text", isObjectId: true }),
```

---

In `model.ts` the `insertElement` method throws errors. Instead of throwing an error it should 
update the `this.errors` with the error message then return false. If there are no errors the method should return true.

Update `model_test.ts` to reflect this change.

---

In `model.ts` the `updateElement` method throws errors. Instead of throwing an error it should update the `this.errors` with the error message then return false. If there are no errors the method should return true.

Update `model_test.ts` to reflect this change.

---

In `model.ts` the `removeElement` method throws an error. Instead of throwing an error it should update the `this.errors` with the error message then return false. If there are no errors the method should return true.

---

In `model.ts` the `renderElement` method throws an error. Instead of throwing an error it should  update the `this.errors` with the error message then return false. If there are no errors the method should return true.

---

In `model.ts` the `fromObject` method throws an error. Instead of throwing an error it should update the `this.errors` with the error message then return false. If there are no errors the method should return true.

---

Running `deno check model.ts` reveals two type errors. Can you fix this.

```
Check file:///Users/rsdoiel/src/github.com/caltechlibrary/models/model.ts
TS2322 [ERROR]: Type '(Element | null)[]' is not assignable to type 'Element[]'.
  Type 'Element | null' is not assignable to type 'Element'.
    Type 'null' is not assignable to type 'Element'.
      this.elements = data.elements.map((elem: any) => {
      ~~~~~~~~~~~~~
    at file:///Users/rsdoiel/src/github.com/caltechlibrary/models/model.ts:365:7

TS2345 [ERROR]: Argument of type 'null' is not assignable to parameter of type 'Element'.
      if (this.elements.includes(null)) {
                                 ~~~~
    at file:///Users/rsdoiel/src/github.com/caltechlibrary/models/model.ts:376:34

Found 2 errors.
```


---

In the `model.ts`, `Model`, the `fromObject` method is not mapping the elements properly. Update the code from

```typescript
    if (Array.isArray(data.elements)) {
      this.elements = data.elements.map((elem: any) => {
        if (elem instanceof Element) {
          return elem;
        } else if (typeof elem === "object") {
          const newElem = new Element(elem);
          newElem.fromObject(elem);
          return newElem;
        }
        this.errors.push("Invalid element data");
        return null;
      }).filter((elem): elem is Element => elem !== null);
      if (this.elements.includes(null)) {
        return false;
      }
    }
```

to

```typescript
    if (Array.isArray(data.elements)) {	  
		for (const elem of data.elements) {
			if (elem instanceof Element) {
				this.elements.push(elem);
			} else if (typeof elem === "object") {
				const newElem = new Element(elem);
				newElem.fromObject(elem);
				this.elements.push(newElem);
			}		
		}
    } else {
		this.errors.push (`data elements are not an array`);
		return false;
	}
```

You can use `this.check()` at the end of the `fromObject` method to check for errors.

---

In `model_test.ts` I found an error in the setup for "TestModelFromObject". Line 210 is listed as

```typescript
      { id: "element1", type: "text", attributes: { name: "element1" } },
```

It should be

```typescript
      { id: "element1", type: "text", attributes: { name: "element1" }, isObjectId: true },
```
