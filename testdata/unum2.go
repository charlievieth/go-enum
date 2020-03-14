// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Unsigned integers - check maximum size

package main

import (
	"encoding/json"
	"fmt"
)

type Unum2 uint8

const (
	Zero Unum2 = iota
	One
	Two
)

func main() {
	ck(Zero, "Zero", false)
	ck(One, "One", false)
	ck(Two, "Two", false)
	ck(3, "Unum2(3)", true)
	ck(255, "Unum2(255)", true)
}

func ck(c Unum2, str string, invalid bool) {
	if fmt.Sprint(c) != str {
		panic("unum2.go: " + str)
	}
	{
		b, err := json.Marshal(c)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("unum2.go: json.Marshal: expected an error for %s", c))
			}
			goto MarshalText
		}
		if err != nil {
			panic("unum2.go: " + err.Error())
		}
		if string(b) != `"`+str+`"` {
			panic(fmt.Sprintf("unum2.go: json.Marshal: got: %s: want: %q", b, str))
		}
		var v Unum2
		if err := json.Unmarshal(b, &v); err != nil {
			panic("unum2.go: json.Unmarshal: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("unum2.go: json.Marshal: got: %s: want: %s", v, c))
		}
	}
MarshalText:
	{
		b, err := c.MarshalText()
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("unum2.go: MarshalText: expected an error for %s", c))
			}
			goto Set
		}
		if err != nil {
			panic("unum2.go: " + err.Error())
		}
		if string(b) != str {
			panic(fmt.Sprintf("unum2.go: MarshalText: got: %s: want: %s", b, str))
		}
		var v Unum2
		if err := v.UnmarshalText(b); err != nil {
			panic("unum2.go: UnmarshalText: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("unum2.go: MarshalText: got: %s: want: %s", v, c))
		}
	}
Set:
	{
		var v Unum2
		err := v.Set(str)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("unum2.go: Set: expected an error for %s", c))
			}
			goto Invalid
		}
		if v != c {
			panic(fmt.Sprintf("unum2.go: Set: got: %s: want: %s", v, c))
		}
	}
Invalid:
	if invalid {
		var v Unum2
		if err := json.Unmarshal([]byte(str), &v); err == nil {
			panic("unum2.go: json.Unmarshal: expected an error for: " + str)
		}
		if err := v.UnmarshalText([]byte(str)); err == nil {
			panic("unum2.go: UnmarshalText: expected an error for: " + str)
		}
	}
}
