package main

import (
	"encoding"
	"encoding/json"
	"fmt"
)

type Interface int

const One Interface = 0

type ValueInterface interface {
	fmt.Stringer
	encoding.TextMarshaler
	json.Marshaler
	Valid() bool
}

type PointerInterface interface {
	encoding.TextUnmarshaler
	json.Unmarshaler
	Set(string) error
}

type ExpectedInterface interface {
	ValueInterface
	PointerInterface
}

var (
	_ ValueInterface    = Interface(0)
	_ PointerInterface  = (*Interface)(nil)
	_ ExpectedInterface = (*Interface)(nil)
)

func main() {
}
