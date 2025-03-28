// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"testing"

	"github.com/ctx42/xtst/internal/affirm"
)

func Test_Len(t *testing.T) {
	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		opt := WithTrail("type.field")

		// --- When ---
		err := Len(1, []int{1, 2}, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected []int length:\n" +
			"  trail: type.field\n" +
			"   want: 1\n" +
			"   have: 2"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_Len_success_tabular(t *testing.T) {
	ch := make(chan int, 4)
	ch <- 0
	ch <- 1
	ch <- 2
	t.Cleanup(func() { <-ch; <-ch; <-ch; close(ch) })

	tt := []struct {
		testN string

		val  any
		want int
	}{
		{"int empty success", []int{}, 0},
		{"int success", []int{1}, 1},
		{"map empty success", map[string]int{}, 0},
		{"map success", map[string]int{"A": 1}, 1},
		{"channel success", ch, 3},
		{"string success", "abc", 3},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			err := Len(tc.want, tc.val)

			// --- Then ---
			affirm.Nil(t, err)
		})
	}
}

func Test_Len_error_tabular(t *testing.T) {
	ch := make(chan int, 4)
	ch <- 0
	ch <- 1
	ch <- 2
	t.Cleanup(func() { <-ch; <-ch; <-ch; close(ch) })

	tt := []struct {
		testN string

		val  any
		want int
		wMsg string
	}{
		{"int empty fail", []int{}, 1, "expected []int length:\n  want: 1\n  have: 0"},
		{"int fail", []int{1}, 2, "expected []int length:\n  want: 2\n  have: 1"},
		{"map empty fail", map[string]int{}, 1, "expected map[string]int length:\n  want: 1\n  have: 0"},
		{"map fail", map[string]int{"A": 1}, 2, "expected map[string]int length:\n  want: 2\n  have: 1"},
		{"invalid type", 1, 0, "cannot execute len(int)"},
		{"chan fail", ch, 4, "expected chan int length:\n  want: 4\n  have: 3"},
		{"string fail", "abc", 4, "expected string length:\n  want: 4\n  have: 3"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			err := Len(tc.want, tc.val)

			// --- Then ---
			affirm.NotNil(t, err)
			affirm.Equal(t, tc.wMsg, err.Error())
		})
	}
}

