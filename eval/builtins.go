package eval

import (
	"fmt"

	"github.com/pandeykartikey/goto/object"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return errorMessageToObject("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.List:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return errorMessageToObject("argument to `len` not supported, got %s", args[0].Type())
			}

		},
	},
	"append": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return errorMessageToObject("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.LIST_OBJ {
				return errorMessageToObject("argument to `append` must be LIST, got %s", args[0].Type())
			}
			list := args[0].(*object.List)
			list.Value = append(list.Value, args[1])
			return NULL
		},
	},
	"print": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}
