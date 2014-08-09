package data

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	data := Data(map[string][]string{
		"name":       []string{"bob", "bill"},
		"profession": []string{"plumber"},
	})

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
	data := Data(map[string][]string{
		"age":    []string{"25", "33"},
		"weight": []string{"155"},
	})

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
	data := Data(map[string][]string{
		"age":    []string{"25.7", "33"},
		"weight": []string{"42"},
	})

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
	data := Data(map[string][]string{
		"retired":         []string{"true", "false"},
		"leftHanded":      []string{"0"},
		"collegeGraduate": []string{"1"},
	})

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
	data := Data(map[string][]string{
		"name":       []string{"bob", "bill"},
		"profession": []string{"plumber"},
	})

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

func TestGetStringsSplit(t *testing.T) {
	data := Data(map[string][]string{
		"children":       []string{"martha,bill,jane", "adam,julia"},
		"favoriteColors": []string{"blue%20green%20fuchsia"},
	})

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

	// Test unmarshaling into data structures separately
	// For maps, both the GetMapFromJSON method and the GetAndMarshalJSON method
	expectedMap := map[string]interface{}{"latitude": 123.456, "longitude": 948.123}
	if got, err := d.GetMapFromJSON("location"); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(got, expectedMap) {
		t.Errorf("location was incorrect. Expected %v, but got %v.\n", expectedMap, got)
	}
	gotMap := map[string]interface{}{}
	if err := d.GetAndMarshalJSON("location", &gotMap); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(gotMap, expectedMap) {
		t.Errorf("location was incorrect. Expected %v, but got %v.\n", expectedMap, gotMap)
	}

	// For slices, both the GetSliceFromJSON method and the GetAndMarshalJSON method
	expectedSlice := []interface{}{"a", "b", "c"}
	if got, err := d.GetSliceFromJSON("things"); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(got, expectedSlice) {
		t.Errorf("things was incorrect. Expected %v, but got %v.\n", expectedSlice, got)
	}
	gotSlice := []interface{}{}
	if err := d.GetAndMarshalJSON("things", &gotSlice); err != nil {
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
