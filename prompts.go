package models

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Prompt holds the elements need to present questions and menus
// to build simple conversational console interfaces
type Prompt struct {
	in   io.Reader
	out  io.Writer
	eout io.Writer

	buf *bufio.Reader
}

// NewPrompt creates a new prompt struct for use with the prompt methods.
// @param in: io.Reader, the source of prompt responses.
// @param out: io.Writer, the place the response is written to (e.g. Stdout)
// @param eout: io.Writer, the place errors are written to (e.g. Stderr)
//
// ```
//
//	prompt := NewPrompt(os.Stdin, os.Stdout, os.Stderr)
//
// ```
func NewPrompt(in io.Reader, out io.Writer, eout io.Writer) *Prompt {
	prompt := new(Prompt)
	prompt.in = in
	prompt.out = out
	prompt.eout = eout
	prompt.buf = bufio.NewReader(in)
	return prompt
}

// GetAnswer display a string and get a response, e.g. "Are you Happy (Y/n)?"
//
// ```
//
//	prompt := NewPrompt(os.Stdin, os.Stdout, os.Stderr)
//	fmt.Println("Are you happy (Y/n)?")
//	answser := prompt.GetAnswer("y", true)
//	if answer == "y" {
//	   fmt.Println("Yeah! We're happy.")
//	} else {
//	   fmt.Println("I'm sad too.")
//	}
//
// ```
func (prompt *Prompt) GetAnswer(defaultAnswer string, lower bool) string {
	prompt.buf.Reset(prompt.in)
	answer, err := prompt.buf.ReadString('\n')
	if err != nil && err != io.EOF {
		return ""
	}
	answer = strings.TrimSpace(answer)
	if answer == "" {
		answer = defaultAnswer
	}
	if lower {
		return strings.ToLower(answer)
	}
	return answer
}

// GetAnswers returns an answer which has an initial verb and an predicate separted
// by a space. E.g. "sleeping now" -> "sleeping" "now" or
// "sleeping at 10:00am" -> "sleeping" "at 10:00am"
//
// ```
//
//	prompt := NewPrompt(os.Stdin, os.Stdout, os.Stderr)
//	fmt.Println("Enter an action and object")
//	verb, object := prompt.GetAnswers("sleeping", "now", false)
//	fmt.Printf("Are you %q? %q?", verb, object)
//
// ```
func (prompt *Prompt) GetAnswers(defaultAnswer string, defaultValue string, lower bool) (string, string) {
	var (
		answer1 string
		answer2 string
	)
	rawAnswer := prompt.GetAnswer(defaultAnswer, false)
	if strings.Contains(rawAnswer, " ") {
		parts := strings.SplitN(rawAnswer, " ", 2)
		answer1, answer2 = parts[0], parts[1]
	} else {
		answer1 = rawAnswer
	}
	answer1 = strings.TrimSpace(answer1)
	answer2 = strings.TrimSpace(answer2)
	if answer1 == "" {
		answer1 = defaultAnswer
	}
	if lower {
		return strings.ToLower(answer1), answer2
	}
	return answer1, answer2
}

// Menu displays header, a list of choices and footer to form a menu.
//
// ```
//
//	choices := []string{
//	   "[1] Vanilla",
//	   "[2] Strawberry",
//	   "[3] Coffee",
//	   "",
//	   "[q]uit to decline ice cream",
//	}
//	prompt := NewPrompt(os.Stdin, os.Stdout, os.Stderr)
//	prompt.Menu(MenuHeader, MenuFooter, choices)
//
// ```
func (prompt *Prompt) Menu(header string, footer string, choices []string) {
	if header != "" {
		fmt.Fprintf(prompt.out, "%s\n", header)
	}
	for _, s := range choices {
		fmt.Fprintf(prompt.out, "\t%s\n", s)
	}
	if footer != "" {
		fmt.Fprintf(prompt.out, "%s\n", footer)
	}
}

