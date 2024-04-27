package tools

import (
	"reflect"
	"time"
)

func CheckModel(model interface{}) bool {
	val := reflect.ValueOf(model)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if field.Kind() == reflect.Ptr && field.Type().Elem() == reflect.TypeOf(time.Time{}) {
			if field.IsNil() {
				return false
			}
			continue
		}

		switch field.Kind() {
		case reflect.Struct:
			if !CheckModel(field.Interface()) {
				return false
			}
		case reflect.Ptr:
			if field.IsNil() {
				return false
			}

			if field.Elem().Kind() == reflect.Struct {
				if !CheckModel(field.Elem().Interface()) {
					return false
				}
			}
		case reflect.Slice:
			for j := 0; j < field.Len(); j++ {
				if field.Index(j).Kind() == reflect.Struct {
					if !CheckModel(field.Index(j).Interface()) {
						return false
					}
				}
			}
		}
	}

	return true
}
