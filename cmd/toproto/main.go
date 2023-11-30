package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/aggronmagi/wplugins/buildpb"
	"github.com/aggronmagi/wplugins/gen"
	"github.com/aggronmagi/wplugins/options"
	"github.com/aggronmagi/wplugins/utils"
	"github.com/aggronmagi/wplugins/utils/plugin"
)

func main() {
	plugin.MainRoot(toProto)
}

func toProto(rq *buildpb.BuildRQ) (rs *buildpb.BuildRS, err error) {
	rs = &buildpb.BuildRS{}
	var data []byte
	for _, file := range rq.Files {
		desc, ok := rq.Programs[file]
		if !ok {
			err = fmt.Errorf("file desc not found [%s]", file)
			return
		}

		pkg, _ := desc.Options.Options[options.ProtoGoPkg]
		_, lowercase := desc.Options.Options[options.ProtoLowerCase]
		syntax, _ := desc.Options.Options[options.ProtoSyntax]

		g := gen.New(
			gen.WithIndent("  "),
			gen.WithKeyTitle(!lowercase),
		)
		proto2Version := false

		g.P("// Generate by wctl plugin(toproto). DO NOT EDIT.")
		// 语法版本
		if syntax != nil && syntax.IntValue == 2 {
			proto2Version = true
			g.P(`syntax = "proto2";`)
		} else {
			g.P(`syntax = "proto3";`)
		}
		if proto2Version {
			log.Println("暂未支持 proto2 语法,忽略", file)
			continue
		}

		// 包名
		g.Doc(desc.Pkg.Doc)
		g.P("package ", desc.Pkg.Package, ";")

		// 选项信息
		if pkg != nil && pkg.Value != "" {
			if pkg.Doc != nil && len(pkg.Doc.Doc) > 0 {
				for _, v := range pkg.Doc.Doc {
					g.P(v)
				}
			}
			g.P(`option go_package = "`, pkg.Value, `";`)
		}
		g.P()
		// import 其它文件
		for _, v := range desc.Imports {
			// ids, ok := rq.Programs[v.File]
			// if !ok {
			// 	log.Println("error: invalid import. ", v.String())
			// 	err = errors.New("invlaid import")
			// 	return
			// }
			g.Pf(`import "%s";`, strings.Replace(v.File, filepath.Ext(v.File), ".proto", -1))
		}
		for _, e := range desc.Enums {
			g.Doc(e.Doc)
			g.P("enum ", g.Key(e.Name), "{")
			g.In()
			for _, v := range e.Values {
				g.Doc(v.Doc)
				g.Pf("%s = %d;", g.Key(v.Name), v.Value)
			}
			g.Out()
			g.P("}")
		}
		// 消息转换
		var typ string
		for _, msg := range desc.Msgs {
			g.Doc(msg.Doc)
			g.P("message ", g.Key(msg.Name), " {")
			g.In()
			for _, field := range msg.Fields {
				g.Doc(field.Doc)
				typ, err = getTypeName(field.Type, lowercase)
				g.P(typ, " ", g.Key(field.Name), " = ", field.No, ";", msg.Doc.GetTailDoc())
			}

			g.Out()
			g.P("}")
		}
		// 输出结果
		data, err = g.Bytes()
		if err != nil {
			log.Println("error:", err)
			return
		}
		rs.Result = append(rs.Result, &buildpb.BuildOutput{
			File: strings.TrimSuffix(file, filepath.Ext(file)) + ".proto",
			Data: data,
		})
	}

	return
}

