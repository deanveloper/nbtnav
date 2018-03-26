package nbt

import "errors"

// Represents a tag in NBT Format.
// If unnamed, Name will be the empty string
// If valueless, Value will be nil
type Tag struct {
	Type  TagType
	Name  string
	Value interface{}
}

// TagType implements fmt.Stringer
type TagType byte

const (
	// Represents the TAG_End type
	TagEnd TagType = iota

	// Represents the TAG_Byte type. Called Int8 because TAG_Byte is signed.
	TagInt8

	// Represents the TAG_Short type. Called Int16 for consistiency.
	TagInt16

	// Represents the TAG_Int type. Called Int32 for consistiency.
	TagInt32

	// Represents the TAG_Long type. Called Int64 for consistiency.
	TagInt64

	// Represents the TAG_Float type. Called Float32 to match Go types.
	TagFloat32

	// Represents the TAG_Double type. Called Float64 to match Go types.
	TagFloat64

	// Represents the TAG_Byte_Array type. Called Int8 because the bytes are signed.
	TagInt8Array

	// Represents the TAG_String type.
	TagString

	// Represents the TAG_List type.
	TagList

	// Represents the TAG_Compound type.
	TagCompound

	// Represents the TAG_Int_Array type. Called Int32 for consistency.
	TagInt32Array
)

func (t TagType) String() string {
	switch t {
	case TagEnd:
		return "End"
	case TagInt8:
		return "Byte"
	case TagInt16:
		return "Short"
	case TagInt32:
		return "Int"
	case TagInt64:
		return "Long"
	case TagFloat32:
		return "Float"
	case TagFloat64:
		return "Double"
	case TagInt8Array:
		return "ByteArray"
	case TagString:
		return "String"
	case TagList:
		return "List"
	case TagCompound:
		return "Compound"
	case TagInt32Array:
		return "IntArray"
	}
	panic(errors.New("Oh no"))
}
