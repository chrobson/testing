// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package assert

import (
	"testing"

	"github.com/ctx42/xtst/internal/affirm"
	"github.com/ctx42/xtst/pkg/check"
	"github.com/ctx42/xtst/pkg/tester"
)

func Test_ChannelWillClose(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t).Close()
		c := make(chan struct{})
		done := make(chan struct{})
		var have bool

		// --- When ---
		go func() { have = ChannelWillClose(tspy, "1s", c); close(done) }()

		// --- Then ---
		close(c)
		<-done
		affirm.True(t, have)
	})

	t.Run("error", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.IgnoreLogs()
		tspy.Close()

		c := make(chan struct{})
		defer close(c)

		// --- When ---
		have := ChannelWillClose(tspy, "1s", c)

		// --- Then ---
		tspy.Finish().AssertExpectations()
		affirm.False(t, have)
	})

	t.Run("log message with trail", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectFail()
		tspy.ExpectLogContain("\t trail: type.field")
		tspy.Close()

		c := make(chan struct{})
		defer close(c)
		opt := check.WithTrail("type.field")

		// --- When ---
		have := ChannelWillClose(tspy, "1s", c, opt)

		// --- Then ---
		tspy.Finish().AssertExpectations()
		affirm.False(t, have)
	})
}
