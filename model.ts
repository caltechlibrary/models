/**
 * model.ts holds the Model class and related functions.
 *
 * @author R. S. Doiel
 */

import { Element } from "./element.ts";
import { isValidVarname } from "./util.ts";

// Define the RenderFunc type
type RenderFunc = (writer: WritableStream, model: Model) => Promise<void>;

// Define the ValidateFunc type
type ValidateFunc = (element: Element, value: string) => boolean;

// Define the GenElementFunc type
type GenElementFunc = () => Element;

// Define the Model class
export class Model {
  // Id is a required field for model, it maps to the HTML element id and name
  id: string = "";

  // This is a Newt specific set of attributes to place in the form element of HTML. I.e., it could
  // be form "class", "method", "action", "encoding". It is not defined in the GitHub YAML issue template syntax
  // (optional)
  attributes: { [key: string]: string } = {};

  // Description, A description for the issue form template, which appears in the template chooser interface.
  // (required)
  description: string = "";

  // Elements, Definition of the input types in the form.
  // (required)
  elements: Element[] = [];

  // Title, A default title that will be pre-populated in the issue submission form.
  // (optional) only there for compatibility with GitHub YAML Issue Templates
  title?: string;

  // isChanged is an internal state used by the modeler to know when a model has changed
  private isChanged: boolean = false;

  // errors array to store validation errors
  errors: string[] = [];

  // renderer is a map of names to RenderFunc functions. A RenderFunc is that take a WritableStream and the model object as parameters then
  // return an error type. This allows for many renderers to be used with Model by
  // registering the function then invoking render with the name registered.
  renderer: { [key: string]: RenderFunc } = {};

  // genElements holds a map to the "type" pointing to an element generator
  genElements: { [key: string]: GenElementFunc } = {};

  // validators holds a list of validate functions associated with types. Key is type name.
  validators: { [key: string]: ValidateFunc } = {};

  constructor(data?: Partial<Model>) {
    Object.assign(this, data);
  }

  // GenElementType takes an element type and returns an Element struct populated for that type or null if type is not supported.
  genElementType(typeName: string): Element | null {
    const fn = this.genElements[typeName];
    if (fn) {
      return fn();
    }
    return null;
  }

  // Validate form data expressed as map[string]string.
  validate(formData: { [key: string]: string }): boolean {
    const ids = this.getElementIds();
    if (ids.length !== Object.keys(formData).length) {
      this.errors.push("Form data does not match model element ids.");
      return false;
    }
    for (const [k, v] of Object.entries(formData)) {
      const elem = this.getElementById(k);
      if (elem) {
        const validator = this.validators[elem.type];
        if (validator && !validator(elem, v)) {
          this.errors.push(`Failed to validate elem.Id ${elem.id}, elem.Type ${elem.type}, value ${v}`);
          return false;
        }
      } else {
        this.errors.push(`Element with id ${k} not found in model.`);
        return false;
      }
    }
    return true;
  }

  // ValidateMapInterface normalizes the map interface values before calling
  // the element's validator function.
  validateMapInterface(data: { [key: string]: any }): boolean {
    const ids = this.getElementIds();
    if (ids.length !== Object.keys(data).length) {
      this.errors.push(`Expected len(ids) ${ids.length}, got len(data) ${Object.keys(data).length}`);
      return false;
    }
    for (const [k, v] of Object.entries(data)) {
      let val: string;
      switch (typeof v) {
        case "string":
          val = v;
          break;
        case "number":
          val = v.toString();
          break;
        case "boolean":
          val = v.toString();
          break;
        default:
          val = JSON.stringify(v);
      }
      const elem = this.getElementById(k);
      if (elem) {
        const validator = this.validators[elem.type];
        if (validator && !validator(elem, val)) {
          this.errors.push(`Failed to validate elem.Id ${elem.id}, value ${val}`);
          return false;
        }
      } else {
        this.errors.push(`Element with id ${k} not found in model.`);
        return false;
      }
    }
    return true;
  }

  // HasChanges checks if the model's elements have changed
  hasChanges(): boolean {
    if (this.isChanged) {
      return true;
    }
    return this.elements.some(e => e.hasChanged());
  }

  // Changed sets the change state
  changed(state: boolean): void {
    this.isChanged = state;
    this.elements.forEach(e => e.changed(state));
  }

  // HasElement checks if the model has a given element id
  hasElement(elementId: string): boolean {
    return this.elements.some(e => e.id === elementId);
  }

  // HasElementType checks if an element type matches given type.
  hasElementType(elementType: string): boolean {
    return this.elements.some(e => e.type.toLowerCase() === elementType.toLowerCase());
  }

  // GetModelIdentifier returns the element which describes the model identifier.
  // Returns the element or null if not found.
  getModelIdentifier(): Element | null {
    return this.elements.find(e => e.isObjectId) || null;
  }

  // GetAttributeIds returns a slice of attribute ids found in the model's .Elements
  getAttributeIds(): string[] {
    return Object.keys(this.attributes);
  }

  // GetElementIds returns a slice of element ids found in the model's .Elements
  getElementIds(): string[] {
    return this.elements.filter(elem => elem.id !== "").map(elem => elem.id);
  }

  // GetPrimaryId returns the primary id
  getPrimaryId(): string {
    const element = this.elements.find(elem => elem.isObjectId);
    return element ? element.id : "";
  }

  // GetGeneratedTypes returns a map of element id and value held by .Generator
  getGeneratedTypes(): { [key: string]: string } {
    const gt: { [key: string]: string } = {};
    this.elements.forEach(elem => {
      if (elem.generator !== "") {
        gt[elem.id] = elem.generator;
      }
    });
    return gt;
  }

  // GetElementById returns an Element from the model's .Elements.
  getElementById(id: string): Element | null {
    return this.elements.find(elem => elem.id === id) || null;
  }

