package utils

import (
	"reflect"
	"strings"
	"unsafe"
)

func Title(in string) (out string) {
	list := strings.Split(in, "_")
	for _, v := range list {
		out += strings.Title(v)
	}
	return
}

// StringToBytes converts string to byte slice without a memory allocation.
func StringToBytes(s string) (b []byte) {
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return b
}

// BytesToString converts byte slice to string without a memory allocation.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func MergeNotSameSlice(l1, l2 []string) (res []string) {
	for _, v := range l1 {
		if v == "" {
			continue
		}
		res = append(res, strings.ToLower(v))
	}
	for _, v := range l2 {
		if v == "" {
			continue
		}
		find := false
		for _, v2 := range res {
			if v == v2 {
				find = true
				break
			}
		}
		if find {
			continue
		}
		res = append(res, strings.ToLower(v))
	}
	return
}
