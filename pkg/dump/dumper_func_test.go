// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ctx42/xtst/internal/affirm"
)

func Test_funcDumper_tabular(t *testing.T) {
	tt := []struct {
		testN string

		val  any
		want string
	}{
		{"func0", func() {}, "<func>(<addr>)"},
		{"func1", func(int) error { return nil }, "<func>(<addr>)"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			val := reflect.ValueOf(tc.val)

			// --- When ---
			have := funcDumper(Dump{}, 0, val)

			// --- Then ---
			affirm.Equal(t, tc.want, have)
		})
	}
}

func Test_funcDumper(t *testing.T) {
	t.Run("nil function", func(t *testing.T) {
		// --- Given ---
		dmp := New(NewConfig(WithPtrAddr))

		var fn func()
		val := reflect.ValueOf(fn)

		// --- When ---
		have := funcDumper(dmp, 0, val)

		// --- Then ---
		affirm.Equal(t, "<func>(<0x0>)", have)
	})

	t.Run("usage error", func(t *testing.T) {
		// --- Given ---
		dmp := New(NewConfig())
		val := reflect.ValueOf(1234)

		// --- When ---
		have := funcDumper(dmp, 0, val)

		// --- Then ---
		affirm.Equal(t, valErrUsage, have)
	})

	t.Run("print pointer address", func(t *testing.T) {
		// --- Given ---
		dmp := New(NewConfig(WithPtrAddr))
		fn := func() {}
		val := reflect.ValueOf(fn)
		want := fmt.Sprintf("<func>(<0x%x>)", val.Pointer())

		// --- When ---
		have := funcDumper(dmp, 0, val)

		// --- Then ---
		affirm.Equal(t, want, have)
	})

	t.Run("uses indent and level", func(t *testing.T) {
		// --- Given ---
		dmp := New(NewConfig(WithIndent(2)))
		val := reflect.ValueOf(1234)

		// --- When ---
		have := funcDumper(dmp, 1, val)

		// --- Then ---
		affirm.Equal(t, "\t\t\t"+valErrUsage, have)
	})
}
