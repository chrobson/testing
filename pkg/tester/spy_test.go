// Copyright (c) 2025 Rafal Zajac
// SPDX-License-Identifier: MIT

package tester

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/ctx42/xtst/internal/affirm"
)

func Test_find_match(t *testing.T) {
	t.Run("equal success", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: Equal, want: "abc def ghi"}

		// --- When ---
		have := ent.match("abc def ghi")

		// --- Then ---
		affirm.True(t, have)
	})

	t.Run("error equal", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: Equal, want: "abc def ghi"}

		// --- When ---
		have := ent.match("def")

		// --- Then ---
		affirm.False(t, have)
	})

	t.Run("contains success", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: Contains, want: "def"}

		// --- When ---
		have := ent.match("abc def ghi")

		// --- Then ---
		affirm.True(t, have)
	})

	t.Run("error contains", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: Contains, want: "def"}

		// --- When ---
		have := ent.match("abc xyz ghi")

		// --- Then ---
		affirm.False(t, have)
	})

	t.Run("not contains success", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: NotContains, want: "jkl"}

		// --- When ---
		have := ent.match("abc def ghi")

		// --- Then ---
		affirm.True(t, have)
	})

	t.Run("error not contains", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: NotContains, want: "def"}

		// --- When ---
		have := ent.match("abc def ghi")

		// --- Then ---
		affirm.False(t, have)
	})

	t.Run("regexp success", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: Regexp, want: "[a-z ]+"}

		// --- When ---
		have := ent.match("abc def ghi")

		// --- Then ---
		affirm.True(t, have)
	})

	t.Run("error regexp", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: Regexp, want: "abc def ghi"}

		// --- When ---
		have := ent.match("xyz")

		// --- Then ---
		affirm.False(t, have)
	})

	t.Run("unknown as equal success", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: "unknown", want: "abc def ghi"}

		// --- When ---
		have := ent.match("abc def ghi")

		// --- Then ---
		affirm.True(t, have)
	})

	t.Run("error unknown as equal", func(t *testing.T) {
		// --- Given ---
		ent := &find{strategy: "unknown", want: "abc def ghi"}

		// --- When ---
		have := ent.match("def")

		// --- Then ---
		affirm.False(t, have)
	})
}

func Test_New(t *testing.T) {
	t.Run("initial values", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		// --- When ---
		spy := New(ti)

		// --- Then ---
		affirm.False(t, spy.helperCntSet)
		affirm.Equal(t, -1, spy.wantHelperCnt)
		affirm.Equal(t, 0, spy.haveHelperCnt)
		affirm.Equal(t, 0, spy.wantNamesCnt)
		affirm.Equal(t, 0, spy.haveNamesCnt)
		affirm.Equal(t, 0, spy.wantTempDirCnt)
		affirm.Equal(t, 0, len(spy.haveTempDirs))
		affirm.Equal(t, 0, len(spy.wantEnv))
		affirm.Equal(t, 0, len(spy.haveEnv))
		affirm.False(t, spy.closed)
		affirm.False(t, spy.finished)
		affirm.False(t, spy.wantFailed)
		affirm.False(t, spy.wantError)
		affirm.False(t, spy.haveError)
		affirm.False(t, spy.wantFatal)
		affirm.False(t, spy.haveFatal)
		affirm.False(t, spy.wantSkipped)
		affirm.False(t, spy.haveSkipped)
		affirm.False(t, spy.panicked)
		affirm.Equal(t, 0, spy.hadCleanupsCnt)
		affirm.Equal(t, 0, len(spy.haveCleanups))
		affirm.Equal(t, 0, len(spy.savedMgs))
		affirm.Equal(t, 0, len(spy.haveLogMgs))
		affirm.False(t, spy.ignoreLog)
		affirm.False(t, spy.runningFinish)
	})

	t.Run("set expected Helper calls", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		// --- When ---
		spy := New(ti, 1)

		// --- Then ---
		affirm.True(t, spy.helperCntSet)
		affirm.Equal(t, 1, spy.wantHelperCnt)
		affirm.Equal(t, 0, spy.haveHelperCnt)
	})

	t.Run("set invalid expected Helper calls panics", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		// --- Then ---
		want := "ExpectHelpers cnt must be greater or equal to minus one"
		affirm.Panic(t, want, func() { New(ti, -2) })
	})
}

func Test_New_finish_and_assert_called_automatically_at_test_end(t *testing.T) {
	// --- Given ---
	var runs []int
	var spy *Spy

	// --- When ---
	t.Run("inner test", func(t *testing.T) {
		// --- Given ---
		spy = New(t, 0)
		spy.ExpectCleanups(2)
		spy.Close()

		// --- When ---
		spy.Cleanup(func() { runs = append(runs, 0) })
		spy.Cleanup(func() { runs = append(runs, 1) })
	})

	// --- Then ---
	affirm.Equal(t, 2, len(runs))
	affirm.Equal(t, 1, runs[0])
	affirm.Equal(t, 0, runs[1])
	affirm.True(t, spy.tt == nil)
}

