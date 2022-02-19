package lorm

type ArgArray []interface{}

func ArrayOf(v ...interface{}) ArgArray {
	return v
}
