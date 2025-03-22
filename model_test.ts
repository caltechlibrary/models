import { Model } from "./model.ts";
import { Element } from "./element.ts";
import { assert, assertEquals } from "@std/assert";
import * as yaml from "@std/yaml";

// Define the test suite
Deno.test("TestModel", async () => {
  let model: Model | null = new Model();
  assert(!model.hasChanges(), "A new empty model should not have changed yet");
  assert(
    !model.hasElement("id"),
    "A new empty model should not have an id yet",
  );
  const elem = model.getModelIdentifier();
  assert(
    elem === null,
    "A new model should not have a identifier assigned yet",
  );
  const attrIds = model.getAttributeIds();
  assert(attrIds.length === 0, "A new model should not have attributes yet");
  let elemIds = model.getElementIds();
  assert(elemIds.length === 0, "A new model should not have element ids yet");
  const elemById = model.getElementById("name");
  assert(
    elemById === null,
    "A new model should not have an element called 'name'",
  );

  const txt = `
id: test_model
attributes:
  method: GET
  action: ./
elements:
  - id: id
    type: text
    attributes:
      required: true
      name: id
    is_object_id: true
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
    attributes:
      name: created
    generator: created_timestamp
`;

  const yamlData = await yaml.parse(txt) as { [key: string]: any };
  model = Model.newModel("test_model");
  assert(model !== null, "Failed to create new model");
  model.fromObject(yamlData);

  assert(model.check(), `Expected valid model, ${model.errors}`);

  const expectedAttr = ["method", "action", "elements"];
  for (const attr of model.getAttributeIds()) {
    assert(
      expectedAttr.includes(attr),
      `Expected ${attr} to be in attribute list`,
    );
  }

  const expectedElemIds = ["id", "name", "msg", "updated"];
  elemIds = model.getElementIds();
  for (const elemId of expectedElemIds) {
    assert(
      elemIds.includes(elemId),
      `Expected element id ${elemId} to be in list`,
    );
  }

  const primaryId = model.getPrimaryId();
  assertEquals(
    primaryId,
    "id",
    `Expected primary id to be "id", got "${primaryId}"`,
  );

  const generatedTypes = model.getGeneratedTypes();
  assertEquals(
    Object.keys(generatedTypes).length,
    2,
    "Expected 2 generator type elements",
  );
  assertEquals(
    generatedTypes["updated"],
    "current_timestamp",
    "Expected 'updated' to be 'current_timestamp'",
  );
  assertEquals(
    generatedTypes["created"],
    "created_timestamp",
    "Expected 'created' to be 'created_timestamp'",
  );
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
