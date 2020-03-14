// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Gaps and an offset.

package main

import (
	"encoding/json"
	"fmt"
)

type Gap int

const (
	Two    Gap = 2
	Three  Gap = 3
	Five   Gap = 5
	Six    Gap = 6
	Seven  Gap = 7
	Eight  Gap = 8
	Nine   Gap = 9
	Eleven Gap = 11
)

func main() {
	ck(0, "Gap(0)", true)
	ck(1, "Gap(1)", true)
	ck(Two, "Two", false)
	ck(Three, "Three", false)
	ck(4, "Gap(4)", true)
	ck(Five, "Five", false)
	ck(Six, "Six", false)
	ck(Seven, "Seven", false)
	ck(Eight, "Eight", false)
	ck(Nine, "Nine", false)
	ck(10, "Gap(10)", true)
	ck(Eleven, "Eleven", false)
	ck(12, "Gap(12)", true)
}

func ck(c Gap, str string, invalid bool) {
	if fmt.Sprint(c) != str {
		panic("gap.go: " + str)
	}
	{
		b, err := json.Marshal(c)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("gap.go: json.Marshal: expected an error for %s", c))
			}
			goto MarshalText
		}
		if err != nil {
			panic("gap.go: " + err.Error())
		}
		if string(b) != `"`+str+`"` {
			panic(fmt.Sprintf("gap.go: json.Marshal: got: %s: want: %q", b, str))
		}
		var v Gap
		if err := json.Unmarshal(b, &v); err != nil {
			panic("gap.go: json.Unmarshal: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("gap.go: json.Marshal: got: %s: want: %s", v, c))
		}
	}
MarshalText:
	{
		b, err := c.MarshalText()
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("gap.go: MarshalText: expected an error for %s", c))
			}
			goto Set
		}
		if err != nil {
			panic("gap.go: " + err.Error())
		}
		if string(b) != str {
			panic(fmt.Sprintf("gap.go: MarshalText: got: %s: want: %s", b, str))
		}
		var v Gap
		if err := v.UnmarshalText(b); err != nil {
			panic("gap.go: UnmarshalText: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("gap.go: MarshalText: got: %s: want: %s", v, c))
		}
	}
Set:
	{
		var v Gap
		err := v.Set(str)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("gap.go: Set: expected an error for %s", c))
			}
			goto Invalid
		}
		if v != c {
			panic(fmt.Sprintf("gap.go: Set: got: %s: want: %s", v, c))
		}
	}
Invalid:
	if invalid {
		var v Gap
		if err := json.Unmarshal([]byte(str), &v); err == nil {
			panic("gap.go: json.Unmarshal: expected an error for: " + str)
		}
		if err := v.UnmarshalText([]byte(str)); err == nil {
			panic("gap.go: UnmarshalText: expected an error for: " + str)
		}
	}
}
