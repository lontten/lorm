package lorm

import (
	"database/sql"
	"fmt"
	"github.com/lontten/lorm/types/jsuuid"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
)

type K struct {
	Name *string
	Ha   jsuuid.NullUUID
	Hb   sql.NullBool
}

func Test_baseStructValue(t *testing.T) {
	type args struct {
		v reflect.Value
	}
	tests := []struct {
		name            string
		args            args
		wantStructValue reflect.Value
		wantErr         bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStructValue, err := baseStructValue(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("baseStructValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotStructValue, tt.wantStructValue) {
				t.Errorf("baseStructValue() gotStructValue = %v, want %v", gotStructValue, tt.wantStructValue)
			}
		})
	}
}

func Test_getStructTableName(t *testing.T) {
	as := assert.New(t)
	type args struct {
		dest   interface{}
		config OrmConfig
	}

	type User struct {
		Name string `tableName:"kk"`
		Age  string `tableName:"kkage"`
	}

	tableName := "kk"

	user := User{Name: "s"}
	users := make([]User, 0)
	f := func(structName string, dest interface{}) string {
		log.Println(structName)
		return "user"
	}
	println(f)
	config := OrmConfig{
		TableNamePrefix: "t_",
		TableNameFun:    nil,
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "获取 v 的 tableName",
			args: args{
				dest:   user,
				config: config,
			},
			want:    tableName,
			wantErr: false,
		},

		{
			name: "获取 v 的 tableName",
			args: args{
				dest:   &user,
				config: config,
			},
			want:    tableName,
			wantErr: false,
		},

		{
			name: "获取 v 的 tableName",
			args: args{
				dest:   users,
				config: config,
			},
			want:    tableName,
			wantErr: false,
		},

		{
			name: "获取 v 的 tableName",
			args: args{
				dest:   &users,
				config: config,
			},
			want:    tableName,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getStructTableName(tt.args.dest, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("getStructTableName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			as.Equal(got, tt.want, "bu")
			//if got != tt.want {
			//	t.Errorf("getStructTableName() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func Test_baseSlic2ePtrType(t *testing.T) {
	switch 2 {
	case 2:
		fmt.Println(2)
		fallthrough
	case 4:
		fmt.Println(4)
	case 5:
		fmt.Println(5)
	default:
		fmt.Println(6)
	}
	fmt.Println(0)

}
func Test_baseSlicePtrType(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	type a struct {
	}
	s := reflect.TypeOf(a{})
	sp := reflect.TypeOf(&a{})
	as := make([]a, 0)
	ass := reflect.TypeOf(as)
	assp := reflect.TypeOf(&as)

	aps := make([](*a), 0)
	apss := reflect.TypeOf(aps)
	apssp := reflect.TypeOf(&aps)

	tests := []struct {
		name           string
		args           args
		wantTyp        int
		wantStructType reflect.Type
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			name:           "struct",
			args:           args{t: s},
			wantTyp:        0,
			wantStructType: s,
			wantErr:        false,
		},

		{
			name:           "struct ptr",
			args:           args{t: sp},
			wantTyp:        1,
			wantStructType: s,
			wantErr:        false,
		},

		{
			name:           "struct slice",
			args:           args{t: ass},
			wantTyp:        2,
			wantStructType: s,
			wantErr:        false,
		},

		{
			name:           "struct slice ptr",
			args:           args{t: assp},
			wantTyp:        0,
			wantStructType: nil,
			wantErr:        true,
		},

		{
			name:           "struct  ptr slice",
			args:           args{t: apss},
			wantTyp:        0,
			wantStructType: nil,
			wantErr:        true,
		},

		{
			name:           "struct  ptr slice ptr",
			args:           args{t: apssp},
			wantTyp:        0,
			wantStructType: nil,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTyp, gotStructType, err := baseStructTypeSliceOrPtr(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("baseStructTypeSliceOrPtr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTyp != tt.wantTyp {
				t.Errorf("baseStructTypeSliceOrPtr() gotTyp = %v, want %v", gotTyp, tt.wantTyp)
			}
			if !reflect.DeepEqual(gotStructType, tt.wantStructType) {
				t.Errorf("baseStructTypeSliceOrPtr() gotStructType = %v, want %v", gotStructType, tt.wantStructType)
			}
		})
	}
}

func Test_baseSliceType(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	type a struct {
	}
	s := reflect.TypeOf(a{})
	sp := reflect.TypeOf(&a{})
	as := make([]a, 0)
	ass := reflect.TypeOf(as)
	assp := reflect.TypeOf(&as)

	aps := make([](*a), 0)
	apss := reflect.TypeOf(aps)
	apssp := reflect.TypeOf(&aps)

	tests := []struct {
		name           string
		args           args
		wantStructType reflect.Type
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			name:           "struct",
			args:           args{t: s},
			wantStructType: nil,
			wantErr:        true,
		},

		{
			name:           "struct ptr",
			args:           args{t: sp},
			wantStructType: nil,
			wantErr:        true,
		},

		{
			name:           "struct slice",
			args:           args{t: ass},
			wantStructType: ass,
			wantErr:        false,
		},

		{
			name:           "struct slice ptr",
			args:           args{t: assp},
			wantStructType: ass,
			wantErr:        false,
		},

		{
			name:           "struct  ptr slice",
			args:           args{t: apss},
			wantStructType: apss,
			wantErr:        false,
		},

		{
			name:           "struct  ptr slice ptr",
			args:           args{t: apssp},
			wantStructType: apss,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStructType, err := baseSliceTypePtr(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("baseSliceType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotStructType, tt.wantStructType) {
				t.Errorf("baseSliceType() gotStructType = %v, want %v", gotStructType, tt.wantStructType)
			}
		})
	}
}