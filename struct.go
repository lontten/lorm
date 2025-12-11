//  Copyright 2025 lontten lontten@163.com
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package lorm

import (
	"errors"
	"reflect"
	"sort"
	"sync"
	"unicode"

	"github.com/lontten/lorm/field"
	"github.com/lontten/lorm/softdelete"
	"github.com/lontten/lorm/utils"
)

type StructValidFieldValueMap map[string]any

// --------------------comp-field-valuer---------
// v03 检查一个 struct/map 是否合法,valuer
var compFieldVCache = make(map[reflect.Type]struct{})
var mutexCompFieldVCache sync.Mutex

func checkCompFieldV(typ reflect.Type) error {
	mutexCompFieldVCache.Lock()
	defer mutexCompFieldVCache.Unlock()

	_, ok := compFieldVCache[typ]
	if ok {
		return nil
	}

	kind := typ.Kind()

	//struct
	if kind == reflect.Struct {
		err := checkStructFieldV(typ)
		if err != nil {
			return err
		} else {
			compFieldVCache[typ] = struct{}{}
			return nil
		}
	}
	//map
	if kind == reflect.Map {
		err := checkMapFieldV(typ)
		if err != nil {
			return err
		} else {
			compFieldVCache[typ] = struct{}{}
			return nil
		}
	}
	return errors.New("need a struct or map-scan")
}

// --------------------comp-field-valuer-scanner---------
// 检查一个 comp 是否合法
var compFieldVSCache = make(map[reflect.Type]reflect.Value)
var mutexCompFieldVSCache sync.Mutex

func checkCompFieldVS(va reflect.Value) error {
	mutexCompFieldVSCache.Lock()
	defer mutexCompFieldVSCache.Unlock()

	typ := va.Type()
	_, ok := compFieldVSCache[typ]
	if ok {
		return nil
	}

	kind := typ.Kind()

	//struct
	if kind == reflect.Struct {
		err := checkStructFieldVS(typ)
		if err == nil {
			compFieldVSCache[typ] = va
			return nil
		}
		return err
	}
	//map
	if kind == reflect.Map {
		is := checkMapFieldValue(va)
		if is {
			compFieldVSCache[typ] = va
			return nil
		}
	}
	return errors.New("checkCompFieldVS err;need a struct or map")
}

var fieldCache = make(map[reflect.Type][]compC)
var fieldMutex sync.Mutex

// 主键ID，转化为id
// tag== lrom:-  跳过
// 过滤掉首字母小写的字段
// 获取model对应的数据字段名：和其在model中的字段名
func getStructC(t reflect.Type) []compC {
	fields, ok := fieldCache[t]
	if ok {
		return fields
	}
	fieldMutex.Lock()
	defer fieldMutex.Unlock()
	fields, ok = fieldCache[t]
	if ok {
		return fields
	}
	return _getStructC(t, "")
}

func getStructCFMap(t reflect.Type) map[string]string {
	list := _getStructC(t, "")
	m := make(map[string]string, 0)
	for _, c := range list {
		m[c.columnName] = c.fieldName
	}
	return m
}

// struct 所有的字段名列表
func getStructCAllList(t reflect.Type) []string {
	list := _getStructC(t, "")
	m := make([]string, 0)
	for _, c := range list {
		m = append(m, c.columnName)
	}
	return m
}

// 排除 软删除字段
// struct 字段名列表
func getStructCList(t reflect.Type) []string {
	list := _getStructC(t, "")
	m := make([]string, 0)
	for _, c := range list {
		if c.isSoftDel {
			continue
		}
		m = append(m, c.columnName)
	}
	return m
}

// 主键ID，转化为id
// tag== db:-  跳过
// 过滤掉首字母小写的字段
// 获取model对应的数据字段名：和其在model中的字段名
// lormName tag中定义的字段对应的数据库字段名
func _getStructC(t reflect.Type, lormName string) (list []compC) {
	numField := t.NumField()
	for i := 0; i < numField; i++ {
		cc := compC{}

		structField := t.Field(i)
		if structField.Anonymous {
			tag := structField.Tag.Get("db")
			if tag == "-" {
				continue
			}
			data := _getStructC(structField.Type, tag)
			list = append(list, data...)
			continue
		}

		name := structField.Name

		// 过滤掉首字母小写的字段
		if unicode.IsLower([]rune(name)[0]) {
			continue
		}

		tag := structField.Tag.Get("db")
		if tag == "-" {
			continue
		}

		cc.kind = structField.Type.Kind()
		if cc.kind == reflect.Ptr {
			cc.canNull = true
			cc.kind = structField.Type.Elem().Kind()
		} else {
			canNull, isScanner := checkHandleNull(structField.Type)
			cc.canNull = canNull
			cc.isScanner = isScanner
		}

		if name == "ID" {
			cc.fieldName = "ID"
			cc.columnName = "id"
			list = append(list, cc)
			continue
		}

		if tag != "" {
			cc.fieldName = name
			cc.columnName = tag
			list = append(list, cc)
			continue
		}

		delType, has := softdelete.SoftDelTypeMap[t]
		if has {
			cc.fieldName = name
			if lormName == "" {
				value := softdelete.SoftDelTypeYesFVMap[delType]
				cc.columnName = value.Name
			} else {
				cc.columnName = lormName
			}
			list = append(list, cc)
			continue
		}

		cc.fieldName = name
		cc.columnName = utils.Camel2Case(name)
		list = append(list, cc)
	}
	return
}

