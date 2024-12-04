package soft_delete

import (
	"github.com/gofrs/uuid"
	"github.com/lontten/lorm/field"
	"gorm.io/plugin/soft_delete"
)
import (
	"github.com/lontten/lorm/types"
	"reflect"
)

type SoftDelType int

const (
	None SoftDelType = iota
	DelTimeNil
	DelBool
	DelBoolNil
	DelUUID
	DelUUIDNil
	DelInt32
	DelInt64
	DelGormSecond
	DelTimeGormMilli
	DelTimeGormNano
	DelGormFlag
)

var (
	SoftDelTypeMap = map[reflect.Type]SoftDelType{
		reflect.TypeOf(DeleteTimeNil{}):    DelTimeNil,
		reflect.TypeOf(DeleteBool{}):       DelBool,
		reflect.TypeOf(DeleteBoolNil{}):    DelBoolNil,
		reflect.TypeOf(DeleteUUID{}):       DelUUID,
		reflect.TypeOf(DeleteUUIDNil{}):    DelUUIDNil,
		reflect.TypeOf(DeleteInt32{}):      DelInt32,
		reflect.TypeOf(DeleteInt64{}):      DelInt64,
		reflect.TypeOf(DeleteGormSecond{}): DelGormSecond,
		reflect.TypeOf(DeleteGormMilli{}):  DelTimeGormMilli,
		reflect.TypeOf(DeleteGormNano{}):   DelTimeGormNano,
		reflect.TypeOf(DeleteGormFlag{}):   DelGormFlag,
	}

	SoftDelTypeNoFVMap = map[SoftDelType]field.FValue{
		DelTimeNil:       {Type: field.Null},
		DelBool:          {Name: "deleted_flag", Type: field.Expression, Value: "false"},
		DelBoolNil:       {Type: field.Null},
		DelUUID:          {Name: "deleted_uuid", Type: field.Val, Value: uuid.UUID{}},
		DelUUIDNil:       {Type: field.Null},
		DelInt32:         {Name: "deleted_flag", Type: field.Expression, Value: "0"},
		DelInt64:         {Name: "deleted_flag", Type: field.Expression, Value: "0"},
		DelGormSecond:    {Name: "deleted_at", Type: field.Expression, Value: "0"},
		DelTimeGormMilli: {Name: "deleted_at", Type: field.Expression, Value: "0"},
		DelTimeGormNano:  {Name: "deleted_at", Type: field.Expression, Value: "0"},
		DelGormFlag:      {Name: "deleted_flag", Type: field.Expression, Value: "0"},
	}
	SoftDelTypeYesFVMap = map[SoftDelType]field.FValue{
		DelTimeNil:       {Name: "deleted_at", Type: field.Now},
		DelBool:          {Name: "deleted_flag", Type: field.Expression, Value: "true"},
		DelBoolNil:       {Name: "deleted_flag", Type: field.Expression, Value: "true"},
		DelUUID:          {Name: "deleted_uuid", Type: field.Val, Value: types.V4()},
		DelUUIDNil:       {Name: "deleted_uuid", Type: field.Val, Value: types.V4()},
		DelInt32:         {Name: "deleted_flag", Type: field.ID},
		DelInt64:         {Name: "deleted_flag", Type: field.ID},
		DelGormSecond:    {Name: "deleted_at", Type: field.UnixSecond},
		DelTimeGormMilli: {Name: "deleted_at", Type: field.UnixMilli},
		DelTimeGormNano:  {Name: "deleted_at", Type: field.UnixNano},
		DelGormFlag:      {Name: "deleted_flag", Type: field.Expression, Value: "1"},
	}
)

// 软删除的删除时间通常是不看的，不用在意格式化问题

// gorm 秒软删除，储存类型为uint32,可储存272年
// 未删除时，为零值，删除时，为删除时间
type DeleteGormSecond struct {
	DeletedAt soft_delete.DeletedAt
}

// gorm 毫秒软删除，储存类型为uint64
// 未删除时，为零值，删除时，为删除时间
type DeleteGormMilli struct {
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:milli"`
}

// gorm 纳秒软删除，储存类型为uint64,可储存584年
// 未删除时，为零值，删除时，为删除时间
type DeleteGormNano struct {
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:nano"`
}

// gorm flag软删除
// 未删除时，为零值，删除时，为1
type DeleteGormFlag struct {
	IsDel soft_delete.DeletedAt `gorm:"softDelete:flag"`
}

// 未删除时，为nil，删除时，为删除时间
type DeleteTimeNil struct {
	DeletedAt *types.DateTime
}

// 未删除时，为false，删除时，为true
type DeleteBool struct {
	DeletedFlag bool
}

// 未删除时，为nil，删除时，为true
type DeleteBoolNil struct {
	DeletedFlag *bool
}

// 未删除时，为uuid.zero，删除时，为非zero
type DeleteUUID struct {
	DeletedUUID types.UUID
}

// 未删除时，为nil，删除时，为非nil
type DeleteUUIDNil struct {
	DeletedUUID *types.UUID
}

// 未删除时，为0，删除时，为i32 ID
type DeleteInt32 struct {
	DeletedFlag int32
}

// 未删除时，为0，删除时，为i64 ID
type DeleteInt64 struct {
	DeletedFlag int64
}
