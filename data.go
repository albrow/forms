package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Data holds data obtained from the request body and url query parameters.
// Because Data is built from multiple sources, sometimes there will be more
// than one value for a given key. You can use Get, Set, Add, and Del to access
// the first element for a given key or access the map directly to access additional
// elements for a given key. You can also use helper methods to convert the first
// value for a given key to a different type (e.g. bool or int).
type Data url.Values

// Parse parses the request body and url query parameters into
// Data. The content in the body of the request has a higher priority,
// will be added to Data first, and will be the result of any operation
// which gets the first element for a given key (e.g. Get, GetInt, or GetBool).
func Parse(req *http.Request) (Data, error) {
	values := url.Values{}
	contentType := req.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		if err := req.ParseMultipartForm(2048); err != nil {
			return nil, err
		}
		for key, vals := range req.MultipartForm.Value {
			for _, val := range vals {
				values.Add(key, val)
			}
		}
	} else if strings.Contains(contentType, "form-urlencoded") {
		if err := req.ParseForm(); err != nil {
			return nil, err
		}
		for key, vals := range req.PostForm {
			for _, val := range vals {
				values.Add(key, val)
			}
		}
	} else if strings.Contains(contentType, "application/json") {
		if err := parseJSON(values, req); err != nil {
			return nil, err
		}
	}
	for key, vals := range req.URL.Query() {
		for _, val := range vals {
			values.Add(key, val)
		}
	}
	return Data(values), nil
}

func parseJSON(values url.Values, req *http.Request) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	rawData := map[string]interface{}{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return err
	}
	// Whatever the underlying type is, we need to convert it to a
	// string. There are only a few possible types, so we can just
	// do a type switch over the possibilities.
	for key, val := range rawData {
		switch val.(type) {
		case string, bool, float64:
			values.Add(key, fmt.Sprint(val))
		case nil:
			values.Add(key, "")
		case map[string]interface{}, []interface{}:
			// for more complicated data structures, convert back to
			// a JSON string and let user decide how to unmarshal
			jsonVal, err := json.Marshal(val)
			if err != nil {
				return err
			}
			values.Add(key, string(jsonVal))
		}
	}
	return nil
}

// Add adds the value to key. It appends to any existing values associated with key.
func (d Data) Add(key string, value string) {
	url.Values(d).Add(key, value)
}

// Del deletes the values associated with key.
func (d Data) Del(key string) {
	url.Values(d).Del(key)
}

// Encode encodes the values into “URL encoded” form ("bar=baz&foo=quux") sorted by key.
func (d Data) Encode() string {
	return url.Values(d).Encode()
}

// Get gets the first value associated with the given key. If there are no values
// associated with the key, Get returns the empty string. To access multiple values,
// use the map directly.
func (d Data) Get(key string) string {
	return url.Values(d).Get(key)
}

// Set sets the key to value. It replaces any existing values.
func (d Data) Set(key string, value string) {
	url.Values(d).Set(key, value)
}

// KeyExists returns true iff data[key] exists. If the value of data[key] is an empty
// string, it is still considered to be in existence.
func (d Data) KeyExists(key string) bool {
	_, found := d[key]
	return found
}

// GetInt returns the first element in data[key] converted to an int.
func (d Data) GetInt(key string) int {
	if !d.KeyExists(key) || len(d[key]) == 0 {
		return 0
	}
	str := d.Get(key)
	if result, err := strconv.Atoi(str); err != nil {
		panic(err)
	} else {
		return result
	}
}

// GetFloat returns the first element in data[key] converted to a float.
func (d Data) GetFloat(key string) float64 {
	if !d.KeyExists(key) || len(d[key]) == 0 {
		return 0.0
	}
	str := d.Get(key)
	if result, err := strconv.ParseFloat(str, 64); err != nil {
		panic(err)
	} else {
		return result
	}
}

// GetBool returns the first element in data[key] converted to a bool.
func (d Data) GetBool(key string) bool {
	if !d.KeyExists(key) || len(d[key]) == 0 {
		return false
	}
	str := d.Get(key)
	if result, err := strconv.ParseBool(str); err != nil {
		panic(err)
	} else {
		return result
	}
}

// GetBytes returns the first element in data[key] converted to a slice of bytes.
func (d Data) GetBytes(key string) []byte {
	return []byte(d.Get(key))
}

// GetStringsSplit returns the first element in data[key] split into a slice delimited by delim.
func (d Data) GetStringsSplit(key string, delim string) []string {
	if !d.KeyExists(key) || len(d[key]) == 0 {
		return nil
	}
	return strings.Split(d[key][0], delim)
}

// GetMapFromJSON assumes that the first element in data[key] is a json string, attempts to
// unmarshal it into a map[string]interface{}, and if successful, returns the result. If
// unmarshaling was not successful, returns an error.
func (d Data) GetMapFromJSON(key string) (map[string]interface{}, error) {
	if !d.KeyExists(key) || len(d[key]) == 0 {
		return nil, nil
	}
	result := map[string]interface{}{}
	if err := json.Unmarshal([]byte(d.Get(key)), &result); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

// GetSliceFromJSON assumes that the first element in data[key] is a json string, attempts to
// unmarshal it into a []interface{}, and if successful, returns the result. If unmarshaling
// was not successful, returns an error.
func (d Data) GetSliceFromJSON(key string) ([]interface{}, error) {
	if !d.KeyExists(key) || len(d[key]) == 0 {
		return nil, nil
	}
	result := []interface{}{}
	if err := json.Unmarshal([]byte(d.Get(key)), &result); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

// GetAndMarshalJSON assumes that the first element in data[key] is a json string and
// attempts to unmarshal it into v. If unmarshaling was not successful, returns an error.
// v should be a pointer to some data structure.
func (d Data) GetAndMarshalJSON(key string, v interface{}) error {
	if err := json.Unmarshal([]byte(d.Get(key)), v); err != nil {
		return err
	}
	return nil
}

// Validator returns a Validator which can be used to easily validate data.
func (d Data) Validator() *Validator {
	return &Validator{
		data: d,
	}
}
