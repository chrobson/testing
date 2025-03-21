// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"github.com/ctx42/xtst/pkg/notice"
)

// True checks "have" is true. Returns nil if it's, otherwise it returns an
// error with a message indicating the expected and actual values.
func True(have bool, opts ...Option) error {
	if !have {
		ops := DefaultOptions().set(opts)
		return notice.New("expected value to be true").
			Trail(ops.Trail)
	}
	return nil
}

// False checks "have" is false. Returns nil if it's, otherwise it returns an
// error with a message indicating the expected and actual values.
func False(have bool, opts ...Option) error {
	if have {
		ops := DefaultOptions().set(opts)
		return notice.New("expected value to be false").
			Trail(ops.Trail)
	}
	return nil
}
