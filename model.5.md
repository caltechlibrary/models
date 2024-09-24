%models(5) user manual | version 0.0.1 14f6d2f
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
: This is a list of elements that describe the data attributes of your model.

## Elements

The elements attribute holds a list of elements. You can think of these as HTML5 form elements described in YAML.
They will also be used to infer SQLite 3 column types.

Each element is made from the following properties.

type
: (required) This is a string and maps to the input element types available in HTML5[^1]. 

id
: (optional) This is the element's identifier. It should be unique with in the model. While optional it is used to retrieve an element from a model. If is
also required when rendering column definitions in SQLite 3. A model that includes a submit or reset button would examples of when to leave it blank.

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

[^1]: See <https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input> for details.

[^2]: See <https://developer.mozilla.org/en-US/docs/Web/HTML/Attributes/pattern> for details of how patterns are used in validation.

## Data type support in models package

The models package starts from the premise of supporting a YAML description of a web form that then can be used to render HTML and SQL Schema. It also needs to be able to be a thin layer in a Go API that can validate the elements of a model just like they are validated browser side by the [HTML5 input types](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input).  The following are all implemented by the models package using a naive validation approach[^3].  

[^3]: "naive" in this case means overly simplistic validation, e.g. min max ranges don't validate against step attributes. 


- [button](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/button)
- [checkbox](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/checkbox)
- [color](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/color)
- [date](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/date)
- [datetime-local](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/datetime-local)
- [email](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/email)
- [hidden](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/hidden)
- [image](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/image)
- [month](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/month)
- [number](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/number)
- [password](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/password)
- [radio](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/radio)
- [range](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/range)
- [reset](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/reset)
- [search](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/search)
- [submit](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/submit)
- [tel](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/tel)
- [text](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/text)
- [time](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/time)
- [url](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/url)
- [week](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/week)

Additional data types[^4] can be defined by using the `Model.Define` function provided in this package. You need to provide a name for the new type as well as the func's name. The "defined" data types are applied before the default types. This allows for improvements to the defaults while retaining a fallback. Hopefully this mechanism can prove useful to expanding the data types supported by models.

[^4]: The validation function is used server side only because it is written in Go. E.g. by Dataset's JSON API.

NOTE: As the models package evolves the validation methods provided out of the box will evolve too. Some may even be dropped if they prove problematic[^5].

[^5]: E.g. "week" input type is not widely used and is poorly supported by browsers in 2024. "image" doesn't make a whole lot of sense.

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










