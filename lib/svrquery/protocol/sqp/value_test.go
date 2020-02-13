package sqp

import (
	"reflect"
	"testing"
)

func TestDynamicValue_MarshalJSON(t *testing.T) {
	type fields struct {
		Type  DataType
		Value interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "byte value",
			fields: fields{
				Type:  Byte,
				Value: byte(128),
			},
			want: []byte(`128`),
		},
		{
			name: "uint16 value",
			fields: fields{
				Type:  Uint16,
				Value: uint16(500),
			},
			want: []byte(`500`),
		},
		{
			name: "uint32 value",
			fields: fields{
				Type:  Uint32,
				Value: uint32(100000),
			},
			want: []byte(`100000`),
		},
		{
			name: "uint64 value",
			fields: fields{
				Type:  Uint64,
				Value: uint64(1000000000),
			},
			want: []byte(`1000000000`),
		},
		{
			name: "string value",
			fields: fields{
				Type:  String,
				Value: "a string",
			},
			want: []byte(`"a string"`),
		},
		{
			name: "unknown type value",
			fields: fields{
				Type:  255,
				Value: "unknown type",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dv := &DynamicValue{
				Type:  tt.fields.Type,
				Value: tt.fields.Value,
			}
			got, err := dv.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("DynamicValue.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DynamicValue.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
