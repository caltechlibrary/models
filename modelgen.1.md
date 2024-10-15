%modelgen(1) user manual | version 0.0.5 af4ccd4
% R. S. Doiel
% 2024-10-15

# NAME

modelgen 

# SYNOPSIS

modelgen [OPTIONS] ACTION [MODEL_NAME] [OUT_NAME]

# DESCRIPTION

modelgen is a demonstration of the models package for Go.  It can read
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
: Display modelgen version.

-license
: Display modelgen license.

# EXAMPLE

In this example we create a new model YAML file interactively using
the "model" action. Then create an HTML page followed by SQL file
holding the SQL schema for SQLite 3.

~~~
modelgen model guestbook.yaml
modelgen html guestbook.yaml guestbook.html
modelgen sqlite guestbook.yaml guestbook.sql
~~~


