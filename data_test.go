// Copyright 2015 Alex Browne.  All rights reserved.
// Use of this source code is governed by the MIT
// license, which can be found in the LICENSE file.

package forms

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	data := newData()
	data.Values = map[string][]string{
		"name":       []string{"bob", "bill"},
		"profession": []string{"plumber"},
	}

	table := []struct {
		key      string
		expected string
	}{
		{
			key:      "name",
			expected: "bob",
		},
		{
			key:      "profession",
			expected: "plumber",
		},
		{
			key:      "favoriteColor",
			expected: "",
		},
	}

	for _, test := range table {
		got := data.Get(test.key)
		if got != test.expected {
			t.Errorf("%s was incorrect. Expected %s, but got %s.\n", test.key, test.expected, got)
		}
	}
}

func TestGetInt(t *testing.T) {
	data := newData()
	data.Values = map[string][]string{
		"age":    []string{"25", "33"},
		"weight": []string{"155"},
	}

	table := []struct {
		key      string
		expected int
	}{
		{
			key:      "age",
			expected: 25,
		},
		{
			key:      "weight",
			expected: 155,
		},
		{
			key:      "height",
			expected: 0,
		},
	}

	for _, test := range table {
		got := data.GetInt(test.key)
		if got != test.expected {
			t.Errorf("%s was incorrect. Expected %d, but got %d.\n", test.key, test.expected, got)
		}
	}
}

func TestGetFloat(t *testing.T) {
	data := newData()
	data.Values = map[string][]string{
		"age":    []string{"25.7", "33"},
		"weight": []string{"42"},
	}

	table := []struct {
		key      string
		expected float64
	}{
		{
			key:      "age",
			expected: 25.7,
		},
		{
			key:      "weight",
			expected: 42.0,
		},
		{
			key:      "height",
			expected: 0.0,
		},
	}

	for _, test := range table {
		got := data.GetFloat(test.key)
		if got != test.expected {
			t.Errorf("%s was incorrect. Expected %f, but got %f.\n", test.key, test.expected, got)
		}
	}
}

func TestGetBool(t *testing.T) {
	data := newData()
	data.Values = map[string][]string{
		"retired":         []string{"true", "false"},
		"leftHanded":      []string{"0"},
		"collegeGraduate": []string{"1"},
	}

	table := []struct {
		key      string
		expected bool
	}{
		{
			key:      "retired",
			expected: true,
		},
		{
			key:      "leftHanded",
			expected: false,
		},
		{
			key:      "collegeGraduate",
			expected: true,
		},
		{
			key:      "sagittarius",
			expected: false,
		},
	}

	for _, test := range table {
		got := data.GetBool(test.key)
		if got != test.expected {
			t.Errorf("%s was incorrect. Expected %t, but got %t.\n", test.key, test.expected, got)
		}
	}
}

func TestBytes(t *testing.T) {
	data := newData()
	data.Values = map[string][]string{
		"name":       []string{"bob", "bill"},
		"profession": []string{"plumber"},
	}

	table := []struct {
		key      string
		expected []byte
	}{
		{
			key:      "name",
			expected: []byte("bob"),
		},
		{
			key:      "profession",
			expected: []byte("plumber"),
		},
		{
			key:      "favoriteColor",
			expected: []byte(""),
		},
	}

	for _, test := range table {
		got := data.GetBytes(test.key)
		if len(got) == 0 && len(test.expected) == 0 {
			// do nothing
			// reflect.DeepEqual doesn't like when both lengths are zero, but it should pass.
		} else if !reflect.DeepEqual(got, test.expected) {
			t.Errorf("%s was incorrect. Expected %v, but got %v.\n", test.key, test.expected, got)
		}
	}
}

func TestCreateFromMap(t *testing.T) {
	m := map[string]string{
		"name":          "bob",
		"age":           "25",
		"favoriteColor": "fuchsia",
	}
	data := CreateFromMap(m)

	table := []struct {
		key      string
		expected string
	}{
		{
			key:      "name",
			expected: "bob",
		},
		{
			key:      "age",
			expected: "25",
		},
		{
			key:      "dreamJob",
			expected: "",
		},
	}

	for _, test := range table {
		got := data.Get(test.key)
		if got != test.expected {
			t.Errorf("%s was incorrect. Expected %s, but got %s.\n", test.key, test.expected, got)
		}
	}
}

