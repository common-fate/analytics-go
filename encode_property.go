package analytics

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/common-fate/analytics-go/acore"
)

func eventToProperties(e Event) acore.Properties {
	v := reflect.ValueOf(e)
	props := acore.NewProperties()

	if v.Kind() == reflect.Pointer {
		v = reflect.Indirect(v)
	}

	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i)

		typeField := v.Type().Field(i)
		tag := typeField.Tag
		jt, ok := tag.Lookup("json")
		if !ok {
			continue
		}
		prefix, opts := parseTag(jt)
		if prefix == "" || prefix == "-" {
			continue
		}
		if opts.Contains("omitempty") && isEmptyValue(fieldVal) {
			continue
		}

		if val, ok := hashStructField(typeField, fieldVal); ok {
			props.Set(prefix, val)
			continue
		}

		fieldInterface := fieldVal.Interface()
		props.Set(prefix, fieldInterface)
	}
	return props
}

func hashStructField(typeField reflect.StructField, fieldVal reflect.Value) (val string, ok bool) {
	tag := typeField.Tag
	at, ok := tag.Lookup("analytics")
	if !ok {
		return
	}

	prefix, _, _ := strings.Cut(at, ",")
	if prefix == "" {
		return
	}

	fieldInterface := fieldVal.Interface()
	return hash(fieldInterface, prefix)
}

func hash(in any, prefix string) (val string, ok bool) {
	fieldBytes, err := json.Marshal(in)
	if err != nil {
		return
	}
	if string(fieldBytes) == `""` || string(fieldBytes) == "" || fieldBytes == nil {
		return
	}

	hashed := sha256.Sum256(fieldBytes)
	hashstr := base64.RawURLEncoding.EncodeToString(hashed[:])

	val = prefix + "_" + hashstr
	ok = true
	return
}

// tagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	tag, opt, _ := strings.Cut(tag, ",")
	return tag, tagOptions(opt)
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var name string
		name, s, _ = strings.Cut(s, ",")
		if name == optionName {
			return true
		}
	}
	return false
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Pointer:
		return v.IsNil()
	}
	return false
}
