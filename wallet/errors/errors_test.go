// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package errors

import (
	std "errors"
	"testing"
)

func depth(err error) int {
	if err == nil {
		return 0
	}
	e, ok := err.(*Error)
	if !ok {
		return 1
	}
	return 1 + depth(e.Err)
}

func eq(e0, e1 *Error) bool {
	if e0.Op != e1.Op {
		return false
	}
	if e0.Kind != e1.Kind {
		return false
	}
	if e0.Err != e1.Err {
		return false
	}
	return true
}

func TestCollapse(t *testing.T) {
	e0 := E(Op("abc"))
	e0 = E(e0, Passphrase)
	if depth(e0) != 1 {
		t.Fatal("e0 was not collapsed")
	}

	e1 := E(Op("abc"), Passphrase)
	if !eq(e0.(*Error), e1.(*Error)) {
		t.Fatal("e0 was not collapsed to e1")
	}
}

func TestIs(t *testing.T) {
	base := std.New("base error")
	e := E(base, Op("operation"), Permission)
	if !Is(e, base) {
		t.Fatal("no match on base errors")
	}
	if Is(e, E("different error")) {
		t.Fatal("match on different error strings")
	}
	if !Is(e, E(Op("operation"))) {
		t.Fatal("no match on operation")
	}
	if Is(e, E(Op("different operation"))) {
		t.Fatal("match on different operation")
	}
	if !Is(e, E(Permission)) {
		t.Fatal("no match on kind")
	}
	if Is(e, E(Invalid)) {
		t.Fatal("match on different kind")
	}
}

func TestDoubleWrappedErrorWithKind(t *testing.T) {
	err := E(Invalid, "abc")
	// Wrap the error again
	err = E(err, "def")
	// Now try to match against the kind
	if !Is(err, Invalid) {
		t.Errorf("Is returned false for error object: %T %+[1]v", err)
		t.Errorf("Wrapped error: %T %+[1]v", err.(*Error).Unwrap())
	}
}
