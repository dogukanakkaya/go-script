package internal

type Builtin struct {
	Name string
	Fn   func(args ...interface{}) interface{}
}

type Value interface{}

type Object map[string]Value

type Array []Value

type ReturnValue struct {
	Value Value
}
