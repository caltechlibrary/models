
# Tutorial

The models package provides a means of describing a data model based on HTML elements expressed as YAML. A given model can then be rendered to HTML or SQLite 3 SQL Schema. In the following tutorial you will create a simple "guestbook" model and render it to HTML and SQLite 3 SQL scheme.

## Steps

The models package includes a demonstration program called `modelgen`. We'll use it in this tutorial.

1. Generate a model YAML file using "model" action
2. Generate an HTML web form using the "html" action
3. Generate SQL schema for SQLite 3 using the "sqlite" action

## Generating our model YAML file

This is done using the `modelgen` command line program providing the "model" action along with a YAML filename to create (e.g. "guestbook.yaml")
The "model" action is interactive. You will be presented with a series of text menus. Pressing the enter key without selecting a menu choice will
move you to the next menu step.  Most menus you will either type a single letter or digit to then be prompted to complete the task.

~~~
modelgen model guestbook.yaml
~~~

When you type in the command above you will be presented with the following menu.

~~~
Manage Model Metadata
	id: "guestbook"
	description: "... description of \"guestbook\" goes here ..."
Menu [i]d, [d]escription or press enter when done
~~~

Let's change the description to "A guestbook demo".  To do this press "d" followed by the enter key.

You should now be prompted to enter a new description. Type "A guestbook demo" (without quotes), then press the enter key.

You should see something like the following displayed.

~~~
Manage Model Metadata
	id: "guestbook"
	description: "... description of \"guestbook\" goes here ..."
Menu [i]d, [d]escription or press enter when done
d
Enter Description: A guestbook demo
Manage Model Metadata
	id: "guestbook"
	description: "A guestbook demo"
Menu [i]d, [d]escription or press enter when done
~~~

Notice that the screen is NOT cleared between menus and prompts. The reason for this is it allows you to see what has gone before as well as ensure it'll run in most terminals regardless of operating system.

Press enter again. You should now see the following appear at the bottom of the screen.

~~~

Manage guestbook elements
	id
Choices [a]dd, [m]odify, [r]emove or press enter when done
~~~

There are three actions that can be taken on your element list. You can add an element by typing "a" and pressing enter. You can modify and element by typing "m" and pressing enter. You can remove an element by typing "r" and enter.

The menu will then prompt you for more information such as the element identifier or name.

When you "add" an element a plain old "text" element is appended to the list.  If you want to some other type of element then you "modify" after adding it.

The following element types are currently supported in the models package and are based on their HTML input types.

- [checkbox](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/checkbox)
- [color](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/color)
- [date](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/date)
- [datetime-local](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/datetime-local)
- [email](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/email)
- [month](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/month)
- [number](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/number)
- [password](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/password)
- [radio](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/radio)
- [range](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/range)
- [tel](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/tel)
- [text](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/text)
- [textarea](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/textarea)
- [time](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/time)
- [url](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/url)

Two additional types are defined, "orcid" and "isni" these use the HTML pattern attribute to control form.  In a Go program using the models package a generator and validator is available for each of these types. The generator populates the `Element` type in the model while the validator will compare a string with the properties described in the element.  The demo program modelgen does not use these features.

When you create a model an "id" element is always generated. If you do not need it you can use "r" to remove it.

Let's add two text elements, "name" and "msg". First add the name, type "a" and press the enter key. At the prompt add "name" (without the quotes).

The result should look like the following.

~~~
a
Enter element id to add: name
Manage guestbook elements
	id
	name
Choices [a]dd, [m]odify, [r]emove or press enter when done
~~~

Now type "a" and press the enter key again, this time at the prompt add "msg" (without the quotes).

~~~
a
Enter element id to add: msg
Manage guestbook elements
	id
	name
	msg
Choices [a]dd, [m]odify, [r]emove or press enter when done
~~~

As mentioned by default a "text" input type is added. This is fine for `name` and `msg`. A name in the guest book should be required so we need to "modify" "name".  Type "m" and press enter, then type "name" (without quotes) and press enter. You should now see something like this.

