package data

import (
	"fmt"
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
	data    Data
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

// Require will add an error to the Validator if data[field]
// does not exist, is an empty string, or consists of only
// whitespace.
func (v *Validator) Require(field string) *ValidationResult {
	if strings.TrimSpace(v.data.Get(field)) == "" {
		return v.addRequiredError(field)
	} else {
		return &ValidationResult{Ok: true}
	}
}

func (v *Validator) addRequiredError(field string) *ValidationResult {
	msg := fmt.Sprintf("%s is required.", field)
	return v.AddError(field, msg)
}

// MinLength will add an error to the Validator if data[field]
// is shorter than length (if data[field] has less than
// length characters), not counting leading or trailing
// whitespace.
func (v *Validator) MinLength(field string, length int) *ValidationResult {
	val := v.data.Get(field)
	trimmed := strings.TrimSpace(val)
	if len(trimmed) < length {
		return v.addMinLengthError(field, length)
	} else {
		return &ValidationResult{Ok: true}
	}
}

func (v *Validator) addMinLengthError(field string, length int) *ValidationResult {
	msg := fmt.Sprintf("%s must be at least %d characters long.", field, length)
	return v.AddError(field, msg)
}

// MaxLength will add an error to the Validator if data[field]
// is longer than length (if data[field] has more than
// length characters), not counting leading or trailing
// whitespace.
func (v *Validator) MaxLength(field string, length int) *ValidationResult {
	val := v.data.Get(field)
	trimmed := strings.TrimSpace(val)
	if len(trimmed) > length {
		return v.addMaxLengthError(field, length)
	} else {
		return &ValidationResult{Ok: true}
	}
}

func (v *Validator) addMaxLengthError(field string, length int) *ValidationResult {
	msg := fmt.Sprintf("%s cannot be more than %d characters long.", field, length)
	return v.AddError(field, msg)
}

// LengthRange will add an error to the Validator if data[field]
// is shorter than min (if data[field] has less than
// min characters) or if data[field] is longer than max
// (if data[field] has more than max characters), not
// counting leading or trailing whitespace.
func (v *Validator) LengthRange(field string, min int, max int) *ValidationResult {
	if val := v.data.Get(field); len(val) < min || len(val) > max {
		return v.addLengthRangeError(field, min, max)
	} else {
		return &ValidationResult{Ok: true}
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
		return &ValidationResult{Ok: true}
	}
}

func (v *Validator) addEqualError(field1 string, field2 string) *ValidationResult {
	// note: "match" is a more natural colloquial term than "be equal"
	// not to be confused with "matching" a regular expression
	msg := fmt.Sprintf("%s and %s must match.", field1, field2)
	return v.AddError(field2, msg)
}

// Match will add an error to the Validator if data[field] does
// not match the regular expression regex.
func (v *Validator) Match(field string, regex *regexp.Regexp) *ValidationResult {
	if !regex.MatchString(v.data.Get(field)) {
		return v.addMatchError(field)
	} else {
		return &ValidationResult{Ok: true}
	}
}

// MatchEmail will add an error to the Validator if data[field]
// does not match the formatting expected of an email.
func (v *Validator) MatchEmail(field string, msg ...string) *ValidationResult {
	regex := regexp.MustCompile("^[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?$")
	return v.Match(field, regex)
}

func (v *Validator) addMatchError(field string) *ValidationResult {
	msg := fmt.Sprintf("%s must be correctly formatted.", field)
	return v.AddError(field, msg)
}

// TypeInt will add an error to the Validator if the first
// element of data[field] cannot be converted to an int.
func (v *Validator) TypeInt(field string) *ValidationResult {
	if _, err := strconv.Atoi(v.data.Get(field)); err != nil {
		return v.addTypeError(field, "integer")
	} else {
		return &ValidationResult{Ok: true}
	}
}

// TypeFloat will add an error to the Validator if the first
// element of data[field] cannot be converted to an float64.
func (v *Validator) TypeFloat(field string) *ValidationResult {
	if _, err := strconv.ParseFloat(v.data.Get(field), 64); err != nil {
		// note: "number" is a more natural colloquial term than "float"
		return v.addTypeError(field, "number")
	} else {
		return &ValidationResult{Ok: true}
	}
}

// TypeBool will add an error to the Validator if the first
// element of data[field] cannot be converted to a bool.
func (v *Validator) TypeBool(field string) *ValidationResult {
	if _, err := strconv.ParseBool(v.data.Get(field)); err != nil {
		// note: "true or false" is a more natural colloquial term than "bool"
		return v.addTypeError(field, "true or false")
	} else {
		return &ValidationResult{Ok: true}
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
