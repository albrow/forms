Forms
=====

[![GoDoc](https://godoc.org/github.com/albrow/forms?status.svg)](https://godoc.org/github.com/albrow/forms)

Forms is a lightweight, but incredibly useful go library for parsing
form data from an http.Request. It supports multipart forms, url-encoded
forms, json data, and url query parameters. It also provides helper methods
for converting data into other types and a Validator object which can be
used to validate the data. Forms is framework-agnostic and works directly
with the http package.

Version 0.3.2


Development Status
------------------

Forms is being actively developed and is well-tested. However, since it is still
a young library, it is not recommended for use in mission-critical production
applications at this time. It is probably fine to use for low-traffic hobby
sites, and in fact we encourage its use in those settings to help polish the API
and find missing features and hidden bugs. Pull requests and issue reports are
much appreciated :)

Forms follows semantic versioning but offers no guarantees of backwards
compatibility until version 1.0. Keep in mind that breaking changes might occur.
We will do our best to make the community aware of any non-trivial breaking
changes beforehand. We recommend using a dependency vendoring tool such as
[godep](https://github.com/tools/godep) to ensure that breaking changes will not
break your application.

Installation
------------

Install like you would any other package:
```
go get github.com/albrow/forms
```

Then include the package in your import statements:
``` go
import "github.com/albrow/forms"
```

Example Usage
-------------

Meant to be used inside the body of an http.HandlerFunc or any function that
has access to an http.Request.

``` go
func CreateUserHandler(res http.ResponseWriter, req *http.Request) {
	// Parse request data.
	userData, err := forms.Parse(req)
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

Contributing
------------

See [CONTRIBUTING.md](https://github.com/albrow/forms/blob/master/CONTRIBUTING.md)

License
-------

Forms is licensed under the MIT License. See the LICENSE file for more information.
