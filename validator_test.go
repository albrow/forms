// Copyright 2015 Alex Browne.  All rights reserved.
// Use of this source code is governed by the MIT
// license, which can be found in the LICENSE file.

package forms

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

func TestCustomMessage(t *testing.T) {
	data := newData()
	val := data.Validator()
	customMsg := "You forgot to include name!"
	val.Require("name").Message(customMsg)

	if !val.HasErrors() {
		t.Error("Expected an error but got none.")
	} else if val.Messages()[0] != customMsg {
		t.Errorf("Expected custom error message \"%s\" but got \"%s\"", customMsg, val.Messages()[0])
	}
}

func TestCustomField(t *testing.T) {
	data := newData()
	val := data.Validator()
	customField := "person.name"
	val.Require("name").Field(customField)

	if !val.HasErrors() {
		t.Error("Expected an error but got none.")
	} else if val.Fields()[0] != customField {
		t.Errorf("Expected custom field name \"%s\" but got \"%s\"", customField, val.Fields()[0])
	}
}

func TestRequire(t *testing.T) {
	data := newData()
	data.Add("name", "Bob")
	data.Add("age", "25")
	data.Add("color", "")

	val := data.Validator()
	val.Require("name")
	val.Require("age")
	if val.HasErrors() {
		t.Errorf("Expected no errors but got errors: %v", val.Messages())
	}

	val.Require("color")
	val.Require("a")
	if len(val.Messages()) != 2 {
		t.Errorf("Expected 2 validation errors but got %d.", len(val.Messages()))
	}
}

func TestRequireFile(t *testing.T) {
	data := newData()
	val := data.Validator()
	val.RequireFile("file")
	if !val.HasErrors() {
		t.Error("Expected val to have errors because file was not included but got none.")
	}

	fileHeader, err := createTestFileHeader("test_file.txt", []byte{})
	if err != nil {
		t.Error(err)
	}
	data.AddFile("file", fileHeader)
	val = data.Validator()
	val.RequireFile("file")
	if len(val.ErrorMap()) != 1 {
		t.Errorf("Expected val to have exactly one error because file was empty but got %d.", len(val.ErrorMap()))
	} else {
		msg := val.ErrorMap()["file"][0]
		if !strings.Contains(msg, "empty") {
			t.Errorf("Expected message to say file was empty but got: %s.", msg)
		}
	}

	// Create the multipart file header
	// Write actual content to it this time
	fileHeaderWithContent, err := createTestFileHeader("test_file.txt", []byte("Hello!\n"))
	if err != nil {
		t.Error(err)
	}
	data.AddFile("file", fileHeaderWithContent)
	val = data.Validator()
	val.RequireFile("file")
	if val.HasErrors() {
		t.Errorf("Expected val to have no errors but got: %v\n", val.ErrorMap())
	}
}

func createTestFileHeader(filename string, content []byte) (*multipart.FileHeader, error) {
	body := bytes.NewBuffer([]byte{})
	partWriter := multipart.NewWriter(body)
	fileWriter, err := partWriter.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err := fileWriter.Write(content); err != nil {
		return nil, err
	}
	if err := partWriter.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+partWriter.Boundary())
	_, fileHeader, err := req.FormFile("file")
	if err != nil {
		return nil, err
	}
	return fileHeader, nil
}

func TestMinLength(t *testing.T) {
	data := newData()
	data.Add("one", "A")
	data.Add("three", "ABC")
	data.Add("five", "ABC")

	val := data.Validator()
	val.MinLength("one", 1)
	val.MinLength("three", 3)
	if val.HasErrors() {
		t.Error("Expected no errors but got errors: %v", val.Messages())
	}

	val.MinLength("five", 5)
	if len(val.Messages()) != 1 {
		t.Error("Expected a validation error.")
	}
}

func TestMaxLength(t *testing.T) {
	data := newData()
	data.Add("one", "A")
	data.Add("three", "ABC")
	data.Add("five", "ABCDEF")
	val := data.Validator()
	val.MaxLength("one", 1)
	val.MaxLength("three", 3)
	if val.HasErrors() {
		t.Errorf("Expected no errors but got errors: %v", val.Messages())
	}

	val.MaxLength("five", 5)
	if len(val.Messages()) != 1 {
		t.Error("Expected a validation error.")
	}
}