func Test_Spy_ExpectCleanups(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectCleanups(2)

		// --- Then ---
		affirm.Equal(t, 2, spy.wantCleanupsCnt)
		affirm.True(t, spy == have)
	})

	t.Run("panics when called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- Then ---
		affirm.Panic(t, errExpectOnClosed, func() { spy.ExpectCleanups(1) })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errExpectOnFinished, func() { spy.ExpectCleanups(1) })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_Cleanup(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		var runs []int
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		spy.Cleanup(func() { runs = append(runs, 0) })
		spy.Cleanup(func() { runs = append(runs, 1) })

		// --- Then ---
		affirm.Equal(t, 2, spy.hadCleanupsCnt)
		affirm.Equal(t, 2, len(spy.haveCleanups))
		affirm.Equal(t, 0, len(runs))

		spy.haveCleanups[0]()
		affirm.Equal(t, 1, len(runs))
		affirm.Equal(t, 0, runs[0])

		spy.haveCleanups[1]()
		affirm.Equal(t, 2, len(runs))
		affirm.Equal(t, 0, runs[0])
		affirm.Equal(t, 1, runs[1])
	})

	t.Run("panics when called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- Then ---
		affirm.Panic(t, errMockOnNotClosed, func() { spy.Cleanup(func() {}) })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errActionOnFinished, func() { spy.Cleanup(func() {}) })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_ExpectError(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectError()

		// --- Then ---
		affirm.True(t, spy.wantError)
		affirm.True(t, spy == have)
	})

	t.Run("panics when called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- Then ---
		affirm.Panic(t, errExpectOnClosed, func() { spy.ExpectError() })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errExpectOnFinished, func() { spy.ExpectError() })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called after ExpectFail", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectFail()

		// --- Then ---
		want := "cannot use ExpectError and ExpectFail in the same time"
		affirm.Panic(t, want, func() { spy.ExpectError() })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_Error(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		spy.Error("msg", 0)

		// --- Then ---
		affirm.Equal(t, 1, len(spy.haveLogMgs))
		affirm.Equal(t, "msg 0", spy.haveLogMgs[0])
		affirm.True(t, spy.haveError)
	})

	t.Run("multi call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		spy.Error("msg", 0)
		spy.Error("msg", 1)

		// --- Then ---
		affirm.Equal(t, 2, len(spy.haveLogMgs))
		affirm.Equal(t, "msg 0", spy.haveLogMgs[0])
		affirm.Equal(t, "msg 1", spy.haveLogMgs[1])
		affirm.True(t, spy.haveError)
	})

	t.Run("panics when called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- Then ---
		affirm.Panic(t, errMockOnNotClosed, func() { spy.Error("msg", 0) })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errActionOnFinished, func() { spy.Error("msg", 0) })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_Errorf(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		spy.Errorf("msg %d", 0)

		// --- Then ---
		affirm.Equal(t, 1, len(spy.haveLogMgs))
		affirm.Equal(t, "msg 0", spy.haveLogMgs[0])
		affirm.True(t, spy.haveError)
	})

	t.Run("multi call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		spy.Errorf("msg %d", 0)
		spy.Errorf("msg %d", 1)

		// --- Then ---
		affirm.Equal(t, 2, len(spy.haveLogMgs))
		affirm.Equal(t, "msg 0", spy.haveLogMgs[0])
		affirm.Equal(t, "msg 1", spy.haveLogMgs[1])
		affirm.True(t, spy.haveError)
	})

	t.Run("panics when called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- Then ---
		affirm.Panic(t, errMockOnNotClosed, func() { spy.Errorf("msg %d", 0) })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errActionOnFinished, func() { spy.Errorf("msg %d", 0) })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_ExpectFatal(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectFatal()

		// --- Then ---
		affirm.True(t, spy.wantFatal)
		affirm.True(t, spy == have)
	})

	t.Run("panics when called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- Then ---
		affirm.Panic(t, errExpectOnClosed, func() { spy.ExpectFatal() })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errExpectOnFinished, func() { spy.ExpectFatal() })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called after ExpectFail", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectFail()

		// --- Then ---
		want := "cannot use ExpectFatal and ExpectFail in the same time"
		affirm.Panic(t, want, func() { spy.ExpectFatal() })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_Fatal(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		affirm.Panic(t, FailNowMsg, func() { spy.Fatal("msg", 0) })

		// --- Then ---
		affirm.Equal(t, 1, len(spy.haveLogMgs))
		affirm.Equal(t, "msg 0", spy.haveLogMgs[0])
		affirm.True(t, spy.haveFatal)
	})

	t.Run("multi call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		affirm.Panic(t, FailNowMsg, func() { spy.Fatal("msg", 0) })
		affirm.Panic(t, FailNowMsg, func() { spy.Fatal("msg", 1) })

		// --- Then ---
		affirm.Equal(t, 2, len(spy.haveLogMgs))
		affirm.Equal(t, "msg 0", spy.haveLogMgs[0])
		affirm.Equal(t, "msg 1", spy.haveLogMgs[1])
		affirm.True(t, spy.haveFatal)
	})

	t.Run("panics when called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- Then ---
		affirm.Panic(t, errMockOnNotClosed, func() { spy.Fatal("msg", 0) })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errActionOnFinished, func() { spy.Fatal("msg", 0) })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_Fatalf(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		affirm.Panic(t, FailNowMsg, func() { spy.Fatalf("msg %d", 0) })

		// --- Then ---
		affirm.Equal(t, 1, len(spy.haveLogMgs))
		affirm.Equal(t, "msg 0", spy.haveLogMgs[0])
		affirm.True(t, spy.haveFatal)
	})

	t.Run("multi call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		affirm.Panic(t, FailNowMsg, func() { spy.Fatalf("msg %d", 0) })
		affirm.Panic(t, FailNowMsg, func() { spy.Fatalf("msg %d", 1) })

		// --- Then ---
		affirm.Equal(t, 2, len(spy.haveLogMgs))
		affirm.Equal(t, "msg 0", spy.haveLogMgs[0])
		affirm.Equal(t, "msg 1", spy.haveLogMgs[1])
		affirm.True(t, spy.haveFatal)
	})

	t.Run("panics when called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- Then ---
		affirm.Panic(t, errMockOnNotClosed, func() { spy.Fatalf("msg %d", 0) })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errActionOnFinished, func() { spy.Fatalf("msg %d", 0) })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_FailNow(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- Then ---
		affirm.Panic(t, FailNowMsg, func() { spy.FailNow() })
		affirm.False(t, spy.panicked)
		affirm.True(t, spy.haveFatal)
	})

	t.Run("panics when called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- Then ---
		affirm.Panic(t, errMockOnNotClosed, func() { spy.Fatal("msg", 0) })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errActionOnFinished, func() { spy.Fatal("msg", 0) })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_Failed(t *testing.T) {
	t.Run("not failed", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		have := spy.Failed()

		// --- Then ---
		affirm.False(t, have)
	})

	t.Run("have error", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()
		spy.Error("msg", 0)

		// --- When ---
		have := spy.Failed()

		// --- Then ---
		affirm.True(t, have)
	})

	t.Run("have fatal", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()
		affirm.Panic(t, FailNowMsg, func() { spy.Fatal("msg", 0) })

		// --- When ---
		have := spy.Failed()

		// --- Then ---
		affirm.True(t, have)
	})

	t.Run("have fatal and error", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()
		affirm.Panic(t, FailNowMsg, func() { spy.Fatal("msg", 0) })
		spy.Error("msg", 1)

		// --- When ---
		have := spy.Failed()

		// --- Then ---
		affirm.True(t, have)
	})

	t.Run("panics when called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- Then ---
		affirm.Panic(t, errMockOnNotClosed, func() { spy.Failed() })
		affirm.True(t, spy.panicked)
	})

	t.Run("does not panic when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()
		spy.Finish()

		// --- When ---
		have := spy.Failed()

		// --- Then ---
		affirm.False(t, have)
	})
}

