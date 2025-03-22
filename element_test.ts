import { Element } from "./element.ts";
import { assert, assertEquals } from "@std/assert";

// Define the test suite
Deno.test("TestNewElement", () => {
  const element = Element.newElement("testElement");
  assert(element !== null, "Element should not be null");
  assertEquals(element?.id, "testElement", "Element ID should be 'testElement'");
  assertEquals(element?.type, "text", "Element type should be 'text'");
  assertEquals(element?.label, "TestElement", `Element label should be 'TestElement', got "${element?.label || ''}"`);
  assertEquals(element?.attributes["name"], "testElement", "Element attribute 'name' should be 'testElement'");
  assertEquals(element?.isObjectId, false, "Element isObjectId should be false");
  assertEquals(element?.hasChanged(), true, "Element should be marked as changed");
});

Deno.test("TestElementCheck", () => {
  const element = new Element({ id: "validElement", type: "text" });
  assertEquals(element.check(), true, "Valid element should pass the check");

  const invalidElement = new Element({ id: "", type: "text" });
  assertEquals(invalidElement.check(), false, "Element with missing ID should fail the check, errors " + invalidElement.errors.join('\n\t'));
  assertEquals(invalidElement.errors[0], "element missing id", "Error message should be 'element missing id'");

  const invalidElement2 = new Element({ id: "validElement", type: "" });
  assertEquals(invalidElement2.check(), false, "Element with missing type should fail the check");
  assertEquals(invalidElement2.errors[0], "element, validElement, missing type", "Error message should be 'element, validElement, missing type'");
});

Deno.test("TestElementHasChanged", () => {
  const element = new Element({ id: "testElement", type: "text" });
  assertEquals(element.hasChanged(), false, "Newly created element should not be marked as changed");

  element.changed(true);
  assertEquals(element.hasChanged(), true, "Element should be marked as changed after calling changed(true)");

  element.changed(false);
  assertEquals(element.hasChanged(), false, "Element should not be marked as changed after calling changed(false)");
});

Deno.test("TestIsValidVarname", () => {
  assertEquals(Element.isValidVarname("validName123"), true, "Valid variable name should return true");
  assertEquals(Element.isValidVarname("invalid name"), false, "Invalid variable name should return false");
  assertEquals(Element.isValidVarname("123invalid"), false, "Variable name starting with a number should return false");
});

Deno.test("TestElementFromObject", () => {
  const element = new Element();
  const data = {
    id: "testElement",
    type: "text",
    attributes: { name: "testElement", required: "true" },
    pattern: "pattern",
    options: [{ key: "option1" }, { key: "option2" }],
    generator: "uuid",
    label: "Test Element",
    is_primary_id: true,
  };

  element.fromObject(data);

  assertEquals(element.id, "testElement", "Element ID should be 'testElement'");
  assertEquals(element.type, "text", "Element type should be 'text'");
  assertEquals(element.attributes["name"], "testElement", "Element attribute 'name' should be 'testElement'");
  assertEquals(element.attributes["required"], "true", "Element attribute 'required' should be 'true'");
  assertEquals(element.pattern, "pattern", "Element pattern should be 'pattern'");
  assertEquals(element.options, [{ key: "option1" }, { key: "option2" }], "Element options should be [{ key: 'option1' }, { key: 'option2' }]");
  assertEquals(element.generator, "uuid", "Element generator should be 'uuid'");
  assertEquals(element.label, "Test Element", "Element label should be 'Test Element'");
  assertEquals(element.isObjectId, true, "Element isObjectId should be true");
});