func TestLengthRange(t *testing.T) {
	data := newData()
	data.Add("one-two", "a")
	data.Add("two-three", "abc")
	data.Add("three-four", "ab")
	data.Add("four-five", "abcdef")

	val := data.Validator()
	val.LengthRange("one-two", 1, 2)
	val.LengthRange("two-three", 2, 3)
	if val.HasErrors() {
		t.Errorf("Expected no errors but got errors: %v", val.Messages())
	}

	val.LengthRange("three-four", 3, 4)
	val.LengthRange("four-five", 4, 5)
	if len(val.Messages()) != 2 {
		t.Errorf("Expected 2 validation errors but got %d.", len(val.Messages()))
	}
}

func TestEqual(t *testing.T) {
	data := newData()
	data.Add("password", "password123")
	data.Add("confirmPassword", "password123")
	data.Add("nonMatching", "password1234")

	val := data.Validator()
	val.Equal("password", "confirmPassword")
	if val.HasErrors() {
		t.Errorf("Expected no errors but got errors: %v", val.Messages())
	}

	val.Equal("password", "nonMatching")
	if len(val.Messages()) != 1 {
		t.Error("Expected a validation error.")
	}
}

func TestMatch(t *testing.T) {
	data := newData()
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
		t.Errorf("Expected no errors but got errors: %v", val.Messages())
	}

	val.Match("not-numeric", numericRegex)
	val.Match("not-alpha", alphaRegex)
	if len(val.Messages()) != 2 {
		t.Errorf("Expected 2 validation errors but got %d.", len(val.Messages()))
	}
}

func TestMatchEmail(t *testing.T) {
	data := newData()
	data.Add("email", "abc@example.com")
	data.Add("not-email", "abc.com")
	val := data.Validator()
	val.MatchEmail("email")
	if val.HasErrors() {
		t.Errorf("Expected no errors but got errors: %v", val.Messages())
	}

	val.MatchEmail("not-email")
	val.MatchEmail("nothing")
	if len(val.Messages()) != 2 {
		t.Errorf("Expected 2 validation errors but got %d.", len(val.Messages()))
	}
}

func TestTypeInt(t *testing.T) {
	data := newData()
	data.Add("age", "23")
	data.Add("weight", "not a number")
	val := data.Validator()
	val.TypeInt("age")
	if val.HasErrors() {
		t.Errorf("Expected no errors but got errors: %v", val.Messages())
	}

	val.TypeInt("weight")
	if len(val.Messages()) != 1 {
		t.Errorf("Expected 1 validation errors but got %d.", len(val.Messages()))
	}
}

func TestTypeFloat(t *testing.T) {
	data := newData()
	data.Add("age", "23")
	data.Add("weight", "155.8")
	data.Add("favoriteNumber", "not a number")
	val := data.Validator()
	val.TypeFloat("age")
	val.TypeFloat("weight")
	if val.HasErrors() {
		t.Errorf("Expected no errors but got errors: %v", val.Messages())
	}

	val.TypeFloat("favoriteNumber")
	if len(val.Messages()) != 1 {
		t.Errorf("Expected 1 validation errors but got %d.", len(val.Messages()))
	}
}

func TestTypeBool(t *testing.T) {
	data := newData()
	data.Add("cool", "true")
	data.Add("fun", "false")
	data.Add("yes", "not a boolean")
	val := data.Validator()
	val.TypeBool("cool")
	val.TypeBool("fun")
	if val.HasErrors() {
		t.Errorf("Expected no errors but got errors: %v", val.Messages())
	}

	val.TypeBool("yes")
	if len(val.Messages()) != 1 {
		t.Errorf("Expected 1 validation errors but got %d.", len(val.Messages()))
	}
}

func TestGreater(t *testing.T) {
	data := newData()
	data.Add("one", "1")
	data.Add("three", "3")
	val := data.Validator()
	val.Greater("one", -1)
	val.Greater("three", 2)
	if val.HasErrors() {
		t.Errorf("Expected no errors but got errors: %v", val.Messages())
	}

	val.Greater("one", 1)
	val.Greater("three", 4)
	if len(val.Messages()) != 2 {
		t.Errorf("Expected 2 validation errors but got %d.", len(val.Messages()))
	}
}

func TestGreaterOrEqual(t *testing.T) {
	data := newData()
	data.Add("one", "1")
	data.Add("three", "3")
	val := data.Validator()
	val.GreaterOrEqual("one", 1)
	val.GreaterOrEqual("three", 2)
	if val.HasErrors() {
		t.Errorf("Expected no errors but got errors: %v", val.Messages())
	}

	val.GreaterOrEqual("one", 2)
	val.GreaterOrEqual("three", 4)
	if len(val.Messages()) != 2 {
		t.Errorf("Expected 2 validation errors but got %d.", len(val.Messages()))
	}
}