// Type = BaseType,CustomType
// 使用 Key
// Type = ListType
// 使用 Key 根据ElemCustom判断 数组元素是自定义类型还是基础类型
// Type = MapType
// Key 是基础类型. Value 根据ElemCustom判断 数组元素是自定义类型还是基础类型
func getTypeName(typ *buildpb.TypeDesc, lowercase bool) (name string, err error) {
	switch typ.Type {
	case buildpb.FieldType_BaseType:
		return baseTypeName(typ.Key)
	case buildpb.FieldType_CustomType:
		if lowercase {
			return typ.Key, nil
		}
		lst := strings.Split(typ.Key, ".")
		if len(lst) > 1 {
			return lst[0] + "." + utils.Title(lst[1]), nil
		} else {
			return utils.Title(typ.Key), nil
		}
	case buildpb.FieldType_ListType:
		if typ.ElemCustom {
			if lowercase {
				return "repeated " + typ.Key, nil
			}
			lst := strings.Split(typ.Key, ".")
			if len(lst) > 1 {
				return "repeated " + lst[0] + "." + utils.Title(lst[1]), nil
			} else {
				return "repeated " + utils.Title(lst[0]), nil
			}
		}
		name, err = baseTypeName(typ.Key)
		if err != nil {
			return
		}
		return "repeated " + name, nil
	case buildpb.FieldType_MapType:
		var key, val string
		key, err = baseTypeName(typ.Key)
		if err != nil {
			return
		}
		if typ.ElemCustom {
			if lowercase {
				val = typ.Value
			}
			lst := strings.Split(typ.Value, ".")
			if len(lst) > 1 {
				val = lst[0] + "." + utils.Title(lst[1])
			} else {
				val = utils.Title(lst[0])
			}

		} else {
			val, err = baseTypeName(typ.Value)
		}
		if err != nil {
			return
		}
		return "map<" + key + "," + val + ">", nil
	}

	return
}

func baseTypeName(in string) (out string, err error) {
	switch in {
	case "string":
		out = in
	case "bytes", "binary":
		out = "bytes"
	case "int8", "int16", "int32", "int":
		out = "int32"
	case "uint8", "uint16", "uint32", "uint":
		out = "uint32"
	case "int64":
		out = in
	case "uint64":
		out = in
	case "float", "float32":
		out = "float"
	case "double", "float64":
		out = "double"
	case "bool":
		out = "bool"
	default:
		err = fmt.Errorf("invalid base type %s", in)
	}
	return
}

// func baseTypeName(typ *buildpb.TypeDesc) string {
// 	switch typ {
// 	case ast.BaseTypeInt8, ast.BaseTypeInt16, ast.BaseTypeInt32:
// 		return "int32"
// 	case ast.BaseTypeUint8, ast.BaseTypeUint16, ast.BaseTypeUint32:
// 		return "uint32"
// 	case ast.BaseTypeInt64:
// 		return "int64"
// 	case ast.BaseTypeUint64:
// 		return "uint64"
// 	case ast.BaseTypeString:
// 		return "string"
// 	case ast.BaseTypeBinary:
// 		return "string"
// 	case ast.BaseTypeBool:
// 		return "bool"
// 	default:
// 		fmt.Println(typ)
// 		panic("unkown type")
// 	}

// 	case buildpb:
// 		return baseTypeName(typ.YTBaseType)
// 	case typ.YTCustomType != nil:
// 		return typ.YTCustomType.Name
// 	case typ.YTListType != nil:
// 		if typ.YTListType.YTBaseType != nil {
// 			return "repeated " + baseTypeName(typ.YTListType.YTBaseType)
// 		} else if typ.YTListType.YTCustomType != nil {
// 			return "repeated " + typ.YTListType.YTCustomType.Name
// 		}
// 	case typ.YTMapTypee != nil:

// 		if typ.YTMapTypee.Value.YTBaseType != nil {
// 			return "map<" + baseTypeName(typ.YTMapTypee.Key) + "," +
// 				baseTypeName(typ.YTMapTypee.Value.YTBaseType) + ">"
// 		} else if typ.YTMapTypee.Value.YTCustomType != nil {
// 			return "map<" + baseTypeName(typ.YTMapTypee.Key) + "," +
// 				typ.YTMapTypee.Value.YTCustomType.Name + ">"
// 		}
// 	}
// 	fmt.Printf("%#v\n", typ)
// 	panic("invalid type")
// }
