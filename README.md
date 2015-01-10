Go Data Parser
==============

This library parses data from the body and url query paramaters
of an http.Request and combines it into a single Data object. It
supports multipart forms, url-encoded forms, and json data in the request.
It also provides helper methods for converting data into other types
and a Validator object which can be used to validate the data. It is
written to be framework-agnostic and works directly with the http
package.

[Full Documentation](http://godoc.org/github.com/albrow/go-data-parser)

Installation
------------

Install like you would any other package:
```
go get github.com/albrow/go-data-parser
```

Then include the package in your import statements:
``` go
import "github.com/albrow/go-data-parser"
```

Example Usage
-------------

Meant to be used inside the body of an http.HandlerFunc or any function that
has access to an http.Request.

``` go
func CreateUserHandler(res http.ResponseWriter, req *http.Request) {
	// Parse request data.
	userData, err := data.Parse(req)
	if err != nil {
		// Handle err
		// ...
	}

	// Validate
	val := userData.Validator()
	val.Require("username")
	val.LengthRange("username", 4, 16)
	val.Require("email")
	val.MatchEmail("email")
	val.Require("password")
	val.MinLength("password", 8)
	val.Require("confirmPassword")
	val.Equal("password", "confirmPassword")
	val.RequireFile("profileImage")
	val.AcceptFileExts("profileImage", "jpg", "png", "gif")
	if val.HasErrors() {
		// Write the errors to the response
		// Maybe this means formatting the errors as json
		// or re-rendering the form with an error message
		// ...
	}

	// Use data to create a user object
	user := &models.User {
		Username: userData.Get("username"),
		Email: userData.Get("email"),
		HashedPassword: hash(userData.Get("password")),
	}

	// Continue by saving the user to the database and writing
	// to the response
	// ...


	// Get the contents of the profileImage file
	imageBytes, err := userData.GetFileBytes("profileImage")
	if err != nil {
	  // Handle err
	}
	// Now you can either copy the file over to your server using io.Copy,
	// upload the file to something like amazon S3, or do whatever you want
	// with it.
}
```