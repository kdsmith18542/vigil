// Copyright (c) 2015-2017 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txscript

import (
	"fmt"
)

const (
	maxInt32 = 1<<31 - 1
	minInt32 = -1 << 31

	// MathOpCodeMaxScriptNumLen is the maximum number of bytes data being
	// interpreted as an integer may be for the majority of op codes.
	MathOpCodeMaxScriptNumLen = 4

	// CltvMaxScriptNumLen is the maximum number of bytes data being interpreted
	// as an integer may be for by-time and by-height locks as interpreted by
	// CHECKLOCKTIMEVERIFY.
	//
	// The value comes from the fact that the current transaction locktime
	// is a uint32 resulting in a maximum locktime of 2^32-1 (the year
	// 2106).  However, script numbers are signed and therefore a standard
	// 4-byte ScriptNum would only support up to a maximum of 2^31-1 (the
	// year 2038).  Thus, a 5-byte ScriptNum is needed since it will support
	// up to 2^39-1 which allows dates beyond the current locktime limit.
	CltvMaxScriptNumLen = 5

	// CsvMaxScriptNumLen is the maximum number of bytes data being interpreted
	// as an integer may be for by-time and by-height locks as interpreted by
	// CHECKSEQUENCEVERIFY.
	//
	// The value comes from the fact that the current transaction sequence
	// is a uint32 resulting in a maximum sequence of 2^32-1.  However,
	// ScriptNums are signed and therefore a standard 4-byte ScriptNum would
	// only support up to a maximum of 2^31-1.  Thus, a 5-byte ScriptNum is
	// needed since it will support up to 2^39-1 which allows sequences
	// beyond the current sequence limit.
	CsvMaxScriptNumLen = 5

	// altSigSuitesMaxscriptNumLen is the maximum number of bytes for the
	// type of alternative signature suite.
	altSigSuitesMaxscriptNumLen = 1
)

// ScriptNum represents a numeric value used in the scripting engine with
// special handling to deal with the subtle semantics required by consensus.
//
// All numbers are stored on the data and alternate stacks encoded as little
// endian with a sign bit.  All numeric opcodes such as OP_ADD, OP_SUB,
// and OP_MUL, are only allowed to operate on 4-byte integers in the range
// [-2^31 + 1, 2^31 - 1], however the results of numeric operations may overflow
// and remain valid so long as they are not used as inputs to other numeric
// operations or otherwise interpreted as an integer.
//
// For example, it is possible for OP_ADD to have 2^31 - 1 for its two operands
// resulting 2^32 - 2, which overflows, but is still pushed to the stack as the
// result of the addition.  That value can then be used as input to OP_VERIFY
// which will succeed because the data is being interpreted as a boolean.
// However, if that same value were to be used as input to another numeric
// opcode, such as OP_SUB, it must fail.
//
// This type handles the aforementioned requirements by storing all numeric
// operation results as an int64 to handle overflow and provides the Bytes
// method to get the serialized representation (including values that overflow).
//
// Then, whenever data is interpreted as an integer, it is converted to this
// type by using the MakeScriptNum function which will return an error if the
// number is out of range or not minimally encoded depending on parameters.
// Since all numeric opcodes involve pulling data from the stack and
// interpreting it as an integer, it provides the required behavior.
type ScriptNum int64

// checkMinimalDataEncoding returns whether or not the passed byte array adheres
// to the minimal encoding requirements.
func checkMinimalDataEncoding(v []byte) error {
	if len(v) == 0 {
		return nil
	}

	// Check that the number is encoded with the minimum possible
	// number of bytes.
	//
	// If the most-significant-byte - excluding the sign bit - is zero
	// then we're not minimal.  Note how this test also rejects the
	// negative-zero encoding, [0x80].
	if v[len(v)-1]&0x7f == 0 {
		// One exception: if there's more than one byte and the most
		// significant bit of the second-most-significant-byte is set
		// it would conflict with the sign bit.  An example of this case
		// is +-255, which encode to 0xff00 and 0xff80 respectively.
		// (big-endian).
		if len(v) == 1 || v[len(v)-2]&0x80 == 0 {
			str := fmt.Sprintf("numeric value encoded as %x is "+
				"not minimally encoded", v)
			return scriptError(ErrMinimalData, str)
		}
	}

	return nil
}

