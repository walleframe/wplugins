package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/walleframe/wplugins/buildpb"
	"github.com/walleframe/wplugins/gen"
	"github.com/walleframe/wplugins/utils/plugin"
)

func main() {
	plugin.MainOneByOne(generateWalleZapLog)
}

func generateWalleZapLog(prog *buildpb.FileDesc, depend map[string]*buildpb.FileDesc) (out []*buildpb.BuildOutput, err error) {

	g := gen.New(
		gen.WithFormat(gen.GoFormat2),
		gen.WithIndent("    "),
		gen.WithKeyTitle(true),
	)

	g.P("// Generate by wctl plugin(wzap). DO NOT EDIT.")
	g.P("package ", prog.Pkg.Package, ";")
	g.P()
	g.P("import (")
	g.In()
	g.P(`"go.uber.org/zap/zapcore"`)
	// 用于binary数据打印
	g.P(`"encoding/base64"`)
	// 用于map key 打印
	g.P(`"strconv"`)
	g.Out()
	g.P(")")

	g.P()
	g.P()

	for _, msg := range prog.Msgs {
		g.Doc(msg.Doc)
		g.P("func (x *", g.Key(msg.Name), ") MarshalLogObject(enc zapcore.ObjectEncoder) error {")
		g.In()
		for _, field := range msg.Fields {
			typ := field.Type
			fname := g.Key(field.Name)
			switch typ.Type {
			case buildpb.FieldType_BaseType:

				tf, err := typFuncName(typ.KeyBase)
				if err != nil {
					return nil, fmt.Errorf("message %s field %s %w", msg.Name, field.Name, err)
				}
				g.Pf(`enc.Add%s("%s", %s)`, tf, fname, typValue(typ.KeyBase, fmt.Sprintf("x.%s", fname)))
			case buildpb.FieldType_CustomType:
				g.Pf(`enc.AddObject("%[1]s", x.%[1]s)`, fname)
			case buildpb.FieldType_ListType:
				if typ.ElemCustom {
					pkg := ""
					stName := g.Key(typ.Key)
					if pl := strings.Split(typ.Key, "."); len(pl) > 1 {
						pkg = pl[0] + "."
						stName = g.Key(pl[1])
					}
					log.Printf("---[%s][%s][%s]\n", fname, pkg, stName)
					g.Pf(`enc.AddArray("%[1]s", %sZapArray%s(x.%[1]s))`, fname, pkg, stName)
				} else {
					g.Pf(`enc.AddArray("%s", zapcore.ArrayMarshalerFunc(func(ae zapcore.ArrayEncoder) error {`, fname)
					g.In()
					g.Pf(`for _, v := range x.%s {`, fname)
					g.In()

					tf, err := typFuncName(typ.KeyBase)
					if err != nil {
						return nil, fmt.Errorf("message %s field %s %w", msg.Name, field.Name, err)
					}

					g.Pf(`ae.Append%s(%s)`, tf, typValue(typ.KeyBase, "v"))
					g.Out()
					g.P("}")
					g.P("return nil")
					g.Out()
					g.P("}))")
				}
			case buildpb.FieldType_MapType:
				g.Pf(`enc.AddObject("%s", zapcore.ObjectMarshalerFunc(func(oe zapcore.ObjectEncoder) error {`, fname)
				g.In()
				g.Pf(`for k, v := range x.%s {`, fname)
				g.In()
				// oe.AddInt32(strconv.FormatInt(int64(k), 10), v)
				//strconv.FormatFloat(0.1, fmt byte, prec int, bitSize int)

				if typ.ElemCustom {
					pkg := ""
					stName := g.Key(field.Type.Value)
					if pl := strings.Split(field.Type.Value, "."); len(pl) > 1 {
						pkg = pl[0] + "."
						stName = g.Key(pl[1])
					}
					g.Pf(`oe.AddObject(%s, %s%s(x.%s))`, typMapKey(typ.KeyBase, "k"), pkg, stName, fname)
				} else {
					tf, err := typFuncName(typ.ValueBase)
					if err != nil {
						return nil, fmt.Errorf("message %s field %s %w", msg.Name, field.Name, err)
					}
					g.Pf(`oe.Add%s(%s, %s)`, tf, typMapKey(typ.KeyBase, "k"), typValue(typ.ValueBase, "v"))
				}

				g.Out()
				g.P("}")
				g.P("return nil")
				g.Out()
				g.P("}))")
			}
		}
		g.P("return nil")
		g.Out()
		g.P("}")
		g.Pf(`type ZapArray%[1]s []*%[1]s

func (x ZapArray%[1]s) MarshalLogArray(ae zapcore.ArrayEncoder) error {
	for _, v := range x {
		ae.AppendObject(v)
	}
	return nil
}`, g.Key(msg.Name))
	}

	data, err := g.Bytes()
	if err != nil {
		log.Println("format code failed.", err)
		return nil, fmt.Errorf("format failed %w", err)
		// data = g.Buffer.Bytes()
		// err = nil
	}
	out = append(out, &buildpb.BuildOutput{
		File: strings.TrimSuffix(prog.File, filepath.Ext(prog.File)) + ".zap.go",
		Data: data,
	})
	return
}

