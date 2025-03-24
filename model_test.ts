import { assert, assertEquals, assertThrows } from "@std/assert";
import { Model } from "./model.ts";
import { Element } from "./element.ts";
import * as yaml from "@std/yaml";
import { crypto } from "https://deno.land/std@0.152.0/crypto/mod.ts";

// Define the test suite for Model
Deno.test("TestNewModel", () => {
  const model = Model.newModel("testModel");
  assert(model !== null, "Model should not be null");
  assertEquals(model.id, "testModel", "Model ID should be 'testModel'");
  assertEquals(model.description, "... description of testModel goes here ...", "Model description should be set");
  assertEquals(model.elements.length, 1, "Model should have one element by default");
  assertEquals(model.elements[0].id, "id", "Default element ID should be 'id'");
  assertEquals(model.elements[0].isObjectId, true, "Default element should be the object identifier");
});

Deno.test("TestModelCheck", () => {
  const model = new Model();
  model.id = "validModel";
  model.description = "A valid model";
  model.elements = [
    new Element({ id: "element1", type: "text", isObjectId: true }),
    new Element({ id: "element2", type: "textarea" }),
  ];

  assertEquals(model.check(), true, `Valid model should pass the check, errors ${model.errors.join(', ')}`);

  model.id = "invalid model id";
  assertEquals(model.check(), false, "Model with invalid ID should fail the check");
  assertEquals(model.errors[0], "Invalid model id, invalid model id", "Error message should indicate invalid model ID");

  model.id = "validModel";
  model.elements = [];
  assertEquals(model.check(), false, "Model with no elements should fail the check");
  assertEquals(model.errors[0], "Missing model elements.", "Error message should indicate missing elements");

  model.elements = [
    new Element({ id: "element1", type: "text" }),
    new Element({ id: "element2", type: "textarea" }),
    new Element({ id: "element3", type: "text", isObjectId: true }),
    new Element({ id: "element4", type: "text", isObjectId: true }),
  ];
  assertEquals(model.check(), false, "Model with duplicate object identifiers should fail the check");
  assertEquals(model.errors[0], "Duplicate model identifier element (3) validModel.element4", "Error message should indicate duplicate object identifier");
});

Deno.test("TestModelHasElement", () => {
  const model = new Model();
  model.elements = [
    new Element({ id: "element1", type: "text" }),
    new Element({ id: "element2", type: "textarea" }),
  ];

  assertEquals(model.hasElement("element1"), true, "Model should have element with ID 'element1'");
  assertEquals(model.hasElement("element3"), false, "Model should not have element with ID 'element3'");
});

Deno.test("TestModelHasElementType", () => {
  const model = new Model();
  model.elements = [
    new Element({ id: "element1", type: "text" }),
    new Element({ id: "element2", type: "textarea" }),
  ];

  assertEquals(model.hasElementType("text"), true, "Model should have element of type 'text'");
  assertEquals(model.hasElementType("button"), false, "Model should not have element of type 'button'");
});

Deno.test("TestModelGetModelIdentifier", () => {
  const model = new Model();
  model.elements = [
    new Element({ id: "element1", type: "text" }),
    new Element({ id: "element2", type: "textarea", isObjectId: true }),
  ];

  const identifier = model.getModelIdentifier();
  assert(identifier !== null, "Model identifier should not be null");
  assertEquals(identifier?.id, "element2", "Model identifier should be the element with isObjectId set to true");
});

Deno.test("TestModelGetAttributeIds", () => {
  const model = new Model();
  model.attributes = { attr1: "value1", attr2: "value2" };

  const attributeIds = model.getAttributeIds();
  assertEquals(attributeIds.length, 2, "Model should have 2 attribute IDs");
  assertEquals(attributeIds.includes("attr1"), true, "Model attributes should include 'attr1'");
  assertEquals(attributeIds.includes("attr2"), true, "Model attributes should include 'attr2'");
});

Deno.test("TestModelGetElementIds", () => {
  const model = new Model();
  model.elements = [
    new Element({ id: "element1", type: "text" }),
    new Element({ id: "element2", type: "textarea" }),
  ];

  const elementIds = model.getElementIds();
  assertEquals(elementIds.length, 2, "Model should have 2 element IDs");
  assertEquals(elementIds.includes("element1"), true, "Model elements should include 'element1'");
  assertEquals(elementIds.includes("element2"), true, "Model elements should include 'element2'");
});

Deno.test("TestModelGetPrimaryId", () => {
  const model = new Model();
  model.elements = [
    new Element({ id: "element1", type: "text" }),
    new Element({ id: "element2", type: "textarea", isObjectId: true }),
  ];

  const primaryId = model.getPrimaryId();
  assertEquals(primaryId, "element2", "Primary ID should be the ID of the element with isObjectId set to true");
});