func Test_Spy_ExpectHelpers(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti)

		// --- When ---
		have := spy.ExpectHelpers(2)

		// --- Then ---
		affirm.Equal(t, 2, spy.wantHelperCnt)
		affirm.True(t, spy.helperCntSet)
		affirm.True(t, spy == have)
	})

	t.Run("set at least once", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti)

		// --- When ---
		have := spy.ExpectHelpers(-1)

		// --- Then ---
		affirm.Equal(t, -1, spy.wantHelperCnt)
		affirm.True(t, spy.helperCntSet)
		affirm.True(t, spy == have)
	})

	t.Run("argument must be greater or equal minus one", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti)

		// --- Then ---
		want := "ExpectHelpers cnt must be greater or equal to minus one"
		affirm.Panic(t, want, func() { spy.ExpectHelpers(-2) })
		affirm.True(t, spy.panicked)
	})

	t.Run("may be set only once", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}
		spy := New(ti)
		spy.ExpectHelpers(1)

		// --- Then ---
		want := "ExpectHelpers may be called only once"
		affirm.Panic(t, want, func() { spy.ExpectHelpers(1) })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- Then ---
		affirm.Panic(t, errExpectOnClosed, func() { spy.ExpectHelpers(1) })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errExpectOnFinished, func() { spy.ExpectHelpers(1) })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_Helper(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		spy.Helper()

		// --- Then ---
		affirm.Equal(t, 1, spy.haveHelperCnt)
	})

	t.Run("multiple calls", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		spy.Helper()
		spy.Helper()

		// --- Then ---
		affirm.Equal(t, 2, spy.haveHelperCnt)
	})

	t.Run("panics when called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- Then ---
		affirm.Panic(t, errMockOnNotClosed, func() { spy.Helper() })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errActionOnFinished, func() { spy.Helper() })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_ExpectSetenv(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectSetenv("k0", "v0")

		// --- Then ---
		affirm.Equal(t, 1, len(spy.wantEnv))
		affirm.Equal(t, "v0", spy.wantEnv["k0"])
		affirm.True(t, spy == have)
	})

	t.Run("panics when called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- Then ---
		affirm.Panic(t, errExpectOnClosed, func() { spy.ExpectSetenv("k0", "v0") })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errExpectOnFinished, func() { spy.ExpectSetenv("k0", "v0") })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_Setenv(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		spy.Setenv("k0", "v0")

		// --- Then ---
		affirm.Equal(t, 1, len(spy.haveEnv))
		affirm.Equal(t, "v0", spy.haveEnv["k0"])
	})

	t.Run("set multiple", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		spy.Setenv("k0", "v0")
		spy.Setenv("k1", "v1")

		// --- Then ---
		affirm.Equal(t, 2, len(spy.haveEnv))
		affirm.Equal(t, "v0", spy.haveEnv["k0"])
		affirm.Equal(t, "v1", spy.haveEnv["k1"])
	})

	t.Run("panics when called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- Then ---
		affirm.Panic(t, errMockOnNotClosed, func() { spy.Setenv("k0", "v0") })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errActionOnFinished, func() { spy.Setenv("k0", "v0") })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_ExpectSkipped(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectSkipped()

		// --- Then ---
		affirm.True(t, spy.wantSkipped)
		affirm.True(t, spy == have)
	})

	t.Run("panics when called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- Then ---
		affirm.Panic(t, errExpectOnClosed, func() { spy.ExpectSkipped() })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errExpectOnFinished, func() { spy.ExpectSkipped() })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_Skip(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		spy.Skip("msg", 0)

		// --- Then ---
		affirm.Equal(t, 1, len(spy.haveLogMgs))
		affirm.Equal(t, "msg 0", spy.haveLogMgs[0])
		affirm.True(t, spy.haveSkipped)
	})

	t.Run("multi call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		spy.Skip("msg", 0)
		spy.Skip("msg", 1)

		// --- Then ---
		affirm.Equal(t, 2, len(spy.haveLogMgs))
		affirm.Equal(t, "msg 0", spy.haveLogMgs[0])
		affirm.Equal(t, "msg 1", spy.haveLogMgs[1])
		affirm.True(t, spy.haveSkipped)
	})

	t.Run("panics when called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- Then ---
		affirm.Panic(t, errMockOnNotClosed, func() { spy.Skip("msg", 0) })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errActionOnFinished, func() { spy.Skip("msg", 0) })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_ExpectTempDir(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectTempDir(1)

		// --- Then ---
		affirm.Equal(t, 1, spy.wantTempDirCnt)
		affirm.True(t, spy == have)
	})

	t.Run("panics when called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- Then ---
		affirm.Panic(t, errExpectOnClosed, func() { spy.ExpectTempDir(1) })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errExpectOnFinished, func() { spy.ExpectTempDir(1) })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_GetTempDir(t *testing.T) {
	t.Run("get directories", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectTempDir(2)
		spy.Close()
		want0 := spy.TempDir()
		want1 := spy.TempDir()
		spy.Finish()

		// --- When ---
		have0 := spy.GetTempDir(0)
		have1 := spy.GetTempDir(1)

		// --- Then ---
		affirm.Equal(t, want0, have0)
		affirm.Equal(t, want1, have1)
	})

	t.Run("get directories when calls to TempDir do not matter", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectTempDir(-1)
		spy.Close()
		want0 := spy.TempDir()
		want1 := spy.TempDir()
		spy.Finish()

		// --- When ---
		have0 := spy.GetTempDir(0)
		have1 := spy.GetTempDir(1)

		// --- Then ---
		affirm.Equal(t, want0, have0)
		affirm.Equal(t, want1, have1)
	})

	t.Run("GetTempDir called before ExpectTempDir", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()
		spy.Finish()

		// --- When ---
		have := spy.GetTempDir(0)

		// --- Then ---
		affirm.Equal(t, "", have)
		affirm.False(t, spy.Failed())
		affirm.True(t, ti.Failed())
		assertSpyHasMgs(t, spy, 1)
		msg := "ExpectTempDir method must be called before GetTempDir"
		assertSpyHasMsg(t, spy, 0, msg)
	})

	t.Run("get not existing index", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectTempDir(2)
		spy.Close()
		spy.TempDir()
		spy.TempDir()
		spy.Finish()

		// --- When ---
		have := spy.GetTempDir(2)

		// --- Then ---
		affirm.Equal(t, "", have)
		affirm.False(t, spy.Failed())
		affirm.True(t, ti.Failed())
		assertSpyHasMgs(t, spy, 1)
		msg := "temp directory with index 2 does not exist"
		assertSpyHasMsg(t, spy, 0, msg)
	})
}

