package utils

import (
	soft_del "github.com/lontten/lorm/soft-delete"
	"reflect"
)

func GetSoftDelType(t reflect.Type) soft_del.SoftDelType {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			delType, has := soft_del.SoftDelTypeMap[field.Type]
			if has {
				return delType
			}
			return GetSoftDelType(field.Type)
		}
	}
	return soft_del.None
}

func IsSoftDelFieldType(t reflect.Type) bool {
	_, has := soft_del.SoftDelTypeMap[t]
	return has
}
