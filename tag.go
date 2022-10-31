package analytics

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"strings"
)

// hashValues looks for fields with an 'analytics' struct tag and hashes
// the data using SHA256.
//
// The expected format of this struct tag is as follows:
//
//	type myEvent struct {
//		UserID string `analytics:"usr"`
//	}
//	// hashValues(&myEvent{UserID: "something"}) will replace the UserID field with usr_kh6XMrBb8gxpzREC4mAOSNs862lIy8tjE9fNDBWrjRE
func hashValues[T any](e T) T {
	v := reflect.ValueOf(e)

	if v.Kind() == reflect.Pointer {
		v = reflect.Indirect(v)
	}

	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		tag := typeField.Tag
		at, ok := tag.Lookup("analytics")
		if !ok {
			continue
		}

		prefix, _, _ := strings.Cut(at, ",")
		if prefix == "" {
			continue
		}

		fieldVal := v.Field(i)
		fieldInterface := fieldVal.Interface()
		fieldBytes, err := json.Marshal(fieldInterface)
		if err != nil {
			continue
		}
		if string(fieldBytes) == `""` || string(fieldBytes) == "" || fieldBytes == nil {
			continue
		}

		hashed := sha256.Sum256(fieldBytes)
		hashstr := base64.RawURLEncoding.EncodeToString(hashed[:])

		val := prefix + "_" + hashstr

		if fieldVal.Kind() != reflect.String {
			continue
		}
		if !fieldVal.CanSet() {
			fieldVal.Interface()
			// fieldVal = fieldVal.Elem()
			fieldVal = reflect.ValueOf(&e).Elem().Elem().Field(i)
		}

		fieldVal.SetString(val)
	}

	return e
}
