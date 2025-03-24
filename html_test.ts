import { assert, assertEquals } from "@std/assert";
import { Model } from "./model.ts";
import { Element } from "./element.ts";
import { modelToHTML } from "./html.ts";
import { parse } from "@std/yaml";

// Define the test suite
Deno.test("TestModelToHTML", async () => {
  // YAML document representing the form structure
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

  // Parse the YAML document
  const parsedObject = parse(yamlDocument) as { [key: string]: any };

  // Create a Model instance and populate it using fromObject
  const model = new Model();
  model.fromObject(parsedObject);

  // Create a WritableStream to capture the output
  const chunks: Uint8Array[] = [];
  const output = new WritableStream({
    write(chunk) {
      chunks.push(chunk);
    },
    close() {
      // Do nothing
    },
  });

  // Render the model to HTML
  await modelToHTML(output, model);

  // Decode the output chunks into a string
  const decoder = new TextDecoder();
  const outputString = chunks.map(chunk => decoder.decode(chunk)).join("");

  // Validate the output
  const expectedOutput = `<!-- testModel: ... description of testModel goes here ... -->
<form id="testModel">
  <div class="testmodel-key"><label class="testmodel-key" for="key">Key</label> <input class="testmodel-key" type="text" id="key" name="key" required></div>
  <div class="testmodel-textinput"><label class="testmodel-textinput" for="textInput">TextInput</label> <input class="testmodel-textinput" type="text" id="textInput" name="textInput" required></div>
  <div class="testmodel-textarea"><label class="testmodel-textarea" for="textArea">TextArea</label> <textarea class="testmodel-textarea" id="textArea" name="textArea" required></textarea></div>
  <div class="testmodel-submitbutton"><label class="testmodel-submitbutton" for="submitButton">SubmitButton</label> <input class="testmodel-submitbutton" type="submit" id="submitButton" value="Submit"></div>
</form>
`;

  assertEquals(outputString, expectedOutput, "The rendered HTML should match the expected output");
});
