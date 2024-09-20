
# models

This is a Go package used to describe data models aligned with the HTML5 data types. The model can be expressed in YAML or JSON. The YAML (or JSON) data structure is patterned after the HTML5 form elements. A single model can be used to generate HTML web forms or used to validate a map that confirms to the model. In princple generators can be written to express the model in other forms, e.g. SQL.

It is important to note that is not an Object Relational Mapper (ORM).  The purpose of the model package is to facilitate describing simple data models using YAML then beable to reuse the models in other Go based projects (e.g. [dataset](http://github.com/caltechlibrary/dataset), [Newt](https://github.com/caltechlibrary/newt)).

This Go package assumes Go version 1.23.1 or better.

