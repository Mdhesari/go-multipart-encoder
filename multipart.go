package multipart

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// DefaultFileExtension is used when a file's extension cannot be detected
const DefaultFileExtension = ""

// Encode converts a struct into multipart form-data format.
// It returns a buffer containing the encoded data, the content type string,
// and any error that occurred during encoding.
//
// Struct fields are encoded based on their types:
// - String fields are encoded as form fields
// - Numeric fields (int, uint, float) are converted to strings and encoded as form fields
// - Boolean fields are converted to strings and encoded as form fields
// - []byte fields are encoded as files with automatic extension detection
// - Other slices are encoded as multiple form fields with the same name
// - Struct fields are JSON-encoded and sent as form fields
//
// Tags can be used to customize encoding:
// - `form:"fieldname"` sets the form field name (defaults to lowercase field name)
// - `form:"-"` skips the field
// - `filename:"custom.ext"` sets custom filename for []byte fields
func Encode(req any) (*bytes.Buffer, string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	v := reflect.ValueOf(req)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, "", fmt.Errorf("req must be a struct")
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		fieldName := fieldType.Tag.Get("form")
		if fieldName == "" {
			fieldName = strings.ToLower(fieldType.Name)
		}
		if fieldName == "-" {
			continue
		}

		var (
			fw  io.Writer
			err error
		)
		switch field.Kind() {
		case reflect.String:
			if field.String() != "" {
				fw, err = w.CreateFormField(fieldName)
				if err == nil {
					_, err = fw.Write([]byte(field.String()))
				}
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fw, err = w.CreateFormField(fieldName)
			if err == nil {
				_, err = fw.Write([]byte(strconv.FormatInt(field.Int(), 10)))
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fw, err = w.CreateFormField(fieldName)
			if err == nil {
				_, err = fw.Write([]byte(strconv.FormatUint(field.Uint(), 10)))
			}

		case reflect.Float32, reflect.Float64:
			fw, err = w.CreateFormField(fieldName)
			if err == nil {
				_, err = fw.Write([]byte(strconv.FormatFloat(field.Float(), 'f', -1, 64)))
			}

		case reflect.Bool:
			fw, err = w.CreateFormField(fieldName)
			if err == nil {
				_, err = fw.Write([]byte(strconv.FormatBool(field.Bool())))
			}

		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.Uint8 && !field.IsNil() {
				filename := fieldType.Tag.Get("filename")
				if filename == "" {
					ext := getExtensionFromContent(field.Bytes())
					if ext == "" {
						ext = DefaultFileExtension
					}
					filename = fieldName + ext
				}

				fw, err = w.CreateFormFile(fieldName, filename)
				if err == nil {
					_, err = fw.Write(field.Bytes())
				}
			} else if !field.IsNil() {
				// Handle slice of primitive values (as multiple form fields with the same name)
				for j := 0; j < field.Len(); j++ {
					elem := field.Index(j)
					fw, err = w.CreateFormField(fieldName)
					if err == nil {
						switch elem.Kind() {
						case reflect.String:
							_, err = fw.Write([]byte(elem.String()))
						case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
							_, err = fw.Write([]byte(strconv.FormatInt(elem.Int(), 10)))
						case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
							_, err = fw.Write([]byte(strconv.FormatUint(elem.Uint(), 10)))
						}
					}
				}
			}

		case reflect.Struct:
			var jsonData []byte
			jsonData, err = json.Marshal(field.Interface())
			if err == nil {
				fw, err = w.CreateFormField(fieldName)
				if err == nil {
					_, err = fw.Write(jsonData)
				}
			}
		}

		if err != nil {
			return nil, "", err
		}
	}

	if err := w.Close(); err != nil {
		return nil, "", fmt.Errorf("failed to close multipart writer %w", err)
	}

	return &b, w.FormDataContentType(), nil
}

func getExtensionFromContent(data []byte) string {
	t := http.DetectContentType(data)
	s, err := mime.ExtensionsByType(t)
	if err != nil || len(s) == 0 {

		return ""
	}

	if s[0] == ".jpe" {

		return ".jpg"
	}

	return s[0]
}
