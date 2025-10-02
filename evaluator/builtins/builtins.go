package builtins

import (
	"go-script/evaluator/builtins/fetch"
	"go-script/evaluator/builtins/json"
	"go-script/evaluator/builtins/print"
)

type Builtin struct {
	Name string
	Fn   func(args ...interface{}) interface{}
}

var builtins = map[string]*Builtin{
	"print": {Name: print.Print.Name, Fn: print.Print.Fn},
	"fetch": {Name: fetch.Fetch.Name, Fn: fetch.Fetch.Fn},
}

var jsonNamespace = make(map[string]*Builtin)

func init() {
	for key, builtin := range json.JSON {
		jsonNamespace[key] = &Builtin{Name: builtin.Name, Fn: builtin.Fn}
	}
}

func Get(name string) (*Builtin, bool) {
	builtin, ok := builtins[name]
	return builtin, ok
}

func GetJSON() map[string]*Builtin {
	return jsonNamespace
}
