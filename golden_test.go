// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains simple golden tests for various examples.
// Besides validating the results when the implementation changes,
// it provides a way to look at the generated code without having
// to execute the print statements in one's head.

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charlievieth/go-enum/internal/testenv"
)

// Golden represents a test case.
type Golden struct {
	name        string
	trimPrefix  string
	lineComment bool
	input       string // input; the package clause is provided when running the test.
	output      string // expected output.
}

var golden = []Golden{
	{"day", "", false, day_in, day_out},
	{"offset", "", false, offset_in, offset_out},
	{"gap", "", false, gap_in, gap_out},
	{"num", "", false, num_in, num_out},
	{"unum", "", false, unum_in, unum_out},
	{"unumpos", "", false, unumpos_in, unumpos_out},
	{"prime", "", false, prime_in, prime_out},
	{"prefix", "Type", false, prefix_in, prefix_out},
	{"tokens", "", true, tokens_in, tokens_out},
}

// Each example starts with "type XXX [u]int", with a single space separating them.

// Simple test: enumeration of type int starting at 0.
const day_in = `type Day int
const (
	Monday Day = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)
`

const day_out = `
const _Day_name = "MondayTuesdayWednesdayThursdayFridaySaturdaySunday"

var _Day_index = [...]uint8{0, 6, 13, 22, 30, 36, 44, 50}

func (i Day) String() string {
	if i < 0 || i >= Day(len(_Day_index)-1) {
		return "Day(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Day_name[_Day_index[i]:_Day_index[i+1]]
}

func (i Day) Valid() bool {
	return !(i < 0 || i >= Day(len(_Day_index)-1))
}

func (i Day) MarshalText() ([]byte, error) {
	if i < 0 || i >= Day(len(_Day_index)-1) {
		return nil, errors.New("invalid Day: " + strconv.FormatInt(int64(i), 10))
	}
	return []byte(_Day_name[_Day_index[i]:_Day_index[i+1]]), nil
}

func (i *Day) Set(s string) (err error) {
	switch s {
	case _Day_name[0:6]:
		*i = Monday
	case _Day_name[6:13]:
		*i = Tuesday
	case _Day_name[13:22]:
		*i = Wednesday
	case _Day_name[22:30]:
		*i = Thursday
	case _Day_name[30:36]:
		*i = Friday
	case _Day_name[36:44]:
		*i = Saturday
	case _Day_name[44:50]:
		*i = Sunday
	default:
		if len(s) <= 32 {
			err = errors.New("malformed Day: " + string(s))
		} else {
			err = errors.New("malformed Day: " + string(s[0:29]) + "...")
		}
	}
	return err
}

func (i *Day) UnmarshalText(s []byte) (err error) {
	switch string(s) {
	case _Day_name[0:6]:
		*i = Monday
	case _Day_name[6:13]:
		*i = Tuesday
	case _Day_name[13:22]:
		*i = Wednesday
	case _Day_name[22:30]:
		*i = Thursday
	case _Day_name[30:36]:
		*i = Friday
	case _Day_name[36:44]:
		*i = Saturday
	case _Day_name[44:50]:
		*i = Sunday
	default:
		if len(s) <= 32 {
			err = errors.New("malformed Day: " + string(s))
		} else {
			err = errors.New("malformed Day: " + string(s[0:29]) + "...")
		}
	}
	return err
}
`

// Enumeration with an offset.
// Also includes a duplicate.
const offset_in = `type Number int
const (
	_ Number = iota
	One
	Two
	Three
	AnotherOne = One  // Duplicate; note that AnotherOne doesn't appear below.
)
`

