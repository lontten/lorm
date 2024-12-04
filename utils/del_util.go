package utils

import (
	"github.com/lontten/lorm/soft-delete"
	"reflect"
)

func GetSoftDelType(t reflect.Type) soft_delete.SoftDelType {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			delType, has := soft_delete.SoftDelTypeMap[field.Type]
			if has {
				return delType
			}
			return GetSoftDelType(field.Type)
		}
	}
	return soft_delete.None
}

func IsSoftDelFieldType(t reflect.Type) bool {
	_, has := soft_delete.SoftDelTypeMap[t]
	return has
}
