import { Model } from "./model.ts";
import { Element } from "./element.ts";

// ModelToHTML takes a model and renders an input form. The form is not
// populated with values, but that could be done easily via JavaScript and DOM calls.
export async function modelToHTML(out: WritableStream, model: Model): Promise<void> {
  const writer = out.getWriter();
  const encoder = new TextEncoder();

  // Include the description as an HTML comment.
  // Write opening form element
  if (model.id !== "") {
    await writer.write(encoder.encode(`<!-- ${model.id}: ${model.description} -->\n`));
    await writer.write(encoder.encode(`<form id="${model.id}"`));
  } else {
    await writer.write(encoder.encode(`<!-- ${model.description} -->\n`));
    await writer.write(encoder.encode(`<form`));
  }

  for (const [k, v] of Object.entries(model.attributes)) {
    switch (k) {
      case "checked":
        await writer.write(encoder.encode(" checked"));
        break;
      case "required":
        await writer.write(encoder.encode(" required"));
        break;
      default:
        await writer.write(encoder.encode(` ${k}="${v}"`));
    }
  }

  const cssBaseClass = model.id.toLowerCase().replaceAll(" ", "_");
  await writer.write(encoder.encode(">\n"));

  for (const elem of model.elements) {
    await elementToHTML(writer, encoder, cssBaseClass, elem);
  }

  if (!model.hasElementType("submit")) {
    const cssName = `${cssBaseClass}-submit`;
    await writer.write(encoder.encode(`  <div class="${cssName}"><input class="${cssName}" type="submit" value="submit"> <input class="${cssName}" type="reset" value="cancel"></div>\n`));
  }

  // Write closing form element
  await writer.write(encoder.encode("</form>\n"));
  await writer.close();
}

// ElementToHTML renders an individual element as HTML, including the label as well as the input element.
async function elementToHTML(writer: WritableStreamDefaultWriter, encoder: TextEncoder, cssBaseClass: string, elem: Element): Promise<void> {
  const cssClass = `${cssBaseClass}-${elem.id.toLowerCase()}`;
  await writer.write(encoder.encode(`  <div class="${cssClass}">`));

  // Generate a default label if not provided
  const labelText = elem.label || elem.id.charAt(0).toUpperCase() + elem.id.slice(1);
  const name = elem.attributes["name"] || elem.id;
  await writer.write(encoder.encode(`<label class="${cssClass}" for="${name}">${labelText}</label> `));

  switch (elem.type.toLowerCase()) {
    case "textarea":
      await writer.write(encoder.encode(`<textarea class="${cssClass}"`));
      break;
    case "button":
      await writer.write(encoder.encode(`<button class="${cssClass}"`));
      break;
    default:
      await writer.write(encoder.encode(`<input class="${cssClass}" type="${elem.type}"`));
  }

  if (elem.id !== "") {
    await writer.write(encoder.encode(` id="${elem.id}"`));
  }

  for (const [k, v] of Object.entries(elem.attributes)) {
    switch (k) {
      case "checked":
        await writer.write(encoder.encode(" checked"));
        break;
      case "required":
        await writer.write(encoder.encode(" required"));
        break;
      default:
        await writer.write(encoder.encode(` ${k}="${v}"`));
    }
  }

  switch (elem.type.toLowerCase()) {
    case "button":
      await writer.write(encoder.encode(`>${elem.label}</button>`));
      break;
    case "textarea":
      await writer.write(encoder.encode("></textarea>"));
      break;
    default:
      await writer.write(encoder.encode(">"));
  }

  await writer.write(encoder.encode("</div>\n"));
}
