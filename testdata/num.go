// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Signed integers spanning zero.

package main

import (
	"encoding/json"
	"fmt"
)

type Num int

const (
	m_2 Num = -2 + iota
	m_1
	m0
	m1
	m2
)

func main() {
	ck(-3, "Num(-3)", true)
	ck(m_2, "m_2", false)
	ck(m_1, "m_1", false)
	ck(m0, "m0", false)
	ck(m1, "m1", false)
	ck(m2, "m2", false)
	ck(3, "Num(3)", true)
}

func ck(c Num, str string, invalid bool) {
	if fmt.Sprint(c) != str {
		panic("num.go: " + str)
	}
	{
		b, err := json.Marshal(c)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("num.go: json.Marshal: expected an error for %s", c))
			}
			goto MarshalText
		}
		if err != nil {
			panic("num.go: " + err.Error())
		}
		if string(b) != `"`+str+`"` {
			panic(fmt.Sprintf("num.go: json.Marshal: got: %s: want: %q", b, str))
		}
		var v Num
		if err := json.Unmarshal(b, &v); err != nil {
			panic("num.go: json.Unmarshal: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("num.go: json.Marshal: got: %s: want: %s", v, c))
		}
	}
MarshalText:
	{
		b, err := c.MarshalText()
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("num.go: MarshalText: expected an error for %s", c))
			}
			goto Set
		}
		if err != nil {
			panic("num.go: " + err.Error())
		}
		if string(b) != str {
			panic(fmt.Sprintf("num.go: MarshalText: got: %s: want: %s", b, str))
		}
		var v Num
		if err := v.UnmarshalText(b); err != nil {
			panic("num.go: UnmarshalText: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("num.go: MarshalText: got: %s: want: %s", v, c))
		}
	}
Set:
	{
		var v Num
		err := v.Set(str)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("num.go: Set: expected an error for %s", c))
			}
			goto Invalid
		}
		if v != c {
			panic(fmt.Sprintf("num.go: Set: got: %s: want: %s", v, c))
		}
	}
Invalid:
	if invalid {
		var v Num
		if err := json.Unmarshal([]byte(str), &v); err == nil {
			panic("num.go: json.Unmarshal: expected an error for: " + str)
		}
		if err := v.UnmarshalText([]byte(str)); err == nil {
			panic("num.go: UnmarshalText: expected an error for: " + str)
		}
	}
}
