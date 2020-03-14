// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Check that constants defined as a conversion are accepted.

package main

import (
	"encoding/json"
	"fmt"
)

type Other int // Imagine this is in another package.

const (
	alpha Other = iota
	beta
	gamma
	delta
)

type Conv int

const (
	Alpha = Conv(alpha)
	Beta  = Conv(beta)
	Gamma = Conv(gamma)
	Delta = Conv(delta)
)

func main() {
	ck(Alpha, "Alpha", false)
	ck(Beta, "Beta", false)
	ck(Gamma, "Gamma", false)
	ck(Delta, "Delta", false)
	ck(42, "Conv(42)", true)
}

func ck(c Conv, str string, invalid bool) {
	if fmt.Sprint(c) != str {
		panic("conv.go: " + str)
	}
	{
		b, err := json.Marshal(c)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("conv.go: json.Marshal: expected an error for %s", c))
			}
			goto MarshalText
		}
		if err != nil {
			panic("conv.go: " + err.Error())
		}
		if string(b) != `"`+str+`"` {
			panic(fmt.Sprintf("conv.go: json.Marshal: got: %s: want: %q", b, str))
		}
		var v Conv
		if err := json.Unmarshal(b, &v); err != nil {
			panic("conv.go: json.Unmarshal: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("conv.go: json.Marshal: got: %s: want: %s", v, c))
		}
	}
MarshalText:
	{
		b, err := c.MarshalText()
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("conv.go: MarshalText: expected an error for %s", c))
			}
			goto Set
		}
		if err != nil {
			panic("conv.go: " + err.Error())
		}
		if string(b) != str {
			panic(fmt.Sprintf("conv.go: MarshalText: got: %s: want: %s", b, str))
		}
		var v Conv
		if err := v.UnmarshalText(b); err != nil {
			panic("conv.go: UnmarshalText: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("conv.go: MarshalText: got: %s: want: %s", v, c))
		}
	}
Set:
	{
		var v Conv
		err := v.Set(str)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("conv.go: Set: expected an error for %s", c))
			}
			goto Invalid
		}
		if v != c {
			panic(fmt.Sprintf("conv.go: Set: got: %s: want: %s", v, c))
		}
	}
Invalid:
	if invalid {
		var v Conv
		if err := json.Unmarshal([]byte(str), &v); err == nil {
			panic("conv.go: json.Unmarshal: expected an error for: " + str)
		}
		if err := v.UnmarshalText([]byte(str)); err == nil {
			panic("conv.go: UnmarshalText: expected an error for: " + str)
		}
	}
}
