// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file is the same as day.go except the constants have different values.

package main

import (
	"encoding/json"
	"fmt"
)

type Day int

const (
	Sunday Day = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func main() {
	ck(Monday, "Monday", false)
	ck(Tuesday, "Tuesday", false)
	ck(Wednesday, "Wednesday", false)
	ck(Thursday, "Thursday", false)
	ck(Friday, "Friday", false)
	ck(Saturday, "Saturday", false)
	ck(Sunday, "Sunday", false)
	ck(-127, "Day(-127)", true)
	ck(127, "Day(127)", true)
}

func ck(c Day, str string, invalid bool) {
	if fmt.Sprint(c) != str {
		panic("vary_day.go: " + str)
	}
	{
		b, err := json.Marshal(c)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("vary_day.go: json.Marshal: expected an error for %s", c))
			}
			goto MarshalText
		}
		if err != nil {
			panic("vary_day.go: " + err.Error())
		}
		if string(b) != `"`+str+`"` {
			panic(fmt.Sprintf("vary_day.go: json.Marshal: got: %s: want: %q", b, str))
		}
		var v Day
		if err := json.Unmarshal(b, &v); err != nil {
			panic("vary_day.go: json.Unmarshal: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("vary_day.go: json.Marshal: got: %s: want: %s", v, c))
		}
	}
MarshalText:
	{
		b, err := c.MarshalText()
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("vary_day.go: MarshalText: expected an error for %s", c))
			}
			goto Set
		}
		if err != nil {
			panic("vary_day.go: " + err.Error())
		}
		if string(b) != str {
			panic(fmt.Sprintf("vary_day.go: MarshalText: got: %s: want: %s", b, str))
		}
		var v Day
		if err := v.UnmarshalText(b); err != nil {
			panic("vary_day.go: UnmarshalText: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("vary_day.go: MarshalText: got: %s: want: %s", v, c))
		}
	}
Set:
	{
		var v Day
		err := v.Set(str)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("vary_day.go: Set: expected an error for %s", c))
			}
			goto Invalid
		}
		if v != c {
			panic(fmt.Sprintf("vary_day.go: Set: got: %s: want: %s", v, c))
		}
	}
Invalid:
	if invalid {
		var v Day
		if err := json.Unmarshal([]byte(str), &v); err == nil {
			panic("vary_day.go: json.Unmarshal: expected an error for: " + str)
		}
		if err := v.UnmarshalText([]byte(str)); err == nil {
			panic("vary_day.go: UnmarshalText: expected an error for: " + str)
		}
	}
}
