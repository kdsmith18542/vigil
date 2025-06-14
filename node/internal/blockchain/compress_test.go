// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// hexToBytes converts the passed hex string into bytes and will panic if there
// is an error.  This is only provided for the hard-coded constants so errors in
// the source code can be detected. It will only (and must only) be called with
// hard-coded values.
func hexToBytes(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic("invalid hex in source file: " + s)
	}
	return b
}

// TestVLQ ensures the variable length quantity serialization, deserialization,
// and size calculation works as expected.
func TestVLQ(t *testing.T) {
	t.Parallel()

	tests := []struct {
		val        uint64
		serialized []byte
	}{
		{0, hexToBytes("00")},
		{1, hexToBytes("01")},
		{127, hexToBytes("7f")},
		{128, hexToBytes("8000")},
		{129, hexToBytes("8001")},
		{255, hexToBytes("807f")},
		{256, hexToBytes("8100")},
		{16383, hexToBytes("fe7f")},
		{16384, hexToBytes("ff00")},
		{16511, hexToBytes("ff7f")}, // Max 2-byte value
		{16512, hexToBytes("808000")},
		{16513, hexToBytes("808001")},
		{16639, hexToBytes("80807f")},
		{32895, hexToBytes("80ff7f")},
		{2113663, hexToBytes("ffff7f")}, // Max 3-byte value
		{2113664, hexToBytes("80808000")},
		{270549119, hexToBytes("ffffff7f")}, // Max 4-byte value
		{270549120, hexToBytes("8080808000")},
		{2147483647, hexToBytes("86fefefe7f")},
		{2147483648, hexToBytes("86fefeff00")},
		{4294967295, hexToBytes("8efefefe7f")}, // Max uint32, 5 bytes
	}

	for _, test := range tests {
		// Ensure the function to calculate the serialized size without
		// actually serializing the value is calculated properly.
		gotSize := serializeSizeVLQ(test.val)
		if gotSize != len(test.serialized) {
			t.Errorf("serializeSizeVLQ: did not get expected size "+
				"for %d - got %d, want %d", test.val, gotSize,
				len(test.serialized))
			continue
		}

		// Ensure the value serializes to the expected bytes.
		gotBytes := make([]byte, gotSize)
		gotBytesWritten := putVLQ(gotBytes, test.val)
		if !bytes.Equal(gotBytes, test.serialized) {
			t.Errorf("putVLQUnchecked: did not get expected bytes "+
				"for %d - got %x, want %x", test.val, gotBytes,
				test.serialized)
			continue
		}
		if gotBytesWritten != len(test.serialized) {
			t.Errorf("putVLQUnchecked: did not get expected number "+
				"of bytes written for %d - got %d, want %d",
				test.val, gotBytesWritten, len(test.serialized))
			continue
		}

		// Ensure the serialized bytes deserialize to the expected
		// value.
		gotVal, gotBytesRead := deserializeVLQ(test.serialized)
		if gotVal != test.val {
			t.Errorf("deserializeVLQ: did not get expected value "+
				"for %x - got %d, want %d", test.serialized,
				gotVal, test.val)
			continue
		}
		if gotBytesRead != len(test.serialized) {
			t.Errorf("deserializeVLQ: did not get expected number "+
				"of bytes read for %d - got %d, want %d",
				test.serialized, gotBytesRead,
				len(test.serialized))
			continue
		}
	}
}