// Bytes returns the number serialized as a little endian with a sign bit.
//
// Example encodings:
//
//	   127 -> [0x7f]
//	  -127 -> [0xff]
//	   128 -> [0x80 0x00]
//	  -128 -> [0x80 0x80]
//	   129 -> [0x81 0x00]
//	  -129 -> [0x81 0x80]
//	   256 -> [0x00 0x01]
//	  -256 -> [0x00 0x81]
//	 32767 -> [0xff 0x7f]
//	-32767 -> [0xff 0xff]
//	 32768 -> [0x00 0x80 0x00]
//	-32768 -> [0x00 0x80 0x80]
func (n ScriptNum) Bytes() []byte {
	// Zero encodes as an empty byte slice.
	if n == 0 {
		return nil
	}

	// Take the absolute value and keep track of whether it was originally
	// negative.
	nu64 := uint64(n)
	isNegative := n < 0
	if isNegative {
		nu64 = uint64(-n)
	}

	// Encode to little endian.  The maximum number of encoded bytes is 9
	// (8 bytes for max int64 plus a potential byte for sign extension).
	result := make([]byte, 0, 9)
	for nu64 > 0 {
		result = append(result, byte(nu64&0xff))
		nu64 >>= 8
	}

	// When the most significant byte already has the high bit set, an
	// additional high byte is required to indicate whether the number is
	// negative or positive.  The additional byte is removed when converting
	// back to an integral and its high bit is used to denote the sign.
	//
	// Otherwise, when the most significant byte does not already have the
	// high bit set, use it to indicate the value is negative, if needed.
	if result[len(result)-1]&0x80 != 0 {
		extraByte := byte(0x00)
		if isNegative {
			extraByte = 0x80
		}
		result = append(result, extraByte)
	} else if isNegative {
		result[len(result)-1] |= 0x80
	}

	return result
}

// Int32 returns the script number clamped to a valid int32.  That is to say
// when the script number is higher than the max allowed int32, the max int32
// value is returned and vice versa for the minimum value.  Note that this
// behavior is different from a simple int32 cast because that truncates
// and the consensus rules dictate numbers which are directly cast to ints
// provide this behavior.
//
// In practice, for most opcodes, the number should never be out of range since
// it will have been created with MakeScriptNum using the defaultScriptLen
// value, which rejects them.  In case something in the future ends up calling
// this function against the result of some arithmetic, which IS allowed to be
// out of range before being reinterpreted as an integer, this will provide the
// correct behavior.
func (n ScriptNum) Int32() int32 {
	if n > maxInt32 {
		return maxInt32
	}

	if n < minInt32 {
		return minInt32
	}

	return int32(n)
}

// MakeScriptNum interprets the passed serialized bytes as an encoded integer
// and returns the result as a script number.
//
// Since the consensus rules dictate that serialized bytes interpreted as ints
// are only allowed to be in the range determined by a maximum number of bytes,
// on a per opcode basis, an error will be returned when the provided bytes
// would result in a number outside of that range.  In particular, the range for
// the vast majority of opcodes dealing with numeric values are limited to 4
// bytes and therefore will pass that value to this function resulting in an
// allowed range of [-2^31 + 1, 2^31 - 1].
//
// An error will be returned if the encoding is not represented with the
// smallest possible number of bytes or is the negative 0 encoding, [0x80].  For
// example, consider the number 127.  It could be encoded as [0x7f], [0x7f
// 0x00], [0x7f 0x00 0x00 ...], etc.  All forms except [0x7f] will return an
// error.
//
// The scriptNumLen is the maximum number of bytes the encoded value can be
// before an ErrNumOutOfRange is returned.  This effectively limits the range of
// allowed values.
//
// WARNING: Great care should be taken if passing a value larger than
// defaultScriptNumLen, which could lead to addition and multiplication
// overflows.
//
// See the Bytes function documentation for example encodings.
func MakeScriptNum(v []byte, scriptNumLen int) (ScriptNum, error) {
	// Interpreting data requires that it is not larger than
	// the passed scriptNumLen value.
	if len(v) > scriptNumLen {
		str := fmt.Sprintf("numeric value encoded as %x is %d bytes "+
			"which exceeds the max allowed of %d", v, len(v),
			scriptNumLen)
		return 0, scriptError(ErrNumOutOfRange, str)
	}

	// Enforce minimal encoding.
	if err := checkMinimalDataEncoding(v); err != nil {
		return 0, err
	}

	// Zero is encoded as an empty byte slice.
	if len(v) == 0 {
		return 0, nil
	}

	// Decode from little endian.
	var result int64
	for i, val := range v {
		result |= int64(val) << uint8(8*i)
	}

	// When the most significant byte of the input bytes has the sign bit
	// set, the result is negative.  So, remove the sign bit from the result
	// and make it negative.
	if v[len(v)-1]&0x80 != 0 {
		// The maximum length of v has already been determined to be 4
		// above, so uint8 is enough to cover the max possible shift
		// value of 24.
		result &= ^(int64(0x80) << uint8(8*(len(v)-1)))
		return ScriptNum(-result), nil
	}

	return ScriptNum(result), nil
}