// SelectMenu display a menu and return answer(s).
//
// ```
//
//	prompt := NewPrompt(os.Stdin, os.Stdout, os.Stderr)
//	choices := []string{
//	    "[v]anilla",
//	    "[c]hocolate",
//	    "[s]trawberry",
//	}
//	flavor, scoops := prompt.SelectMenuItem(
//	    "Menu: enter flavor letter and number of scoops",
//	    "E.g. three vanilla scoops, "c 3",
//	    choices, "", "", true)
//	fmt.Printf("You selected %q of flavor %q.\n", scoops, flavor)
//
// ```
func (prompt *Prompt) SelectMenu(header string, footer string, choices []string, defaultAnswer string, defaultValue string, lower bool) (string, string) {
	prompt.Menu(header, footer, choices)
	return prompt.GetAnswers(defaultAnswer, defaultValue, lower)
}

// GetDigit returns a numeric answer that is greater than or
// equal than zero and less than the number of choices.
//
// The response integer is normalized to a zero based array.
// E.g. a user inputs the number three it is normalized
// to two so it aligns with the choices present.
//
// ```
//
//	choices := []string{
//	   "Vanilla",
//	   "Strawberry",
//	   "Coffee",
//	   "",
//	   "Decline ice cream",
//	}
//	prompt := NewPrompt(os.Stdin, os.Stdout, os.Stderr)
//	prompt.NumberedMenu("Enter number for choice", "Pick your flavor:", choices)
//	if x, ok := prompt.GetDigit(choices); ok {
//	   fmt.Printf("You selected %d\n", x)
//	} else {
//	   fmt.Printf("No selection\n")
//	}
//
// ```
func (prompt *Prompt) GetDigit(choices []string) (int, bool) {
	maxChoices := 0
	for _, s := range choices {
		if s != "" {
			maxChoices++
		}
	}
	answer := prompt.GetAnswer("", true)
	if answer != "" {
		pos, err := strconv.Atoi(answer)
		if err == nil {
			// Adjust input to zero based array address.
			pos--
			if pos >= 0 && pos < maxChoices {
				return pos, true
			}
		}
	}
	return -1, false
}

// NumberedMenu displays header, a list of choices and footer. The choices
// are numbered when forming the menu. Only non-empty strings are numbered.
//
// ```
//
//	choices := []string{
//	   "Vanilla",
//	   "Strawberry",
//	   "Coffee",
//	   "",
//	   "Decline ice cream",
//	}
//	prompt := NewPrompt(os.Stdin, os.Stdout, os.Stderr)
//	prompt.Menu("Enter a number and press enter", "Pick your flavor:", choices)
//
// ```
func (prompt *Prompt) NumberedMenu(header string, footer string, choices []string) {
	if header != "" {
		fmt.Fprintf(prompt.out, "%s\n", header)
	}
	x := 0
	for _, s := range choices {
		if s != "" {
			x++
			fmt.Fprintf(prompt.out, "\t%d. %s\n", x, s)
		} else {
			fmt.Fprintln(prompt.out, "")
		}
	}
	if footer != "" {
		fmt.Fprintf(prompt.out, "%s\n", footer)
	}
}

// SelectNumberedMenu display a numbered menu and return the intereger and OK status.
//
// ```
//
//	prompt := NewPrompt(os.Stdin, os.Stdout, os.Stderr)
//	choices := []string{
//	    "vanilla",
//	    "chocolate",
//	    "strawberry",
//	    "",
//	    "decline ice cream",
//	}
//	i, ok := prompt.SelectMenuNumber(
//	    "Menu: enter flavor number",
//	    "E.g. 1 for vanilla",
//	    choices, "", "", true)
//	fmt.Printf("You selected %d, status %t\n", i, ok)
//
// ```
func (prompt *Prompt) SelectNumberedMenu(header string, footer string, choices []string) (int, bool) {
	prompt.NumberedMenu(header, footer, choices)
	return prompt.GetDigit(choices)
}
