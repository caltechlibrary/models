
# Tutorial

The models package provides a means of describing a data model based on HTML elements expressed as YAML. A given model can then be rendered to HTML or SQLite 3 SQL Schema. In the followng totorial you will create a simple "guestbook" model and render it to HTML and SQLite 3 SQL scheme.

## Steps

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

You should now be prompted to enter a new description. If you just press enter then the description will not be changed.
When prompted type in "A guestbook demo" then press the enter key.  You should see something like the following displayed.

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

Notice that the screen is NOT cleared between actions.  The interactive interface is based around the simple notion of prompt and
response. It should run easily in most POSIX terminals because it doesn't do anything fancy. Now press the enter key again and 
you will see a new menu added at the bottom of the screen.

~~~
Manage guestbook attributes (none required)


Menu [a]dd, [m]odify, [r]emove or press enter when done
~~~

A model may have attributes. These are rendered into the HTML generated form but are not required. They are present primarily for compatibility
with GitHub YAML Issue Templates, the inspiration for the models. package. Press enter again and the new Menu should appear to define the elements
of your data model.

~~~
Manage guestbook elements

	1: id

Menu [a]dd, [m]odify, [r]emove or press enter when done
~~~

There are three actions that can be taken on your element list. You can add an element "a", modify "m" and element and remove "r" an element. The action
"m", "r" can be followed by a space and the number of the element you wish to change.

When you "add" an element a plain old "text" element is appended to the list.  If you want to something different then you "modify" the element to meet your needs.

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

Two additional types are defined, "orcid" and "isni" these use the HTML pattern attribute to control form.  In a Go program using the models package a
generator and validator is available for each of these types. The generator populates the `Element` type in the model while the validator will compare
a string with the properties described in the element.  The demo program modelgen does not use these features.

When you create a model a "id" element is always generated. If you do not need it you can use "r" to remove it.

Let's add two text elements, "name" and "msg". First add the name, type "e" and press the enter key. At the prompt add "name" (without the quotes).
The result should look like the following.

~~~
a
Enter element id to add: name
Manage guestbook elements

	1: id
	2: name

Menu [a]dd, [m]odify, [r]emove or press enter when done
~~~

Now type "a" and press the enter key again, this time at the prompt add "msg" (without the quotes).

~~~
a
Enter element id to add: msg
Manage guestbook elements

	1: id
	2: name
	3: msg

Menu [a]dd, [m]odify, [r]emove or press enter when done
~~~

As mentioned by default a "text" input type is added. This is fine for `name` and `msg` but we should require a name when 
adding an entry to the guest book. We do that by adding an attribute to name. To do that we need to modify name. Type "m" and then press the entery key.
You are then prompted for the element you want to modify. In our case it is "name" (without the quotes) followed by pressing the enter key.
You should now see something like the following.

~~~
m
Enter element id to modify: name
Manage guestbook.name element

	id name
	type text
	label Name
	pattern
	attributes name
	object identifier? false

Menu [t]ype, [l]abel, [o]bject identifier, [p]attern, [a]ttributes, or press enter when done
~~~

We want to add an attribute so the next step is to type "a" followed by the enter key. We will then be shown the list of attributes that are defined.

~~~
a
Manage guestbook.name attributes

	1: name -> "name"

Menu [a]dd, [m]odify, [r]emove or press enter when done
~~~

Since we need to add our "required" attribute type "a" followed by the enter key. At the prompt type in "required" (without the quotes).

~~~
a
Enter attribute name: required
Manage guestbook.name attributes

	1: name -> "name"
	2: required -> ""

Menu [a]dd, [m]odify, [r]emove or press enter when done
~~~

Notice that "required" was added to the list but it's value is an empty string. We want the value to be "true". We modify it by typing in "m" and pressing enter. Whenm prompted we type in "required" and press enter. We are then prompted for a value, type in "true" (without the quotes) and press enter.

~~~
Enter attribute name: required
Enter required's value: true
Manage guestbook.name attributes

	1: name -> "name"
	2: required -> "true"

Menu [a]dd, [m]odify, [r]emove or press enter when done
~~~

We're done making changes to our model.  When you're done you generally just press the enter key by itself. Give it a try.
You should now see something like the following.

~~~
Manage guestbook.name element

	id name
	type text
	label Name
	pattern
	attributes name,
		required
	object identifier? false

Menu [t]ype, [l]abel, [o]bject identifier, [p]attern, [a]ttributes, or press enter when done
~~~

Press enter again and you should not see the follow (we're backing out of the menus).

~~~
Manage guestbook elements

	1: id
	2: name
	3: msg

Menu [a]dd, [m]odify, [r]emove or press enter when done
~~~

Press enter again and you're be prompted to save the model. The default answer is "Y" so you can just press enter again.

~~~

Save guestbook.yaml (Y/n)?
~~~

We've created our first model YAML file. Congraduation.

You can see the model using the `cat` command on macOS and Linux or the `type` command on Windows.
Your guestbook.yaml should look something like this.

~~~
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

## Generating HTML

A simple HTML form can be generated with the following command now that we have our model YAML file.

~~~
modelgen html guestbook.yaml
~~~

Here's an example of the output.

~~~html
<form id="guestbook">
  <div class="guestbook-id"><input class="guestbook-id" type="text" id="id" readonly="true" ></div>
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
