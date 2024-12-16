package utils

import (
	"github.com/lontten/lorm/softdelete"
	"reflect"
)

func GetSoftDelType(t reflect.Type) softdelete.SoftDelType {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			delType, has := softdelete.SoftDelTypeMap[field.Type]
			if has {
				return delType
			}
			return GetSoftDelType(field.Type)
		}
	}
	return softdelete.None
}

func IsSoftDelFieldType(t reflect.Type) bool {
	delType := GetSoftDelType(t)
	return delType != softdelete.None
}