const offset_out = `
const _Number_name = "OneTwoThree"

var _Number_index = [...]uint8{0, 3, 6, 11}

func (i Number) String() string {
	i -= 1
	if i < 0 || i >= Number(len(_Number_index)-1) {
		return "Number(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Number_name[_Number_index[i]:_Number_index[i+1]]
}

func (i Number) Valid() bool {
	i -= 1
	return !(i < 0 || i >= Number(len(_Number_index)-1))
}

func (i Number) MarshalText() ([]byte, error) {
	i -= 1
	if i < 0 || i >= Number(len(_Number_index)-1) {
		return nil, errors.New("invalid Number: " + strconv.FormatInt(int64(i+1), 10))
	}
	return []byte(_Number_name[_Number_index[i]:_Number_index[i+1]]), nil
}

func (i *Number) Set(s string) (err error) {
	switch s {
	case _Number_name[0:3]:
		*i = One
	case _Number_name[3:6]:
		*i = Two
	case _Number_name[6:11]:
		*i = Three
	default:
		if len(s) <= 32 {
			err = errors.New("malformed Number: " + string(s))
		} else {
			err = errors.New("malformed Number: " + string(s[0:29]) + "...")
		}
	}
	return err
}

func (i *Number) UnmarshalText(s []byte) (err error) {
	switch string(s) {
	case _Number_name[0:3]:
		*i = One
	case _Number_name[3:6]:
		*i = Two
	case _Number_name[6:11]:
		*i = Three
	default:
		if len(s) <= 32 {
			err = errors.New("malformed Number: " + string(s))
		} else {
			err = errors.New("malformed Number: " + string(s[0:29]) + "...")
		}
	}
	return err
}
`

// Gaps and an offset.
const gap_in = `type Gap int
const (
	Two Gap = 2
	Three Gap = 3
	Five Gap = 5
	Six Gap = 6
	Seven Gap = 7
	Eight Gap = 8
	Nine Gap = 9
	Eleven Gap = 11
)
`

const gap_out = `
const (
	_Gap_name_0 = "TwoThree"
	_Gap_name_1 = "FiveSixSevenEightNine"
	_Gap_name_2 = "Eleven"
)

var (
	_Gap_index_0 = [...]uint8{0, 3, 8}
	_Gap_index_1 = [...]uint8{0, 4, 7, 12, 17, 21}
)

func (i Gap) String() string {
	switch {
	case 2 <= i && i <= 3:
		i -= 2
		return _Gap_name_0[_Gap_index_0[i]:_Gap_index_0[i+1]]
	case 5 <= i && i <= 9:
		i -= 5
		return _Gap_name_1[_Gap_index_1[i]:_Gap_index_1[i+1]]
	case i == 11:
		return _Gap_name_2
	default:
		return "Gap(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}

func (i Gap) Valid() bool {
	switch {
	case 2 <= i && i <= 3:
	case 5 <= i && i <= 9:
	case i == 11:
	default:
		return false
	}
	return true
}

func (i Gap) MarshalText() ([]byte, error) {
	if i.Valid() {
		return []byte(i.String()), nil
	}
	return nil, errors.New("invalid Gap: " + strconv.FormatInt(int64(i), 10))
}

func (i *Gap) Set(s string) (err error) {
	switch s {
	case _Gap_name_0[0:3]:
		*i = Two
	case _Gap_name_0[3:8]:
		*i = Three
	case _Gap_name_1[0:4]:
		*i = Five
	case _Gap_name_1[4:7]:
		*i = Six
	case _Gap_name_1[7:12]:
		*i = Seven
	case _Gap_name_1[12:17]:
		*i = Eight
	case _Gap_name_1[17:21]:
		*i = Nine
	case _Gap_name_2:
		*i = Eleven
	default:
		if len(s) <= 32 {
			err = errors.New("malformed Gap: " + string(s))
		} else {
			err = errors.New("malformed Gap: " + string(s[0:29]) + "...")
		}
	}
	return err
}

func (i *Gap) UnmarshalText(s []byte) (err error) {
	switch string(s) {
	case _Gap_name_0[0:3]:
		*i = Two
	case _Gap_name_0[3:8]:
		*i = Three
	case _Gap_name_1[0:4]:
		*i = Five
	case _Gap_name_1[4:7]:
		*i = Six
	case _Gap_name_1[7:12]:
		*i = Seven
	case _Gap_name_1[12:17]:
		*i = Eight
	case _Gap_name_1[17:21]:
		*i = Nine
	case _Gap_name_2:
		*i = Eleven
	default:
		if len(s) <= 32 {
			err = errors.New("malformed Gap: " + string(s))
		} else {
			err = errors.New("malformed Gap: " + string(s[0:29]) + "...")
		}
	}
	return err
}
`

// Signed integers spanning zero.
const num_in = `type Num int
const (
	m_2 Num = -2 + iota
	m_1
	m0
	m1
	m2
)
`