func _getStructC_columnNameMap(t reflect.Type, lormName string) map[string]compC {
	cm := make(map[string]compC)
	list := _getStructC(t, lormName)
	for _, c := range list {
		cm[c.columnName] = c
	}
	return cm
}

type compC struct {
	fieldName   string // 字段名字
	columnName  string // 数据库字段名字
	columnIndex int
	isSoftDel   bool // 是否是软删除字段
	canNull     bool // 可以直接接收null；指针、实现了Valuer接口并处理了nil的结构体；（基础类型需要手动处理nil）
	isScanner   bool // 是否是Scanner
	kind        reflect.Kind
}

type compCV struct {
	fieldName  string // 字段名字
	columnName string // 数据库字段名字
	isSoftDel  bool   // 是否是软删除字段
	value      field.Value
	isZero     bool // 是否是零值
}

//--------------------------------- value -----------------------------------------

// 获取 struct 值
// 返回值类型有 Val,一种
func getStructCV(v reflect.Value) (list []compCV) {
	t := v.Type()
	cs := getStructC(t)
	for _, c := range cs {
		cv := compCV{
			fieldName:  c.fieldName,
			columnName: c.columnName,
			isSoftDel:  c.isSoftDel,
		}

		fieldV := v.FieldByName(c.fieldName)
		inter := getFieldInterZero(fieldV)
		if inter != nil {
			cv.value = field.Value{
				Type:  field.Val,
				Value: inter,
			}
		} else {
			cv.isZero = true
		}

		list = append(list, cv)
	}
	return list
}

type colName2fieldNameMap map[string]string
type compCVMap struct {
	//有效字段列表
	columns []string
	//有效值列表
	columnValues []field.Value

	modelZeroColumnNames      []string //零值字段列表
	modelNoSoftDelColumnNames []string // model 所有字段列表- 忽略软删除字段
	modelAllColumnNames       []string //所有字段列表

	//所有字段 dbName:fieldName
	modelAllCFNameMap colName2fieldNameMap
}

func getStructCVMap(v reflect.Value) (m compCVMap) {
	m = compCVMap{
		columns:                   make([]string, 0),
		columnValues:              make([]field.Value, 0),
		modelZeroColumnNames:      make([]string, 0),
		modelAllColumnNames:       make([]string, 0),
		modelNoSoftDelColumnNames: make([]string, 0),
		modelAllCFNameMap:         colName2fieldNameMap{},
	}
	list := getStructCV(v)
	for _, cv := range list {
		m.modelAllCFNameMap[cv.columnName] = cv.fieldName
		m.modelAllColumnNames = append(m.modelAllColumnNames, cv.columnName)
		if cv.isZero {
			m.modelZeroColumnNames = append(m.modelZeroColumnNames, cv.columnName)
		}
		if !cv.isSoftDel {
			m.modelNoSoftDelColumnNames = append(m.modelNoSoftDelColumnNames, cv.columnName)
		}
		if !cv.isZero && !cv.isSoftDel {
			m.columns = append(m.columns, cv.columnName)
			m.columnValues = append(m.columnValues, cv.value)
		}
	}
	return m
}

// 获取map[string]any
// 返回值类型有 None,Null,Val,三种
func getMapCV(v reflect.Value) (list []compCV) {
	keys := v.MapKeys()

	// 按键名字典序排序
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].String() < keys[j].String()
	})

	for _, k := range keys {
		list = append(list, compCV{
			columnName: k.String(),
			value:      getFieldInter(v.MapIndex(k)),
		})
	}
	return
}

// 获取map[string]any
// 返回值类型有 None,Null,Val,三种
func getMapCVMap(v reflect.Value) (m compCVMap) {
	list := getMapCV(v)
	for _, cv := range list {
		m.columns = append(m.columns, cv.columnName)
		m.columnValues = append(m.columnValues, cv.value)
	}
	return m
}