// TestScriptCompression ensures the domain-specific script compression and
// decompression works as expected.
func TestScriptCompression(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		scriptVersion uint16
		uncompressed  []byte
		compressed    []byte
	}{
		{
			name:          "nil",
			scriptVersion: 0,
			uncompressed:  nil,
			compressed:    hexToBytes("40"),
		},
		{
			name:          "pay-to-pubkey-hash 1",
			scriptVersion: 0,
			uncompressed:  hexToBytes("76a9141018853670f9f3b0582c5b9ee8ce93764ac32b9388ac"),
			compressed:    hexToBytes("001018853670f9f3b0582c5b9ee8ce93764ac32b93"),
		},
		{
			name:          "pay-to-pubkey-hash 2",
			scriptVersion: 0,
			uncompressed:  hexToBytes("76a914e34cce70c86373273efcc54ce7d2a491bb4a0e8488ac"),
			compressed:    hexToBytes("00e34cce70c86373273efcc54ce7d2a491bb4a0e84"),
		},
		{
			name:          "pay-to-script-hash 1",
			scriptVersion: 0,
			uncompressed:  hexToBytes("a914da1745e9b549bd0bfa1a569971c77eba30cd5a4b87"),
			compressed:    hexToBytes("01da1745e9b549bd0bfa1a569971c77eba30cd5a4b"),
		},
		{
			name:          "pay-to-script-hash 2",
			scriptVersion: 0,
			uncompressed:  hexToBytes("a914f815b036d9bbbce5e9f2a00abd1bf3dc91e9551087"),
			compressed:    hexToBytes("01f815b036d9bbbce5e9f2a00abd1bf3dc91e95510"),
		},
		{
			name:          "pay-to-pubkey compressed 0x02",
			scriptVersion: 0,
			uncompressed:  hexToBytes("2102192d74d0cb94344c9569c2e77901573d8d7903c3ebec3a957724895dca52c6b4ac"),
			compressed:    hexToBytes("02192d74d0cb94344c9569c2e77901573d8d7903c3ebec3a957724895dca52c6b4"),
		},
		{
			name:          "pay-to-pubkey compressed 0x03",
			scriptVersion: 0,
			uncompressed:  hexToBytes("2103b0bd634234abbb1ba1e986e884185c61cf43e001f9137f23c2c409273eb16e65ac"),
			compressed:    hexToBytes("03b0bd634234abbb1ba1e986e884185c61cf43e001f9137f23c2c409273eb16e65"),
		},
		{
			name:          "pay-to-pubkey uncompressed 0x04 even",
			scriptVersion: 0,
			uncompressed:  hexToBytes("4104192d74d0cb94344c9569c2e77901573d8d7903c3ebec3a957724895dca52c6b40d45264838c0bd96852662ce6a847b197376830160c6d2eb5e6a4c44d33f453eac"),
			compressed:    hexToBytes("04192d74d0cb94344c9569c2e77901573d8d7903c3ebec3a957724895dca52c6b4"),
		},
		{
			name:          "pay-to-pubkey uncompressed 0x04 odd",
			scriptVersion: 0,
			uncompressed:  hexToBytes("410411db93e1dcdb8a016b49840f8c53bc1eb68a382e97b1482ecad7b148a6909a5cb2e0eaddfb84ccf9744464f82e160bfa9b8b64f9d4c03f999b8643f656b412a3ac"),
			compressed:    hexToBytes("0511db93e1dcdb8a016b49840f8c53bc1eb68a382e97b1482ecad7b148a6909a5c"),
		},
		{
			name:          "pay-to-pubkey invalid pubkey",
			scriptVersion: 0,
			uncompressed:  hexToBytes("3302aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaac"),
			compressed:    hexToBytes("633302aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaac"),
		},
		{
			name:          "null data",
			scriptVersion: 0,
			uncompressed:  hexToBytes("6a200102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"),
			compressed:    hexToBytes("626a200102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"),
		},
		{
			name:          "requires 2 size bytes - data push 200 bytes",
			scriptVersion: 0,
			uncompressed:  append(hexToBytes("4cc8"), bytes.Repeat([]byte{0x00}, 200)...),
			// [0x80, 0x50] = 208 as a variable length quantity
			// [0x4c, 0xc8] = OP_PUSHDATA1 200
			compressed: append(hexToBytes("810a4cc8"), bytes.Repeat([]byte{0x00}, 200)...),
		},
	}

	for _, test := range tests {
		// Ensure the function to calculate the serialized size without
		// actually serializing the value is calculated properly.
		gotSize := compressedScriptSize(test.scriptVersion, test.uncompressed)
		if gotSize != len(test.compressed) {
			t.Errorf("compressedScriptSize (%s): did not get "+
				"expected size - got %d, want %d", test.name,
				gotSize, len(test.compressed))
			continue
		}

		// Ensure the script compresses to the expected bytes.
		gotCompressed := make([]byte, gotSize)
		gotBytesWritten := putCompressedScript(gotCompressed, test.scriptVersion,
			test.uncompressed)
		if !bytes.Equal(gotCompressed, test.compressed) {
			t.Errorf("putCompressedScript (%s): did not get "+
				"expected bytes - got %x, want %x", test.name,
				gotCompressed, test.compressed)
			continue
		}
		if gotBytesWritten != len(test.compressed) {
			t.Errorf("putCompressedScript (%s): did not get "+
				"expected number of bytes written - got %d, "+
				"want %d", test.name, gotBytesWritten,
				len(test.compressed))
			continue
		}

		// Ensure the compressed script size is properly decoded from
		// the compressed script.
		gotDecodedSize := decodeCompressedScriptSize(test.compressed)
		if gotDecodedSize != len(test.compressed) {
			t.Errorf("decodeCompressedScriptSize (%s): did not get "+
				"expected size - got %d, want %d", test.name,
				gotDecodedSize, len(test.compressed))
			continue
		}

		// Ensure the script decompresses to the expected bytes.
		gotDecompressed := decompressScript(test.compressed)
		if !bytes.Equal(gotDecompressed, test.uncompressed) {
			t.Errorf("decompressScript (%s): did not get expected "+
				"bytes - got %x, want %x", test.name,
				gotDecompressed, test.uncompressed)
			continue
		}
	}
}

