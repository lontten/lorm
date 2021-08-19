package lorm

import (
	"database/sql/driver"
	"github.com/lontten/lorm/types"
	"github.com/lontten/lorm/utils"
	"github.com/pkg/errors"
	"reflect"
	"sync"
)

// * struct
// * slice
// slice

// 检查数据类型，
//获取基础 value
//bool 为 是否为 slice类型
func (c *OrmContext) initTargetDest(dest interface{}) {
	v := reflect.ValueOf(dest)
	arr := make([]reflect.Value, 0)
	isPtr, baseValue := basePtrValue(v)
	is, base := baseStructValue(baseValue)
	if isPtr && is {
		arr = append(arr, base)
		c.isSlice = false
		c.destBaseValueArr = arr
		c.destBaseValue = base
		return
	}
	//slice
	is, has, base := baseSliceValue(baseValue)
	if !has {
		c.err = ErrContainEmpty
		return
	}
	if !is {
		c.err = errors.New("need ptr or slice")
		return
	}
	e := base.Index(0)
	if e.Kind() != reflect.Struct {
		c.err = errors.New("need a slice struct type")
		return
	}
	c.isSlice = true
	c.destBaseValueArr = utils.ToSliceValue(baseValue)
	c.destBaseValue = e
}

var targetDestValidCache = make(map[reflect.Type]error)
var mutextargetDestValidCache sync.Mutex

func (c *OrmContext) checkTargetDestField() {
	v := c.destBaseValue
	mutextargetDestValidCache.Lock()
	defer mutextargetDestValidCache.Unlock()

	typ := v.Type()
	b, ok := targetDestValidCache[typ]
	if ok {
		c.err = b
		return
	}

	numField := v.NumField()
	for i := 0; i < numField; i++ {
		field := v.Field(i)

		typ, validField, ok := baseStructValidField(field)
		if !ok {
			c.err = errors.New("struct field " + field.String() + " need field is ptr slice struct")
			return
		}
		//为 struct类型
		if typ == 3 {
			_, ok = validField.Interface().(driver.Valuer)
			if !ok {
				c.err = errors.New("struct field " + field.String() + " need imp sql Value")
				return
			}
			_, ok = validField.Interface().(types.NullEr)
			if !ok {
				c.err = errors.New("struct field " + field.String() + " need imp core NullEr ")
				return
			}

		}
	}
	structValidCache[typ] = nil
	return

}
