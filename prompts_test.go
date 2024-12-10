package models

import (
	"bytes"
	"testing"
)

// TestGetAnswers tests the interactive prompting methods.
func TestGetAnswers(t *testing.T) {
	txt := "a required true"
	in := bytes.NewBuffer([]byte(txt))
	out := bytes.NewBuffer([]byte{})
	eout := bytes.NewBuffer([]byte{})
	prompt := NewPrompt(in, out, eout)
	answer1, answer2 := prompt.GetAnswers("", "", false)
	if answer1 != "a" {
		t.Errorf("expected %q, got %q (from %q)", "a", answer1, txt)
	}
	if answer2 != "required true" {
		t.Errorf("expected %q, got %q (from %q)", "required true", answer2, txt)
	}
}