  // NewModel makes sure model id is valid, populates a Model with the identifier element providing
  // returns a Model.
  static newModel(modelId: string): Model {
    const model = new Model();
    model.id = modelId;
    model.description = `... description of ${modelId} goes here ...`;
    model.attributes = {};
    model.elements = [];
    // Make the required element ...
    const element = Element.newElement("id");
    if (element === null) {
      model.errors.push(`Failed to create element with id "id"`);
      return model;
    }
    element.isObjectId = true;
    element.type = "text";
    element.attributes = { required: "true" };
    model.insertElement(0, element);
    model.check();
    return model;
  }

  // Check analyzes the model and makes sure at least one element exists and the
  // model has a single identifier (e.g., "identifier")
  check(): boolean {
    this.errors = [];
    if (!isValidVarname(this.id)) {
      this.errors.push(`Invalid model id, ${this.id}`);
    }
    if (this.elements.length === 0) {
      this.errors.push("Missing model elements.");
      return false;
    }
    let hasModelId = false;
    for (const [i, e] of this.elements.entries()) {
      // Check to make sure each element is valid
      if (!e.check()) {
        this.errors.push(`Error for ${this.id}.${e.id}`);
      }
      if (e.isObjectId) {
        if (hasModelId) {
          this.errors.push(`Duplicate model identifier element (${i}) ${this.id}.${e.id}`);
        }
        hasModelId = true;
      }
    }
    if (!hasModelId) {
      this.errors.push(`Missing required object identifier for model ${this.id}`);
    }
    return this.errors.length === 0;
  }

  // InsertElement will add a new element to model.Elements in the position indicated,
  // It will also set isChanged to true on addition.
  insertElement(pos: number, element: Element): boolean {
    if (!isValidVarname(element.id)) {
      this.errors.push(`Element id is not valid: ${element.id}`);
      return false;
    }
    if (this.hasElement(element.id)) {
      this.errors.push(`Duplicate element id: ${element.id}`);
      return false;
    }
    if (pos < 0) {
      pos = 0;
    }
    if (pos >= this.elements.length) {
      this.elements.push(element);
    } else {
      this.elements.splice(pos, 0, element);
    }
    this.changed(true);
    return true;
  }

  // UpdateElement will update an existing element with element id with the new element.
  updateElement(elementId: string, element: Element): boolean {
    const index = this.elements.findIndex(e => e.id === elementId);
    if (index === -1) {
      this.errors.push(`Element id ${elementId} not found`);
      return false;
    }
    this.elements[index] = element;
    this.changed(true);
    return true;
  }

  // RemoveElement removes an element by id from the model.Elements
  removeElement(elementId: string): boolean {
    const index = this.elements.findIndex(e => e.id === elementId);
    if (index === -1) {
      this.errors.push(`Element id ${elementId} not found`);
      return false;
    }
    this.elements.splice(index, 1);
    this.changed(true);
    return true;
  }

  // RenderElement renders the model using the specified render function.
  async renderElement(name: string, out: WritableStream): Promise<boolean> {
    const fn = this.renderer[name];
    if (fn) {
      await fn(out, this);
      return true;
    } else {
      this.errors.push(`${name} is not a registered rendering function`);
      return false;
    }
  }

  // ToSQLiteScheme takes a model and tries to render a SQLite3 SQL create statement.
  async toSQLiteScheme(out: WritableStream): Promise<void> {
    // Implement the rendering logic here
  }

  // ToHTML takes a model and tries to render an HTML web form
  async toHTML(out: WritableStream): Promise<void> {
    // Implement the rendering logic here
  }

  // ModelToYAML takes a model and interactively prompts to create
  // a YAML model file.
  async modelToYAML(out: WritableStream): Promise<void> {
    // Implement the rendering logic here
  }

  // Register takes a name (string) and a RenderFunc and registers it with the model.
  // Registered names can then be invoked by the register name.
  register(name: string, fn: RenderFunc): void {
    this.renderer[name] = fn;
  }

  // Render takes a registered render WritableStream and register name invoking the function
  // with the model.
  async render(out: WritableStream, name: string): Promise<void> {
    const fn = this.renderer[name];
    if (fn) {
      await fn(out, this);
    } else {
      this.errors.push(`${name} is not a registered rendering function`);
      throw new Error(`${name} is not a registered rendering function`);
    }
  }

  // IsSupportedElementType checks if the element type is supported by Newt, returns true if OK, false if it is not
  isSupportedElementType(eType: string): boolean {
    return Object.keys(this.genElements).includes(eType);
  }

  // Define attaches a type definition (an element generator) and validator for the named type
  define(typeName: string, genElementFn: GenElementFunc, validateFn: ValidateFunc): void {
    this.genElements[typeName] = genElementFn;
    this.validators[typeName] = validateFn;
  }

  // fromObject takes a parameter of {[key: string]: any} and maps it into the attributes of a Model
  fromObject(data: { [key: string]: any }): boolean {
    if (data.id && typeof data.id === "string") {
      this.id = data.id;
    }
    if (data.attributes && typeof data.attributes === "object") {
      this.attributes = data.attributes;
    }
    if (data.description && typeof data.description === "string") {
      this.description = data.description;
    }
    if (data.title && typeof data.title === "string") {
      this.title = data.title;
    }
    if (Array.isArray(data.elements)) {
      this.elements = [];
      for (const elem of data.elements) {
        if (elem instanceof Element) {
          this.elements.push(elem);
        } else if (typeof elem === "object") {
          const newElem = new Element(elem);
          newElem.fromObject(elem);
          this.elements.push(newElem);
        }
      }
    } else {
      this.errors.push("data elements are not an array");
      return false;
    }
    return this.check();
  }
}
