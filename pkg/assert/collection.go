// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"cmp"

	"github.com/ctx42/testing/internal/core"
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/tester"
)

// Len asserts "have" has "want" elements. Returns true if it is, otherwise it
// marks the test as failed, writes error message to test log and returns false.
func Len(t tester.T, want int, have any, opts ...check.Option) bool {
	t.Helper()
	if e := check.Len(want, have, opts...); e != nil {
		cnt, _ := core.Len(have)
		if want > cnt {
			t.Fatal(e)
		} else {
			t.Error(e)
		}
		return false
	}
	return true
}

// Has asserts slice has "want" value. Returns true if it does, otherwise marks
// the test as failed, writes error message to test log and returns false.
func Has[T comparable](t tester.T, want T, bag []T, opts ...check.Option) bool {
	t.Helper()
	if e := check.Has(want, bag, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// HasNo asserts slice does not have "want" value. Returns true if it does not,
// otherwise marks the test as failed, writes error message to test log and
// returns false.
func HasNo[T comparable](t tester.T, want T, bag []T, opts ...check.Option) bool {
	t.Helper()
	if e := check.HasNo(want, bag, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// HasKey asserts map has a key. Returns true if it does, otherwise marks the
// test as failed, writes error message to test log and returns false.
func HasKey[K comparable, V any](t tester.T, key K, set map[K]V, opts ...check.Option) (V, bool) {
	t.Helper()
	val, e := check.HasKey(key, set, opts...)
	if e != nil {
		t.Error(e)
		return val, false
	}
	return val, true
}

// HasNoKey asserts map has no key. Returns true if it doesn't, otherwise marks
// the test as failed, writes error message to test log and returns false.
func HasNoKey[K comparable, V any](t tester.T, key K, set map[K]V, opts ...check.Option) bool {
	t.Helper()
	if e := check.HasNoKey(key, set, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// HasKeyValue asserts map has a key with given value. Returns true if it
// doesn't, otherwise marks the test as failed, writes error message to test
// log and returns false.
func HasKeyValue[K, V comparable](
	t tester.T,
	key K,
	want V,
	set map[K]V,
	opts ...check.Option,
) bool {

	t.Helper()
	if e := check.HasKeyValue(key, want, set, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// SliceSubset checks the "want" is a subset "have". In other words all values
// in "want" slice must be in "have" slice. Returns nil if they are, otherwise
// returns an error with a message indicating the expected and actual values.
func SliceSubset[T comparable](t tester.T, want, have []T, opts ...check.Option) bool {
	t.Helper()
	if e := check.SliceSubset(want, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}

// MapSubset asserts the "want" is a subset "have". In other words all keys and
// their corresponding values in "want" map must be in "have" map. It is not an
// error when "have" map has some other keys. Returns true if "want is a subset
// of "have", otherwise marks the test as failed, writes error message to test
// log and returns false.
func MapSubset[K cmp.Ordered, V any](
	t tester.T,
	want, have map[K]V,
	opts ...check.Option,
) bool {

	t.Helper()
	if e := check.MapSubset(want, have, opts...); e != nil {
		t.Error(e)
		return false
	}
	return true
}
