%modelgen(1) user manual | version 0.0.1 3061bf5
% R. S. Doiel
% 2024-09-20

# NAME

modelgen 

# SYNOPSIS

modelgen [OPTIONS] html|sqlite3 [MODEL_NAME] [OUT_NAME]

# DESCRIPTION

modelgen is a demonstration of the models package for Go.  It can read
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
: Display modelgen version.

-license
: Display modelgen license.

# EXAMPLE

~~~
modelgen html guestbook.yaml guestbook.html
modelgen sqlite3 guestbook.yaml guestbook.sql
~~~


