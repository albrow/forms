// Copyright 2015 Alex Browne.  All rights reserved.
// Use of this source code is governed by the MIT
// license, which can be found in the LICENSE file.

package forms

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Validator has methods for validating its underlying Data.
// A Validator stores any errors that occurred during validation,
// and they can be accessed directly. In a typical workflow, you
// will create a Validator from some Data, call some methods on
// that validator (e.g. Require), check if the validator
// has errors, then do something with the errors if it does.
type Validator struct {
	data    *Data
	results []*ValidationResult
}

// ValidationResult is returned from every validation method and can
// be used to override the default field name or error message. If
// you want to use the default fields and messages, simply discard
// the ValidationResult.
type ValidationResult struct {
	Ok      bool
	field   string
	message string
}

var validationOk = &ValidationResult{Ok: true}

// Field changes the field name associated with the validation result.
func (vr *ValidationResult) Field(field string) *ValidationResult {
	vr.field = field
	return vr
}

// Message changes the error message associated with the validation
// result. msg should typically be a user-readable sentence, such as
// "username is required."
func (vr *ValidationResult) Message(msg string) *ValidationResult {
	vr.message = msg
	return vr
}

// AddError adds an error associated with field to the validator. msg
// should typically be a user-readable sentence, such as "username
// is required."
func (v *Validator) AddError(field string, msg string) *ValidationResult {
	result := &ValidationResult{
		field:   field,
		message: msg,
	}
	v.results = append(v.results, result)
	return result
}

// HasErrors returns true iff the Validator has errors, i.e.
// if any validation methods called on the Validator failed.
func (v *Validator) HasErrors() bool {
	return len(v.results) > 0
}

// Messages returns the messages for all validation results for
// the Validator, in order.
func (v *Validator) Messages() []string {
	msgs := []string{}
	for _, vr := range v.results {
		msgs = append(msgs, vr.message)
	}
	return msgs
}

// Fields returns the fields for all validation results for
// the Validator, in order.
func (v *Validator) Fields() []string {
	fields := []string{}
	for _, vr := range v.results {
		fields = append(fields, vr.field)
	}
	return fields
}

// ErrorMap reutrns all the fields and error messages for
// the validator in the form of a map. The keys of the map
// are field names, and the values are any error messages
// associated with that field name.
func (v *Validator) ErrorMap() map[string][]string {
	errMap := map[string][]string{}
	for _, vr := range v.results {
		if _, found := errMap[vr.field]; found {
			errMap[vr.field] = append(errMap[vr.field], vr.message)
		} else {
			errMap[vr.field] = []string{vr.message}
		}
	}
	return errMap
}

// Require will add an error to the Validator if data.Values[field]
// does not exist, is an empty string, or consists of only
// whitespace.
func (v *Validator) Require(field string) *ValidationResult {
	if strings.TrimSpace(v.data.Get(field)) == "" {
		return v.addRequiredError(field)
	} else {
		return validationOk
	}
}

// RequireFile will add an error to the Validator if data.Files[field]
// does not exist or is an empty file
func (v *Validator) RequireFile(field string) *ValidationResult {
	if !v.data.FileExists(field) {
		return v.addRequiredError(field)
	}
	bytes, err := v.data.GetFileBytes(field)
	if err != nil {
		return v.AddError(field, "Could not read file.")
	}
	if len(bytes) == 0 {
		return v.addFileEmptyError(field)
	}
	return validationOk
}

func (v *Validator) addRequiredError(field string) *ValidationResult {
	msg := fmt.Sprintf("%s is required.", field)
	return v.AddError(field, msg)
}

func (v *Validator) addFileEmptyError(field string) *ValidationResult {
	msg := fmt.Sprintf("%s is required and cannot be an empty file.", field)
	return v.AddError(field, msg)
}

// MinLength will add an error to the Validator if data.Values[field]
// is shorter than length (if data.Values[field] has less than
// length characters), not counting leading or trailing
// whitespace.
func (v *Validator) MinLength(field string, length int) *ValidationResult {
	val := v.data.Get(field)
	trimmed := strings.TrimSpace(val)
	if len(trimmed) < length {
		return v.addMinLengthError(field, length)
	} else {
		return validationOk
	}
}