const num_out = `
const _Num_name = "m_2m_1m0m1m2"

var _Num_index = [...]uint8{0, 3, 6, 8, 10, 12}

func (i Num) String() string {
	i -= -2
	if i < 0 || i >= Num(len(_Num_index)-1) {
		return "Num(" + strconv.FormatInt(int64(i+-2), 10) + ")"
	}
	return _Num_name[_Num_index[i]:_Num_index[i+1]]
}

func (i Num) Valid() bool {
	i -= -2
	return !(i < 0 || i >= Num(len(_Num_index)-1))
}

func (i Num) MarshalText() ([]byte, error) {
	i -= -2
	if i < 0 || i >= Num(len(_Num_index)-1) {
		return nil, errors.New("invalid Num: " + strconv.FormatInt(int64(i+-2), 10))
	}
	return []byte(_Num_name[_Num_index[i]:_Num_index[i+1]]), nil
}

func (i *Num) Set(s string) (err error) {
	switch s {
	case _Num_name[0:3]:
		*i = m_2
	case _Num_name[3:6]:
		*i = m_1
	case _Num_name[6:8]:
		*i = m0
	case _Num_name[8:10]:
		*i = m1
	case _Num_name[10:12]:
		*i = m2
	default:
		if len(s) <= 32 {
			err = errors.New("malformed Num: " + string(s))
		} else {
			err = errors.New("malformed Num: " + string(s[0:29]) + "...")
		}
	}
	return err
}

func (i *Num) UnmarshalText(s []byte) (err error) {
	switch string(s) {
	case _Num_name[0:3]:
		*i = m_2
	case _Num_name[3:6]:
		*i = m_1
	case _Num_name[6:8]:
		*i = m0
	case _Num_name[8:10]:
		*i = m1
	case _Num_name[10:12]:
		*i = m2
	default:
		if len(s) <= 32 {
			err = errors.New("malformed Num: " + string(s))
		} else {
			err = errors.New("malformed Num: " + string(s[0:29]) + "...")
		}
	}
	return err
}
`

// Unsigned integers spanning zero.
const unum_in = `type Unum uint
const (
	m_2 Unum = iota + 253
	m_1
)

const (
	m0 Unum = iota
	m1
	m2
)
`

const unum_out = `
const (
	_Unum_name_0 = "m0m1m2"
	_Unum_name_1 = "m_2m_1"
)

var (
	_Unum_index_0 = [...]uint8{0, 2, 4, 6}
	_Unum_index_1 = [...]uint8{0, 3, 6}
)

func (i Unum) String() string {
	switch {
	case i <= 2:
		return _Unum_name_0[_Unum_index_0[i]:_Unum_index_0[i+1]]
	case 253 <= i && i <= 254:
		i -= 253
		return _Unum_name_1[_Unum_index_1[i]:_Unum_index_1[i+1]]
	default:
		return "Unum(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}

func (i Unum) Valid() bool {
	switch {
	case i <= 2:
	case 253 <= i && i <= 254:
	default:
		return false
	}
	return true
}

func (i Unum) MarshalText() ([]byte, error) {
	if i.Valid() {
		return []byte(i.String()), nil
	}
	return nil, errors.New("invalid Unum: " + strconv.FormatInt(int64(i), 10))
}

func (i *Unum) Set(s string) (err error) {
	switch s {
	case _Unum_name_0[0:2]:
		*i = m0
	case _Unum_name_0[2:4]:
		*i = m1
	case _Unum_name_0[4:6]:
		*i = m2
	case _Unum_name_1[0:3]:
		*i = m_2
	case _Unum_name_1[3:6]:
		*i = m_1
	default:
		if len(s) <= 32 {
			err = errors.New("malformed Unum: " + string(s))
		} else {
			err = errors.New("malformed Unum: " + string(s[0:29]) + "...")
		}
	}
	return err
}

func (i *Unum) UnmarshalText(s []byte) (err error) {
	switch string(s) {
	case _Unum_name_0[0:2]:
		*i = m0
	case _Unum_name_0[2:4]:
		*i = m1
	case _Unum_name_0[4:6]:
		*i = m2
	case _Unum_name_1[0:3]:
		*i = m_2
	case _Unum_name_1[3:6]:
		*i = m_1
	default:
		if len(s) <= 32 {
			err = errors.New("malformed Unum: " + string(s))
		} else {
			err = errors.New("malformed Unum: " + string(s[0:29]) + "...")
		}
	}
	return err
}
`

