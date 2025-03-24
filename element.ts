/**
 * element.ts holds the Element class and related functions. Elemement are used by the Model class.
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
/**
 * element.ts holds the Element class and related functions.
 *
 * @author R. S. Doiel
 */

import { isValidVarname } from "./util.ts";

// Define the Element class
export class Element {
  // Type, The type of element that you want to input. It is required. Valid values are
  // checkboxes, dropdown, input, markdown and textarea.
  //
  // The input type corresponds to the native input types defined for HTML 5. E.g. text, textarea,
  // email, phone, date, url, checkbox, radio, button, submit, cancel, select
  // See MDN developer docs for input, <https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input>
  type: string = "text";

  // Id for the element, except when type is set to markdown. Can only use alpha-numeric characters,
  // -, and _. Must be unique in the form definition. If provided, the id is the canonical identifier
  // for the field in URL query parameter prefill.
  id: string = "";

  // Attributes, a set of key-value pairs that define the properties of the element.
  // This is a required element as it holds the "value" attribute when expressing
  // HTML content. Other commonly used attributes
  attributes: { [key: string]: string } = {};

  // Pattern holds a validation pattern. When combined with an input type (or input type alias, e.g. orcid)
  // produces a form element that sports a specific client-side validation expectation. This intern can be used
  // to generate appropriate validation code server-side.
  pattern: string = "";

  // Options holds a list of values and their labels used for HTML select elements in rendering their option child elements
  options?: { [key: string]: string }[] = [];

  // IsObjectId (i.e., is the identifier of the object) used by the modeled data.
  // It is used in calculating routes and templates where the object identifier is required.
  isObjectId: boolean = false;

  // Generator indicates the type of automatic population of a field. It is used to
  // indicate auto-increment and UUIDs for primary keys and timestamps for datetime-oriented fields.
  generator: string = "";

  // Label is used when rendering an HTML form as a label element tied to the input element via the set attribute and
  // the element's id.
  label: string = "";

  // errors array to store validation errors
  errors: string[] = [];

  // Private attribute to track changes
  private isChanged: boolean = false;

  constructor(data?: Partial<Element>) {
    Object.assign(this, data);
  }

  // HasChanged checks to see if the Element has been changed.
  hasChanged(): boolean {
    return this.isChanged;
  }

  // Changed sets the change state on element
  changed(state: boolean): void {
    this.isChanged = state;
  }

  // Check reviews an Element to make sure it is valid.
  check(): boolean {
    this.errors = [];
    if (this.id === "") {
      this.errors.push("element missing id");
    }
    if (!isValidVarname(this.id)) {
      this.errors.push(`Invalid element id, ${this.id}`);
    }
    if (this.type === "") {
      this.errors.push(`element, ${this.id}, missing type`);
    }
    return this.errors.length === 0;
  }

  // fromObject takes a parameter of type {[key: string]: any} and maps it into the attributes of an Element
  fromObject(data: { [key: string]: any }): void {
    if (data.id && typeof data.id === "string") {
      this.id = data.id;
    }
    if (data.type && typeof data.type === "string") {
      this.type = data.type;
    }
    if (data.attributes && typeof data.attributes === "object") {
      this.attributes = data.attributes;
    }
    if (data.pattern && typeof data.pattern === "string") {
      this.pattern = data.pattern;
    }
    if (data.options && Array.isArray(data.options)) {
      this.options = data.options;
    }
    if (data.generator && typeof data.generator === "string") {
      this.generator = data.generator;
    }
    if (data.label && typeof data.label === "string") {
      this.label = data.label;
    }
    if (data.is_primary_id !== undefined || data.is_object_id !== undefined) {
      this.isObjectId = Boolean(data.is_primary_id || data.is_object_id);
    }
  }

  // NewElement makes sure element id is valid, populates an element as a basic input type.
  // The new element has the attribute "name" and label set to default values.
  static newElement(elementId: string): Element {
    const element = new Element();
    element.id = elementId;
    element.attributes = { name: elementId };
    element.type = "text";
    element.label = elementId.charAt(0).toUpperCase() + elementId.slice(1);
    element.isObjectId = false;
    element.changed(true);
    element.check();
    return element;
  }
}
