// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Enumeration with an offset.
// Also includes a duplicate.

package main

import (
	"encoding/json"
	"fmt"
)

type Number int

const (
	_ Number = iota
	One
	Two
	Three
	AnotherOne = One // Duplicate; note that AnotherOne doesn't appear below.
)

func main() {
	ck(One, "One", false)
	ck(Two, "Two", false)
	ck(Three, "Three", false)
	ck(AnotherOne, "One", false)
	ck(127, "Number(127)", true)
}

func ck(c Number, str string, invalid bool) {
	if fmt.Sprint(c) != str {
		panic("number.go: " + str)
	}
	{
		b, err := json.Marshal(c)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("number.go: json.Marshal: expected an error for %s", c))
			}
			goto MarshalText
		}
		if err != nil {
			panic("number.go: " + err.Error())
		}
		if string(b) != `"`+str+`"` {
			panic(fmt.Sprintf("number.go: json.Marshal: got: %s: want: %q", b, str))
		}
		var v Number
		if err := json.Unmarshal(b, &v); err != nil {
			panic("number.go: json.Unmarshal: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("number.go: json.Marshal: got: %s: want: %s", v, c))
		}
	}
MarshalText:
	{
		b, err := c.MarshalText()
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("number.go: MarshalText: expected an error for %s", c))
			}
			goto Set
		}
		if err != nil {
			panic("number.go: " + err.Error())
		}
		if string(b) != str {
			panic(fmt.Sprintf("number.go: MarshalText: got: %s: want: %s", b, str))
		}
		var v Number
		if err := v.UnmarshalText(b); err != nil {
			panic("number.go: UnmarshalText: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("number.go: MarshalText: got: %s: want: %s", v, c))
		}
	}
Set:
	{
		var v Number
		err := v.Set(str)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("number.go: Set: expected an error for %s", c))
			}
			goto Invalid
		}
		if v != c {
			panic(fmt.Sprintf("number.go: Set: got: %s: want: %s", v, c))
		}
	}
Invalid:
	if invalid {
		var v Number
		if err := json.Unmarshal([]byte(str), &v); err == nil {
			panic("number.go: json.Unmarshal: expected an error for: " + str)
		}
		if err := v.UnmarshalText([]byte(str)); err == nil {
			panic("number.go: UnmarshalText: expected an error for: " + str)
		}
	}
}
