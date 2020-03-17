package main

import (
	"encoding"
	"fmt"
)

type Interface int

const One Interface = 0

type ValueInterface interface {
	fmt.Stringer
	encoding.TextMarshaler
	Valid() bool
}

type PointerInterface interface {
	encoding.TextUnmarshaler
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
