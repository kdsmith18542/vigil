// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package VGLjson

import (
	"errors"
	"reflect"
	"sort"
	"testing"
)

// Register methods for testing purposes.  This does not conflict with
// registration performed by external packages as they are done in separate
// builds.
func init() {
	MustRegister("getblock", (*testGetBlockCmd)(nil), 0)
	MustRegister("getblockcount", (*testGetBlockCountCmd)(nil), 0)
	MustRegister("session", (*testSessionCmd)(nil), UFWebsocketOnly)
	MustRegister("help", (*testHelpCmd)(nil), 0)
}

type testGetBlockCmd struct {
	Hash      string
	Verbose   *bool `jsonrpcdefault:"true"`
	VerboseTx *bool `jsonrpcdefault:"false"`
}
type testGetBlockCountCmd struct{}
type testSessionCmd struct{}
type testHelpCmd struct {
	Command *string
}

// TestUsageFlagStringer tests the stringized output for the UsageFlag type.
func TestUsageFlagStringer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   UsageFlag
		want string
	}{
		{0, "0x0"},
		{1, "0x0"}, // was UFWalletOnly
		{UFWebsocketOnly, "UFWebsocketOnly"},
		{UFNotification, "UFNotification"},
		{UFWebsocketOnly | UFNotification, "UFWebsocketOnly|UFNotification"},
		{1 | UFWebsocketOnly | UFNotification | (1 << 31),
			"UFWebsocketOnly|UFNotification|0x80000000"},
	}

	// Detect additional usage flags that don't have the stringer added.
	numUsageFlags := 0
	highestBit := highestUsageFlagBit
	for highestBit > 1 {
		numUsageFlags++
		highestBit >>= 1
	}
	if len(tests)-3 != numUsageFlags {
		t.Errorf("It appears a usage flag was added without adding " +
			"an associated stringer test")
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		result := test.in.String()
		if result != test.want {
			t.Errorf("String #%d\n got: %s want: %s", i, result,
				test.want)
			continue
		}
	}
}

// TestRegisterCmdErrors ensures the RegisterCmd function returns the expected
// error when provided with invalid types.
func TestRegisterCmdErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		method  string
		cmdFunc func() interface{}
		flags   UsageFlag
		err     error
	}{
		{
			name:   "duplicate method",
			method: "getblock",
			cmdFunc: func() interface{} {
				return struct{}{}
			},
			err: ErrDuplicateMethod,
		},
		{
			name:   "invalid usage flags",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				return 0
			},
			flags: highestUsageFlagBit,
			err:   ErrInvalidUsageFlags,
		},
		{
			name:   "invalid type",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				return 0
			},
			err: ErrInvalidType,
		},
		{
			name:   "invalid type 2",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				return &[]string{}
			},
			err: ErrInvalidType,
		},
		{
			name:   "embedded field",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				// nolint: unused
				type test struct{ int }
				return (*test)(nil)
			},
			err: ErrEmbeddedType,
		},
		{
			name:   "unexported field",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				// nolint: structcheck, unused
				type test struct{ a int }
				return (*test)(nil)
			},
			err: ErrUnexportedField,
		},
		{
			name:   "unsupported field type 1",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct{ A **int }
				return (*test)(nil)
			},
			err: ErrUnsupportedFieldType,
		},
		{
			name:   "unsupported field type 2",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct{ A chan int }
				return (*test)(nil)
			},
			err: ErrUnsupportedFieldType,
		},
		{
			name:   "unsupported field type 3",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct{ A complex64 }
				return (*test)(nil)
			},
			err: ErrUnsupportedFieldType,
		},
		{
			name:   "unsupported field type 4",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct{ A complex128 }
				return (*test)(nil)
			},
			err: ErrUnsupportedFieldType,
		},
		{
			name:   "unsupported field type 5",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct{ A func() }
				return (*test)(nil)
			},
			err: ErrUnsupportedFieldType,
		},
		{
			name:   "unsupported field type 6",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct{ A interface{} }
				return (*test)(nil)
			},
			err: ErrUnsupportedFieldType,
		},
		{
			name:   "required after optional",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct {
					A *int
					B int
				}
				return (*test)(nil)
			},
			err: ErrNonOptionalField,
		},
		{
			name:   "non-optional with default",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct {
					A int `jsonrpcdefault:"1"`
				}
				return (*test)(nil)
			},
			err: ErrNonOptionalDefault,
		},
		{
			name:   "mismatched default",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct {
					A *int `jsonrpcdefault:"1.7"`
				}
				return (*test)(nil)
			},
			err: ErrMismatchedDefault,
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		err := Register(test.method, test.cmdFunc(), test.flags)
		if !errors.Is(err, test.err) {
			t.Errorf("Test #%d (%s): mismatched error - got %v, "+
				"want %v", i, test.name, err, test.err)
			continue
		}
	}
}

// TestMustRegisterCmdPanic ensures the MustRegisterCmd function panics when
// used to register an invalid type.
func TestMustRegisterCmdPanic(t *testing.T) {
	t.Parallel()

	// Setup a defer to catch the expected panic to ensure it actually
	// paniced.
	defer func() {
		if err := recover(); err == nil {
			t.Error("MustRegisterCmd did not panic as expected")
		}
	}()

	// Intentionally try to register an invalid type to force a panic.
	MustRegister("panicme", 0, 0)
}

// TestRegisteredCmdMethods tests the RegisteredCmdMethods function ensure it
// works as expected.
func TestRegisteredCmdMethods(t *testing.T) {
	t.Parallel()

	// Ensure the registered methods for plain string methods are returned.
	methods := RegisteredMethods("")
	if len(methods) == 0 {
		t.Fatal("RegisteredCmdMethods: no methods")
	}

	// Ensure the returned methods are sorted.
	sortedMethods := make([]string, len(methods))
	copy(sortedMethods, methods)
	sort.Strings(sortedMethods)
	if !reflect.DeepEqual(sortedMethods, methods) {
		t.Fatal("RegisteredCmdMethods: methods are not sorted")
	}
}
