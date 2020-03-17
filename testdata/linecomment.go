package main

import (
	"encoding/json"
	"fmt"
)

type Linecomment int

const (
	Backslash               Linecomment = 1   // \B\a\c\k\s\l\a\s\h
	Tab                     Linecomment = 3   // Tab	Tab
	BadEscape               Linecomment = 5   // \b\f\n\r\t\u
	Quoted                  Linecomment = 7   // "Quoted"
	SingleQuote             Linecomment = 9   // 'SingleQuote'
	Backtick                Linecomment = 11  // `Backtick`
	EscapedQuote            Linecomment = 13  // \"EscapedQuote\"
	Null                    Linecomment = 15  // Null \u0000
	NullEscaped             Linecomment = 17  // NullEscaped \\u0000
	ContinuationByte        Linecomment = 20  // ContinuationByte \u0080
	ContinuationByteEscaped Linecomment = 22  // ContinuationByteEscaped \u0080
	UTF8                    Linecomment = 100 // 日a本b語ç日ð本Ê語þ日¥本¼語i日©
)

func main() {
	ck(Backslash, "\\B\\a\\c\\k\\s\\l\\a\\s\\h", false)
	ck(Tab, "Tab\tTab", false)
	ck(BadEscape, "\\b\\f\\n\\r\\t\\u", false)
	ck(Quoted, "\"Quoted\"", false)
	ck(SingleQuote, "'SingleQuote'", false)
	ck(Backtick, "`Backtick`", false)
	ck(EscapedQuote, "\\\"EscapedQuote\\\"", false)
	ck(Null, "Null \\u0000", false)
	ck(NullEscaped, "NullEscaped \\\\u0000", false)
	ck(ContinuationByte, "ContinuationByte \\u0080", false)
	ck(ContinuationByteEscaped, "ContinuationByteEscaped \\u0080", false)
	ck(UTF8, "日a本b語ç日ð本Ê語þ日¥本¼語i日©", false)
}

func ck(c Linecomment, str string, invalid bool) {
	if fmt.Sprint(c) != str {
		panic(fmt.Sprintf("linecomment.go: got: %s want: %s", c, str))
	}
	{
		b, err := json.Marshal(c)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("linecomment.go: json.Marshal: expected an error for %s", c))
			}
			goto MarshalText
		}
		if err != nil {
			panic("linecomment.go: " + err.Error())
		}
		// Unlike other tests the JSON encoded version of str != "str".
		exp, err := json.Marshal(c.String())
		if err != nil {
			panic(fmt.Sprintf("linecomment.go: json.Marshal: expected an error for %s", c))
		}
		if string(b) != string(exp) {
			panic(fmt.Sprintf("linecomment.go: json.Marshal: got: %s: want: %q", b, str))
		}
		var v Linecomment
		if err := json.Unmarshal(b, &v); err != nil {
			panic("linecomment.go: json.Unmarshal: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("linecomment.go: json.Marshal: got: %s: want: %s", v, c))
		}
	}
MarshalText:
	{
		b, err := c.MarshalText()
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("linecomment.go: MarshalText: expected an error for %s", c))
			}
			goto Set
		}
		if err != nil {
			panic("linecomment.go: " + err.Error())
		}
		if string(b) != str {
			panic(fmt.Sprintf("linecomment.go: MarshalText: got: %s: want: %s", b, str))
		}
		var v Linecomment
		if err := v.UnmarshalText(b); err != nil {
			panic("linecomment.go: UnmarshalText: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("linecomment.go: MarshalText: got: %s: want: %s", v, c))
		}
	}
Set:
	{
		var v Linecomment
		err := v.Set(str)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("linecomment.go: Set: expected an error for %s", c))
			}
			goto Invalid
		}
		if v != c {
			panic(fmt.Sprintf("linecomment.go: Set: got: %s: want: %s", v, c))
		}
	}
Invalid:
	if invalid {
		var v Linecomment
		if err := json.Unmarshal([]byte(str), &v); err == nil {
			panic("linecomment.go: json.Unmarshal: expected an error for: " + str)
		}
		if err := v.UnmarshalText([]byte(str)); err == nil {
			panic("linecomment.go: UnmarshalText: expected an error for: " + str)
		}
	}
}
