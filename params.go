package sparkle

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
)

var ErrNilUnmarshalTarget = errors.New("Passed nil object as unmarshal target")
var ErrNotPointerTarget = errors.New("Target is not a pointer")
var ErrNotStructTarget = errors.New("Target pointer does not point to struct")
var ErrCouldNotGetContextRequest = errors.New("Could not obtain request from context")
var ErrUnsupportedType = errors.New("Can not unmarshall to field type as it is unsupported")

// A ValueProvider that provides values from the current request query string and
// form inputs
var currentValueProvider ParameterValueProvider
var maxMemoryForMultipartForm int64

func init() {
	currentValueProvider = &FormValueProvider{}
	// Default to 32 MB
	maxMemoryForMultipartForm = 32 * 1024 * 1024
}

// A map type that has both string keys and values
type KeyValueStringMap map[string]string

// ParameterValueProvider is an interface which when implemented provides
// the mechinism neccisary for (*Context).UnmarshalParameters to get data
// for use in unmarshalling parameters.
type ParameterValueProvider interface {
	Values(c *Context) (KeyValueStringMap, error)
}

// SetParameterValueProvider sets the ParameterValueProvider that sparkle 
// will use as data for unmarshalling paramters
func SetParameterValueProvider(p ParameterValueProvider) {
	currentValueProvider = p
}

// FormValueProvider is a ParameterValueProvider that obtains data from
// the request query string and form paramaeters
type FormValueProvider struct{}

// Values provides a KeyValueStringMap of the Query String and Form parameters
func (p *FormValueProvider) Values(c *Context) (KeyValueStringMap, error) {
	request := c.Request()
	if request == nil {
		// Should never happen
		return nil, ErrCouldNotGetContextRequest
	}

	if err := parseFormData(request); err != nil {
		return nil, err
	}

	return request.Form, nil
}

// Sets the Max Memory to be used when UnmarshalParameters encounters a
// multipart form.
func SetMaxMemory(maxMemory int64) {
	maxMemoryForMultipartForm = maxMemory
}

func canUnmarshal(v interface{}) error {
	if v == nil {
		return ErrNilUnmarshalTarget
	}

	it := reflect.ValueOf(v)
	if it.Kind() != reflect.Ptr {
		return ErrNotPointerTarget
	}

	if rt := it.Elem(); rt.Kind() != reflect.Struct {
		return ErrNotStructTarget
	}

	return nil
}

func parseFormData(r *http.Request) error {
	if r.Header.Get("Content-Type") == "multipart/form-data" {
		return r.ParseMultipartForm(maxMemoryForMultipartForm)
	}

	return r.ParseForm()
}

func getFieldAndKind(rt reflect.Value, fieldName string) (reflect.Value, reflect.Kind) {
	// Look for property on structure with the same name
	fieldValue := rt.FieldByName(fieldName)
	fieldKind := fieldValue.Kind()
	if fieldKind == reflect.Ptr {
		// Step one in if it's a pointer
		fieldValue = fieldValue.Elem()
		fieldKind = fieldValue.Kind()
	}

	return fieldValue, fieldKind
}

// Unmarshals Query and Post parameters to an object supplied in v
// 
// If a parameter doesn't have a corresponding field in the struct, it is 
// ignored. If the struct has a field that doesn't have a corresponding 
// parameter, then the struct field will not be altered.
//
// If successful returns nil with the values in v set accordingly
// Returns ErrNilUnmarshalTarget if v was nil
// Returns ErrNotPointerTarget if v was not a pointer
// Returns ErrNotStructTarget if v was not a pointer to a struct
func (c *Context) UnmarshalParameters(v interface{}) error {

	if err := canUnmarshal(v); err != nil {
		return err
	}

	rt := reflect.ValueOf(v).Elem()

	valueMap, err := currentValueProvider(c)
	if err != nil {
		return err
	}

	// Okay, so we should step through the values in the form now
	for key, value := range valueMap {
		fieldValue, fieldKind := getFieldAndKind(rt, key)

		// if we can't set the result, then ignore it
		if !fieldValue.CanSet() {
			continue
		}

		switch fieldKind {
		case reflect.String:
			fieldValue.SetString(value[0])

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if v, err := strconv.ParseUint(value[0], 10, 64); err != nil {
				return err
			} else {
				fieldValue.SetUint(v)
			}
		case reflect.Int:
			if v, err := strconv.ParseInt(value[0], 10, 64); err != nil {
				return err
			} else {
				fieldValue.SetInt(v)
			}
		case reflect.Bool:
			if v, err := strconv.ParseBool(value[0]); err != nil {
				return err
			} else {
				fieldValue.SetBool(v)
			}
		default:
			return ErrUnsupportedType
		}
	}

	return nil
}
