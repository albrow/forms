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
	data   Data
	Keys   []string
	Errors []string
}

// Error adds an error associated with key to the validator. err
// should typically be a user-readable sentence, such as "username
// is required."
func (v *Validator) Error(key string, err string) {
	v.Keys = append(v.Keys, key)
	v.Errors = append(v.Errors, err)
}

// HasErrors returns true iff the Validator has errors, i.e.
// if any validation methods called on the Validator failed.
func (v *Validator) HasErrors() bool {
	return len(v.Errors) > 0
}

// Require will add an error to the Validator if data[key]
// does not exist, is an empty string, or consists of only
// whitespace.
func (v *Validator) Require(key string, msg ...string) {
	if strings.TrimSpace(v.data.Get(key)) == "" {
		v.requiredError(key, msg...)
	}
}

func (v *Validator) requiredError(key string, msg ...string) {
	if len(msg) != 0 {
		v.Error(key, msg[0])
	} else {
		err := fmt.Sprintf("%s is required.", key)
		v.Error(key, err)
	}
}

// MinLength will add an error to the Validator if data[key]
// is shorter than length (if data[key] has less than
// length characters), not counting leading or trailing
// whitespace.
func (v *Validator) MinLength(key string, length int, msg ...string) {
	val := v.data.Get(key)
	trimmed := strings.TrimSpace(val)
	if len(trimmed) < length {
		v.minLengthError(key, length, msg...)
	}
}

func (v *Validator) minLengthError(key string, length int, msg ...string) {
	if len(msg) != 0 {
		v.Error(key, msg[0])
	} else {
		err := fmt.Sprintf("%s must be at least %d characters long.", key, length)
		v.Error(key, err)
	}
}

// MaxLength will add an error to the Validator if data[key]
// is longer than length (if data[key] has more than
// length characters), not counting leading or trailing
// whitespace.
func (v *Validator) MaxLength(key string, length int, msg ...string) {
	val := v.data.Get(key)
	trimmed := strings.TrimSpace(val)
	if len(trimmed) > length {
		v.maxLengthError(key, length, msg...)
	}
}

func (v *Validator) maxLengthError(key string, length int, msg ...string) {
	if len(msg) != 0 {
		v.Error(key, msg[0])
	} else {
		err := fmt.Sprintf("%s cannot be more than %d characters long.", key, length)
		v.Error(key, err)
	}
}

// LengthRange will add an error to the Validator if data[key]
// is shorter than min (if data[key] has less than
// min characters) or if data[key] is longer than max
// (if data[key] has more than max characters), not
// counting leading or trailing whitespace.
func (v *Validator) LengthRange(key string, min int, max int, msg ...string) {
	if val := v.data.Get(key); len(val) < min || len(val) > max {
		v.lengthRangeError(key, min, max, msg...)
	}
}

func (v *Validator) lengthRangeError(key string, min int, max int, msg ...string) {
	if len(msg) != 0 {
		v.Error(key, msg[0])
	} else {
		err := fmt.Sprintf("%s must be between %d and %d characters long.", key, min, max)
		v.Error(key, err)
	}
}

// Equal will add an error to the Validator if data[key1]
// is not equal to data[key2].
func (v *Validator) Equal(key1 string, key2 string, msg ...string) {
	val1 := v.data.Get(key1)
	val2 := v.data.Get(key2)
	if val1 != val2 {
		v.equalError(key1, key2, msg...)
	}
}

func (v *Validator) equalError(key1 string, key2 string, msg ...string) {
	if len(msg) != 0 {
		v.Error(key2, msg[0])
	} else {
		// note: "match" is a more natural colloquial term than "be equal"
		// not to be confused with "matching" a regular expression
		err := fmt.Sprintf("%s and %s must match.", key1, key2)
		v.Error(key2, err)
	}
}

// Match will add an error to the Validator if data[key] does
// not match the regular expression regex.
func (v *Validator) Match(key string, regex *regexp.Regexp, msg ...string) {
	if !regex.MatchString(v.data.Get(key)) {
		v.matchError(key, msg...)
	}
}

// MatchEmail will add an error to the Validator if data[key]
// does not match the formatting expected of an email.
func (v *Validator) MatchEmail(key string, msg ...string) {
	regex := regexp.MustCompile("^[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?$")
	v.Match(key, regex, msg...)
}

func (v *Validator) matchError(key string, msg ...string) {
	if len(msg) != 0 {
		v.Error(key, msg[0])
	} else {
		err := fmt.Sprintf("%s must be correctly formatted.", key)
		v.Error(key, err)
	}
}

// TypeInt will add an error to the Validator if the first
// element of data[key] cannot be converted to an int.
func (v *Validator) TypeInt(key string, msg ...string) {
	if _, err := strconv.Atoi(v.data.Get(key)); err != nil {
		v.typeError(key, "integer", msg...)
	}
}

// TypeFloat will add an error to the Validator if the first
// element of data[key] cannot be converted to an float64.
func (v *Validator) TypeFloat(key string, msg ...string) {
	if _, err := strconv.ParseFloat(v.data.Get(key), 64); err != nil {
		// note: "number" is a more natural colloquial term than "float"
		v.typeError(key, "number", msg...)
	}
}

// TypeBool will add an error to the Validator if the first
// element of data[key] cannot be converted to a bool.
func (v *Validator) TypeBool(key string, msg ...string) {
	if _, err := strconv.ParseBool(v.data.Get(key)); err != nil {
		// note: "true or false" is a more natural colloquial term than "bool"
		v.typeError(key, "true or false", msg...)
	}
}

func (v *Validator) typeError(key string, typ string, msg ...string) {
	if len(msg) != 0 {
		v.Error(key, msg[0])
	} else {
		article := "a"
		if strings.Contains("aeiou", string(typ[0])) {
			article = "an"
		}
		err := fmt.Sprintf("%s must be %s %s", key, article, typ)
		v.Error(key, err)
	}
}
