package lorm

import (
	"database/sql"
	"github.com/lontten/lorm/types/jsuuid"
	"reflect"
	"testing"
)
 

type K struct {
	Name *string
	Ha jsuuid.NullUUID
	Hb sql.NullBool
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