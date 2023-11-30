package keyarg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type EmptyArg struct{}

// 需要导入的包
func (x *EmptyArg) Imports() []string {
	return nil
}

// 是否是构造参数
func (x *EmptyArg) ConstructArg() bool {
	return false
}

// 参数类型
func (x *EmptyArg) ArgType() (arg string) {
	return
}

func (x *EmptyArg) ArgName() string {
	return ""
}

// 格式化代码
func (x *EmptyArg) FormatCode(obj string) (code string) {
	return
}

// SourceArg 原始字符串参数
type SourceArg struct {
	EmptyArg
	source string
}

// 格式化代码
func (x *SourceArg) FormatCode(obj string) (code string) {
	return fmt.Sprintf(`%s.WriteString("%s")`, obj, x.source)
}

// SourceMatch 原始字符串匹配
type SourceMatch struct{}

func (SourceMatch) Match(st KeyStructer, k ArgIndex, arg string) (KeyArg, error) {
	if err := checkValidKeyChar(arg); err != nil {
		return nil, err
	}
	return &SourceArg{source: arg}, nil
}

func checkValidKeyChar(key string) error {
	for _, v := range key {
		if (v >= '0' && v <= '9') || (v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') {
			continue
		}
		// 出现了非 0-9 a-z A-Z的字符,不匹配
		return errors.New(fmt.Sprintf("[%s] contain invalid char.", key))
	}
	return nil
}

// MergeSourceArg 合并source,减少write次数
func MergeSourceArg(in []KeyArg) (outs []KeyArg) {
	var pre *SourceArg
	for _, v := range in {
		src, ok := v.(*SourceArg)
		// 当前得是source
		if ok {
			// 没有前一个source,保存当前的
			if pre == nil {
				pre = src
				continue
			}
			// 拼接前一个source
			pre.source += ArgSplit + src.source
			continue
		}
		// 之前有source,先放source
		if pre != nil {
			outs = append(outs, pre)
			pre = nil
		}
		outs = append(outs, v)
	}
	if pre != nil {
		outs = append(outs, pre)
		pre = nil
	}
	return
}

func splitKeyOp(in string) (fun, op string, err error) {
	split := ""
	switch {
	case strings.Contains(in, "+"):
		split = "+"
	case strings.Contains(in, "%"):
		split = "%"
	case strings.Contains(in, "-"):
		split = "-"
	default:
		fun = in
		err = checkValidKeyChar(in)
		return
	}

	lists := strings.Split(in, split)
	if len(lists) != 2 {
		err = fmt.Errorf("%s split op failed,not 2 parts.", in)
		return
	}
	fun = trimParentheses(lists[0])
	err = checkValidKeyChar(fun)
	if err != nil {
		return
	}

	number, err := strconv.ParseInt(trimParentheses(lists[1]), 10, 64)
	if err != nil {
		return
	}

	op = fmt.Sprintf(" %s %d", split, number)

	return
}

func trimParentheses(in string) string {
	return strings.TrimSpace(
		strings.TrimPrefix(
			strings.TrimSpace(
				strings.TrimSuffix(
					strings.TrimSpace(in),
					")",
				),
			),
			"(",
		),
	)
}