// Unsigned positive integers.
const unumpos_in = `type Unumpos uint
const (
	m253 Unumpos = iota + 253
	m254
)

const (
	m1 Unumpos = iota + 1
	m2
	m3
)
`

const unumpos_out = `
const (
	_Unumpos_name_0 = "m1m2m3"
	_Unumpos_name_1 = "m253m254"
)

var (
	_Unumpos_index_0 = [...]uint8{0, 2, 4, 6}
	_Unumpos_index_1 = [...]uint8{0, 4, 8}
)

func (i Unumpos) String() string {
	switch {
	case 1 <= i && i <= 3:
		i -= 1
		return _Unumpos_name_0[_Unumpos_index_0[i]:_Unumpos_index_0[i+1]]
	case 253 <= i && i <= 254:
		i -= 253
		return _Unumpos_name_1[_Unumpos_index_1[i]:_Unumpos_index_1[i+1]]
	default:
		return "Unumpos(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}

func (i Unumpos) Valid() bool {
	switch {
	case 1 <= i && i <= 3:
	case 253 <= i && i <= 254:
	default:
		return false
	}
	return true
}

func (i Unumpos) MarshalText() ([]byte, error) {
	if i.Valid() {
		return []byte(i.String()), nil
	}
	return nil, errors.New("invalid Unumpos: " + strconv.FormatInt(int64(i), 10))
}

func (i *Unumpos) Set(s string) (err error) {
	switch s {
	case _Unumpos_name_0[0:2]:
		*i = m1
	case _Unumpos_name_0[2:4]:
		*i = m2
	case _Unumpos_name_0[4:6]:
		*i = m3
	case _Unumpos_name_1[0:4]:
		*i = m253
	case _Unumpos_name_1[4:8]:
		*i = m254
	default:
		if len(s) <= 32 {
			err = errors.New("malformed Unumpos: " + string(s))
		} else {
			err = errors.New("malformed Unumpos: " + string(s[0:29]) + "...")
		}
	}
	return err
}

func (i *Unumpos) UnmarshalText(s []byte) (err error) {
	switch string(s) {
	case _Unumpos_name_0[0:2]:
		*i = m1
	case _Unumpos_name_0[2:4]:
		*i = m2
	case _Unumpos_name_0[4:6]:
		*i = m3
	case _Unumpos_name_1[0:4]:
		*i = m253
	case _Unumpos_name_1[4:8]:
		*i = m254
	default:
		if len(s) <= 32 {
			err = errors.New("malformed Unumpos: " + string(s))
		} else {
			err = errors.New("malformed Unumpos: " + string(s[0:29]) + "...")
		}
	}
	return err
}
`

// Enough gaps to trigger a map implementation of the method.
// Also includes a duplicate to test that it doesn't cause problems
const prime_in = `type Prime int
const (
	p2 Prime = 2
	p3 Prime = 3
	p5 Prime = 5
	p7 Prime = 7
	// Skip duplicates: we don't support them
	// p77 Prime = 7 // Duplicate; note that p77 doesn't appear below.
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
`

