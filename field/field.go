package field

type ValueType int

const (
	None       ValueType = iota // 无,未定义
	Default                     // 字段默认值
	Null                        // 设为null
	Now                         // 当前时间
	UnixSecond                  // 秒时间戳
	UnixMilli                   // 毫秒时间戳
	UnixNano                    // 纳秒时间戳
	Val                         // 自定义值
	Increment                   // 字段自增
	Expression                  // 表达式
	ID                          // 主键id
)

type Value struct {
	Type  ValueType // 值类型
	Value any       // 值
}

type FValue struct {
	Name  string    // 字段名
	Type  ValueType // 值类型
	Value any       // 值
}

func (fv FValue) ToValue() Value {
	return Value{
		Type:  fv.Type,
		Value: fv.Value,
	}
}