func (v *Validator) addMinLengthError(field string, length int) *ValidationResult {
	msg := fmt.Sprintf("%s must be at least %d characters long.", field, length)
	return v.AddError(field, msg)
}

// MaxLength will add an error to the Validator if data.Values[field]
// is longer than length (if data.Values[field] has more than
// length characters), not counting leading or trailing
// whitespace.
func (v *Validator) MaxLength(field string, length int) *ValidationResult {
	val := v.data.Get(field)
	trimmed := strings.TrimSpace(val)
	if len(trimmed) > length {
		return v.addMaxLengthError(field, length)
	} else {
		return validationOk
	}
}

func (v *Validator) addMaxLengthError(field string, length int) *ValidationResult {
	msg := fmt.Sprintf("%s cannot be more than %d characters long.", field, length)
	return v.AddError(field, msg)
}

// LengthRange will add an error to the Validator if data.Values[field]
// is shorter than min (if data.Values[field] has less than
// min characters) or if data.Values[field] is longer than max
// (if data.Values[field] has more than max characters), not
// counting leading or trailing whitespace.
func (v *Validator) LengthRange(field string, min int, max int) *ValidationResult {
	if val := v.data.Get(field); len(val) < min || len(val) > max {
		return v.addLengthRangeError(field, min, max)
	} else {
		return validationOk
	}
}

func (v *Validator) addLengthRangeError(field string, min int, max int) *ValidationResult {
	msg := fmt.Sprintf("%s must be between %d and %d characters long.", field, min, max)
	return v.AddError(field, msg)
}

// Equal will add an error to the Validator if data[field1]
// is not equal to data[field2].
func (v *Validator) Equal(field1 string, field2 string) *ValidationResult {
	val1 := v.data.Get(field1)
	val2 := v.data.Get(field2)
	if val1 != val2 {
		return v.addEqualError(field1, field2)
	} else {
		return validationOk
	}
}

func (v *Validator) addEqualError(field1 string, field2 string) *ValidationResult {
	// note: "match" is a more natural colloquial term than "be equal"
	// not to be confused with "matching" a regular expression
	msg := fmt.Sprintf("%s and %s must match.", field1, field2)
	return v.AddError(field2, msg)
}

// Match will add an error to the Validator if data.Values[field] does
// not match the regular expression regex.
func (v *Validator) Match(field string, regex *regexp.Regexp) *ValidationResult {
	if !regex.MatchString(v.data.Get(field)) {
		return v.addMatchError(field)
	} else {
		return validationOk
	}
}

// MatchEmail will add an error to the Validator if data.Values[field]
// does not match the formatting expected of an email.
func (v *Validator) MatchEmail(field string) *ValidationResult {
	regex := regexp.MustCompile("^[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?$")
	return v.Match(field, regex)
}

func (v *Validator) addMatchError(field string) *ValidationResult {
	msg := fmt.Sprintf("%s must be correctly formatted.", field)
	return v.AddError(field, msg)
}

// TypeInt will add an error to the Validator if the first
// element of data.Values[field] cannot be converted to an int.
func (v *Validator) TypeInt(field string) *ValidationResult {
	if _, err := strconv.Atoi(v.data.Get(field)); err != nil {
		return v.addTypeError(field, "integer")
	} else {
		return validationOk
	}
}

// TypeFloat will add an error to the Validator if the first
// element of data.Values[field] cannot be converted to an float64.
func (v *Validator) TypeFloat(field string) *ValidationResult {
	if _, err := strconv.ParseFloat(v.data.Get(field), 64); err != nil {
		// note: "number" is a more natural colloquial term than "float"
		return v.addTypeError(field, "number")
	} else {
		return validationOk
	}
}

// TypeBool will add an error to the Validator if the first
// element of data.Values[field] cannot be converted to a bool.
func (v *Validator) TypeBool(field string) *ValidationResult {
	if _, err := strconv.ParseBool(v.data.Get(field)); err != nil {
		// note: "true or false" is a more natural colloquial term than "bool"
		return v.addTypeError(field, "true or false")
	} else {
		return validationOk
	}
}

