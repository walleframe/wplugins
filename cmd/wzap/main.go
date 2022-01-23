package main

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/aggronmagi/wplugins/buildpb"
	"github.com/aggronmagi/wplugins/gen"
	"github.com/aggronmagi/wplugins/utils/plugin"
	"go.uber.org/zap/zapcore"
)

func main() {
	plugin.MainOneByOne(generateWalleZapLog)
}

type AddRs struct {
	Params []int64
	I64    int64
	Mv     map[int32]int32
}

func (x *AddRs) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddArray("a", zapcore.ArrayMarshalerFunc(func(ae zapcore.ArrayEncoder) error {
		for _, v := range x.Params {
			ae.AppendInt64(v)
		}
		return nil
	}))
	enc.AddInt64("i64", x.I64)
	enc.AddObject("mv", zapcore.ObjectMarshalerFunc(func(oe zapcore.ObjectEncoder) error {
		for k, v := range x.Mv {
			oe.AddInt32(strconv.FormatInt(int64(k), 10), v)
		}
		return nil
	}))
	// enc.AddString()
	// enc.AddFloat32()
	return nil
}

func generateWalleZapLog(prog *buildpb.FileDesc, depend map[string]*buildpb.FileDesc) (out []*buildpb.BuildOutput, err error) {
	g := gen.New(gen.WithGoFmt(true), gen.WithIndent("    "), gen.WithKeyTitle(true))

	g.P("// Generate by wctl plugin(wzap). DO NOT EDIT.")
	g.P("package ", prog.Pkg.Package, ";")
	g.P()
	g.P("import (")
	g.In()
	g.P(`"go.uber.org/zap/zapcore"`)
	g.P(`"encoding/base64"`)
	g.P(`"strconv"`)
	g.Out()
	g.P(")")

	g.P()
	g.P()
	for _, msg := range prog.Msgs {
		g.Doc(msg.Doc)
		g.Pf("func (x *", g.Key(msg.Name), ") MarshalLogObject(enc zapcore.ObjectEncoder) error {")
		g.In()
		for _, field := range msg.Fields {
			typ := field.Type
			fname := g.Key(field.Name)
			switch typ.Type {
			case buildpb.FieldType_BaseType:
				f, k, v, err := baseTypeFuncName(typ.KeyBase, fname)
				if err != nil {
					return nil, fmt.Errorf("message %s field %s %w", msg.Name, field.Name, err)
				}
				g.Pf(`enc.%s(%s,%s)`, f, k, v)
			case buildpb.FieldType_CustomType:
				g.Pf(`enc.AddObject("%s", x.%s)`, field.Name, fname)
			case buildpb.FieldType_ListType:
				if typ.ElemCustom {
					g.Pf(`enc.AddArray("%s", xx.ZapArray(x.%s))`, field.Name)
				} else {
					
				}
			case buildpb.FieldType_MapType:
				
			}
		}
		g.Out()
		g.P("}")
	}

	data, err := g.Bytes()
	if err != nil {
		return nil, err
	}
	out = append(out, &buildpb.BuildOutput{
		File: filepath.Base(prog.File) + ".zap.go",
		Data: data,
	})
	return
}

func baseTypeFuncName(typ buildpb.BaseTypeDesc, fn string) (f, k, v string, err error) {
	v = "x." + fn
	k = fmt.Sprintf(`"%s"`, fn)
	switch typ {
	case buildpb.BaseTypeDesc_String:
		f = "AddString"
	case buildpb.BaseTypeDesc_Binary:
		f = "AddString"
		v = fmt.Sprintf("base64.StdEncoding.EncodeToString([]byte(v.%s))", fn)
	case "int8", "int16", "int32", "int":
		f = "AddInt32"
	case "uint8", "uint16", "uint32", "uint":
		f = "AddUint32"
	case "int64":
		f = "AddInt64"
	case "uint64":
		f = "AddUint64"
	case "float", "float32":
		f = "AddFloat32"
	case buildpb.BaseTypeDesc_Int16:
		f = "AddFloat32"
	case buildpb.BaseTypeDesc_Bool:
		f = "AddBool"
	default:
		err = fmt.Errorf("invalid base type %s", in)
	}
	return
}
