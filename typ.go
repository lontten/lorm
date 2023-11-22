package lorm

// todo 下面未重构--------------
type ArgArray []interface{}

func ArrayOf(v ...interface{}) ArgArray {
	return v
}
