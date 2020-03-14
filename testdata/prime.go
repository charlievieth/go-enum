// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Enough gaps to trigger a map implementation of the method.
// Also includes a duplicate to test that it doesn't cause problems

package main

import (
	"encoding/json"
	"fmt"
)

type Prime int

const (
	p2 Prime = 2
	p3 Prime = 3
	p5 Prime = 5
	p7 Prime = 7

	// Skip duplicates: we don't support them
	// xp77 Prime = 7 // Duplicate; note that p77 doesn't appear below.

	p11 Prime = 11
	p13 Prime = 13
	p17 Prime = 17
	p19 Prime = 19
	p23 Prime = 23
	p29 Prime = 29
	p37 Prime = 31
	p41 Prime = 41
	p43 Prime = 43
)

func main() {
	ck(0, "Prime(0)", true)
	ck(1, "Prime(1)", true)
	ck(p2, "p2", false)
	ck(p3, "p3", false)
	ck(4, "Prime(4)", true)
	ck(p5, "p5", false)
	ck(p7, "p7", false)

	// Skip duplicates: we don't support them
	// ck(p77, "p7")

	ck(p11, "p11", false)
	ck(p13, "p13", false)
	ck(p17, "p17", false)
	ck(p19, "p19", false)
	ck(p23, "p23", false)
	ck(p29, "p29", false)
	ck(p37, "p37", false)
	ck(p41, "p41", false)
	ck(p43, "p43", false)
	ck(44, "Prime(44)", true)
}

func ck(c Prime, str string, invalid bool) {
	if fmt.Sprint(c) != str {
		panic("prime.go: " + str)
	}
	{
		b, err := json.Marshal(c)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("prime.go: json.Marshal: expected an error for %s", c))
			}
			goto MarshalText
		}
		if err != nil {
			panic("prime.go: " + err.Error())
		}
		if string(b) != `"`+str+`"` {
			panic(fmt.Sprintf("prime.go: json.Marshal: got: %s: want: %q", b, str))
		}
		var v Prime
		if err := json.Unmarshal(b, &v); err != nil {
			panic("prime.go: json.Unmarshal: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("prime.go: json.Marshal: got: %s: want: %s", v, c))
		}
	}
MarshalText:
	{
		b, err := c.MarshalText()
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("prime.go: MarshalText: expected an error for %s", c))
			}
			goto Set
		}
		if err != nil {
			panic("prime.go: " + err.Error())
		}
		if string(b) != str {
			panic(fmt.Sprintf("prime.go: MarshalText: got: %s: want: %s", b, str))
		}
		var v Prime
		if err := v.UnmarshalText(b); err != nil {
			panic("prime.go: UnmarshalText: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("prime.go: MarshalText: got: %s: want: %s", v, c))
		}
	}
Set:
	{
		var v Prime
		err := v.Set(str)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("prime.go: Set: expected an error for %s", c))
			}
			goto Invalid
		}
		if v != c {
			panic(fmt.Sprintf("prime.go: Set: got: %s: want: %s", v, c))
		}
	}
Invalid:
	if invalid {
		var v Prime
		if err := json.Unmarshal([]byte(str), &v); err == nil {
			panic("prime.go: json.Unmarshal: expected an error for: " + str)
		}
		if err := v.UnmarshalText([]byte(str)); err == nil {
			panic("prime.go: UnmarshalText: expected an error for: " + str)
		}
	}
}
