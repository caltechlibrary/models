
# models

This is a Go package used to describe data models aligned with the HTML5 data types. The model can be expressed in YAML or JSON. The YAML (or JSON) data structure is patterned after the HTML5 form elements. A single model can be used to generate HTML web forms or used to validate a map that confirms to the model. In princple generators can be written to express the model in other forms, e.g. SQL.

It is important to note that is not an Object Relational Mapper (ORM).  The purpose of the model package is to facilitate describing simple data models using YAML then beable to reuse the models in other Go based projects (e.g. [dataset](http://github.com/caltechlibrary/dataset), [Newt](https://github.com/caltechlibrary/newt)).

This Go package assumes Go version 1.23.1 or better.

# Oberservation: Web forms describe a simple data structure

The models package grew out of an observation that if you can define the elements of an HTML5 web form you can also describe a simple data model or schema. The problem is HTML5 is combersum to type, read and work with.  On the other hand it lends itself to expression in simpler representations.

YAML can be used to represent a web form in a clear and concise way. From that description you can extrapulate HTML and SQL Schema. You can also use it as a guide to data validation for web form submissions.

Our common use cases.

1. Web form as YAML can be used to generate HTML web forms
2. Web form elements can be used to inferring the SQL column type
3. Web form as YAML is a guide to validating web form submissions

# A simple example

A "guest book" model.

~~~yaml
id: guest_book_entry
attributes:
  action: ./signbook.html
  method: POST
  x-success: ./thankyou.html
  x-failure: ./oops.html
elements:
  - id: record_id
    type: text
    pattern: [a-z0-9]+\.[a-z0-9]+
    attributes:
      name: record_id
      placeholder: A unique record id
      required: true
  - id: name
    type: text
    attributes:
      name: name
      placeholder: E.g. Doe, Jane
      required: true
  - id: message
    type: text
    attributes:
      placeholder: Please leave a message
      name: message
  - id: signed
    type: date
    attributes:
      name: signed
      required: true
~~~

This "model" describes JSON data that might look like the following.

~~~json
{ 
    "record_id": "jane.doe",
    "Doe, Jane",
    "signed": "2024-09-10"
}
~~~

The model could be used to generate the web form and validate the data. It implies an SQL Scheme.  The model package provides the means of working with a model and to validate the model's content. By normalizing your data elements to throse supported by HTML5 you also can easily generate the code you need (e.g. HTML form or SQL scheme).

The package doesn't provide the extrapolated forms but does provide the functions and method to make it easy to build them.


