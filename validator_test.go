package data

import (
	"regexp"
	"testing"
)

func TestRequire(t *testing.T) {
	data := Data{}
	data.Add("name", "Bob")
	data.Add("age", "25")
	data.Add("color", "")

	val := data.Validator()
	val.Require("name")
	val.Require("age")
	if val.HasErrors() {
		t.Errorf("Expected no errors but got errors: %v", val.Errors)
	}

	val.Require("color")
	val.Require("a")
	if len(val.Errors) != 2 {
		t.Errorf("Expected 2 validation errors but got %d.", len(val.Errors))
	}
}

func TestMinLength(t *testing.T) {
	data := Data{}
	data.Add("one", "A")
	data.Add("three", "ABC")
	data.Add("five", "ABC")

	val := data.Validator()
	val.MinLength("one", 1)
	val.MinLength("three", 3)
	if val.HasErrors() {
		t.Error("Expected no errors but got errors: %v", val.Errors)
	}

	val.MinLength("five", 5)
	if len(val.Errors) != 1 {
		t.Error("Expected a validation error.")
	}
}

func TestMaxLength(t *testing.T) {
	data := Data{}
	data.Add("one", "A")
	data.Add("three", "ABC")
	data.Add("five", "ABCDEF")
	val := data.Validator()
	val.MaxLength("one", 1)
	val.MaxLength("three", 3)
	if val.HasErrors() {
		t.Errorf("Expected no errors but got errors: %v", val.Errors)
	}

	val.MaxLength("five", 5)
	if len(val.Errors) != 1 {
		t.Error("Expected a validation error.")
	}
}

func TestLengthRange(t *testing.T) {
	data := Data{}
	data.Add("one-two", "a")
	data.Add("two-three", "abc")
	data.Add("three-four", "ab")
	data.Add("four-five", "abcdef")

	val := data.Validator()
	val.LengthRange("one-two", 1, 2)
	val.LengthRange("two-three", 2, 3)
	if val.HasErrors() {
		t.Errorf("Expected no errors but got errors: %v", val.Errors)
	}

	val.LengthRange("three-four", 3, 4)
	val.LengthRange("four-five", 4, 5)
	if len(val.Errors) != 2 {
		t.Errorf("Expected 2 validation errors but got %d.", len(val.Errors))
	}
}

func TestEqual(t *testing.T) {
	data := Data{}
	data.Add("password", "password123")
	data.Add("confirmPassword", "password123")
	data.Add("nonMatching", "password1234")

	val := data.Validator()
	val.Equal("password", "confirmPassword")
	if val.HasErrors() {
		t.Errorf("Expected no errors but got errors: %v", val.Errors)
	}

	val.Equal("password", "nonMatching")
	if len(val.Errors) != 1 {
		t.Error("Expected a validation error.")
	}
}

func TestMatch(t *testing.T) {
	data := Data{}
	data.Add("numeric", "123")
	data.Add("alpha", "abc")
	data.Add("not-numeric", "123a")
	data.Add("not-alpha", "abc1")

	val := data.Validator()
	numericRegex := regexp.MustCompile("^[0-9]+$")
	alphaRegex := regexp.MustCompile("^[a-zA-Z]+$")
	val.Match("numeric", numericRegex)
	val.Match("alpha", alphaRegex)
	if val.HasErrors() {
		t.Errorf("Expected no errors but got errors: %v", val.Errors)
	}

	val.Match("not-numeric", numericRegex)
	val.Match("not-alpha", alphaRegex)
	if len(val.Errors) != 2 {
		t.Errorf("Expected 2 validation errors but got %d.", len(val.Errors))
	}
}

func TestMatchEmail(t *testing.T) {
	data := Data{}
	data.Add("email", "abc@example.com")
	data.Add("not-email", "abc.com")
	val := data.Validator()
	val.MatchEmail("email")
	if val.HasErrors() {
		t.Errorf("Expected no errors but got errors: %v", val.Errors)
	}

	val.MatchEmail("not-email")
	val.MatchEmail("nothing")
	if len(val.Errors) != 2 {
		t.Errorf("Expected 2 validation errors but got %d.", len(val.Errors))
	}
}