func (v *Validator) addTypeError(field string, typ string) *ValidationResult {
	article := "a"
	if strings.Contains("aeiou", string(typ[0])) {
		article = "an"
	}
	msg := fmt.Sprintf("%s must be %s %s", field, article, typ)
	return v.AddError(field, msg)
}

// Greater will add an error to the Validator if the first
// element of data.Values[field] is not greater than value or if the first
// element of data.Values[field] cannot be converted to a number.
func (v *Validator) Greater(field string, value float64) *ValidationResult {
	return v.inequality(field, value, greater, "greater than")
}

// GreaterOrEqual will add an error to the Validator if the first
// element of data.Values[field] is not greater than or equal to value or if
// the first element of data.Values[field] cannot be converted to a number.
func (v *Validator) GreaterOrEqual(field string, value float64) *ValidationResult {
	return v.inequality(field, value, greaterOrEqual, "greater than or equal to")
}

// Less will add an error to the Validator if the first
// element of data.Values[field] is not less than value or if the first
// element of data.Values[field] cannot be converted to a number.
func (v *Validator) Less(field string, value float64) *ValidationResult {
	return v.inequality(field, value, less, "less than")
}

// LessOrEqual will add an error to the Validator if the first
// element of data.Values[field] is not less than or equal to value or if
// the first element of data.Values[field] cannot be converted to a number.
func (v *Validator) LessOrEqual(field string, value float64) *ValidationResult {
	return v.inequality(field, value, lessOrEqual, "less than or equal to")
}

type conditional func(given float64, target float64) bool

var greater conditional = func(given float64, target float64) bool {
	return given > target
}

var greaterOrEqual conditional = func(given float64, target float64) bool {
	return given >= target
}

var less conditional = func(given float64, target float64) bool {
	return given < target
}

var lessOrEqual conditional = func(given float64, target float64) bool {
	return given <= target
}

func (v *Validator) inequality(field string, value float64, condition conditional, explanation string) *ValidationResult {
	if valFloat, err := strconv.ParseFloat(v.data.Get(field), 64); err != nil {
		// note: "number" is a more natural colloquial term than "float"
		return v.addTypeError(field, "number")
	} else {
		if !condition(valFloat, value) {
			return v.AddError(field, fmt.Sprintf("%s must be %s %f.", field, explanation, value))
		} else {
			return validationOk
		}
	}
}

// AcceptFileExts will add an error to the Validator if the extension
// of the file identified by field is not in exts. exts should be one ore more
// allowed file extensions, not including the preceding ".". If the file does not
// exist, it does not add an error to the Validator.
func (v *Validator) AcceptFileExts(field string, exts ...string) *ValidationResult {
	if !v.data.FileExists(field) {
		return validationOk
	}
	header := v.data.GetFile(field)
	gotExt := filepath.Ext(header.Filename)
	for _, ext := range exts {
		if ext == gotExt[1:] {
			return validationOk
		}
	}
	return v.addFileExtError(field, gotExt, exts...)
}

func (v *Validator) addFileExtError(field string, gotExt string, allowedExts ...string) *ValidationResult {
	msg := fmt.Sprintf("The file extension %s is not allowed. Allowed extensions include: ", gotExt)

	// Append each allowed extension to the message, in a human-readable list
	// e.g. "x, y, and z"
	for i, ext := range allowedExts {
		if i == len(allowedExts)-1 {
			// special case for the last element
			switch len(allowedExts) {
			case 1:
				msg += ext
			default:
				msg += fmt.Sprintf("and %s", ext)
			}
		} else {
			// default case for middle elements
			// we only reach here if there is at least
			// one element
			switch len(allowedExts) {
			case 2:
				msg += fmt.Sprintf("%s ", ext)
			default:
				msg += fmt.Sprintf("%s, ", ext)
			}
		}
	}
	return v.AddError(field, msg)
}
