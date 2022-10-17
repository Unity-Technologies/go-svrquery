package sqp

import "encoding/json"

// DynamicValue represents dynamically typed values
type DynamicValue struct {
	Type  DataType
	Value interface{}
}

// NewDynamicValue returns a DynamicValue loaded from a packetReader.
// As this route determines the type dynamically, the read count is
// one byte greater to account for the type byte.
func NewDynamicValue(r *packetReader) (int64, *DynamicValue, error) {
	dv := &DynamicValue{}
	dt, err := r.ReadByte()
	if err != nil {
		return 0, nil, err
	}

	n, err := dynamicValue(dv, DataType(dt), r)
	if err != nil {
		return 1 + n, nil, err
	}

	return 1 + n, dv, nil
}

// NewDynamicValueWithType returns a DynamicValue of a given type loaded from a packetReader
// As no type is read via this route, the read count is exactly the size of the value read.
func NewDynamicValueWithType(r *packetReader, dt DataType) (int64, *DynamicValue, error) {
	dv := &DynamicValue{}

	n, err := dynamicValue(dv, dt, r)
	if err != nil {
		return n, nil, err
	}

	return n, dv, nil
}

func dynamicValue(dv *DynamicValue, dt DataType, r *packetReader) (int64, error) {
	var err error
	dv.Type = dt
	switch dv.Type {
	case Byte:
		dv.Value, err = r.ReadByte()
		return int64(Byte.Size()), err
	case Uint16:
		dv.Value, err = r.ReadUint16()
		return int64(Uint16.Size()), err
	case Uint32:
		dv.Value, err = r.ReadUint32()
		return int64(Uint32.Size()), err
	case Uint64:
		dv.Value, err = r.ReadUint64()
		return int64(Uint64.Size()), err
	case String:
		var count int64
		count, dv.Value, err = r.ReadString()
		return count, err
	case Float32:
		dv.Value, err = r.ReadFloat32()
		return int64(Float32.Size()), err
	}

	return 0, ErrUnknownDataType(dv.Type)
}

// Byte returns the value as a byte
func (dv *DynamicValue) Byte() byte {
	return dv.Value.(byte)
}

// Uint16 returns the value as a uint16
func (dv *DynamicValue) Uint16() uint16 {
	return dv.Value.(uint16)
}

// Uint32 returns the value as a uint32
func (dv *DynamicValue) Uint32() uint32 {
	return dv.Value.(uint32)
}

// Uint64 returns the value as a uint64
func (dv *DynamicValue) Uint64() uint64 {
	return dv.Value.(uint64)
}

// String returns the value as a string
func (dv *DynamicValue) String() string {
	return dv.Value.(string)
}

// Float32 returns the value as a float32
func (dv *DynamicValue) Float32() float32 {
	return dv.Value.(float32)
}

// MarshalJSON returns the json marshalled version of the dynamic value
func (dv *DynamicValue) MarshalJSON() ([]byte, error) {
	switch dv.Type {
	case Byte, Uint16, Uint32, Uint64, String, Float32:
		return json.Marshal(dv.Value)
	}
	return nil, ErrUnknownDataType(dv.Type)
}
