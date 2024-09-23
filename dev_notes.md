
# Development Notes

## Naming things

Check
: This refers to "checking" a model to see if it makes sense as a model. An example would be a model read in as hand coded YAML.

Register
: This refers to associating a render method with the model

Render
: This takes the model structure, passes it to the registered method to generate a model representation. An example representation of the model might be an SQL statement or as an HTML form. It is important to note that this renders the model as opposed to the data held by the model.

## Named Ideas to concider implementation

Define
: This refers to associating a type definition available to elements of a model. An example would be the definition of textarea, text, or date type making them available to be model. This is the method you'd use to define specialized types like ORCID, ROR, ISNI, etc.

Validate
: This refers to the confirming correctness of the data the model is holding. This is the only method that will operates on the data the model is holding. Calling Validate on the model will then loop through all the elements in the model invoking their Validate method.

