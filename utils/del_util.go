package utils

import (
	"github.com/lontten/lorm/softdelete"
	"reflect"
)

// GetSoftDelType
// 获取一个 struct 的 type 的软删除类型
func GetSoftDelType(t reflect.Type) softdelete.SoftDelType {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			delType, has := softdelete.SoftDelTypeMap[field.Type]
			if has {
				return delType
			}
			softDelType := GetSoftDelType(field.Type)
			if softDelType != softdelete.None {
				return softDelType
			}
		}
	}
	return softdelete.None
}

func IsSoftDelFieldType(t reflect.Type) bool {
	_, has := softdelete.SoftDelTypeMap[t]
	return has
}
