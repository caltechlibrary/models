/**
 * element.ts holds the Element class and related functions.
 *
 * @author R. S. Doiel
 */

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
    if (!Element.isValidVarname(this.id)) {
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

  // Function to check if a string is a valid variable name
  static isValidVarname(name: string): boolean {
    const validVarnameRegex = /^[a-zA-Z_][a-zA-Z0-9_]*$/;
    return validVarnameRegex.test(name);
  }
}
