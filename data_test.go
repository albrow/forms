package data

import (
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