// TestScriptCompressionErrors ensures calling various functions related to
// script compression with incorrect data returns the expected results.
func TestScriptCompressionErrors(t *testing.T) {
	t.Parallel()

	// A nil script must result in a decoded size of 0.
	if gotSize := decodeCompressedScriptSize(nil); gotSize != 0 {
		t.Fatalf("decodeCompressedScriptSize with nil script did not "+
			"return 0 - got %d", gotSize)
	}

	// A nil script must result in a nil decompressed script.
	if gotScript := decompressScript(nil); gotScript != nil {
		t.Fatalf("decompressScript with nil script did not return nil "+
			"decompressed script - got %x", gotScript)
	}

	// A compressed script for a pay-to-pubkey (uncompressed) that results
	// in an invalid pubkey must result in a nil decompressed script.
	compressedScript := hexToBytes("04012d74d0cb94344c9569c2e77901573d8d" +
		"7903c3ebec3a957724895dca52c6b4")
	if gotScript := decompressScript(compressedScript); gotScript != nil {
		t.Fatalf("decompressScript with compressed pay-to-"+
			"uncompressed-pubkey that is invalid did not return "+
			"nil decompressed script - got %x", gotScript)
	}
}

// TestAmountCompression ensures the domain-specific transaction output amount
// compression and decompression works as expected.
func TestAmountCompression(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		uncompressed uint64
		compressed   uint64
	}{
		{
			name:         "0 VGL (sometimes used in nulldata)",
			uncompressed: 0,
			compressed:   0,
		},
		{
			name:         "546 atoms (current network dust value)",
			uncompressed: 546,
			compressed:   4911,
		},
		{
			name:         "0.00001 VGL (typical transaction fee)",
			uncompressed: 1000,
			compressed:   4,
		},
		{
			name:         "0.0001 VGL (typical transaction fee)",
			uncompressed: 10000,
			compressed:   5,
		},
		{
			name:         "0.12345678 VGL",
			uncompressed: 12345678,
			compressed:   111111101,
		},
		{
			name:         "0.5 VGL",
			uncompressed: 50000000,
			compressed:   48,
		},
		{
			name:         "1 VGL",
			uncompressed: 100000000,
			compressed:   9,
		},
		{
			name:         "5 VGL",
			uncompressed: 500000000,
			compressed:   49,
		},
		{
			name:         "21000000 VGL (max minted coins)",
			uncompressed: 2100000000000000,
			compressed:   21000000,
		},
	}

	for _, test := range tests {
		// Ensure the amount compresses to the expected value.
		gotCompressed := compressTxOutAmount(test.uncompressed)
		if gotCompressed != test.compressed {
			t.Errorf("compressTxOutAmount (%s): did not get "+
				"expected value - got %d, want %d", test.name,
				gotCompressed, test.compressed)
			continue
		}

		// Ensure the value decompresses to the expected value.
		gotDecompressed := decompressTxOutAmount(test.compressed)
		if gotDecompressed != test.uncompressed {
			t.Errorf("decompressTxOutAmount (%s): did not get "+
				"expected value - got %d, want %d", test.name,
				gotDecompressed, test.uncompressed)
			continue
		}
	}
}

