// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"testing"

	"github.com/ctx42/xtst/internal/affirm"
	"github.com/ctx42/xtst/pkg/check"
	"github.com/ctx42/xtst/pkg/tester"
)

func Test_Panic(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := Panic(tspy, func() { panic("test") })

		// --- Then ---
		affirm.True(t, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		// --- When ---
		have := Panic(tspy, func() {})

		// --- Then ---
		affirm.False(t, have)
	})

	t.Run("error with option", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("\ttrail: type.field")
		tspy.Close()

		opt := check.WithTrail("type.field")

		// --- When ---
		have := Panic(tspy, func() {}, opt)

		// --- Then ---
		affirm.False(t, have)
	})
}

func Test_NoPanic(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := NoPanic(tspy, func() {})

		// --- Then ---
		affirm.True(t, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		// --- When ---
		have := NoPanic(tspy, func() { panic("test") })

		// --- Then ---
		affirm.False(t, have)
	})

	t.Run("error with option", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("\t      trail: type.field\n")
		tspy.Close()

		opt := check.WithTrail("type.field")

		// --- When ---
		have := NoPanic(tspy, func() { panic("test") }, opt)

		// --- Then ---
		affirm.False(t, have)
	})
}

func Test_PanicContain(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := PanicContain(tspy, "def", func() { panic("abc def ghi") })

		// --- Then ---
		affirm.True(t, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		// --- When ---
		have := PanicContain(tspy, "xyz", func() { panic("abc def ghi") })

		// --- Then ---
		affirm.False(t, have)
	})

	t.Run("error with option", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("\t      trail: type.field\n")
		tspy.Close()

		opt := check.WithTrail("type.field")

		// --- When ---
		have := PanicContain(tspy, "xyz", func() { panic("abc def ghi") }, opt)

		// --- Then ---
		affirm.False(t, have)
	})
}

func Test_PanicMsg(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		msg := PanicMsg(tspy, func() { panic("abc def ghi") })

		// --- Then ---
		if msg == nil {
			t.Error("expected PanicMsg to return non-nil value")
			return
		}
		affirm.Equal(t, "abc def ghi", *msg)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		// --- When ---
		msg := PanicMsg(tspy, func() {})

		// --- Then ---
		if msg != nil {
			t.Error("expected PanicMsg to return nil value")
		}
	})

	t.Run("error with option", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("\ttrail: type.field")
		tspy.Close()

		opt := check.WithTrail("type.field")

		// --- When ---
		msg := PanicMsg(tspy, func() {}, opt)

		// --- Then ---
		if msg != nil {
			t.Error("expected PanicMsg to return nil value")
		}
	})
}