func Test_Spy_TempDir(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		have := spy.TempDir()

		// --- Then ---
		fi, err := os.Lstat(have)
		affirm.True(t, err == nil)
		affirm.True(t, fi.IsDir())
		affirm.Equal(t, 1, len(spy.haveTempDirs))

		// Dir is added to cleanups and removed when test finishes.
		spy.Finish()
		_, err = os.Lstat(have)
		affirm.True(t, errors.Is(err, os.ErrNotExist))
	})

	t.Run("call multi", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		got0 := spy.TempDir()
		got1 := spy.TempDir()

		// --- Then ---
		affirm.False(t, got0 == got1)

		fi, err := os.Lstat(got0)
		affirm.True(t, err == nil)
		affirm.True(t, fi.IsDir())

		fi, err = os.Lstat(got1)
		affirm.True(t, err == nil)
		affirm.True(t, fi.IsDir())

		affirm.Equal(t, 2, len(spy.haveTempDirs))
	})

	t.Run("panics when called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- Then ---
		affirm.Panic(t, errMockOnNotClosed, func() { spy.TempDir() })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errActionOnFinished, func() { spy.TempDir() })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_Context(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		have := spy.Context()

		// --- Then ---
		affirm.True(t, have != nil)
		affirm.True(t, have.Err() == nil)

		spy.Finish()
		affirm.True(t, errors.Is(have.Err(), context.Canceled))
	})

	t.Run("call multi", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		have0 := spy.Context()
		have1 := spy.Context()

		// --- Then ---
		affirm.True(t, have0 == have1)
	})

	t.Run("panics when called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- Then ---
		affirm.Panic(t, errMockOnNotClosed, func() { spy.Context() })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errActionOnFinished, func() { spy.Context() })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_IgnoreLogs(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- When ---
		have := spy.IgnoreLogs()

		// --- Then ---
		affirm.True(t, spy.ignoreLog)
		affirm.True(t, spy == have)
	})

	t.Run("panics when called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- Then ---
		affirm.Panic(t, errExpectOnClosed, func() { spy.IgnoreLogs() })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errExpectOnFinished, func() { spy.IgnoreLogs() })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when any of the ExpectLog* methods used", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectLogEqual("abc")

		// --- Then ---
		affirm.Panic(t, errIgnoreLogsAfterExpectLog, func() { spy.IgnoreLogs() })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_ExpectLog(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectLog(Contains, "msg %d", 0)

		// --- Then ---
		affirm.Equal(t, 1, len(spy.wantLogMgs))
		ent := spy.wantLogMgs[0]
		affirm.Equal(t, Contains, ent.strategy)
		affirm.Equal(t, "msg 0", ent.want)
		affirm.True(t, spy == have)
	})

	t.Run("does not add empty messages", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectLog(Contains, "")

		// --- Then ---
		affirm.Equal(t, 0, len(spy.wantLogMgs))
		affirm.True(t, spy == have)
	})

	t.Run("set with percentages without args", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectLog(Contains, "msg %d")

		// --- Then ---
		affirm.Equal(t, 1, len(spy.wantLogMgs))
		ent := spy.wantLogMgs[0]
		affirm.Equal(t, Contains, ent.strategy)
		affirm.Equal(t, "msg %d", ent.want)
		affirm.True(t, spy == have)
	})

	t.Run("multi set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- When ---
		spy.ExpectLog(Contains, "msg %d", 0)
		spy.ExpectLog(Equal, "msg %d", 1)
		spy.ExpectLog(Regexp, "msg [0-9]")
		spy.ExpectLog(NotContains, "msg 6")

		// --- Then ---
		affirm.Equal(t, 4, len(spy.wantLogMgs))
		ent := spy.wantLogMgs[0]
		affirm.Equal(t, Contains, ent.strategy)
		affirm.Equal(t, "msg 0", ent.want)

		ent = spy.wantLogMgs[1]
		affirm.Equal(t, Equal, ent.strategy)
		affirm.Equal(t, "msg 1", ent.want)

		ent = spy.wantLogMgs[2]
		affirm.Equal(t, Regexp, ent.strategy)
		affirm.Equal(t, "msg [0-9]", ent.want)

		ent = spy.wantLogMgs[3]
		affirm.Equal(t, NotContains, ent.strategy)
		affirm.Equal(t, "msg 6", ent.want)
	})

	t.Run("panics when called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- Then ---
		affirm.Panic(t, errExpectOnClosed, func() { spy.ExpectLog(Equal, "msg") })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errExpectOnFinished, func() { spy.ExpectLog(Equal, "msg") })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called after IgnoreLogs", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.IgnoreLogs()

		// --- Then ---
		affirm.Panic(t, errExpectLogAfterIgnoreLogs, func() { spy.ExpectLog(Equal, "msg") })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_ExpectLogEqual(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectLogEqual("msg %d", 0)

		// --- Then ---
		affirm.Equal(t, 1, len(spy.wantLogMgs))
		ent := spy.wantLogMgs[0]
		affirm.Equal(t, Equal, ent.strategy)
		affirm.Equal(t, "msg 0", ent.want)
		affirm.True(t, spy == have)
	})

	t.Run("multi set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- When ---
		spy.ExpectLogEqual("msg %d", 0)
		spy.ExpectLogEqual("msg %d", 1)

		// --- Then ---
		affirm.Equal(t, 2, len(spy.wantLogMgs))
		ent := spy.wantLogMgs[0]
		affirm.Equal(t, Equal, ent.strategy)
		affirm.Equal(t, "msg 0", ent.want)

		ent = spy.wantLogMgs[1]
		affirm.Equal(t, Equal, ent.strategy)
		affirm.Equal(t, "msg 1", ent.want)
	})

	t.Run("panics when called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- Then ---
		affirm.Panic(t, errExpectOnClosed, func() { spy.ExpectLogEqual("msg") })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errExpectOnFinished, func() { spy.ExpectLogEqual("msg") })
	})

	t.Run("panics when called after IgnoreLogs", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.IgnoreLogs()

		// --- Then ---
		affirm.Panic(t, errExpectLogAfterIgnoreLogs, func() { spy.ExpectLogEqual("msg") })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_ExpectLogContain(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectLogContain("msg %d", 0)

		// --- Then ---
		affirm.Equal(t, 1, len(spy.wantLogMgs))
		ent := spy.wantLogMgs[0]
		affirm.Equal(t, Contains, ent.strategy)
		affirm.Equal(t, "msg 0", ent.want)
		affirm.True(t, spy == have)
	})

	t.Run("multi set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- When ---
		spy.ExpectLogContain("msg %d", 0)
		spy.ExpectLogContain("msg %d", 1)

		// --- Then ---
		affirm.Equal(t, 2, len(spy.wantLogMgs))
		ent := spy.wantLogMgs[0]
		affirm.Equal(t, Contains, ent.strategy)
		affirm.Equal(t, "msg 0", ent.want)

		ent = spy.wantLogMgs[1]
		affirm.Equal(t, Contains, ent.strategy)
		affirm.Equal(t, "msg 1", ent.want)
	})

	t.Run("panics when called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- Then ---
		affirm.Panic(t, errExpectOnClosed, func() { spy.ExpectLogContain("msg") })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errExpectOnFinished, func() { spy.ExpectLogContain("msg") })
	})

	t.Run("panics when called after IgnoreLogs", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.IgnoreLogs()

		// --- Then ---
		affirm.Panic(t, errExpectLogAfterIgnoreLogs, func() { spy.ExpectLogContain("msg") })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_ExpectLogNotContain(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectLogNotContain("msg %d", 0)

		// --- Then ---
		affirm.Equal(t, 1, len(spy.wantLogMgs))
		ent := spy.wantLogMgs[0]
		affirm.Equal(t, NotContains, ent.strategy)
		affirm.Equal(t, "msg 0", ent.want)
		affirm.True(t, spy == have)
	})

	t.Run("multi set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- When ---
		spy.ExpectLogNotContain("msg %d", 0)
		spy.ExpectLogNotContain("msg %d", 1)

		// --- Then ---
		affirm.Equal(t, 2, len(spy.wantLogMgs))
		ent := spy.wantLogMgs[0]
		affirm.Equal(t, NotContains, ent.strategy)
		affirm.Equal(t, "msg 0", ent.want)

		ent = spy.wantLogMgs[1]
		affirm.Equal(t, NotContains, ent.strategy)
		affirm.Equal(t, "msg 1", ent.want)
	})

	t.Run("panics when called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- Then ---
		affirm.Panic(t, errExpectOnClosed, func() { spy.ExpectLogNotContain("msg") })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errExpectOnFinished, func() { spy.ExpectLogNotContain("msg") })
	})

	t.Run("panics when called after IgnoreLogs", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.IgnoreLogs()

		// --- Then ---
		affirm.Panic(t, errExpectLogAfterIgnoreLogs, func() { spy.ExpectLogNotContain("msg") })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_Log(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		spy.Log("msg", 0)

		// --- Then ---
		affirm.Equal(t, 1, len(spy.haveLogMgs))
		affirm.Equal(t, "msg 0", spy.haveLogMgs[0])
	})

	t.Run("multi call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		spy.Log("msg", 0)
		spy.Log("msg", 1)

		// --- Then ---
		affirm.Equal(t, 2, len(spy.haveLogMgs))
		affirm.Equal(t, "msg 0", spy.haveLogMgs[0])
		affirm.Equal(t, "msg 1", spy.haveLogMgs[1])
	})

	t.Run("panics when called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- Then ---
		affirm.Panic(t, errMockOnNotClosed, func() { spy.Log("msg", 0) })
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errActionOnFinished, func() { spy.Log("msg", 0) })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_Logf(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		spy.Logf("msg %d", 0)

		// --- Then ---
		affirm.Equal(t, 1, len(spy.haveLogMgs))
		affirm.Equal(t, "msg 0", spy.haveLogMgs[0])
	})

	t.Run("multi call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- When ---
		spy.Logf("msg %d", 0)
		spy.Logf("msg %d", 1)

		// --- Then ---
		affirm.Equal(t, 2, len(spy.haveLogMgs))
		affirm.Equal(t, "msg 0", spy.haveLogMgs[0])
		affirm.Equal(t, "msg 1", spy.haveLogMgs[1])
	})

	t.Run("panics when called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- Then ---
		affirm.Panic(t, errMockOnNotClosed, func() { spy.Logf("msg %d", 0) })
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errActionOnFinished, func() { spy.Logf("msg %d", 0) })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_ExpectedNames(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectedNames(2)

		// --- Then ---
		affirm.Equal(t, 2, spy.wantNamesCnt)
		affirm.True(t, spy == have)
	})

	t.Run("panics when called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()

		// --- Then ---
		affirm.Panic(t, errExpectOnClosed, func() { spy.ExpectedNames(1) })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errExpectOnFinished, func() { spy.ExpectedNames(1) })
		affirm.True(t, spy.panicked)
	})
}

