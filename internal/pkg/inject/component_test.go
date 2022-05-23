package inject

import (
	"fmt"
	"github.com/magiconair/properties"
	"io/ioutil"
	"testing"
)

type A struct {
	B *B `inject:"-"`
	C IC `inject:"impl2"`
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
func TestComponent(t *testing.T) {
	a := &A{}
	b := &B{}
	cc := &ICImpl1{}
	d := &ICImpl2{}

	c := new(Component).
		Load(a, "a").
		Load(b, "b").
		Load(cc, "impl1").
		Load(d, "impl2")

	c.Install()

	cc.M()
	fmt.Printf("%v,%v", a.B, b.B)
}

func TestProperties(t *testing.T) {
	pFile, err := ioutil.ReadFile("/home/shenweijie/test.properties")
	if err != nil {
		panic(err)
	}

	p, err := properties.Load(pFile, properties.UTF8)

	for _, k := range p.Keys() {
		v := p.GetString(k, "")
		fmt.Printf("%s = %s\n", k, v)
	}

}