func TestGetStringsSplit(t *testing.T) {
	data := newData()
	data.Values = map[string][]string{
		"children":       []string{"martha,bill,jane", "adam,julia"},
		"favoriteColors": []string{"blue%20green%20fuchsia"},
	}

	table := []struct {
		key       string
		delim     string
		expecteds []string
	}{
		{
			key:       "children",
			delim:     ",",
			expecteds: []string{"martha", "bill", "jane"},
		},
		{
			key:       "favoriteColors",
			delim:     "%20",
			expecteds: []string{"blue", "green", "fuchsia"},
		},
		{
			key:       "height",
			delim:     "-",
			expecteds: []string{},
		},
	}

	for _, test := range table {
		gots := data.GetStringsSplit(test.key, test.delim)
		if len(gots) == 0 && len(test.expecteds) == 0 {
			// do nothing
			// reflect.DeepEqual doesn't like when both lengths are zero, but it should pass.
		} else if !reflect.DeepEqual(gots, test.expecteds) {
			t.Errorf("%s was incorrect. Expected %v, but got %v.\n", test.key, test.expecteds, gots)
		}
	}
}

func TestParseUrlEncoded(t *testing.T) {
	// Construct a urlencoded form request
	// Add some simple key-value params to the form
	fieldData := map[string]string{
		"name":           "Bob",
		"age":            "25",
		"favoriteNumber": "99.99",
		"leftHanded":     "true",
	}
	values := url.Values{}
	for fieldname, value := range fieldData {
		values.Add(fieldname, value)
	}
	req, err := http.NewRequest("POST", "/", strings.NewReader(values.Encode()))
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Parse the request
	d, err := Parse(req)
	if err != nil {
		t.Error(err)
	}
	testBasicFormFields(t, d)
}

func TestParseMultipart(t *testing.T) {
	// Construct a multipart request
	body := bytes.NewBuffer([]byte{})
	form := multipart.NewWriter(body)

	// Add some simple key-value params to the form
	fieldData := map[string]string{
		"name":           "Bob",
		"age":            "25",
		"favoriteNumber": "99.99",
		"leftHanded":     "true",
	}
	for fieldname, value := range fieldData {
		if err := form.WriteField(fieldname, value); err != nil {
			panic(err)
		}
	}

	// Add a file to the form
	testFile, err := os.Open("test_file.txt")
	if err != nil {
		t.Error(err)
	}
	defer testFile.Close()
	fileWriter, err := form.CreateFormFile("file", "test_file.txt")
	if err != nil {
		panic(err)
	}
	if _, err := io.Copy(fileWriter, testFile); err != nil {
		panic(err)
	}
	// Close the form to finish writing
	if err := form.Close(); err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+form.Boundary())

	// Parse the request
	d, err := Parse(req)
	if err != nil {
		t.Error(err)
	}
	testBasicFormFields(t, d)

	// Next test that the file was parsed correctly
	if !d.FileExists("file") {
		t.Error("Expected FileExists() to return true but it returned false.")
	}
	header := d.GetFile("file")
	if header == nil {
		t.Error("Exected GetFile() to return a *multipart.FileHeader but got nil.")
	}
	if header.Filename != "test_file.txt" {
		t.Errorf(`Expected header.Filename to equal "test_file.txt" but got %s`, header.Filename)
	}
	file, err := header.Open()
	if err != nil {
		t.Error(err)
	}
	gotBytes, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error(err)
	}
	if string(gotBytes) != "Hello!" {
		t.Errorf(`Expected file contents when read directly to be "Hello!" but got %s`, string(gotBytes))
	}
	gotBytes, err = d.GetFileBytes("file")
	if err != nil {
		t.Error(err)
	}
	if string(gotBytes) != "Hello!" {
		t.Errorf(`Expected GetFileBytes("file") to return "Hello!" but got %s`, string(gotBytes))
	}
}

// Used for testing multipart and urlencoded form data, since both tests expect the same data
// to be present.
func testBasicFormFields(t *testing.T, d *Data) {
	// use a table for testing
	fields := []struct {
		key      string
		got      interface{}
		expected interface{}
	}{
		{
			key:      "name",
			got:      d.Get("name"),
			expected: "Bob",
		},
		{
			key:      "age",
			got:      d.GetInt("age"),
			expected: 25,
		},
		{
			key:      "favoriteNumber",
			got:      d.GetFloat("favoriteNumber"),
			expected: 99.99,
		},
		{
			key:      "leftHanded",
			got:      d.GetBool("leftHanded"),
			expected: true,
		},
	}
	for _, test := range fields {
		if !reflect.DeepEqual(test.got, test.expected) {
			t.Errorf("%s was incorrect. Expected %v, but got %v.\n", test.key, test.expected, test.got)
		}
	}
}