// TestCompressedTxOut ensures the transaction output serialization and
// deserialization works as expected.
func TestCompressedTxOut(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		amount        uint64
		scriptVersion uint16
		pkScript      []byte
		version       uint32
		compressed    []byte
		hasAmount     bool
	}{
		{
			name:          "nulldata with 0 VGL",
			amount:        0,
			scriptVersion: 0,
			pkScript:      hexToBytes("6a200102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"),
			version:       1,
			compressed:    hexToBytes("00626a200102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"),
			hasAmount:     false,
		},
		{
			name:          "pay-to-pubkey-hash dust, no amount",
			amount:        0,
			scriptVersion: 0,
			pkScript:      hexToBytes("76a9141018853670f9f3b0582c5b9ee8ce93764ac32b9388ac"),
			version:       1,
			compressed:    hexToBytes("00001018853670f9f3b0582c5b9ee8ce93764ac32b93"),
			hasAmount:     false,
		},
		{
			name:          "pay-to-pubkey-hash dust, amount",
			amount:        546,
			scriptVersion: 0,
			pkScript:      hexToBytes("76a9141018853670f9f3b0582c5b9ee8ce93764ac32b9388ac"),
			version:       1,
			compressed:    hexToBytes("a52f00001018853670f9f3b0582c5b9ee8ce93764ac32b93"),
			hasAmount:     true,
		},
		{
			name:          "pay-to-pubkey uncompressed, no amount",
			amount:        0,
			scriptVersion: 0,
			pkScript:      hexToBytes("4104192d74d0cb94344c9569c2e77901573d8d7903c3ebec3a957724895dca52c6b40d45264838c0bd96852662ce6a847b197376830160c6d2eb5e6a4c44d33f453eac"),
			version:       1,
			compressed:    hexToBytes("0004192d74d0cb94344c9569c2e77901573d8d7903c3ebec3a957724895dca52c6b4"),
			hasAmount:     false,
		},
		{
			name:          "pay-to-pubkey uncompressed 1 VGL, amount present",
			amount:        100000000,
			scriptVersion: 0,
			pkScript:      hexToBytes("4104192d74d0cb94344c9569c2e77901573d8d7903c3ebec3a957724895dca52c6b40d45264838c0bd96852662ce6a847b197376830160c6d2eb5e6a4c44d33f453eac"),
			version:       1,
			compressed:    hexToBytes("090004192d74d0cb94344c9569c2e77901573d8d7903c3ebec3a957724895dca52c6b4"),
			hasAmount:     true,
		},
	}

	for _, test := range tests {
		targetSz := compressedTxOutSize(0, test.scriptVersion, test.pkScript,
			test.hasAmount) - 1
		target := make([]byte, targetSz)
		putCompressedScript(target, test.scriptVersion, test.pkScript)

		// Ensure the function to calculate the serialized size without
		// actually serializing the txout is calculated properly.
		gotSize := compressedTxOutSize(test.amount, test.scriptVersion,
			test.pkScript, test.hasAmount)
		if gotSize != len(test.compressed) {
			t.Errorf("compressedTxOutSize (%s): did not get "+
				"expected size - got %d, want %d", test.name,
				gotSize, len(test.compressed))
			continue
		}

		// Ensure the txout compresses to the expected value.
		gotCompressed := make([]byte, gotSize)
		gotBytesWritten := putCompressedTxOut(gotCompressed,
			test.amount, test.scriptVersion, test.pkScript, test.hasAmount)
		if !bytes.Equal(gotCompressed, test.compressed) {
			t.Errorf("compressTxOut (%s): did not get expected "+
				"bytes - got %x, want %x", test.name,
				gotCompressed, test.compressed)
			continue
		}
		if gotBytesWritten != len(test.compressed) {
			t.Errorf("compressTxOut (%s): did not get expected "+
				"number of bytes written - got %d, want %d",
				test.name, gotBytesWritten,
				len(test.compressed))
			continue
		}

		// Ensure the serialized bytes are decoded back to the expected
		// compressed values.
		gotAmount, gotScrVersion, gotScript, gotBytesRead, err :=
			decodeCompressedTxOut(test.compressed, test.hasAmount)
		if err != nil {
			t.Errorf("decodeCompressedTxOut (%s): unexpected "+
				"error: %v", test.name, err)
			continue
		}
		if gotAmount != int64(test.amount) {
			t.Errorf("decodeCompressedTxOut (%s): did not get "+
				"expected amount - got %d, want %d",
				test.name, gotAmount, test.amount)
			continue
		}
		if gotScrVersion != test.scriptVersion {
			t.Errorf("decodeCompressedTxOut (%s): did not get "+
				"expected script version - got %d, want %d",
				test.name, gotScrVersion, test.scriptVersion)
			continue
		}
		if !bytes.Equal(gotScript, test.pkScript) {
			t.Errorf("decodeCompressedTxOut (%s): did not get "+
				"expected script - got %x, want %x",
				test.name, gotScript, test.pkScript)
			continue
		}
		if gotBytesRead != len(test.compressed) {
			t.Errorf("decodeCompressedTxOut (%s): did not get "+
				"expected number of bytes read - got %d, want %d",
				test.name, gotBytesRead, len(test.compressed))
			continue
		}
	}
}

