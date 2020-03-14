// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Unsigned integers spanning zero.

package main

import (
	"encoding/json"
	"fmt"
)

type Unum uint8

const (
	m_2 Unum = iota + 253
	m_1
)

const (
	m0 Unum = iota
	m1
	m2
)

func main() {
	ck(^Unum(0)-3, "Unum(252)", true)
	ck(m_2, "m_2", false)
	ck(m_1, "m_1", false)
	ck(m0, "m0", false)
	ck(m1, "m1", false)
	ck(m2, "m2", false)
	ck(3, "Unum(3)", true)
}

func ck(c Unum, str string, invalid bool) {
	if fmt.Sprint(c) != str {
		panic("unum.go: " + str)
	}
	{
		b, err := json.Marshal(c)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("unum.go: json.Marshal: expected an error for %s", c))
			}
			goto MarshalText
		}
		if err != nil {
			panic("unum.go: " + err.Error())
		}
		if string(b) != `"`+str+`"` {
			panic(fmt.Sprintf("unum.go: json.Marshal: got: %s: want: %q", b, str))
		}
		var v Unum
		if err := json.Unmarshal(b, &v); err != nil {
			panic("unum.go: json.Unmarshal: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("unum.go: json.Marshal: got: %s: want: %s", v, c))
		}
	}
MarshalText:
	{
		b, err := c.MarshalText()
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("unum.go: MarshalText: expected an error for %s", c))
			}
			goto Set
		}
		if err != nil {
			panic("unum.go: " + err.Error())
		}
		if string(b) != str {
			panic(fmt.Sprintf("unum.go: MarshalText: got: %s: want: %s", b, str))
		}
		var v Unum
		if err := v.UnmarshalText(b); err != nil {
			panic("unum.go: UnmarshalText: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("unum.go: MarshalText: got: %s: want: %s", v, c))
		}
	}
Set:
	{
		var v Unum
		err := v.Set(str)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("unum.go: Set: expected an error for %s", c))
			}
			goto Invalid
		}
		if v != c {
			panic(fmt.Sprintf("unum.go: Set: got: %s: want: %s", v, c))
		}
	}
Invalid:
	if invalid {
		var v Unum
		if err := json.Unmarshal([]byte(str), &v); err == nil {
			panic("unum.go: json.Unmarshal: expected an error for: " + str)
		}
		if err := v.UnmarshalText([]byte(str)); err == nil {
			panic("unum.go: UnmarshalText: expected an error for: " + str)
		}
	}
}
