package main

import "fmt"

type A struct {
	B *B `inject:"-"`     // no id set, use type autoWare
	C IC `inject:"impl2"` // multiple implements, use id 'impl2' autoWare
}

type B struct {
	B *A `inject:"-"`
	C IC `inject:"impl1"`
}

type ICImpl1 struct {
}

func (c *ICImpl1) M() {
	fmt.Printf("ICImpl1: M")
}

type IC interface {
	M()
}

type ICImpl2 struct {
}

func (d *ICImpl2) M() {
	fmt.Printf("ICImpl2: M")
}
