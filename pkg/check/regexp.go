// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"fmt"
	"regexp"

	"github.com/ctx42/xtst/pkg/notice"
)

// Regexp checks that "want" regexp matches "have". Returns nil if it does,
// otherwise, it returns an error with a message indicating the expected and
// actual values.
//
// The "want" can be either regular expression string or instance of
// [regexp.Regexp]. The [fmt.Sprint] is used to get string representation of
// have argument.
func Regexp(want, have any, opts ...Option) error {
	match, err := matchRegexp(want, have)
	if err != nil {
		return notice.New("expected valid regexp").Append("error", "%q", err)
	}
	if match {
		return nil
	}
	ops := DefaultOptions().set(opts)
	return notice.New("expected regexp to match").
		Trail(ops.Trail).
		Append("regexp", "%s", want).
		Have("%q", have)
}

// matchRegexp return true if a specified regexp matches a string.
func matchRegexp(rx, have any) (bool, error) {
	var r *regexp.Regexp
	if rr, ok := rx.(*regexp.Regexp); ok {
		r = rr
	} else {
		var err error
		rxs := fmt.Sprint(rx)
		if r, err = regexp.Compile(rxs); err != nil {
			return false, err
		}
	}
	return r.FindStringIndex(fmt.Sprint(have)) != nil, nil
}