func TestLess(t *testing.T) {
	data := newData()
	data.Add("one", "1")
	data.Add("three", "3")
	val := data.Validator()
	val.Less("one", 2)
	val.Less("three", 4)
	if val.HasErrors() {
		t.Errorf("Expected no errors but got errors: %v", val.Messages())
	}

	val.Less("one", -1)
	val.Less("three", 3)
	if len(val.Messages()) != 2 {
		t.Errorf("Expected 2 validation errors but got %d.", len(val.Messages()))
	}
}

func TestLessOrEqual(t *testing.T) {
	data := newData()
	data.Add("one", "1")
	data.Add("three", "3")
	val := data.Validator()
	val.LessOrEqual("one", 1)
	val.LessOrEqual("three", 4)
	if val.HasErrors() {
		t.Errorf("Expected no errors but got errors: %v", val.Messages())
	}

	val.LessOrEqual("one", -1)
	val.LessOrEqual("three", 2)
	if len(val.Messages()) != 2 {
		t.Errorf("Expected 2 validation errors but got %d.", len(val.Messages()))
	}
}

func TestAcceptFileExts(t *testing.T) {
	data := newData()
	fileHeader, err := createTestFileHeader("test_file.txt", []byte{})
	if err != nil {
		t.Error(err)
	}
	data.AddFile("file", fileHeader)
	val := data.Validator()
	val.AcceptFileExts("file", "txt")
	if val.HasErrors() {
		t.Errorf("Expected no errors for the single allowed ext case but got %v\n", val.ErrorMap())
	}
	val = data.Validator()
	val.AcceptFileExts("file", "txt", "jpg")
	if val.HasErrors() {
		t.Errorf("Expected no errors for the multiple allowed ext case but got %v\n", val.ErrorMap())
	}
	val = data.Validator()
	val.AcceptFileExts("foo", "txt")
	if val.HasErrors() {
		t.Errorf("Expected no errors for the not-provided field case but got %v\n", val.ErrorMap())
	}

	// use a table-driven test here
	table := []struct {
		allowedExts             []string
		expectedMessageContains []string
	}{
		{
			allowedExts: []string{"jpg"},
			expectedMessageContains: []string{
				"The file extension .txt is not allowed.",
				"include: jpg",
			},
		},
		{
			allowedExts: []string{"jpg", "png"},
			expectedMessageContains: []string{
				"The file extension .txt is not allowed.",
				"include: jpg and png",
			},
		},
		{
			allowedExts: []string{"jpg", "png", "gif"},
			expectedMessageContains: []string{
				"The file extension .txt is not allowed.",
				"include: jpg, png, and gif",
			},
		},
	}
	for i, test := range table {
		val = data.Validator()
		val.AcceptFileExts("file", test.allowedExts...)
		if !val.HasErrors() {
			t.Errorf("Expected val to have errors for test case %d but got none.", i)
			continue // avoid index out-of-bounds error in proceeding lines
		}
		gotMsg := val.ErrorMap()["file"][0]
		for _, expectedMsg := range test.expectedMessageContains {
			if !strings.Contains(gotMsg, expectedMsg) {
				t.Errorf(`Expected error in case %d to contain "%s" but it did not.%sGot: "%s"`, i, expectedMsg, "\n", gotMsg)
			}
		}
	}
}

func ExampleValidator() {
	// Construct a request object for example purposes only.
	// Typically you would be using this inside a http.HandlerFunc,
	// not constructing your own request.
	req, _ := http.NewRequest("GET", "/", nil)
	values := url.Values{}
	values.Add("name", "Bob")
	values.Add("age", "25")
	req.PostForm = values
	req.Header.Set("Content-Type", "form-urlencoded")

	// Parse the form data.
	data, _ := Parse(req)

	// Validate the data.
	val := data.Validator()
	val.Require("name")
	val.MinLength("name", 4)
	val.Require("age")

	// Here's how you can change the error message or field name
	val.Require("retired").Field("retired_status").Message("Must specify whether or not person is retired.")

	// Check for validation errors and print them if there are any.
	if val.HasErrors() {
		fmt.Printf("%#v\n", val.Messages())
	}

	// Output:
	// []string{"name must be at least 4 characters long.", "Must specify whether or not person is retired."}
}
