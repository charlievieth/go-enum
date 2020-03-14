// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Import "C" shouldn't be imported.

package main

/*
#define HELLO 1
*/
import "C"

import (
	"encoding/json"
	"fmt"
)

type Cgo uint32

const (
	// MustScanSubDirs indicates that events were coalesced hierarchically.
	MustScanSubDirs Cgo = 1 << iota
)

func main() {
	_ = C.HELLO
	ck(MustScanSubDirs, "MustScanSubDirs", false)
}

func ck(c Cgo, str string, invalid bool) {
	if fmt.Sprint(c) != str {
		panic("cgo.go: " + str)
	}
	{
		b, err := json.Marshal(c)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("cgo.go: json.Marshal: expected an error for %s", c))
			}
			goto MarshalText
		}
		if err != nil {
			panic("cgo.go: " + err.Error())
		}
		if string(b) != `"`+str+`"` {
			panic(fmt.Sprintf("cgo.go: json.Marshal: got: %s: want: %q", b, str))
		}
		var v Cgo
		if err := json.Unmarshal(b, &v); err != nil {
			panic("cgo.go: json.Unmarshal: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("cgo.go: json.Marshal: got: %s: want: %s", v, c))
		}
	}
MarshalText:
	{
		b, err := c.MarshalText()
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("cgo.go: MarshalText: expected an error for %s", c))
			}
			goto Set
		}
		if err != nil {
			panic("cgo.go: " + err.Error())
		}
		if string(b) != str {
			panic(fmt.Sprintf("cgo.go: MarshalText: got: %s: want: %s", b, str))
		}
		var v Cgo
		if err := v.UnmarshalText(b); err != nil {
			panic("cgo.go: UnmarshalText: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("cgo.go: MarshalText: got: %s: want: %s", v, c))
		}
	}
Set:
	{
		var v Cgo
		err := v.Set(str)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("cgo.go: Set: expected an error for %s", c))
			}
			goto Invalid
		}
		if v != c {
			panic(fmt.Sprintf("cgo.go: Set: got: %s: want: %s", v, c))
		}
	}
Invalid:
	if invalid {
		var v Cgo
		if err := json.Unmarshal([]byte(str), &v); err == nil {
			panic("cgo.go: json.Unmarshal: expected an error for: " + str)
		}
		if err := v.UnmarshalText([]byte(str)); err == nil {
			panic("cgo.go: UnmarshalText: expected an error for: " + str)
		}
	}
}
