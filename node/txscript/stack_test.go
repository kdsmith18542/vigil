// Copyright (c) 2013-2017 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txscript

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

// TestStack tests that all of the stack operations work as expected.
func TestStack(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		before    [][]byte
		operation func(*stack) error
		err       error
		after     [][]byte
	}{
		{
			"noop",
			[][]byte{{1}, {2}, {3}, {4}, {5}},
			func(s *stack) error {
				return nil
			},
			nil,
			[][]byte{{1}, {2}, {3}, {4}, {5}},
		},
		{
			"peek underflow (byte)",
			[][]byte{{1}, {2}, {3}, {4}, {5}},
			func(s *stack) error {
				_, err := s.PeekByteArray(5)
				return err
			},
			ErrInvalidStackOperation,
			nil,
		},
		{
			"peek underflow (int)",
			[][]byte{{1}, {2}, {3}, {4}, {5}},
			func(s *stack) error {
				_, err := s.PeekInt(5, MathOpCodeMaxScriptNumLen)
				return err
			},
			ErrInvalidStackOperation,
			nil,
		},
		{
			"peek underflow (bool)",
			[][]byte{{1}, {2}, {3}, {4}, {5}},
			func(s *stack) error {
				_, err := s.PeekBool(5)
				return err
			},
			ErrInvalidStackOperation,
			nil,
		},
		{
			"pop",
			[][]byte{{1}, {2}, {3}, {4}, {5}},
			func(s *stack) error {
				val, err := s.PopByteArray()
				if err != nil {
					return err
				}
				if !bytes.Equal(val, []byte{5}) {
					return errors.New("not equal")
				}
				return err
			},
			nil,
			[][]byte{{1}, {2}, {3}, {4}},
		},
		{
			"pop everything",
			[][]byte{{1}, {2}, {3}, {4}, {5}},
			func(s *stack) error {
				for i := 0; i < 5; i++ {
					_, err := s.PopByteArray()
					if err != nil {
						return err
					}
				}
				return nil
			},
			nil,
			nil,
		},
		{
			"pop underflow",
			[][]byte{{1}, {2}, {3}, {4}, {5}},
			func(s *stack) error {
				for i := 0; i < 6; i++ {
					_, err := s.PopByteArray()
					if err != nil {
						return err
					}
				}
				return nil
			},
			ErrInvalidStackOperation,
			nil,
		},
		{
			"pop bool",
			[][]byte{nil},
			func(s *stack) error {
				val, err := s.PopBool()
				if err != nil {
					return err
				}

				if val {
					return errors.New("unexpected value")
				}
				return nil
			},
			nil,
			nil,
		},
		{
			"pop bool",
			[][]byte{{1}},
			func(s *stack) error {
				val, err := s.PopBool()
				if err != nil {
					return err
				}

				if !val {
					return errors.New("unexpected value")
				}
				return nil
			},
			nil,
			nil,
		},
		{
			"pop bool",
			nil,
			func(s *stack) error {
				_, err := s.PopBool()
				return err
			},
			ErrInvalidStackOperation,
			nil,
		},
		{
			"popInt 0",
			[][]byte{nil},
			func(s *stack) error {
				v, err := s.PopInt(MathOpCodeMaxScriptNumLen)
				if err != nil {
					return err
				}
				if v != 0 {
					return errors.New("0 != 0 on PopInt")
				}
				return nil
			},
			nil,
			nil,
		},
		{
			"non-minimal popInt 0",
			[][]byte{{0x0}},
			func(s *stack) error {
				_, err := s.PopInt(MathOpCodeMaxScriptNumLen)
				return err
			},
			ErrMinimalData,
			nil,
		},
		{
			"non-minimal popInt -0",
			[][]byte{{0x80}},
			func(s *stack) error {
				_, err := s.PopInt(MathOpCodeMaxScriptNumLen)
				return err
			},
			ErrMinimalData,
			nil,
		},
		{
			"popInt 1",
			[][]byte{{0x01}},
			func(s *stack) error {
				v, err := s.PopInt(MathOpCodeMaxScriptNumLen)
				if err != nil {
					return err
				}
				if v != 1 {
					return errors.New("1 != 1 on PopInt")
				}
				return nil
			},
			nil,
			nil,
		},
		{
			"non-minimal popInt 1 leading 0",
			[][]byte{{0x01, 0x00, 0x00, 0x00}},
			func(s *stack) error {
				_, err := s.PopInt(MathOpCodeMaxScriptNumLen)
				return err
			},
			ErrMinimalData,
			nil,
		},
		{
			"popInt -1",
			[][]byte{{0x81}},
			func(s *stack) error {
				v, err := s.PopInt(MathOpCodeMaxScriptNumLen)
				if err != nil {
					return err
				}
				if v != -1 {
					return errors.New("-1 != -1 on PopInt")
				}
				return nil
			},
			nil,
			nil,
		},
		{
			"popInt 5 byte",
			[][]byte{{0xff, 0xff, 0xff, 0xff, 0x7f}},
			func(s *stack) error {
				v, err := s.PopInt(5)
				if err != nil {
					return err
				}
				if v != 549755813887 {
					return fmt.Errorf("%v != 549755813887 on PopInt", v)
				}
				return nil
			},
			nil,
			nil,
		},
		{
			"non-minimal popInt 5-byte leading 0",
			[][]byte{{0xff, 0xff, 0xff, 0x7f, 0x00}},
			func(s *stack) error {
				_, err := s.PopInt(5)
				return err
			},
			ErrMinimalData,
			nil,
		},
		{
			"too big popInt 5-byte leading 0",
			[][]byte{{0xff, 0xff, 0xff, 0xff, 0x7f, 0x00}},
			func(s *stack) error {
				_, err := s.PopInt(5)
				return err
			},
			ErrNumOutOfRange,
			nil,
		},
		// Triggers the multibyte case in asInt
		{
			"popInt -513",
			[][]byte{{0x1, 0x82}},
			func(s *stack) error {
				v, err := s.PopInt(MathOpCodeMaxScriptNumLen)
				if err != nil {
					return err
				}
				if v != -513 {
					fmt.Printf("%v != %v\n", v, -513)
					return errors.New("1 != 1 on PopInt")
				}
				return nil
			},
			nil,
			nil,
		},
		// Confirm that the asInt code doesn't modify the base data.
		{
			"peekint nomodify -1",
			[][]byte{{0x81}},
			func(s *stack) error {
				v, err := s.PeekInt(0, MathOpCodeMaxScriptNumLen)
				if err != nil {
					return err
				}
				if v != -1 {
					fmt.Printf("%v != %v\n", v, -1)
					return errors.New("-1 != -1 on PeekInt")
				}
				return nil
			},
			nil,
			[][]byte{{0x81}},
		},
		{
			"PushInt 0",
			nil,
			func(s *stack) error {
				s.PushInt(ScriptNum(0))
				return nil
			},
			nil,
			[][]byte{{}},
		},
		{
			"PushInt 1",
			nil,
			func(s *stack) error {
				s.PushInt(ScriptNum(1))
				return nil
			},
			nil,
			[][]byte{{0x1}},
		},
		{
			"PushInt -1",
			nil,
			func(s *stack) error {
				s.PushInt(ScriptNum(-1))
				return nil
			},
			nil,
			[][]byte{{0x81}},
		},
		{
			"PushInt two bytes",
			nil,
			func(s *stack) error {
				s.PushInt(ScriptNum(256))
				return nil
			},
			nil,
			// little endian.. *sigh*
			[][]byte{{0x00, 0x01}},
		},
		{
			"PushInt leading zeros",
			nil,
			func(s *stack) error {
				// this will have the highbit set
				s.PushInt(ScriptNum(128))
				return nil
			},
			nil,
			[][]byte{{0x80, 0x00}},
		},
		{
			"dup",
			[][]byte{{1}},
			func(s *stack) error {
				return s.DupN(1)
			},
			nil,
			[][]byte{{1}, {1}},
		},
		{
			"dup2",
			[][]byte{{1}, {2}},
			func(s *stack) error {
				return s.DupN(2)
			},
			nil,
			[][]byte{{1}, {2}, {1}, {2}},
		},
		{
			"dup3",
			[][]byte{{1}, {2}, {3}},
			func(s *stack) error {
				return s.DupN(3)
			},
			nil,
			[][]byte{{1}, {2}, {3}, {1}, {2}, {3}},
		},
		{
			"dup0",
			[][]byte{{1}},
			func(s *stack) error {
				return s.DupN(0)
			},
			ErrInvalidStackOperation,
			nil,
		},
		{
			"dup-1",
			[][]byte{{1}},
			func(s *stack) error {
				return s.DupN(-1)
			},
			ErrInvalidStackOperation,
			nil,
		},
		{
			"dup too much",
			[][]byte{{1}},
			func(s *stack) error {
				return s.DupN(2)
			},
			ErrInvalidStackOperation,
			nil,
		},
		{
			"PushBool true",
			nil,
			func(s *stack) error {
				s.PushBool(true)

				return nil
			},
			nil,
			[][]byte{{1}},
		},
		{
			"PushBool false",
			nil,
			func(s *stack) error {
				s.PushBool(false)

				return nil
			},
			nil,
			[][]byte{nil},
		},
		{
			"PushBool PopBool",
			nil,
			func(s *stack) error {
				s.PushBool(true)
				val, err := s.PopBool()
				if err != nil {
					return err
				}
				if !val {
					return errors.New("unexpected value")
				}

				return nil
			},
			nil,
			nil,
		},
		{
			"PushBool PopBool 2",
			nil,
			func(s *stack) error {
				s.PushBool(false)
				val, err := s.PopBool()
				if err != nil {
					return err
				}
				if val {
					return errors.New("unexpected value")
				}

				return nil
			},
			nil,
			nil,
		},
		{
			"PushInt PopBool",
			nil,
			func(s *stack) error {
				s.PushInt(ScriptNum(1))
				val, err := s.PopBool()
				if err != nil {
					return err
				}
				if !val {
					return errors.New("unexpected value")
				}

				return nil
			},
			nil,
			nil,
		},
		{
			"PushInt PopBool 2",
			nil,
			func(s *stack) error {
				s.PushInt(ScriptNum(0))
				val, err := s.PopBool()
				if err != nil {
					return err
				}
				if val {
					return errors.New("unexpected value")
				}

				return nil
			},
			nil,
			nil,
		},
		{
			"Nip top",
			[][]byte{{1}, {2}, {3}},
			func(s *stack) error {
				return s.NipN(0)
			},
			nil,
			[][]byte{{1}, {2}},
		},
		{
			"Nip middle",
			[][]byte{{1}, {2}, {3}},
			func(s *stack) error {
				return s.NipN(1)
			},
			nil,
			[][]byte{{1}, {3}},
		},
		{
			"Nip low",
			[][]byte{{1}, {2}, {3}},
			func(s *stack) error {
				return s.NipN(2)
			},
			nil,
			[][]byte{{2}, {3}},
		},
		{
			"Nip too much",
			[][]byte{{1}, {2}, {3}},
			func(s *stack) error {
				// bite off more than we can chew
				return s.NipN(3)
			},
			ErrInvalidStackOperation,
			[][]byte{{2}, {3}},
		},
		{
			"keep on tucking",
			[][]byte{{1}, {2}, {3}},
			func(s *stack) error {
				return s.Tuck()
			},
			nil,
			[][]byte{{1}, {3}, {2}, {3}},
		},
		{
			"a little tucked up",
			[][]byte{{1}}, // too few arguments for tuck
			func(s *stack) error {
				return s.Tuck()
			},
			ErrInvalidStackOperation,
			nil,
		},
		{
			"all tucked up",
			nil, // too few arguments for tuck
			func(s *stack) error {
				return s.Tuck()
			},
			ErrInvalidStackOperation,
			nil,
		},
		{
			"drop 0",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) error {
				return s.DropN(0)
			},
			nil,
			[][]byte{{1}, {2}, {3}, {4}},
		},
		{
			"drop 1",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) error {
				return s.DropN(1)
			},
			nil,
			[][]byte{{1}, {2}, {3}},
		},
		{
			"drop 2",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) error {
				return s.DropN(2)
			},
			nil,
			[][]byte{{1}, {2}},
		},
		{
			"drop 3",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) error {
				return s.DropN(3)
			},
			nil,
			[][]byte{{1}},
		},
		{
			"drop 4",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) error {
				return s.DropN(4)
			},
			nil,
			nil,
		},
		{
			"drop 4/5",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) error {
				return s.DropN(5)
			},
			ErrInvalidStackOperation,
			nil,
		},
		{
			"Rot1",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) error {
				return s.RotN(1)
			},
			nil,
			[][]byte{{1}, {3}, {4}, {2}},
		},
		{
			"Rot2",
			[][]byte{{1}, {2}, {3}, {4}, {5}, {6}},
			func(s *stack) error {
				return s.RotN(2)
			},
			nil,
			[][]byte{{3}, {4}, {5}, {6}, {1}, {2}},
		},
		{
			"Rot too little",
			[][]byte{{1}, {2}},
			func(s *stack) error {
				return s.RotN(1)
			},
			ErrInvalidStackOperation,
			nil,
		},
		{
			"Rot0",
			[][]byte{{1}, {2}, {3}},
			func(s *stack) error {
				return s.RotN(0)
			},
			ErrInvalidStackOperation,
			nil,
		},
		{
			"Swap1",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) error {
				return s.SwapN(1)
			},
			nil,
			[][]byte{{1}, {2}, {4}, {3}},
		},
		{
			"Swap2",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) error {
				return s.SwapN(2)
			},
			nil,
			[][]byte{{3}, {4}, {1}, {2}},
		},
		{
			"Swap too little",
			[][]byte{{1}},
			func(s *stack) error {
				return s.SwapN(1)
			},
			ErrInvalidStackOperation,
			nil,
		},
		{
			"Swap0",
			[][]byte{{1}, {2}, {3}},
			func(s *stack) error {
				return s.SwapN(0)
			},
			ErrInvalidStackOperation,
			nil,
		},
		{
			"Over1",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) error {
				return s.OverN(1)
			},
			nil,
			[][]byte{{1}, {2}, {3}, {4}, {3}},
		},
		{
			"Over2",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) error {
				return s.OverN(2)
			},
			nil,
			[][]byte{{1}, {2}, {3}, {4}, {1}, {2}},
		},
		{
			"Over too little",
			[][]byte{{1}},
			func(s *stack) error {
				return s.OverN(1)
			},
			ErrInvalidStackOperation,
			nil,
		},
		{
			"Over0",
			[][]byte{{1}, {2}, {3}},
			func(s *stack) error {
				return s.OverN(0)
			},
			ErrInvalidStackOperation,
			nil,
		},
		{
			"Pick1",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) error {
				return s.PickN(1)
			},
			nil,
			[][]byte{{1}, {2}, {3}, {4}, {3}},
		},
		{
			"Pick2",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) error {
				return s.PickN(2)
			},
			nil,
			[][]byte{{1}, {2}, {3}, {4}, {2}},
		},
		{
			"Pick too little",
			[][]byte{{1}},
			func(s *stack) error {
				return s.PickN(1)
			},
			ErrInvalidStackOperation,
			nil,
		},
		{
			"Roll1",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) error {
				return s.RollN(1)
			},
			nil,
			[][]byte{{1}, {2}, {4}, {3}},
		},
		{
			"Roll2",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) error {
				return s.RollN(2)
			},
			nil,
			[][]byte{{1}, {3}, {4}, {2}},
		},
		{
			"Roll too little",
			[][]byte{{1}},
			func(s *stack) error {
				return s.RollN(1)
			},
			ErrInvalidStackOperation,
			nil,
		},
		{
			"Peek bool",
			[][]byte{{1}},
			func(s *stack) error {
				// Peek bool is otherwise pretty well tested,
				// just check it works.
				val, err := s.PeekBool(0)
				if err != nil {
					return err
				}
				if !val {
					return errors.New("invalid result")
				}
				return nil
			},
			nil,
			[][]byte{{1}},
		},
		{
			"Peek bool 2",
			[][]byte{nil},
			func(s *stack) error {
				// Peek bool is otherwise pretty well tested,
				// just check it works.
				val, err := s.PeekBool(0)
				if err != nil {
					return err
				}
				if val {
					return errors.New("invalid result")
				}
				return nil
			},
			nil,
			[][]byte{nil},
		},
		{
			"Peek int",
			[][]byte{{1}},
			func(s *stack) error {
				// Peek int is otherwise pretty well tested,
				// just check it works.
				val, err := s.PeekInt(0, MathOpCodeMaxScriptNumLen)
				if err != nil {
					return err
				}
				if val != 1 {
					return errors.New("invalid result")
				}
				return nil
			},
			nil,
			[][]byte{{1}},
		},
		{
			"Peek int 2",
			[][]byte{nil},
			func(s *stack) error {
				// Peek int is otherwise pretty well tested,
				// just check it works.
				val, err := s.PeekInt(0, MathOpCodeMaxScriptNumLen)
				if err != nil {
					return err
				}
				if val != 0 {
					return errors.New("invalid result")
				}
				return nil
			},
			nil,
			[][]byte{nil},
		},
		{
			"peekInt 5 byte",
			[][]byte{{0xff, 0xff, 0xff, 0xff, 0x7f}, nil},
			func(s *stack) error {
				v, err := s.PeekInt(1, 5)
				if err != nil {
					return err
				}
				if v != 549755813887 {
					return fmt.Errorf("%v != 549755813887 on PeekInt", v)
				}
				return nil
			},
			nil,
			[][]byte{{0xff, 0xff, 0xff, 0xff, 0x7f}, nil},
		},
		{
			"non-minimal peekInt 5-byte leading 0",
			[][]byte{{0xff, 0xff, 0xff, 0x7f, 0x00}, nil},
			func(s *stack) error {
				_, err := s.PeekInt(1, 5)
				return err
			},
			ErrMinimalData,
			nil,
		},
		{
			"too big peekInt 5-byte leading 0",
			[][]byte{{0xff, 0xff, 0xff, 0xff, 0x7f, 0x00}, nil},
			func(s *stack) error {
				_, err := s.PeekInt(1, 5)
				return err
			},
			ErrNumOutOfRange,
			nil,
		},
		{
			"pop int",
			nil,
			func(s *stack) error {
				s.PushInt(ScriptNum(1))
				val, err := s.PopInt(MathOpCodeMaxScriptNumLen)
				if err != nil {
					return err
				}
				if val != 1 {
					return errors.New("invalid result")
				}
				return nil
			},
			nil,
			nil,
		},
		{
			"pop empty",
			nil,
			func(s *stack) error {
				_, err := s.PopInt(MathOpCodeMaxScriptNumLen)
				return err
			},
			ErrInvalidStackOperation,
			nil,
		},
	}

	for _, test := range tests {
		// Setup the initial stack state and perform the test operation.
		s := stack{}
		for i := range test.before {
			s.PushByteArray(test.before[i])
		}
		err := test.operation(&s)

		// Ensure the error matches the value specified in the test instance.
		if !errors.Is(err, test.err) {
			t.Errorf("%s: unexpected error - got %v, want %v", test.name, err,
				test.err)
			continue
		}
		if err != nil {
			continue
		}

		// Ensure the resulting stack is the expected length.
		if int32(len(test.after)) != s.Depth() {
			t.Errorf("%s: stack depth doesn't match expected: %v "+
				"vs %v", test.name, len(test.after),
				s.Depth())
			continue
		}

		// Ensure all items of the resulting stack are the expected
		// values.
		for i := range test.after {
			val, err := s.PeekByteArray(s.Depth() - int32(i) - 1)
			if err != nil {
				t.Errorf("%s: can't peek %dth stack entry: %v",
					test.name, i, err)
				break
			}

			if !bytes.Equal(val, test.after[i]) {
				t.Errorf("%s: %dth stack entry doesn't match "+
					"expected: %v vs %v", test.name, i, val,
					test.after[i])
				break
			}
		}
	}
}