func Test_Has(t *testing.T) {
	t.Run("has", func(t *testing.T) {
		// --- Given ---
		val := []int{1, 2, 3}

		// --- When ---
		err := Has(2, val)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		val := []int{1, 2, 3}

		// --- When ---
		err := Has(42, val)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected slice to have a value:\n" +
			"   want: 42\n" +
			"  slice: \n" +
			"         []int{\n" +
			"           1,\n" +
			"           2,\n" +
			"           3,\n" +
			"         }"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		val := []int{1, 2, 3}
		opt := WithTrail("type.field")

		// --- When ---
		err := Has(42, val, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected slice to have a value:\n" +
			"  trail: type.field\n" +
			"   want: 42\n" +
			"  slice: \n" +
			"         []int{\n" +
			"           1,\n" +
			"           2,\n" +
			"           3,\n" +
			"         }"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_HasNo(t *testing.T) {
	t.Run("doesnt have", func(t *testing.T) {
		// --- Given ---
		val := []int{1, 2, 3}

		// --- When ---
		err := HasNo(4, val)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		val := []int{1, 2, 3}

		// --- When ---
		err := HasNo(2, val)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected slice not to have value:\n" +
			"   want: 2\n" +
			"  index: 1\n" +
			"  slice: \n" +
			"         []int{\n" +
			"           1,\n" +
			"           2,\n" +
			"           3,\n" +
			"         }"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		val := []int{1, 2, 3}
		opt := WithTrail("type.field")

		// --- When ---
		err := HasNo(2, val, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected slice not to have value:\n" +
			"  trail: type.field\n" +
			"   want: 2\n" +
			"  index: 1\n" +
			"  slice: \n" +
			"         []int{\n" +
			"           1,\n" +
			"           2,\n" +
			"           3,\n" +
			"         }"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("nil slice", func(t *testing.T) {
		// --- Given ---
		var val []any

		// --- When ---
		err := HasNo(2, val)

		// --- Then ---
		affirm.Nil(t, err)
	})
}

func Test_HasKey(t *testing.T) {
	t.Run("has", func(t *testing.T) {
		// --- Given ---
		val := map[string]int{"A": 1, "B": 2, "C": 3}

		// --- When ---
		haveValue, err := HasKey("B", val)

		// --- Then ---
		affirm.Equal(t, 2, haveValue)
		affirm.Nil(t, err)
	})

	t.Run("has not", func(t *testing.T) {
		// --- Given ---
		val := map[string]int{"A": 1, "B": 2, "C": 3}

		// --- When ---
		haveValue, err := HasKey("X", val)

		// --- Then ---
		affirm.Equal(t, 0, haveValue)
		affirm.NotNil(t, err)
		wMsg := "expected map to have a key:\n" +
			"  key: \"X\"\n" +
			"  map: \n" +
			"       map[string]int{\n" +
			"         \"A\": 1,\n" +
			"         \"B\": 2,\n" +
			"         \"C\": 3,\n" +
			"       }"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("nil map", func(t *testing.T) {
		// --- Given ---
		var m map[string]any

		// --- When ---
		haveValue, err := HasKey("X", m)

		// --- Then ---
		affirm.Nil(t, haveValue)
		affirm.NotNil(t, err)
		wMsg := "expected map to have a key:\n" +
			"  key: \"X\"\n" +
			"  map: map[string]any(nil)"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		var m map[string]any
		opt := WithTrail("type.field")

		// --- When ---
		haveValue, err := HasKey("X", m, opt)

		// --- Then ---
		affirm.Nil(t, haveValue)
		affirm.NotNil(t, err)
		wMsg := "expected map to have a key:\n" +
			"  trail: type.field\n" +
			"    key: \"X\"\n" +
			"    map: map[string]any(nil)"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_HasNoKey(t *testing.T) {
	t.Run("has not", func(t *testing.T) {
		// --- Given ---
		val := map[string]int{"A": 1, "B": 2, "C": 3}

		// --- When ---
		err := HasNoKey("D", val)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		val := map[string]int{"A": 1, "B": 2, "C": 3}

		// --- When ---
		err := HasNoKey("B", val)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected map not to have a key:\n" +
			"    key: \"B\"\n" +
			"  value: 2\n" +
			"    map: \n" +
			"         map[string]int{\n" +
			"           \"A\": 1,\n" +
			"           \"B\": 2,\n" +
			"           \"C\": 3,\n" +
			"         }"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		val := map[string]int{"A": 1, "B": 2, "C": 3}
		opt := WithTrail("type.field")

		// --- When ---
		err := HasNoKey("B", val, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected map not to have a key:\n" +
			"  trail: type.field\n" +
			"    key: \"B\"\n" +
			"  value: 2\n" +
			"    map: \n" +
			"         map[string]int{\n" +
			"           \"A\": 1,\n" +
			"           \"B\": 2,\n" +
			"           \"C\": 3,\n" +
			"         }"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_HasKeyValue(t *testing.T) {
	t.Run("has key and value matches", func(t *testing.T) {
		// --- Given ---
		val := map[string]int{"A": 1, "B": 2, "C": 3}

		// --- When ---
		err := HasKeyValue("B", 2, val)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("has key but value does not match", func(t *testing.T) {
		// --- Given ---
		val := map[string]int{"A": 1, "B": 2, "C": 3}
		opt := WithTrail("type.field")

		// --- When ---
		err := HasKeyValue("B", 100, val, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected map to have a key with a value:\n" +
			"  trail: type.field\n" +
			"    key: \"B\"\n" +
			"   want: 100\n" +
			"   have: 2"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("has no key", func(t *testing.T) {
		// --- Given ---
		val := map[string]int{"A": 1, "B": 2, "C": 3}
		opt := WithTrail("type.field")

		// --- When ---
		err := HasKeyValue("X", 2, val, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected map to have a key:\n  trail: type.field\n" +
			"    key: \"X\"\n" +
			"    map: \n" +
			"         map[string]int{\n" +
			"           \"A\": 1,\n" +
			"           \"B\": 2,\n" +
			"           \"C\": 3,\n" +
			"         }"
		affirm.Equal(t, wMsg, err.Error())
	})
}

func Test_SliceSubset(t *testing.T) {
	t.Run("slices equal", func(t *testing.T) {
		// --- Given ---
		want := []string{"A", "B", "C"}
		have := []string{"C", "B", "A"}

		// --- When ---
		err := SliceSubset(want, have)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("want slice is subset of have slice", func(t *testing.T) {
		// --- Given ---
		want := []string{"A", "B", "C"}
		have := []string{"D", "C", "B", "A"}

		// --- When ---
		err := SliceSubset(want, have)

		// --- Then ---
		affirm.Nil(t, err)
	})

	t.Run("want slice is not a subset of have slice", func(t *testing.T) {
		// --- Given ---
		want := []int{9, 8, 0, 1, 2}
		have := []int{2, 1, 0}

		// --- When ---
		err := SliceSubset(want, have)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected \"want\" slice to be a subset of \"have\" slice:\n" +
			"  missing values: \n" +
			"                  []int{\n" +
			"                    9,\n" +
			"                    8,\n" +
			"                  }"
		affirm.Equal(t, wMsg, err.Error())
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		want := []int{9, 9, 0, 1, 2}
		have := []int{2, 1, 0}
		opt := WithTrail("type.field")

		// --- When ---
		err := SliceSubset(want, have, opt)

		// --- Then ---
		affirm.NotNil(t, err)
		wMsg := "expected \"want\" slice to be a subset of \"have\" slice:\n" +
			"           trail: type.field\n" +
			"  missing values: \n" +
			"                  []int{\n" +
			"                    9,\n" +
			"                    9,\n" +
			"                  }"
		affirm.Equal(t, wMsg, err.Error())
	})
}
