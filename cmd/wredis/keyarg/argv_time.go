package keyarg

import (
	"fmt"
	"strings"
)

type TimeArg struct {
	EmptyArg
	fun string
	op  string
}

// 需要导入的包
func (x *TimeArg) Imports() []string {
	return []string{"github.com/walleframe/walle/util/wtime"}
}

// 格式化代码
func (x *TimeArg) FormatCode(obj string) (code string) {
	return fmt.Sprintf(`%s.WriteInt64(wtime.%s()%s)`, obj, x.fun, x.op)
}

// TimeMatch 时间参数匹配
type TimeMatch struct{}

func (x TimeMatch) Match(st KeyStructer, k ArgIndex, arg string) (_ KeyArg, err error) {
	// 非@开头，不匹配
	if !strings.HasPrefix(arg, "@") {
		return nil, nil
	}

	fun := trimParentheses(strings.TrimPrefix(arg, "@"))
	op := ""
	fun, op, err = splitKeyOp(fun)
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(fun) {
	case "daystamp":
		fun = "DayStamp"
	case "weekstamp":
		fun = "WeekStamp"
	case "monthstamp":
		fun = "MonthStamp"
	case "yearstamp":
		fun = "YearStamp"
	default:
		err = fmt.Errorf("[%s] is invalid wtime functions.", fun)
		return
	}
	return &TimeArg{fun: fun, op: op}, nil
}