Deno.test("TestModelGetGeneratedTypes", () => {
  const model = new Model();
  model.elements = [
    new Element({ id: "element1", type: "text", generator: "uuid" }),
    new Element({ id: "element2", type: "textarea" }),
  ];

  const generatedTypes = model.getGeneratedTypes();
  assertEquals(generatedTypes["element1"], "uuid", "Generated type for 'element1' should be 'uuid'");
  assertEquals(generatedTypes["element2"], undefined, "Generated type for 'element2' should be undefined");
});

Deno.test("TestModelGetElementById", () => {
  const model = new Model();
  model.elements = [
    new Element({ id: "element1", type: "text" }),
    new Element({ id: "element2", type: "textarea" }),
  ];

  const element = model.getElementById("element1");
  assert(element !== null, "Element with ID 'element1' should not be null");
  assertEquals(element?.id, "element1", "Element ID should be 'element1'");

  const nonExistentElement = model.getElementById("element3");
  assertEquals(nonExistentElement, null, "Non-existent element should be null");
});

Deno.test("TestModelInsertElement", () => {
  const model = new Model();
  const element = new Element({ id: "element1", type: "text" });

  const result = model.insertElement(0, element);
  assertEquals(result, true, "Inserting a valid element should return true");
  assertEquals(model.elements.length, 1, "Model should have 1 element after insertion");
  assertEquals(model.elements[0].id, "element1", "Inserted element ID should be 'element1'");

  const duplicateResult = model.insertElement(0, element);
  assertEquals(duplicateResult, false, "Inserting a duplicate element should return false");
  assertEquals(model.errors[0], `Duplicate element id: element1`, "Error message should indicate duplicate element ID");
});

Deno.test("TestModelUpdateElement", () => {
  const model = new Model();
  const element = new Element({ id: "element1", type: "text" });
  model.insertElement(0, element);

  const updatedElement = new Element({ id: "element1", type: "textarea" });
  const result = model.updateElement("element1", updatedElement);
  assertEquals(result, true, "Updating an existing element should return true");
  assertEquals(model.elements[0].type, "textarea", "Updated element type should be 'textarea'");

  const nonExistentResult = model.updateElement("element2", updatedElement);
  assertEquals(nonExistentResult, false, "Updating a non-existent element should return false");
  assertEquals(model.errors[0], `Element id element2 not found`, "Error message should indicate element ID not found");
});

Deno.test("TestModelRemoveElement", () => {
  const model = new Model();
  const element = new Element({ id: "element1", type: "text" });
  model.insertElement(0, element);

  const result = model.removeElement("element1");
  assertEquals(result, true, "Removing an existing element should return true");
  assertEquals(model.elements.length, 0, "Model should have 0 elements after removal");

  const nonExistentResult = model.removeElement("element1");
  assertEquals(nonExistentResult, false, "Removing a non-existent element should return false");
  assertEquals(model.errors[0], `Element id element1 not found`, "Error message should indicate element ID not found");
});

Deno.test("TestModelRenderElement", async () => {
  const model = new Model();
  const element = new Element({ id: "element1", type: "text" });
  model.insertElement(0, element);

  const result = await model.renderElement("html", new WritableStream());
  assertEquals(result, false, "Rendering with an unregistered function should return false");
  assertEquals(model.errors[0], `html is not a registered rendering function`, "Error message should indicate unregistered rendering function");

  model.register("html", async (out, model) => {
    // Mock render function
  });

  const validResult = await model.renderElement("html", new WritableStream());
  assertEquals(validResult, true, "Rendering with a registered function should return true");
});

Deno.test("TestModelFromObject", () => {
  const model = new Model();
  const data = {
    id: "testModel",
    description: "A test model",
    attributes: { attr1: "value1" },
    elements: [
      { id: "element1", type: "text", attributes: { name: "element1" }, isObjectId: true },
      { id: "element2", type: "textarea", attributes: { name: "element2" } },
    ],
  };

  const result = model.fromObject(data);
  assertEquals(result, true, "fromObject should return true for valid data");
  assertEquals(model.id, "testModel", "Model ID should be 'testModel'");
  assertEquals(model.description, "A test model", "Model description should be 'A test model'");
  assertEquals(model.attributes["attr1"], "value1", "Model attribute 'attr1' should be 'value1'");
  assertEquals(model.elements.length, 2, "Model should have 2 elements");
  assertEquals(model.elements[0].id, "element1", "First element ID should be 'element1'");
  assertEquals(model.elements[1].id, "element2", "Second element ID should be 'element2'");
});

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

  const obj = yaml.parse(yamlDocument) as { [key: string]: any };
  const result = model.fromObject(obj);
  assertEquals(result, true, "fromObject should return true for valid YAML data");
  assertEquals(model.check(), true, `expected true, got false for "${model.errors.join(', ')}"`);
});

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
