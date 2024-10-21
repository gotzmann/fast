package gguf

import (
	"encoding/binary"
	"fmt"
	"os"
)

type MetadataValueType uint32

const (
	UINT8 MetadataValueType = iota
	INT8
	UINT16
	INT16
	UINT32
	INT32
	FLOAT32
	BOOL
	STRING
	ARRAY
	UINT64
	INT64
	FLOAT64
)

func (t MetadataValueType) String() string {
	switch t {
	case UINT8:
		return "UINT8"
	case INT8:
		return "INT8"
	case UINT16:
		return "UINT16"
	case INT16:
		return "INT16"
	case UINT32:
		return "UINT32"
	case INT32:
		return "INT32"
	case FLOAT32:
		return "FLOAT32"
	case BOOL:
		return "BOOL"
	case STRING:
		return "STRING"
	case ARRAY:
		return "ARRAY"
	case UINT64:
		return "UINT64"
	case INT64:
		return "INT64"
	case FLOAT64:
		return "FLOAT64"
	default:
		return "Unknown"
	}
}

type MetadataArray struct {
	Type   MetadataValueType
	Length uint64
	Array  []MetadataValue
}

func (a *MetadataArray) String() string {
	str := fmt.Sprintf("Type: %s, Length: %d, Array:\n", a.Type.String(), a.Length)
	for i, v := range a.Array {
		str += fmt.Sprintf("\tValue %d: %s\n", i, v.String(a.Type))
	}
	return str
}

func (a *MetadataArray) Read(file *os.File) error {
	err := binary.Read(file, binary.LittleEndian, &a.Type)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &a.Length)
	if err != nil {
		return err
	}
	a.Array = make([]MetadataValue, a.Length)
	for i := range a.Array {
		err = a.Array[i].Read(a.Type, file)
		if err != nil {
			return err
		}
	}
	return nil
}

type MetadataValue struct {
	U8  uint8
	I8  int8
	U16 uint16
	I16 int16
	U32 uint32
	I32 int32
	F32 float32
	U64 uint64
	I64 int64
	F64 float64
	B   bool
	STR String
	ARR MetadataArray
}

func (v *MetadataValue) String(forType MetadataValueType) string {
	switch forType {
	case UINT8:
		return fmt.Sprintf("%d", v.U8)
	case INT8:
		return fmt.Sprintf("%d", v.I8)
	case UINT16:
		return fmt.Sprintf("%d", v.U16)
	case INT16:
		return fmt.Sprintf("%d", v.I16)
	case UINT32:
		return fmt.Sprintf("%d", v.U32)
	case INT32:
		return fmt.Sprintf("%d", v.I32)
	case FLOAT32:
		return fmt.Sprintf("%f", v.F32)
	case BOOL:
		return fmt.Sprintf("%t", v.B)
	case STRING:
		return v.STR.String()
	case ARRAY:
		return fmt.Sprintf("%v", v.ARR)
	case UINT64:
		return fmt.Sprintf("%d", v.U64)
	case INT64:
		return fmt.Sprintf("%d", v.I64)
	case FLOAT64:
		return fmt.Sprintf("%f", v.F64)
	default:
		return "Unknown"
	}
}

func (v *MetadataValue) Read(forType MetadataValueType, file *os.File) error {
	switch forType {
	case UINT8:
		return binary.Read(file, binary.LittleEndian, &v.U8)
	case INT8:
		return binary.Read(file, binary.LittleEndian, &v.I8)
	case UINT16:
		return binary.Read(file, binary.LittleEndian, &v.U16)
	case INT16:
		return binary.Read(file, binary.LittleEndian, &v.I16)
	case UINT32:
		return binary.Read(file, binary.LittleEndian, &v.U32)
	case INT32:
		return binary.Read(file, binary.LittleEndian, &v.I32)
	case FLOAT32:
		return binary.Read(file, binary.LittleEndian, &v.F32)
	case BOOL:
		return binary.Read(file, binary.LittleEndian, &v.B)
	case STRING:
		return v.STR.Read(file)
	case ARRAY:
		return v.ARR.Read(file)
	case UINT64:
		return binary.Read(file, binary.LittleEndian, &v.U64)
	case INT64:
		return binary.Read(file, binary.LittleEndian, &v.I64)
	case FLOAT64:
		return binary.Read(file, binary.LittleEndian, &v.F64)
	default:
		return fmt.Errorf("unknown type: %d", forType)
	}
}

type MetadataKV struct {
	Key   String
	Type  MetadataValueType
	Value MetadataValue
}

func (kv *MetadataKV) String() string {
	vStr := kv.Value.String(kv.Type)
	// If value string exceeds 100 characters, only show "too long to display"
	if len(vStr) > 100 {
		vStr = fmt.Sprintf("Value too long to display (%d characters)", len(vStr))
	}
	return fmt.Sprintf("Key: %s, Type: %s, Value: %s", kv.Key.String(), kv.Type.String(), vStr)
}

func (kv *MetadataKV) Read(file *os.File) error {
	err := kv.Key.Read(file)
	if err != nil {
		return err
	}

	err = binary.Read(file, binary.LittleEndian, &kv.Type)
	if err != nil {
		return err
	}

	switch kv.Type {
	case UINT8:
		err = binary.Read(file, binary.LittleEndian, &kv.Value.U8)
	case INT8:
		err = binary.Read(file, binary.LittleEndian, &kv.Value.I8)
	case UINT16:
		err = binary.Read(file, binary.LittleEndian, &kv.Value.U16)
	case INT16:
		err = binary.Read(file, binary.LittleEndian, &kv.Value.I16)
	case UINT32:
		err = binary.Read(file, binary.LittleEndian, &kv.Value.U32)
	case INT32:
		err = binary.Read(file, binary.LittleEndian, &kv.Value.I32)
	case FLOAT32:
		err = binary.Read(file, binary.LittleEndian, &kv.Value.F32)
	case BOOL:
		err = binary.Read(file, binary.LittleEndian, &kv.Value.B)
	case STRING:
		err = kv.Value.STR.Read(file)
	case ARRAY:
		err = binary.Read(file, binary.LittleEndian, &kv.Value.ARR.Type)
		if err != nil {
			return err
		}
		err = binary.Read(file, binary.LittleEndian, &kv.Value.ARR.Length)
		if err != nil {
			return err
		}
		kv.Value.ARR.Array = make([]MetadataValue, kv.Value.ARR.Length)
		for i := range kv.Value.ARR.Array {
			err = kv.Value.ARR.Array[i].Read(kv.Value.ARR.Type, file)
			if err != nil {
				return err
			}
		}
	case UINT64:
		err = binary.Read(file, binary.LittleEndian, &kv.Value.U64)
	case INT64:
		err = binary.Read(file, binary.LittleEndian, &kv.Value.I64)
	case FLOAT64:
		err = binary.Read(file, binary.LittleEndian, &kv.Value.F64)
	}

	return err
}