package keyarg

import (
	"errors"
	"fmt"
	"strings"

	"github.com/walleframe/wplugins/buildpb"
	"go.uber.org/multierr"
)

type KeyStructer interface {
	GetFieldName(idx int) string
	GetFieldType(field string) *buildpb.FieldType
}

type ArgIndex interface {
	Get() int
}

type ArgMatcher interface {
	// 是否匹配当前参数
	Match(st KeyStructer, k ArgIndex, arg string) (KeyArg, error)
}

type KeyArg interface {
	// 需要导入的包
	Imports() []string
	// 是否是构造参数
	ConstructArg() bool
	// 参数类型
	ArgType() (arg string)
	//
	ArgName() string
	// 格式化代码
	FormatCode(obj string) (code string)
}

var ArgSplit string = ":"

var keyArgRegistry []ArgMatcher

//var sourceMatch = &SourceMatch{}

func init() {
	keyArgRegistry = append(keyArgRegistry, // 外部参数支持
		// 时间参数
		TimeMatch{},
		// go类型参数
		GoTypeMatch{},

		// 最后的默认匹配
		SourceMatch{},
	)
}

func MatchKey(argv string, st KeyStructer) (args []KeyArg, err error) {
	lists := strings.Split(argv, ArgSplit)
	var argIdx argIndex
	for _, v := range lists {
		v = strings.TrimSpace(v)
		var arg KeyArg
		for _, m := range keyArgRegistry {
			arg, err = m.Match(st, &argIdx, v)
			if err != nil {
				return
			}
			if arg != nil {
				args = append(args, arg)
				break
			}
		}
		if arg != nil {
			continue
		}
		// // 默认匹配
		// arg = sourceMatch.Match(v)
		// if arg != nil {
		// 	args = append(args, arg)
		// 	continue
		// }
		// 没有匹配上
		err = errors.New(fmt.Sprintf("[%s] not match any KeyArg", v))
		return
	}
	// 自定义参数名称检测
	argMap := make(map[string]int)
	for k, arg := range args {
		if arg.ConstructArg() {
			name := arg.ArgName()
			if last, ok := argMap[name]; ok {
				err = multierr.Append(err, fmt.Errorf("conflict arg name.%d arg name equal to %d[%s].", k, last, name))
				continue
			}
			argMap[name] = k
		}
	}
	if err != nil {
		return
	}
	// 合并source arg，减少实际运行代码中的write操作。
	args = MergeSourceArg(args)
	return
}

func MatchGoTypes(argv string, st KeyStructer) (args []KeyArg, err error) {
	lists := strings.Split(argv, ArgSplit)
	var argIdx argIndex
	var arg KeyArg
	var match GoTypeMatch
	for _, v := range lists {
		v = strings.TrimSpace(v)

		arg, err = match.Match(st, &argIdx, v)
		if err != nil {
			return
		}

		if arg != nil {
			args = append(args, arg)
			continue
		}
		// 没有匹配上
		err = errors.New(fmt.Sprintf("[%s] not match any KeyArg", v))
		return
	}
	// 自定义参数名称检测
	argMap := make(map[string]int)
	for k, arg := range args {
		if arg.ConstructArg() {
			name := arg.ArgName()
			if last, ok := argMap[name]; ok {
				err = multierr.Append(err, fmt.Errorf("conflict arg name.%d arg name equal to %d[%s].", k, last, name))
				continue
			}
			argMap[name] = k
		}
	}
	if err != nil {
		return
	}
	return
}

type argIndex int

func (x *argIndex) Get() int {
	*x++
	return int(*x)
}
