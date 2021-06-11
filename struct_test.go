package lorm

import (
	"database/sql"
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
		Age string `tableName:"kkage"`
	}

	tableName:="kk"

	user := User{Name: "s"}
	users := make([]User, 0)
	f := func(structName string, dest interface{}) string {
		log.Println(structName)
		return "user"
	}
	println(f)
	config := OrmConfig{
		TableNamePrefix: "t_",
		TableNameFun:nil ,
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
			as.Equal(got,tt.want,"bu")
			//if got != tt.want {
			//	t.Errorf("getStructTableName() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
