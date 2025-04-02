// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package dump

import (
	"reflect"
	"strings"
)

// mapDumper requires val to be dereferenced representation of [reflect.Struct]
// and returns its string representation in format defined by [Dump]
// configuration.
func structDumper(dmp Dump, lvl int, val reflect.Value) string {
	prn := NewPrinter(dmp)
	prn.Tab(dmp.Indent + lvl)
	vTyp := val.Type()

	num := val.NumField() // Total number of fields.
	lastPrivate := false
	prn.Write("{").NLI(num)

	for i := 0; i < num; i++ {
		last := i == num-1

		fld := vTyp.Field(i)
		if !fld.IsExported() {
			lastPrivate = last
			continue
		}

		// Field name.
		prn.Tab(dmp.Indent + lvl + 1)
		prn.Write(fld.Name)
		prn.Write(":").Space()

		// Field value.
		dmp.PrintType = true
		sub := dmp.value(lvl+1, val.Field(i))
		sub = strings.TrimLeft(sub, " \t")

		prn.Write(sub)
		prn.Comma(last).Sep(last).NL()
	}

	if lastPrivate && dmp.Flat {
		return strings.TrimRight(prn.String(), ",") + "}"
	}
	prn.Tab(dmp.Indent + lvl).Write("}")

	return prn.String()
}