type jsonData struct {
	Name     string             `json:"name"`
	Age      int                `json:"age"`
	Cool     bool               `json:"cool"`
	Aptitude string             `json:"aptitude"`
	Location map[string]float64 `json:"location"`
	Things   []string           `json:"things"`
}

func TestParseJSON(t *testing.T) {
	// Construct and parse a json request
	input := `{
		"name": "bob",
		"age": 25,
		"cool": true,
		"aptitude": null,
		"location": {"latitude": 123.456, "longitude": 948.123},
		"things": ["a", "b", "c"]
	}`
	body := bytes.NewBuffer([]byte(input))
	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Content-Type", "application/json")
	d, err := Parse(req)
	if err != nil {
		t.Error(err)
	}

	// use a table for testing
	table := []struct {
		key      string
		got      interface{}
		expected interface{}
	}{
		{
			key:      "name",
			got:      d.Get("name"),
			expected: "bob",
		},
		{
			key:      "age",
			got:      d.GetFloat("age"),
			expected: 25.0,
		},
		{
			key:      "cool",
			got:      d.GetBool("cool"),
			expected: true,
		},
		{
			key:      "aptitude",
			got:      d.Get("aptitude"),
			expected: "",
		},
	}
	for _, test := range table {
		if !reflect.DeepEqual(test.got, test.expected) {
			t.Errorf("%s was incorrect. Expected %v, but got %v.\n", test.key, test.expected, test.got)
		}
	}

	// Test unmarshaling the entire body to a data structure.
	expected := jsonData{
		Name:     "bob",
		Age:      25,
		Cool:     true,
		Aptitude: "",
		Location: map[string]float64{"latitude": 123.456, "longitude": 948.123},
		Things:   []string{"a", "b", "c"},
	}
	var got jsonData
	if err := d.BindJSON(&got); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(got, expected) {
		t.Errorf("Result of BindJSON was incorrect. Expected %+v, but got %+v.\n", expected, got)
	}

	// Test unmarshaling into data structures separately
	// For maps, both the GetMapFromJSON method and the GetAndUnmarshalJSON method
	expectedMap := map[string]interface{}{"latitude": 123.456, "longitude": 948.123}
	if got, err := d.GetMapFromJSON("location"); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(got, expectedMap) {
		t.Errorf("location was incorrect. Expected %v, but got %v.\n", expectedMap, got)
	}
	gotMap := map[string]interface{}{}
	if err := d.GetAndUnmarshalJSON("location", &gotMap); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(gotMap, expectedMap) {
		t.Errorf("location was incorrect. Expected %v, but got %v.\n", expectedMap, gotMap)
	}

	// For slices, both the GetSliceFromJSON method and the GetAndUnmarshalJSON method
	expectedSlice := []interface{}{"a", "b", "c"}
	if got, err := d.GetSliceFromJSON("things"); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(got, expectedSlice) {
		t.Errorf("things was incorrect. Expected %v, but got %v.\n", expectedSlice, got)
	}
	gotSlice := []interface{}{}
	if err := d.GetAndUnmarshalJSON("things", &gotSlice); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(gotSlice, expectedSlice) {
		t.Errorf("things was incorrect. Expected %v, but got %v.\n", expectedSlice, gotSlice)
	}
}

func ExampleParse() {
	// Construct a request object for example purposes only.
	// Typically you would be using this inside a http.HandlerFunc,
	// not constructing your own request.
	req, _ := http.NewRequest("GET", "/", nil)
	values := url.Values{}
	values.Add("name", "Bob")
	values.Add("age", "25")
	values.Add("retired", "false")
	req.PostForm = values
	req.Header.Set("Content-Type", "form-urlencoded")

	// Parse the form data.
	data, err := Parse(req)
	if err != nil {
		panic(err)
	}
	name := data.Get("name")
	age := data.GetInt("age")
	retired := data.GetBool("retired")
	if retired {
		fmt.Printf("%s is %d years old and he has retired.", name, age)
	} else {
		fmt.Printf("%s is %d years old and not yet retired.", name, age)
	}
	// Output:
	// Bob is 25 years old and not yet retired.
}
