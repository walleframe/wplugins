package keyarg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type NumberArg struct {
	EmptyArg
	typ string
	arg string
	op  string
}

// 是否是构造参数
func (x *NumberArg) ConstructArg() bool {
	return true
}

// 参数类型
func (x *NumberArg) ArgType() (arg string) {
	return strings.ToLower(x.typ)
}

func (x *NumberArg) ArgName() string {
	return x.arg
}


// 格式化代码
func (x *NumberArg) FormatCode(obj string) (code string) {
	return fmt.Sprintf(`%s.Write%s(%s%s)`, obj, strings.Title(x.typ), x.arg, x.op)
}

type StringArg struct {
	EmptyArg
	arg string
}

// 是否是构造参数
func (x *StringArg) ConstructArg() bool {
	return true
}

// 参数类型
func (x *StringArg) ArgType() (arg string) {
	return "string"
}

func (x *StringArg) ArgName() string {
	return x.arg
}

// 格式化代码
func (x *StringArg) FormatCode(obj string) (code string) {
	return fmt.Sprintf(`%s.WriteString(%s)`, obj, x.arg)
}

// GoTypeMatch go类型参数匹配
type GoTypeMatch struct{}

func (x GoTypeMatch) Match(st KeyStructer, k ArgIndex, arg string) (_ KeyArg, err error) {
	// 非@开头，不匹配
	if !strings.HasPrefix(arg, "$") {
		return nil, nil
	}

	typ := strings.ToLower(trimParentheses(strings.TrimPrefix(arg, "$")))
	name := ""
	op := ""
	if strings.Contains(typ, "=") {
		lists := strings.Split(typ, "=")
		if len(lists) != 2 {
			err = fmt.Errorf("[%s] split = not 2 parts.", arg)
		}
		name = lists[0]
		typ = lists[1]
	}
	typ, op, err = splitKeyOp(typ)
	if err != nil {
		return nil, err
	}
	if name == "" {
		name = "arg" + strconv.Itoa(k.Get())
	}
	switch strings.ToLower(typ) {
	case "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
		return &NumberArg{typ: typ, arg: name, op: op}, nil
	case "bool":
		if op != "" {
			err = errors.New("bool type not support number op.")
			return
		}
		return &NumberArg{typ: typ, arg: name}, nil
	case "string":
		if op != "" {
			err = errors.New("string type not support number op.")
			return
		}
		return &StringArg{arg: name}, nil
	default:
		err = fmt.Errorf("[%s] is invalid golang types.", typ)
		return
	}
}
