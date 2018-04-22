package eval

import (
	"fmt"
	"os"
	"strings"

	"github.com/waits/tang/object"
)

const (
	F_DEFAULT = iota
	F_PERCENT // %
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *object.List:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s",
					args[0].Type())
			}
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
	"format": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) == 0 || args[0].Type() != object.STRING {
				return newError("first argument to `format` must be STRING")
			}

			format, args := args[0].(*object.String).Value, args[1:]
			state := F_DEFAULT
			var b strings.Builder

			for i := 0; i < len(format); i++ {
				switch format[i] {
				case '%':
					if state == F_PERCENT {
						b.WriteRune('%')
						state = F_DEFAULT
						break
					}
					state = F_PERCENT
				case 'v':
					if state == F_PERCENT {
						val := args[0]
						args = args[1:]
						b.WriteString(val.Inspect())
						state = F_DEFAULT
						break
					}
					b.WriteRune('v')
				default:
					if state == F_PERCENT {
						return newError("invalid format verb `%%%s`", string(format[i]))
					}
					b.WriteByte(format[i])
				}
			}

			return &object.String{Value: b.String()}
		},
	},
	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != object.LIST {
				return newError("argument to `first` must be LIST, got %s",
					args[0].Type())
			}

			list := args[0].(*object.List)
			if len(list.Elements) > 0 {
				return list.Elements[0]
			}

			return NULL
		},
	},
	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != object.LIST {
				return newError("argument to `last` must be LIST, got %s",
					args[0].Type())
			}

			list := args[0].(*object.List)
			length := len(list.Elements)

			if len(list.Elements) > 0 {
				return list.Elements[length-1]
			}

			return NULL
		},
	},
	"rest": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != object.LIST {
				return newError("argument to `rest` must be LIST, got %s",
					args[0].Type())
			}

			list := args[0].(*object.List)
			length := len(list.Elements)

			if length > 0 {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, list.Elements[1:length])
				return &object.List{Elements: newElements}
			}
			if len(list.Elements) > 0 {
				return list.Elements[length-1]
			}

			return NULL
		},
	},
	"append": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2",
					len(args))
			}
			if args[0].Type() != object.LIST {
				return newError("argument to `append` must be LIST, got %s",
					args[0].Type())
			}

			list := args[0].(*object.List)
			length := len(list.Elements)

			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, list.Elements)
			newElements[length] = args[1]

			return &object.List{Elements: newElements}
		},
	},
	"exit": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != object.INTEGER {
				return newError("argument to `exit` must be INTEGER, got %s",
					args[0].Type())
			}

			arg := args[0].(*object.Integer)
			code := int(arg.Value)
			os.Exit(code)

			return NULL
		},
	},
	"panic": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			fmt.Printf("PANIC: %s\n", args[0].Inspect())
			os.Exit(1)

			return NULL
		},
	},
}