func Test_Name(t *testing.T) {
	t.Run("call", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		spy.Name()

		// --- Then ---
		affirm.Equal(t, 1, spy.haveNamesCnt)
	})

	t.Run("multiple calls", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- When ---
		spy.Name()
		spy.Name()

		// --- Then ---
		affirm.Equal(t, 2, spy.haveNamesCnt)
	})

	t.Run("returns test name", func(t *testing.T) {
		// --- Given ---
		spy := New(t, 0)
		spy.ExpectedNames(1)
		spy.Close()

		// --- When ---
		have := spy.Name()

		// --- Then ---
		affirm.Equal(t, "Test_Name/returns_test_name", have)
	})

	t.Run("panics when called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- Then ---
		affirm.Panic(t, errMockOnNotClosed, func() { spy.Name() })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errActionOnFinished, func() { spy.Name() })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_ExpectFail(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- When ---
		have := spy.ExpectFail()

		// --- Then ---
		affirm.True(t, spy.wantFailed)
		affirm.True(t, spy == have)
	})

	t.Run("panics when called after ExpectFatal", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.ExpectFatal()

		// --- Then ---
		want := "cannot use ExpectFail and ExpectFatal in the same time"
		affirm.Panic(t, want, func() { spy.ExpectFail() })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called after ExpectError", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.ExpectError()

		// --- When ---

		// --- Then ---
		want := "cannot use ExpectFail and ExpectError in the same time"
		affirm.Panic(t, want, func() { spy.ExpectFail() })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- Then ---
		affirm.Panic(t, errExpectOnClosed, func() { spy.ExpectFail() })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errExpectOnFinished, func() { spy.ExpectFail() })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_Close(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- When ---
		have := spy.Close()

		// --- Then ---
		affirm.True(t, spy.closed)
		affirm.True(t, spy == have)
	})

	t.Run("panics when called on closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// --- Then ---
		affirm.Panic(t, errDoubleClose, func() { spy.Close() })
		affirm.True(t, spy.panicked)
	})

	t.Run("panics when called on finished Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errCloseOnFinished, func() { spy.Close() })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_Finish(t *testing.T) {
	t.Run("finish test", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)

		// --- When ---
		spy.Finish()

		// --- Then ---
		affirm.True(t, spy.finished)
	})

	t.Run("runs cleanups", func(t *testing.T) {
		// --- Given ---
		var runs []int
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()
		spy.Cleanup(func() { runs = append(runs, 0) })
		spy.Cleanup(func() { runs = append(runs, 1) })

		// --- When ---
		spy.Finish()

		// --- Then ---
		affirm.Equal(t, 2, len(runs))
		affirm.Equal(t, 1, runs[0])
		affirm.Equal(t, 0, runs[1])
	})

	t.Run("panics when run twice", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Finish()

		// --- Then ---
		affirm.Panic(t, errDoubleFinish, func() { spy.Finish() })
		affirm.True(t, spy.panicked)
	})
}

