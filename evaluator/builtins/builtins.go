package builtins

import (
	"go-script/evaluator/builtins/fetch"
	"go-script/evaluator/builtins/json"
	"go-script/evaluator/builtins/print"
	"go-script/internal"
)

var builtins = map[string]*internal.Builtin{
	"print": {Name: print.Print.Name, Fn: print.Print.Fn},
	"fetch": {Name: fetch.Fetch.Name, Fn: fetch.Fetch.Fn},
}

var jsonNamespace = make(map[string]*internal.Builtin)

func init() {
	for key, builtin := range json.JSON {
		jsonNamespace[key] = &internal.Builtin{Name: builtin.Name, Fn: builtin.Fn}
	}
}

func Get(name string) (*internal.Builtin, bool) {
	builtin, ok := builtins[name]
	return builtin, ok
}

func GetJSON() map[string]*internal.Builtin {
	return jsonNamespace
}
