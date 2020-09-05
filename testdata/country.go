package main

import (
	"encoding/json"
	"fmt"
)

type Country int

const (
	China  Country = iota // 中国	\China
	India                 // Bhārat /India
	Russia                // Росси́я	Russia
)

func main() {
	ck(China, "中国\t\\China", false)
	ck(India, "Bhārat /India", false)
	ck(Russia, "Росси́я\tRussia", false)
	ck(42, "Country(42)", true)
}

func ck(c Country, str string, invalid bool) {
	if s := fmt.Sprint(c); s != str {
		panic(fmt.Sprintf("country.go: want: %q got: %q", str, s))
	}
	{
		b, err := json.Marshal(c)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("country.go: json.Marshal: expected an error for %s", c))
			}
			goto MarshalText
		}
		if err != nil {
			panic("country.go: " + err.Error())
		}
		exp, err := json.Marshal(str)
		if err != nil {
			panic("country.go: " + err.Error())
		}
		if string(b) != string(exp) {
			panic(fmt.Sprintf("country.go: json.Marshal: got: %s: want: %s", b, exp))
		}
		var v Country
		if err := json.Unmarshal(b, &v); err != nil {
			panic("country.go: json.Unmarshal: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("country.go: json.Marshal: got: %s: want: %s", v, c))
		}
	}
MarshalText:
	{
		b, err := c.MarshalText()
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("country.go: MarshalText: expected an error for %s", c))
			}
			goto Set
		}
		if err != nil {
			panic("country.go: " + err.Error())
		}
		if string(b) != str {
			panic(fmt.Sprintf("country.go: MarshalText: got: %s: want: %s", b, str))
		}
		var v Country
		if err := v.UnmarshalText(b); err != nil {
			panic("country.go: UnmarshalText: " + err.Error())
		}
		if v != c {
			panic(fmt.Sprintf("country.go: MarshalText: got: %s: want: %s", v, c))
		}
	}
Set:
	{
		var v Country
		err := v.Set(str)
		if invalid {
			if err == nil {
				panic(fmt.Sprintf("country.go: Set: expected an error for %s", c))
			}
			goto Invalid
		}
		if v != c {
			panic(fmt.Sprintf("country.go: Set: got: %s: want: %s", v, c))
		}
	}
Invalid:
	if invalid {
		var v Country
		if err := json.Unmarshal([]byte(str), &v); err == nil {
			panic("country.go: json.Unmarshal: expected an error for: " + str)
		}
		if err := v.UnmarshalText([]byte(str)); err == nil {
			panic("country.go: UnmarshalText: expected an error for: " + str)
		}
	}
}