~~~
m
Enter element id to modify: name
Manage guestbook.name element
	id name
	type text
	label Name
	pattern 
	attributes:
		name
	object identifier? false
Choices [t]ype, [l]abel, [o]bject identifier, [p]attern, [a]ttributes, or press enter when done
~~~

You are now able to modify the "name" element.  To make an element required we want to add `require="true"` to
the attribute list.  Type "a" and press enter". You should see something like this.

~~~
a
Modify element guestbook.name attributes
	name -> "name"
Choices [a]dd, [m]odify, [r]emove or press enter when done
~~~

This is a list of `guestbook.name` defined attributes. Currently there is only a `guestbook.name.name` attribute. Type "a" and press enter. Then enter "require" (without quotes) and press enter.

~~~
a
Enter attribute name: required
Modify element guestbook.name attributes
	name -> "name"
	required -> ""
Choices [a]dd, [m]odify, [r]emove or press enter when done
~~~

Notice we now have a "required" attribute but the value is empty. We want the value to be "true". So once again we type "m" and when prompted type "required" before typing "true" (without quotes).

~~~
m
Enter attribute name: required
Enter required's value: true
Modify element guestbook.name attributes
	required -> "true"
	name -> "name"
Choices [a]dd, [m]odify, [r]emove or press enter when done
~~~

The name attribute is now a required element.  If we press enter again
you should see the full element description.

~~~
Manage guestbook.name element
	id name
	type text
	label Name
	pattern 
	attributes:
		name,
		required
	object identifier? false
Choices [t]ype, [l]abel, [o]bject identifier, [p]attern, [a]ttributes, or press enter when done
~~~

Press enter one more time we are back to the elements list.

~~~
Manage guestbook elements
	id
	name
	msg
Choices [a]dd, [m]odify, [r]emove or press enter when done
~~~

We have our basic guestbook model completed.  Press enter and you should be prompted to save the model file.

~~~
Save guestbook.yaml (Y/n)?
~~~

If you reply "y" then enter or just press enter it'll save the changes and exit the program.  If you answer anything else it'll exit the program without saving.

We've created our first model YAML file.

You can see the model using the `cat` command on macOS and Linux or the `type` command on Windows.
Your guestbook.yaml should look something like this.

~~~yaml
id: guestbook
description: A guestbook demo
elements:
    - type: text
      id: id
      attributes:
        required: "true"
      is_primary_id: true
    - type: text
      id: name
      attributes:
        name: name
        required: "true"
      label: Name
    - type: text
      id: msg
      attributes:
        name: msg
      label: Msg
~~~

Congratulations! Model done.

## Generating HTML

A simple HTML form can be generated with the following command now that we have our model YAML file.

~~~
modelgen html guestbook.yaml
~~~

Here's an example of the output.

~~~html
<!-- guestbook: A guestbook demo -->
<form id="guestbook">
  <div class="guestbook-id"><input class="guestbook-id" type="text" id="id" required ></div>
  <div class="guestbook-name"><label class="guestbook-name" set="name">Name</label> <input class="guestbook-name" type="text" id="name" name="name" required ></div>
  <div class="guestbook-msg"><label class="guestbook-msg" set="msg">Msg</label> <input class="guestbook-msg" type="text" id="msg" name="msg" ></div>
  <div class="guestbook-submit"><input class="guestbook-submit" type="submit" value="submit"> <input class="guestbook-submit" type="reset" value="cancel"></div>
</form>
~~~

## Generating SQL for SQLite Schema

A simple SQLite 3 Schema can be create via SQL. Use the following command using the "sqlite" action.

~~~
modedgen sqlite guestbook.yaml
~~~

Here's an example of the output.

~~~sql
--
-- guestbook
--
-- A guestbook demo
create table guestbook if not exists (
  id text primary key,
  name text,
  msg text
);
~~~
