package inject

import (
	"fmt"
	"reflect"
)

const (
	Tag     = "inject"
	TagNone = "-"
)

type Instance struct {
	ID  string
	Obj interface{}
}

type Component struct {
	instances []*Instance
}

func (c *Component) Install() {
	for _, com := range c.instances {
		c.initializeInstance(com)
	}
}

func (c *Component) Load(instance interface{}, id ...string) *Component {
	if !isStructPointer(reflect.TypeOf(instance)) {
		panic("instance must be a pointer type")
	}

	ins := &Instance{Obj: instance}
	if len(id) > 0 {
		ins.ID = id[0]
	}

	c.instances = append(c.instances, ins)
	return c
}

func (c *Component) initializeInstance(instance *Instance) {
	typeOf := reflect.TypeOf(instance.Obj).Elem()
	for i := 0; i < typeOf.NumField(); i++ {
		fieldType := typeOf.Field(i)
		if tag := fieldType.Tag.Get("inject"); tag != "" {
			if !fieldType.IsExported() {
				panic("component field must export")
			}

			fieldValue := reflect.ValueOf(instance.Obj).Elem().Field(i)
			fieldInstance := c.findInstanceByType(fieldType.Type)
			if len(fieldInstance) == 0 {
				panic(fmt.Sprintf("no suitable instance found for type %s", fieldType.Name))
			}

			if len(fieldInstance) > 1 {
				if TagNone == tag {
					panic(fmt.Sprintf("one more instance exists for type %s", fieldType.Name))
				}

				var byID []*Instance
				for _, i := range fieldInstance {
					if i.ID == tag {
						byID = append(byID, i)
					}
				}

				if len(byID) == 0 {
					panic(fmt.Sprintf("no instance found for id %s", tag))
				}

				if len(byID) > 1 {
					panic(fmt.Sprintf("find one more instance for id %s", tag))
				}

				fieldValue.Set(reflect.ValueOf(byID[0].Obj))
			}

			fieldValue.Set(reflect.ValueOf(fieldInstance[0].Obj))
		}
	}
}

func (c *Component) findInstanceByType(t reflect.Type) []*Instance {
	var result []*Instance

	matcher := func(instanceType, t reflect.Type) bool {
		if isStructPointer(t) {
			return instanceType == t
		}

		if isInterface(t) {
			return instanceType.Implements(t)
		}

		return false
	}

	for _, i := range c.instances {
		if matcher(reflect.TypeOf(i.Obj), t) {
			result = append(result, i)
		}
	}

	return result
}

func isStructPointer(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

func isInterface(t reflect.Type) bool {
	return t.Kind() == reflect.Interface
}