const prime_out = `
const _Prime_name = "p2p3p5p7p11p13p17p19p23p29p37p41p43"

var _Prime_map = map[Prime]string{
	2:  _Prime_name[0:2],
	3:  _Prime_name[2:4],
	5:  _Prime_name[4:6],
	7:  _Prime_name[6:8],
	11: _Prime_name[8:11],
	13: _Prime_name[11:14],
	17: _Prime_name[14:17],
	19: _Prime_name[17:20],
	23: _Prime_name[20:23],
	29: _Prime_name[23:26],
	31: _Prime_name[26:29],
	41: _Prime_name[29:32],
	43: _Prime_name[32:35],
}

func (i Prime) String() string {
	if str, ok := _Prime_map[i]; ok {
		return str
	}
	return "Prime(" + strconv.FormatInt(int64(i), 10) + ")"
}

func (i Prime) Valid() bool {
	_, ok := _Prime_map[i]
	return ok
}

func (i Prime) MarshalText() ([]byte, error) {
	if str, ok := _Prime_map[i]; ok {
		return []byte(str), nil
	}
	return nil, errors.New("invalid Prime: " + strconv.FormatInt(int64(i), 10))
}

func (i *Prime) Set(s string) (err error) {
	switch s {
	case _Prime_name[0:2]:
		*i = p2
	case _Prime_name[2:4]:
		*i = p3
	case _Prime_name[4:6]:
		*i = p5
	case _Prime_name[6:8]:
		*i = p7
	case _Prime_name[8:11]:
		*i = p11
	case _Prime_name[11:14]:
		*i = p13
	case _Prime_name[14:17]:
		*i = p17
	case _Prime_name[17:20]:
		*i = p19
	case _Prime_name[20:23]:
		*i = p23
	case _Prime_name[23:26]:
		*i = p29
	case _Prime_name[26:29]:
		*i = p37
	case _Prime_name[29:32]:
		*i = p41
	case _Prime_name[32:35]:
		*i = p43
	default:
		if len(s) <= 32 {
			err = errors.New("malformed Prime: " + string(s))
		} else {
			err = errors.New("malformed Prime: " + string(s[0:29]) + "...")
		}
	}
	return err
}

func (i *Prime) UnmarshalText(s []byte) (err error) {
	switch string(s) {
	case _Prime_name[0:2]:
		*i = p2
	case _Prime_name[2:4]:
		*i = p3
	case _Prime_name[4:6]:
		*i = p5
	case _Prime_name[6:8]:
		*i = p7
	case _Prime_name[8:11]:
		*i = p11
	case _Prime_name[11:14]:
		*i = p13
	case _Prime_name[14:17]:
		*i = p17
	case _Prime_name[17:20]:
		*i = p19
	case _Prime_name[20:23]:
		*i = p23
	case _Prime_name[23:26]:
		*i = p29
	case _Prime_name[26:29]:
		*i = p37
	case _Prime_name[29:32]:
		*i = p41
	case _Prime_name[32:35]:
		*i = p43
	default:
		if len(s) <= 32 {
			err = errors.New("malformed Prime: " + string(s))
		} else {
			err = errors.New("malformed Prime: " + string(s[0:29]) + "...")
		}
	}
	return err
}
`

const prefix_in = `type Type int
const (
	TypeInt Type = iota
	TypeString
	TypeFloat
	TypeRune
	TypeByte
	TypeStruct
	TypeSlice
)
`

const prefix_out = `
const _Type_name = "IntStringFloatRuneByteStructSlice"

var _Type_index = [...]uint8{0, 3, 9, 14, 18, 22, 28, 33}

func (i Type) String() string {
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}

func (i Type) Valid() bool {
	return !(i < 0 || i >= Type(len(_Type_index)-1))
}

func (i Type) MarshalText() ([]byte, error) {
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return nil, errors.New("invalid Type: " + strconv.FormatInt(int64(i), 10))
	}
	return []byte(_Type_name[_Type_index[i]:_Type_index[i+1]]), nil
}

func (i *Type) Set(s string) (err error) {
	switch s {
	case _Type_name[0:3]:
		*i = TypeInt
	case _Type_name[3:9]:
		*i = TypeString
	case _Type_name[9:14]:
		*i = TypeFloat
	case _Type_name[14:18]:
		*i = TypeRune
	case _Type_name[18:22]:
		*i = TypeByte
	case _Type_name[22:28]:
		*i = TypeStruct
	case _Type_name[28:33]:
		*i = TypeSlice
	default:
		if len(s) <= 32 {
			err = errors.New("malformed Type: " + string(s))
		} else {
			err = errors.New("malformed Type: " + string(s[0:29]) + "...")
		}
	}
	return err
}

func (i *Type) UnmarshalText(s []byte) (err error) {
	switch string(s) {
	case _Type_name[0:3]:
		*i = TypeInt
	case _Type_name[3:9]:
		*i = TypeString
	case _Type_name[9:14]:
		*i = TypeFloat
	case _Type_name[14:18]:
		*i = TypeRune
	case _Type_name[18:22]:
		*i = TypeByte
	case _Type_name[22:28]:
		*i = TypeStruct
	case _Type_name[28:33]:
		*i = TypeSlice
	default:
		if len(s) <= 32 {
			err = errors.New("malformed Type: " + string(s))
		} else {
			err = errors.New("malformed Type: " + string(s[0:29]) + "...")
		}
	}
	return err
}
`

const tokens_in = `type Token int
const (
	And Token = iota // &
	Or               // |
	Add              // +
	Sub              // -
	Ident
	Period // .

	// not to be used
	SingleBefore
	// not to be used
	BeforeAndInline // inline
	InlineGeneral /* inline general */
)
`