// TestTxOutCompressionErrors ensures calling various functions related to
// txout compression with incorrect data returns the expected results.
func TestTxOutCompressionErrors(t *testing.T) {
	t.Parallel()

	// A compressed txout with a value and missing compressed script must error.
	compressedTxOut := hexToBytes("00")
	_, _, _, _, err := decodeCompressedTxOut(compressedTxOut, true)
	if !isDeserializeErr(err) {
		t.Fatalf("decodeCompressedTxOut with value and missing "+
			"compressed script did not return expected error type "+
			"- got %T, want errDeserialize", err)
	}

	// A compressed txout without a value and with an empty compressed
	// script returns empty but is valid.
	compressedTxOut = hexToBytes("00")
	_, _, _, _, err = decodeCompressedTxOut(compressedTxOut, false)
	if err != nil {
		t.Fatalf("decodeCompressedTxOut with missing compressed script "+
			"did not return expected error type - got %T, want "+
			"errDeserialize", err)
	}

	// A compressed txout with short compressed script must error.
	compressedTxOut = hexToBytes("0010")
	_, _, _, _, err = decodeCompressedTxOut(compressedTxOut, false)
	if !isDeserializeErr(err) {
		t.Fatalf("decodeCompressedTxOut with short compressed script "+
			"did not return expected error type - got %T, want "+
			"errDeserialize", err)
	}
}
