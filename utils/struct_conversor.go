package utils

import "reflect"

const (
	// tagName é uma tag usada nos structs
	tagName = "apelido"
)

// ConvertStructByAlias é uma função que converte uma interface em outra
func ConvertStructByAlias(src interface{}, dst interface{}) error {
	srcValue := reflect.ValueOf(src).Elem()
	dstValue := reflect.ValueOf(dst).Elem()

	srcType := srcValue.Type()
	dstType := dstValue.Type()

	for i := 0; i < srcType.NumField(); i++ {
		srcField := srcType.Field(i)
		dstField, exists := dstType.FieldByName(srcField.Name)
		if !exists {
			continue
		}

		alias := srcField.Tag.Get(tagName)
		if alias != "" {
			srcFieldValue := srcValue.FieldByName(srcField.Name)
			dstFieldValue := dstValue.FieldByName(dstField.Name)
			dstFieldValue.Set(srcFieldValue)
		}
	}

	return nil
}
