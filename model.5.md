%models(5) user manual | version 0.0.0 b53cae9
% R. S. Doiel
% 2024-09-20

# NAME

models

# SYNOPSIS

~~~
import "github.com/caltechlibrary/models"
~~~

# DESCRIPTION

__models__ is a Go package. A model is expressed in YAML. They are used by `modelgen` to render HTML web forms
or SQLite3 schema.

## Model

id
: The identifier for the model. Is the "id" given to the generated HTML web form.

title
: If provided it will be used to insert a title above your web form.

attributes
: These map to the HTML attributes in a web form. Typical you would include method (e.g. GET, POST) and action (e.g. a URL to the form page).
Attributes are a key/value map of form attributes.

description
: This is simple description of the model. It will be included as a comment in the SQLite3 SQL. This is a text string or block.

elements
: This is a listof elements that describe the data attritubes of your model.

## Elements

The elements attribute holds a list of elements. You can think of these as HTML5 form elements described in YAML.
They will also be used to infer SQLite 3 column types.

Each element is made from the following properties.

type
: (required) This is a string and maps to the input element types available in HTML5[^1]. 

id
: (optional) This is the element's identifier. It should be unique with in the model. While optional it is used to retrieve an element from a model. If is
also required when rendering colum definitions in SQLite 3. A model that includes a submit or reset button would examples of when to leave it blank.

attributes
: (optional) This is a list of key/value pairs that map to HTML5 input elements. Boolean HTML element attributes like "required" and "checked" you are expressed
as `required: true` and `checked: true` in YAML. NOTE: attributes value's are resolved to quoted strings when rendered as HTML.

pattern
: (optional) This is a regular expression pattern that is used to validate the input of the element[^2].

options
: (optional) Are a list of key/value maps used to expression HTML5 select elements. They can be be used in validation of a model's content as well as in render HTML selection elements.

is\_primary\_id
: (optional) If set to true it indicates a given element holds the model's primary identifier. If you are store model content in a SQLite 3 database or Dataset collection this would be the unique identifier used to retrieve the modeled object.

label
: (optional) If set it is used as the text content of the label when rendering a web form.

# Example

This is an example model of a guest book entry used in a Dataset base guest book web application.

~~~yaml
id: test_model
attributes:
  method: GET
  action: ./
  x-success: http://localhost:8000/success.html
  x-fail: http://localhost:8000/failed.html
elements:
  - id: id
    type: text
    attributes:
      name: id
      placeholder: Enter a unique string
      required: true
    is_primary_id: true
  - id: name
    type: text
    attributes:
      name: name
      placeholdertext: Enter your name
      required: true
  - id: msg
    type: textarea
    attributes:
      name: msg
      placehodertext: Enter a short message or comment
~~~





[^1]: See <https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input> for details.

[^2]: See <https://developer.mozilla.org/en-US/docs/Web/HTML/Attributes/pattern> for details of how patterns are used in validation.




