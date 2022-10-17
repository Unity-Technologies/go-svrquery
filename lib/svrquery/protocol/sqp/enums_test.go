package sqp

import "testing"

func TestDataType_Size(t *testing.T) {
	tests := []struct {
		name string
		dt   DataType
		want int
	}{
		{
			name: "Byte size",
			dt:   Byte,
			want: 1,
		},
		{
			name: "Uint16 size",
			dt:   Uint16,
			want: 2,
		},
		{
			name: "Uint32 size",
			dt:   Uint32,
			want: 4,
		},
		{
			name: "Uint64 size",
			dt:   Uint64,
			want: 8,
		},
		{
			name: "String size",
			dt:   String,
			want: -1,
		},
		{
			name: "Float32 size",
			dt:   Float32,
			want: 4,
		},
		{
			name: "Unknown size",
			dt:   DataType(99),
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dt.Size(); got != tt.want {
				t.Errorf("DataType.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