func typFuncName(typ buildpb.BaseTypeDesc) (tf string, err error) {
	switch typ {
	case buildpb.BaseTypeDesc_String,
		buildpb.BaseTypeDesc_Binary:
		tf = "String"
	case buildpb.BaseTypeDesc_Int8, buildpb.BaseTypeDesc_Int16,
		buildpb.BaseTypeDesc_Int32:
		tf = "Int32"
	case buildpb.BaseTypeDesc_Uint8, buildpb.BaseTypeDesc_Uint16,
		buildpb.BaseTypeDesc_Uint32:
		tf = "Uint32"
	case buildpb.BaseTypeDesc_Int64:
		tf = "Int64"
	case buildpb.BaseTypeDesc_Uint64:
		tf = "Uint64"
	case buildpb.BaseTypeDesc_Float32:
		tf = "Float32"
	case buildpb.BaseTypeDesc_Float64:
		tf = "Float64"
	case buildpb.BaseTypeDesc_Bool:
		tf = "Bool"
	default:
		err = fmt.Errorf("invalid base type %+v", typ)
	}
	return
}

func typValue(typ buildpb.BaseTypeDesc, fname string) (value string) {
	if typ == buildpb.BaseTypeDesc_Binary {
		return fmt.Sprintf(fmtBinaryData, fname)
	}
	return fname
}

const (
	fmtBinaryData = "base64.StdEncoding.EncodeToString([]byte(%s))"
)

func typMapKey(typ buildpb.BaseTypeDesc, k string) (tf string) {
	switch typ {
	case buildpb.BaseTypeDesc_String:
		tf = k
	case buildpb.BaseTypeDesc_Binary:
		tf = k
	case buildpb.BaseTypeDesc_Int8, buildpb.BaseTypeDesc_Int16,
		buildpb.BaseTypeDesc_Int32:
		tf = fmt.Sprintf("strconv.FormatInt(int64(%s), 10)", k)
	case buildpb.BaseTypeDesc_Uint8, buildpb.BaseTypeDesc_Uint16,
		buildpb.BaseTypeDesc_Uint32:
		tf = fmt.Sprintf("strconv.FormatUint(uint64(%s), 10)", k)
	case buildpb.BaseTypeDesc_Int64:
		tf = fmt.Sprintf("strconv.FormatInt(%s, 10)", k)
	case buildpb.BaseTypeDesc_Uint64:
		tf = fmt.Sprintf("strconv.FormatUint(%s, 10)", k)
	case buildpb.BaseTypeDesc_Float32:
		tf = fmt.Sprintf("strconv.FormatFloat(%s, 'f', -1, 32)", k)
	case buildpb.BaseTypeDesc_Float64:
		tf = fmt.Sprintf("strconv.FormatFloat(%s, 'f', -1, 64)", k)
	case buildpb.BaseTypeDesc_Bool:
		tf = fmt.Sprintf("strconv.FormatBool(%s)", k)
	}
	return
}
