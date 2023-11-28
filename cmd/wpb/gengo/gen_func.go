package gengo

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"unsafe"

	"google.golang.org/protobuf/encoding/protowire"
)

var UseFuncMap = template.FuncMap{}

func init() {
	UseFuncMap["ValueName"] = func(in ...string) string {
		return strings.Join(in, "")
	}
	UseFuncMap["TagBinary"] = func(num int, typ string) (string, error) {
		t := SwitchTagType(typ)
		if t < 0 {
			return "", fmt.Errorf("%s is invalid protowire.Type", typ)
		}
		return BinaryBytes(protowire.AppendTag(nil, protowire.Number(num), t)), nil
	}
	UseFuncMap["TagByes"] = func(num int, typ string) (string, error) {
		t := SwitchTagType(typ)
		if t < 0 {
			return "", fmt.Errorf("%s is invalid protowire.Type", typ)
		}
		return SourceBytes(protowire.AppendTag(nil, protowire.Number(num), t)), nil
	}
	UseFuncMap["TagSize"] = func(num int) int {
		return protowire.SizeTag(protowire.Number(num))
	}
}

var SwitchTagType = func(typ string) protowire.Type {
	switch typ {
	case "protowire.VarintType", "VarintType":
		return protowire.VarintType
	case "protowire.Fixed32Type", "Fixed32Type":
		return protowire.Fixed32Type
	case "protowire.Fixed64Type", "Fixed64Type":
		return protowire.Fixed64Type
	case "protowire.BytesType", "BytesType":
		return protowire.BytesType
	default:
		return -1
	}
}

var SourceBytes = func(in []byte) string {
	buf := make([]byte, 0, len(in)*5-1)
	for k, v := range in {
		if k > 0 {
			buf = append(buf, ',')
		}
		buf = append(buf, '0', 'x')
		buf = strconv.AppendInt(buf, int64(v), 16)
	}
	return *(*string)(unsafe.Pointer(&buf))
}

var BinaryBytes = func(in []byte) string {
	// buf := make([]byte, 0, len(in)*9-1)
	// for k, v := range in {
	// 	if k > 0 {
	// 		buf = append(buf, ' ')
	// 	}
	// 	buf = strconv.AppendInt(buf, int64(v), 2)
	// }
	// return *(*string)(unsafe.Pointer(&buf))
	buf := strings.Builder{}
	buf.Grow(len(in) * 9)
	for k, v := range in {
		if k > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(fmt.Sprintf("%08b", v))
	}
	return buf.String()
}
