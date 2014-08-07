package data

import (
	"fmt"
	"regexp"
	"strings"
)

type Validator struct {
	data   Data
	Keys   []string
	Errors []string
}

func (v *Validator) Error(key string, err string) {
	v.Keys = append(v.Keys, key)
	v.Errors = append(v.Errors, err)
}

func (v *Validator) HasErrors() bool {
	return len(v.Errors) > 0
}

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

func (v *Validator) MinLength(key string, length int, msg ...string) {
	if len(v.data.Get(key)) < length {
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

func (v *Validator) MaxLength(key string, length int, msg ...string) {
	if len(v.data.Get(key)) > length {
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

func (v *Validator) Match(key string, regex *regexp.Regexp, msg ...string) {
	if !regex.MatchString(v.data.Get(key)) {
		v.matchError(key, msg...)
	}
}

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