func Test_Spy_AssertExpectations(t *testing.T) {
	t.Run("finishes and runs cleanups", func(t *testing.T) {
		// --- Given ---
		var runs []int
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectCleanups(1)
		spy.Close()
		spy.Cleanup(func() { runs = append(runs, 0) })

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.True(t, have)
		affirm.True(t, spy.finished)
		affirm.Equal(t, 1, len(runs))
		affirm.Equal(t, 0, runs[0])
	})

	t.Run("Spy with no expectations is success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.True(t, have)
		affirm.False(t, spy.Failed())
		affirm.False(t, ti.Failed())
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("panics when called on not closed Spy", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)

		// --- Then ---
		affirm.Panic(t, errAssertOnNotClosed, func() { spy.AssertExpectations() })
		affirm.True(t, spy.panicked)
	})

	t.Run("called on panicked instance", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.Close()

		// Make Spy panic.
		spy.Finish()
		affirm.Panic(t, errDoubleFinish, func() { spy.Finish() })
		affirm.True(t, spy.panicked)

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.False(t, have)
		affirm.False(t, spy.Failed())
		affirm.True(t, ti.Failed())
		assertSpyHasMgs(t, spy, 1)
		assertSpyHasMsg(t, spy, 0, "invalid Spy usage")
	})

	t.Run("ExpectSetenv - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectSetenv("key0", "val0")
		spy.ExpectSetenv("key1", "val1")
		spy.Close()
		spy.Setenv("key0", "val0")
		spy.Setenv("key1", "val1")
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.True(t, have)
		affirm.False(t, spy.Failed())
		affirm.False(t, ti.Failed())
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("ExpectSetenv - error - not expected but set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()
		spy.Setenv("key0", "val0")
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.False(t, have)
		affirm.False(t, spy.Failed())
		affirm.True(t, ti.Failed())
		assertSpyHasMgs(t, spy, 1)
		wMsg := "expected HUT not to set environment variable:\n" +
			"\t  have key: \"key0\"\n" +
			"\thave value: \"val0\""
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectSetenv - error - none set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectSetenv("key0", "val0")
		spy.ExpectSetenv("key1", "val1")
		spy.Close()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.False(t, have)
		affirm.False(t, spy.Failed())
		affirm.True(t, ti.Failed())
		assertSpyHasMgs(t, spy, 2)
		wMsg := "expected HUT to set environment variable:\n" +
			"\t  want key: \"key0\"\n" +
			"\twant value: \"val0\""
		assertSpyHasMsg(t, spy, 0, wMsg)
		wMsg = "expected HUT to set environment variable:\n" +
			"\t  want key: \"key1\"\n" +
			"\twant value: \"val1\""
		assertSpyHasMsg(t, spy, 1, wMsg)
	})

	t.Run("ExpectSetenv - error - wrong value set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectSetenv("key0", "val0")
		spy.ExpectSetenv("key1", "val1")
		spy.Close()
		spy.Setenv("key0", "valX")
		spy.Setenv("key1", "valY")
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.False(t, have)
		affirm.False(t, spy.Failed())
		affirm.True(t, ti.Failed())
		assertSpyHasMgs(t, spy, 2)
		wMsg := "expected HUT to set environment variable:\n" +
			"\t  want key: \"key0\"\n" +
			"\twant value: \"val0\"\n" +
			"\thave value: \"valX\""
		assertSpyHasMsg(t, spy, 0, wMsg)
		wMsg = "expected HUT to set environment variable:\n" +
			"\t  want key: \"key1\"\n" +
			"\twant value: \"val1\"\n" +
			"\thave value: \"valY\""
		assertSpyHasMsg(t, spy, 1, wMsg)
	})

	t.Run("ExpectSetenv - error - not all required set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectSetenv("key0", "val0")
		spy.ExpectSetenv("key1", "val1")
		spy.Close()
		spy.Setenv("key1", "val1")
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.False(t, have)
		affirm.False(t, spy.Failed())
		affirm.True(t, ti.Failed())
		assertSpyHasMgs(t, spy, 1)
		wMsg := "expected HUT to set environment variable:\n" +
			"\t  want key: \"key0\"\n" +
			"\twant value: \"val0\""
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectSetenv - error - more than expected set", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectSetenv("key0", "val0")
		spy.ExpectSetenv("key1", "val1")
		spy.Close()
		spy.Setenv("key0", "val0")
		spy.Setenv("key1", "val1")
		spy.Setenv("key2", "val2")
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.False(t, have)
		affirm.False(t, spy.Failed())
		affirm.True(t, ti.Failed())
		assertSpyHasMgs(t, spy, 1)
		wMsg := "expected HUT not to set environment variable:\n" +
			"\t  have key: \"key2\"\n" +
			"\thave value: \"val2\""
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectSkipped - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectSkipped()
		spy.ExpectLogEqual("msg 0")
		spy.Close()
		spy.Skip("msg", 0)
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.True(t, have)
		affirm.False(t, spy.Failed())
		affirm.False(t, ti.Failed())
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("error ExpectSkipped", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectLogEqual("msg 0")
		spy.Close()
		spy.Skip("msg", 0)
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.False(t, have)
		affirm.False(t, spy.Failed())
		affirm.True(t, ti.Failed())
		assertSpyHasMgs(t, spy, 1)
		wMsg := "expected HUT to mark test as skipped:\n" +
			"\twant: false\n" +
			"\thave: true"
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectFail - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectFail()
		spy.ExpectLogEqual("msg 0")
		spy.Close()
		spy.Error("msg", 0)
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.True(t, have)
		affirm.True(t, spy.Failed())
		affirm.False(t, ti.Failed())
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("ExpectFail - fail", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectFail()
		spy.Close()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.False(t, have)
		affirm.False(t, spy.Failed())
		affirm.True(t, ti.Failed())
		wMsg := "expected HUT to call the t.Error* or t.Fatal* methods"
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectError - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectError()
		spy.ExpectLogEqual("msg 0")
		spy.Close()
		spy.Error("msg", 0)
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.True(t, have)
		affirm.True(t, spy.Failed())
		affirm.False(t, ti.Failed())
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("ExpectError - fail", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectError()
		spy.Close()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.False(t, have)
		affirm.False(t, spy.Failed())
		affirm.True(t, ti.Failed())
		wMsg := "expected HUT to call any of the t.Error* methods"
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectFatal - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectFatal()
		spy.ExpectLogEqual("msg 0")
		spy.Close()
		affirm.Panic(t, FailNowMsg, func() { spy.Fatal("msg", 0) })
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.True(t, have)
		affirm.True(t, spy.Failed())
		affirm.False(t, ti.Failed())
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("ExpectFatal - fail", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectFatal()
		spy.Close()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.False(t, have)
		affirm.False(t, spy.Failed())
		affirm.True(t, ti.Failed())
		wMsg := "expected HUT to call any of the t.Fatal* methods"
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectHelpers - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.ExpectHelpers(2)
		spy.Close()
		spy.Helper()
		spy.Helper()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.True(t, have)
		affirm.False(t, spy.Failed())
		affirm.False(t, ti.Failed())
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("ExpectHelpers - fail", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.ExpectHelpers(2)
		spy.Close()
		spy.Helper()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.False(t, have)
		affirm.False(t, spy.Failed())
		affirm.True(t, ti.Failed())
		wMsg := "expected Helper to be called N times:\n" +
			"\twant: 2\n" +
			"\thave: 1"
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectNames - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectedNames(2)
		spy.Close()
		spy.Name()
		spy.Name()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.True(t, have)
		affirm.False(t, spy.Failed())
		affirm.False(t, ti.Failed())
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("ExpectNames - fail", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectedNames(2)
		spy.Close()
		spy.Name()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.False(t, have)
		affirm.False(t, spy.Failed())
		affirm.True(t, ti.Failed())
		wMsg := "expected Name to be called N times:\n" +
			"\twant: 2\n" +
			"\thave: 1"
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectCleanups - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectCleanups(2)
		spy.Close()
		spy.Cleanup(func() {})
		spy.Cleanup(func() {})
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.True(t, have)
		affirm.False(t, spy.Failed())
		affirm.False(t, ti.Failed())
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("ExpectCleanups - fail", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectCleanups(2)
		spy.Close()
		spy.Cleanup(func() {})
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.False(t, have)
		affirm.False(t, spy.Failed())
		affirm.True(t, ti.Failed())
		wMsg := "expected Cleanup to be called N times:\n" +
			"\twant: 2\n" +
			"\thave: 1"
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectTempDir - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectTempDir(2)
		spy.Close()
		spy.TempDir()
		spy.TempDir()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.True(t, have)
		affirm.False(t, spy.Failed())
		affirm.False(t, ti.Failed())
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("ExpectTempDir - does not matter and called", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectTempDir(-1)
		spy.Close()
		spy.TempDir()
		spy.TempDir()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.True(t, have)
		affirm.False(t, spy.Failed())
		affirm.False(t, ti.Failed())
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("ExpectTempDir - does not matter and not called", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectTempDir(-1)
		spy.Close()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.True(t, have)
		affirm.False(t, spy.Failed())
		affirm.False(t, ti.Failed())
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("ExpectTempDir - fail", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectTempDir(2)
		spy.Close()
		spy.TempDir()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.False(t, have)
		affirm.False(t, spy.Failed())
		affirm.True(t, ti.Failed())
		wMsg := "expected TempDir to be called N times:\n" +
			"\twant: 2\n" +
			"\thave: 1"
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectLog - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectLog(Equal, "msg %d", 0)
		spy.Close()
		spy.Logf("msg %d", 0)
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.True(t, have)
		affirm.False(t, spy.Failed())
		affirm.False(t, ti.Failed())
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("ExpectLog - fail", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectLog(Equal, "msg %d", 0)
		spy.Close()
		spy.Log("msg", 1)
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.False(t, have)
		affirm.False(t, spy.Failed())
		affirm.True(t, ti.Failed())
		wMsg := "expected HUT to log message 0:\n" +
			"\tmatcher: equal\n" +
			"\t   want: \"msg 0\"\n" +
			"\t   have: \"msg 1\""
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("ExpectLog - not contains - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectLog(NotContains, "msg %d", 1)
		spy.Close()
		spy.Logf("msg %d", 0)
		spy.Logf("msg %d", 2)
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.True(t, have)
		affirm.False(t, spy.Failed())
		affirm.False(t, ti.Failed())
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("ExpectLog - not contains - fail", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.ExpectLog(NotContains, "msg %d", 0)
		spy.Close()
		spy.Logf("msg %d", 0)
		spy.Logf("msg %d", 1)
		spy.Logf("msg %d", 2)
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.False(t, have)
		affirm.False(t, spy.Failed())
		affirm.True(t, ti.Failed())
		wMsg := "expected HUT to log message 0:\n" +
			"\tmatcher: not-contains\n" +
			"\t   want: \"msg 0\"\n" +
			"\t   have: \"msg 0\\nmsg 1\\nmsg 2\""
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("expected no logs - fail", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()
		spy.Log("msg")
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.False(t, have)
		affirm.False(t, spy.Failed())
		affirm.True(t, ti.Failed())
		wMsg := "expected HUT to log no messages but got:\n" +
			"\thave: \"msg\""
		assertSpyHasMsg(t, spy, 0, wMsg)
	})

	t.Run("do not care if messages were logged", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.IgnoreLogs()
		spy.Close()
		spy.Log("msg")
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.True(t, have)
		affirm.False(t, spy.Failed())
		affirm.False(t, ti.Failed())
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("expect Helper to be called at least once - success", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.ExpectHelpers(-1)
		spy.Close()
		spy.Helper()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.True(t, have)
		affirm.False(t, spy.Failed())
		affirm.False(t, ti.Failed())
		assertSpyHasMgs(t, spy, 0)
	})

	t.Run("expect Helper to be called at least once - fail", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti)
		spy.ExpectHelpers(-1)
		spy.Close()
		spy.Finish()

		// --- When ---
		have := spy.AssertExpectations()

		// --- Then ---
		affirm.False(t, have)
		affirm.False(t, spy.Failed())
		affirm.True(t, ti.Failed())
		wMsg := "expected Helper to be called N times:\n" +
			"\twant: >= 1\n" +
			"\thave: 0"
		assertSpyHasMsg(t, spy, 0, wMsg)
	})
}

func Test_Spy_assertFailed(t *testing.T) {
	t.Run("returns true if wantFailed false", func(t *testing.T) {
		// --- Given ---
		ti := &testing.T{}

		spy := New(ti, 0)
		spy.Close()
		spy.Finish()

		// --- When ---
		have := spy.assertFailed()

		// --- Then ---
		affirm.True(t, have)
	})
}

func Test_Spy_checkState(t *testing.T) {
	tt := []struct {
		testN string

		spy    func(ti *testing.T) *Spy
		action func(*Spy)
		want   string
	}{
		// Expect* calls.
		{
			"expect called on finished",
			func(ti *testing.T) *Spy { return New(ti).Close().Finish() },
			func(spy *Spy) { spy.ExpectError() },
			errExpectOnFinished,
		},
		{
			"expect called on closed",
			func(ti *testing.T) *Spy { return New(ti).Close() },
			func(spy *Spy) { spy.ExpectError() },
			errExpectOnClosed,
		},
		// Close calls.
		{
			"call to Close when finished",
			func(ti *testing.T) *Spy { return New(ti).Finish() },
			func(spy *Spy) { spy.Close() },
			errCloseOnFinished,
		},
		{
			"call to Close when already Closed",
			func(ti *testing.T) *Spy { return New(ti).Close() },
			func(spy *Spy) { spy.Close() },
			errDoubleClose,
		},
		// Mocked calls.
		{
			"call to mocked method when finished",
			func(ti *testing.T) *Spy { return New(ti).Close().Finish() },
			func(spy *Spy) { spy.Log("msg") },
			errActionOnFinished,
		},
		{
			"call to mocked method when not closed",
			func(ti *testing.T) *Spy { return New(ti) },
			func(spy *Spy) { spy.Log("msg") },
			errMockOnNotClosed,
		},
		// Assert* calls.
		{
			"assert* call when not finished",
			func(ti *testing.T) *Spy { return New(ti).Close() },
			func(spy *Spy) { spy.assertExpectations() },
			errAssertOnNotFinished,
		},
		{
			"assert* call when not closed",
			func(ti *testing.T) *Spy { return New(ti).Finish() },
			func(spy *Spy) { spy.assertExpectations() },
			errAssertOnNotClosed,
		},
		// Failed call.
		{
			"Failed call when not closed",
			func(ti *testing.T) *Spy { return New(ti) },
			func(spy *Spy) { spy.Failed() },
			errMockOnNotClosed,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			ti := &testing.T{}

			spy := tc.spy(ti)
			affirm.False(t, ti.Failed()) // Make sure it is not failed.

			// --- Then ---
			affirm.Panic(t, tc.want, func() { tc.action(spy) })
			affirm.True(t, spy.panicked)
			affirm.True(t, ti.Failed())
		})
	}
}

// =============================================================================
// ================================== HELPERS ==================================
// =============================================================================

// assertSpyHasMgs asserts spy has saved a "want" number of messages. On
// failure, it marks the test as failed, returns false, but continues execution.
func assertSpyHasMgs(t *testing.T, spy *Spy, want int) bool {
	t.Helper()
	have := len(spy.savedMgs)
	if want != have {
		format := "expected Spy to have saved messages:\n" +
			"\twant: %d\n" +
			"\thave: %d"
		if have > 0 {
			format += "\n\tmessages:\n"
			for idx, msg := range spy.savedMgs {
				format += fmt.Sprintf("\t\t%d: %q\n", idx, msg)
			}
		}
		t.Errorf(format, want, have)
		return false
	}
	return true
}

func Test_assertSpyHasMgs(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ti := &testing.T{}
		spy := &Spy{}
		spy.savedMgs = []string{"abc", "def"}
		have := assertSpyHasMgs(ti, spy, 2)
		if !have || ti.Failed() {
			t.Error("expected success")
		}
	})
	t.Run("error", func(t *testing.T) {
		ti := &testing.T{}
		spy := &Spy{}
		spy.savedMgs = []string{"abc", "def"}
		have := assertSpyHasMgs(ti, spy, 3)
		if have || !ti.Failed() {
			t.Error("expected failure")
		}
	})
}

// assertSpyHasMsg asserts Spy has saved a message with given idx equal to
// "want". On failure, it marks the test as failed, returns false, but
// continues execution.
func assertSpyHasMsg(t *testing.T, spy *Spy, idx int, want string) bool {
	t.Helper()
	if idx >= len(spy.savedMgs) {
		t.Errorf("expected Spy saved message with index %d to exist", idx)
		return false
	}
	have := spy.savedMgs[idx]
	if want != have {
		format := "expected Spy saved message with index %d to be:\n" +
			"\twant: %q\n" +
			"\thave: %q"
		t.Errorf(format, idx, want, have)
		return false
	}
	return true
}

func Test_assertSpyHasMsg(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ti := &testing.T{}
		spy := &Spy{}
		spy.savedMgs = []string{"abc", "def"}
		have := assertSpyHasMsg(ti, spy, 1, "def")
		if !have || ti.Failed() {
			t.Error("expected success")
		}
	})
	t.Run("error", func(t *testing.T) {
		ti := &testing.T{}
		spy := &Spy{}
		spy.savedMgs = []string{"abc", "def"}
		have := assertSpyHasMsg(ti, spy, 1, "xyz")
		if have || !ti.Failed() {
			t.Error("expected failure")
		}
	})
}