const tokens_out = `
const _Token_name = "&|+-Ident.SingleBeforeinlineinline general"

var _Token_index = [...]uint8{0, 1, 2, 3, 4, 9, 10, 22, 28, 42}

func (i Token) String() string {
	if i < 0 || i >= Token(len(_Token_index)-1) {
		return "Token(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Token_name[_Token_index[i]:_Token_index[i+1]]
}

func (i Token) Valid() bool {
	return !(i < 0 || i >= Token(len(_Token_index)-1))
}

func (i Token) MarshalText() ([]byte, error) {
	if i < 0 || i >= Token(len(_Token_index)-1) {
		return nil, errors.New("invalid Token: " + strconv.FormatInt(int64(i), 10))
	}
	return []byte(_Token_name[_Token_index[i]:_Token_index[i+1]]), nil
}

func (i *Token) Set(s string) (err error) {
	switch s {
	case _Token_name[0:1]:
		*i = And
	case _Token_name[1:2]:
		*i = Or
	case _Token_name[2:3]:
		*i = Add
	case _Token_name[3:4]:
		*i = Sub
	case _Token_name[4:9]:
		*i = Ident
	case _Token_name[9:10]:
		*i = Period
	case _Token_name[10:22]:
		*i = SingleBefore
	case _Token_name[22:28]:
		*i = BeforeAndInline
	case _Token_name[28:42]:
		*i = InlineGeneral
	default:
		if len(s) <= 32 {
			err = errors.New("malformed Token: " + string(s))
		} else {
			err = errors.New("malformed Token: " + string(s[0:29]) + "...")
		}
	}
	return err
}

func (i *Token) UnmarshalText(s []byte) (err error) {
	switch string(s) {
	case _Token_name[0:1]:
		*i = And
	case _Token_name[1:2]:
		*i = Or
	case _Token_name[2:3]:
		*i = Add
	case _Token_name[3:4]:
		*i = Sub
	case _Token_name[4:9]:
		*i = Ident
	case _Token_name[9:10]:
		*i = Period
	case _Token_name[10:22]:
		*i = SingleBefore
	case _Token_name[22:28]:
		*i = BeforeAndInline
	case _Token_name[28:42]:
		*i = InlineGeneral
	default:
		if len(s) <= 32 {
			err = errors.New("malformed Token: " + string(s))
		} else {
			err = errors.New("malformed Token: " + string(s[0:29]) + "...")
		}
	}
	return err
}
`

func TestGolden(t *testing.T) {
	testenv.NeedsTool(t, "go")

	dir, err := ioutil.TempDir("", "stringer")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir)

	tmp, err := ioutil.TempDir("", "stringer-tests")
	if err != nil {
		t.Error(err)
	}

	var failed []string

	for _, test := range golden {
		g := Generator{
			trimPrefix:  test.trimPrefix,
			lineComment: test.lineComment,
		}
		input := "package test\n" + test.input
		file := test.name + ".go"
		absFile := filepath.Join(dir, file)
		err := ioutil.WriteFile(absFile, []byte(input), 0644)
		if err != nil {
			t.Error(err)
		}

		g.parsePackage([]string{absFile}, nil)
		// Extract the name and type of the constant from the first line.
		tokens := strings.SplitN(test.input, " ", 3)
		if len(tokens) != 3 {
			t.Fatalf("%s: need type declaration on first line", test.name)
		}
		g.generate(tokens[1])
		got := string(g.format())
		// ignore trailing whitespace (it's not important)
		got = strings.TrimRight(got, "\n")
		test.output = strings.TrimRight(test.output, "\n")
		if got != test.output {
			if testing.Verbose() {
				t.Errorf("%s: got(%d)\n====\n%s====\nexpected(%d)\n====\n%s====\n",
					test.name, len(got), got, len(test.output), test.output)
			} else {
				t.Error(test.name)
			}
			failed = append(failed, test.name)

			base := filepath.Join(tmp, test.name)
			if err := ioutil.WriteFile(base+".got", []byte(got), 0644); err != nil {
				t.Error(err)
			}
			if err := ioutil.WriteFile(base+".exp", []byte(test.output), 0644); err != nil {
				t.Error(err)
			}
		}
	}

	if t.Failed() {
		t.Errorf("\nTest failures:\n  %s\nFailure diffs are saved at: %s",
			strings.Join(failed, "\n  "), tmp)
	} else {
		os.RemoveAll(tmp)
	}
}
