/**
 * element_test.ts contains testing for element.ts
 * 
 * @author R. S. Doiel, <rsdoiel@caltech.edu>
 *
 * Copyright (c) 2025, Caltech
 * All rights not granted herein are expressly reserved by Caltech.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted provided
 * that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this list of conditions and 
 *    the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions
 *    and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *
 * 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or
 *    promote products derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, 
 * INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
 * WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE
 * USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
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
  assertEquals(invalidElement.check(), false, "Element with missing ID should fail the check");
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
