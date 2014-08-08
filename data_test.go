package data

import (
	"fmt"
	"net/http"
	"net/url"
)

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
