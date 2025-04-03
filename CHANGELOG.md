## v0.5.0 (Thu, 03 Apr 2025 14:50:14 UTC)
- Add internal package diff.
- Add assert and check package.
- Update documentation and examples. Add ability to set dump.Config configuration in check.Options.
- Add check.ErrorIs and assert.ErrorIs.
- Add check.ErrorAs and assert.ErrorAs.
- Add check.ErrorEqual and assert.ErrorEqual.
- Add check.ErrorContain and assert.ErrorContain.
- Add check.Regexp and assert.Regexp.
- Add check.ErrorRegexp and assert.ErrorRegexp.
- Create files for specific check and assert topics.
- Refactor notice path - now it's called trail.
- Add check.FileExist, check.NoFileExist, check.DirExist, check.NoDirExist and assert.NoFileExist, assert.DirExist, assert.NoDirExist.
- Add check.Empty, check.NotEmpty and assert.Empty, assert.NotEmpty.
- Add check.Zero, check.NotZero and assert.Zero, assert.NotZero.
- Add internal helper Same which checks if two generic pointers point to the same memory.
- Add check.Same, check.NotSame and assert.Same, assert.NotSame.
- Add check.Len and assert.Len.
- Add check.True, check.False and assert.True, assert.False.
- Add check.Contain, check.NotContain and assert.Contain, assert.NotContain.
- Add check.Count, check.NotCount and assert.Count, assert.NotCount.
- Add check.Panic, check.NotPanic, check.PanicContain, check.PanicMsg and assert.Panic, assert.NotPanic, assert.PanicContain, assert.PanicMsg.
- Add check.FileContain and assert.FileContain.
- Add check.SameType and assert.SameType.
- Add check.Has, check.HasNo and assert.Has, assert.HasNo.
- Add check.HasKey, check.HasNoKey and assert.HasKey, assert.HasNoKey.
- Add check.HasKeyValue and assert.HasKeyValue.
- Add check.ExitCode and assert.ExitCode.
- Add check.SliceSubset and assert.SliceSubset.
- Change check.Options structure and add default time parsing option.
- Add FormatDates helper.
- Add helper getting time.Time from values of type any, where any may be time.Time, string, int, int64.
- Improve check.getTime error messages.
- Add check.TimeEqual and assert.TimeEqual.
- Add missing test cases to notice.Notice.
- Add check.TimeLoc and assert.TimeLoc. Build options only when needed.
- Add helper check.getDut getting time.Duration from values of type any, where any may be time.Duration, string, int, int64. Rename check and assertion TimeLoc to Zone and TimeEqual to Time.
- Add check.TimeExact and assert.TimeExact.
- Add check.Within and assert.Within.
- Better time check and assert function documentation.
- More readable log messages when asserting / checking comparing dates.
- Remove dead code.
- Update check.getDur helper. Code cleanup. Better error messages.
- Add check.Before and assert.Before.
- Add check.After and assert.After.
- Add check.EqualOrAfter and assert.EqualOrAfter.
- Add check.EqualOrBefore and assert.EqualOrBefore.
- Add ability to configure package-wide defaults.
- Update dump options and tests to use explicit "With" prefixes. Add ability to configure package-wide defaults.
- Add support for configurable recent duration in check package.
- Update documentation, add injectable clock to check.Options.
- Add check.Recent and assert.Recent.
- Better names.
- Improve test coverage.
- Add check.ChannelWillClose and assert.ChannelWillClose.
- Update dump package with better user interface and fixed dump format.
- Code cleanup and style.
- Indent with spaces when dumping values and when rendering notice messages.
- Improve notice package message formating.
- Remove internal diff package.
- Remove internal playground package.
- Update dump package documentation.
- Improve multiple notice messages rendering with the same header.
- Remove internal diff package.
- Add check.Equal which recursively compares two values.
- Fix code linting errors.
- Add assert.Equal.
- Export dump.Printer code printer.
- Refactor dump.Dump - separate dump.Config struct is not needed.
- Rename module.
- Update documentation.
- Use Option slice for equalError helper.
- Use custom byte dumper for equality errors.
- Add check.NotEqual and assert.NotEqual.
- Cleanup.
- Introduce internal core package with helpers used by higher level abstractions.
- Update behaviour for core.DidPanic and simplify affirm.Panic helper.
- Add check.JSON and assert.JSON.
- Add check.MapSubset and assert.MapSubset.
- Add check.MapsSubset and assert.MapsSubset.
- Rename check.SameType and assert.SameType to check.Type and assert.Type.
- Add check.Epsilon and assert.Epsilon.
- Add check.Fields and assert.Fields.
- Fix assert.Epsilon.
- The assert package documentation.
- The assert package documentation.

## v0.4.0 (Tue, 18 Mar 2025 18:20:26 UTC)
- Add Go Report Card and Go Doc icons / links.
- Add Go Test badge.
- Update affirm.Panic helper to handle errors, and other types that might be passed to panic.
- Add internal tstkit package with simple golden file reader.
- Add dump package.
- Update copyright lines to use SPDX-FileCopyrightText.

## v0.3.1 (Sat, 15 Mar 2025 12:00:58 UTC)
- Add work in progress disclaimer and more documentation.

## v0.3.0 (Sat, 15 Mar 2025 11:48:58 UTC)
- Create go.yml.
- Add package must with helper functions which panic on error. Add `affirm.DeeoEqual` helper.
- Add notice package and update documentation.

## v0.2.0 (Fri, 14 Mar 2025 19:26:12 UTC)
- Update documentation.

## v0.1.0 (Fri, 14 Mar 2025 17:37:29 UTC)
- Initial commit.

