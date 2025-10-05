package internal

type Builtin struct {
	Name string
	Fn   func(args ...interface{}) interface{}
}

type Value interface{}

type ReturnValue struct {
	Value Value
}

type Object map[string]Value

type Array []Value

type ArrayLike interface {
	GetElements() Array
}
