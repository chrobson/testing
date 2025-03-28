// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"testing"

	"github.com/ctx42/xtst/internal/affirm"
	"github.com/ctx42/xtst/pkg/check"
	"github.com/ctx42/xtst/pkg/tester"
)

func Test_Len(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		// --- When ---
		have := Len(tspy, 2, []int{0, 1})

		// --- Then ---
		affirm.True(t, have)
	})

	t.Run("fatal when want is less than actual length", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFatal()
		tspy.IgnoreLogs()
		tspy.Close()

		defer func() { _ = recover() }()

		// --- When ---
		have := Len(tspy, 3, []int{0, 1})

		// --- Then ---
		affirm.False(t, have)
	})

	t.Run("error when want is greater than actual length", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		// --- When ---
		have := Len(tspy, 1, []int{0, 1})

		// --- Then ---
		affirm.False(t, have)
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		opt := check.WithTrail("type.field")

		// --- When ---
		have := Len(tspy, 1, []int{0, 1}, opt)

		// --- Then ---
		affirm.False(t, have)
	})
}

func Test_Has(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		val := []int{1, 2, 3}

		// --- When ---
		have := Has(tspy, 2, val)

		// --- Then ---
		affirm.True(t, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		val := []int{1, 2, 3}

		// --- When ---
		have := Has(tspy, 42, val)

		// --- Then ---
		affirm.False(t, have)
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		val := []int{1, 2, 3}
		opt := check.WithTrail("type.field")

		// --- When ---
		have := Has(tspy, 42, val, opt)

		// --- Then ---
		affirm.False(t, have)
	})
}

func Test_HasNo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()
		val := []int{1, 2, 3}

		// --- When ---
		have := HasNo(tspy, 4, val)

		// --- Then ---
		affirm.True(t, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		val := []int{1, 2, 3}

		// --- When ---
		have := HasNo(tspy, 2, val)

		// --- Then ---
		affirm.False(t, have)
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		val := []int{1, 2, 3}
		opt := check.WithTrail("type.field")

		// --- When ---
		have := HasNo(tspy, 2, val, opt)

		// --- Then ---
		affirm.False(t, have)
	})
}

func Test_HasKey(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		val := map[string]int{"A": 1, "B": 2, "C": 3}

		// --- When ---
		haveValue, success := HasKey(tspy, "B", val)

		// --- Then ---
		affirm.Equal(t, 2, haveValue)
		affirm.True(t, success)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		val := map[string]int{"A": 1, "B": 2, "C": 3}

		// --- When ---
		haveValue, success := HasKey(tspy, "X", val)

		// --- Then ---
		affirm.Equal(t, 0, haveValue)
		affirm.False(t, success)
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		val := map[string]int{"A": 1, "B": 2, "C": 3}
		opt := check.WithTrail("type.field")

		// --- When ---
		haveValue, success := HasKey(tspy, "X", val, opt)

		// --- Then ---
		affirm.Equal(t, 0, haveValue)
		affirm.False(t, success)
	})
}

func Test_HasNoKey(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()

		val := map[string]int{"A": 1, "B": 2, "C": 3}

		// --- When ---
		success := HasNoKey(tspy, "D", val)

		// --- Then ---
		affirm.True(t, success)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		val := map[string]int{"A": 1, "B": 2, "C": 3}

		// --- When ---
		success := HasNoKey(tspy, "B", val)

		// --- Then ---
		affirm.False(t, success)
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		val := map[string]int{"A": 1, "B": 2, "C": 3}
		opt := check.WithTrail("type.field")

		// --- When ---
		success := HasNoKey(tspy, "B", val, opt)

		// --- Then ---
		affirm.False(t, success)
	})
}

func Test_HasKeyValue(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()
		val := map[string]int{"A": 1, "B": 2, "C": 3}

		// --- When ---
		have := HasKeyValue(tspy, "B", 2, val)

		// --- Then ---
		affirm.True(t, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		val := map[string]int{"A": 1, "B": 2, "C": 3}

		// --- When ---
		have := HasKeyValue(tspy, "B", 100, val)

		// --- Then ---
		affirm.False(t, have)
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("  trail: type.field\n")
		tspy.Close()

		val := map[string]int{"A": 1, "B": 2, "C": 3}
		opt := check.WithTrail("type.field")

		// --- When ---
		have := HasKeyValue(tspy, "B", 100, val, opt)

		// --- Then ---
		affirm.False(t, have)
	})
}

func Test_SliceSubset(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.Close()

		sWant := []string{"A", "B", "C"}
		sHave := []string{"C", "B", "A"}

		// --- When ---
		have := SliceSubset(tspy, sWant, sHave)

		// --- Then ---
		affirm.True(t, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.IgnoreLogs()
		tspy.Close()

		sWant := []string{"X", "Y", "A", "B", "C"}
		sHave := []string{"C", "B", "A"}

		// --- When ---
		have := SliceSubset(tspy, sWant, sHave)

		// --- Then ---
		affirm.False(t, have)
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectError()
		tspy.ExpectLogContain("           trail: type.field\n")
		tspy.Close()

		sWant := []string{"X", "Y", "A", "B", "C"}
		sHave := []string{"C", "B", "A"}
		opt := check.WithTrail("type.field")

		// --- When ---
		have := SliceSubset(tspy, sWant, sHave, opt)

		// --- Then ---
		affirm.False(t, have)
	})
}
