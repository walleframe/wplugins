package genparse

import (
	"fmt"
	"log"
	"strings"

	"github.com/walleframe/wplugins/buildpb"
	"github.com/walleframe/wplugins/cmd/wpb/gengo"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var genzapTemplate = `
{{ $i := Import "go.uber.org/zap/zapcore" "ObjectEncoder" }} {{$i := Import "strconv" "FormatInt"}} {{$_ := Import "encoding/base64" "encode"}}
func (x *{{.GoName}}) MarshalLogObject(enc zapcore.ObjectEncoder) error { {{- range $i,$field := .Fields }}
	{{ZapFieldGen $field}}{{end}}
	return nil 
}

type ZapArray{{.GoName}} []*{{.GoName}}
func (x ZapArray{{.GoName}}) MarshalLogArray(ae zapcore.ArrayEncoder) error{
	for _, v := range x {
		ae.AppendObject(v)
	}
	return nil 
}

func LogArray{{.GoName}}(name string, v []*{{.GoName}}) zap.Field { {{$i := Import "go.uber.org/zap" "Array"}}
	return zap.Array(name, ZapArray{{.GoName}}(v))
}
`

func init() {
	gengo.RegisterCustomModule(&gengo.CustomModule{
		Templates: [][2]string{{"genzap", genzapTemplate}},
		Funcs: map[string]interface{}{
			"ZapFieldGen": genZapField,
		},
	})
}

func genZapField(field *gengo.GenerateField) (_ string) {
	typ := field.WType
	switch typ.Type {
	case buildpb.FieldType_BaseType:
		return fmt.Sprintf(`enc.Add%s("%s", x.%s)`, getZapBaseTypeFunc(field), field.DescName, field.GoName)
	case buildpb.FieldType_CustomType:
		return fmt.Sprintf(`enc.AddObject("%s", x.%s)`, field.DescName, field.GoName)
	case buildpb.FieldType_ListType:
		if typ.ElemCustom {
			return fmt.Sprintf(`enc.AddArray("%s", x.%s)`, field.DescName, field.GoName)
		}
		return fmt.Sprintf(
			`enc.AddArray("%s", zapcore.ArrayMarshalerFunc(func(ae zapcore.ArrayEncoder) error {
	for _,v := range x.%s {
		ae.Append%s
	}
	return nil 
}))`, field.DescName, field.GoName, getZapArrayFunc(field))
	case buildpb.FieldType_MapType:
		return fmt.Sprintf(`enc.AddObject("%s", zapcore.ObjectMarshalerFunc(func(oe zapcore.ObjectEncoder) error {
	for k,v := range x.%s {
		oe.Add%s
	}
	return nil 
}))`, field.DescName, field.GoName, getZapMapFunc(field))
	}
	return
}

func getZapBaseTypeFunc(field *gengo.GenerateField) (f string) {
	switch field.GoType {
	case "[]byte":
		return "Binary"
	default:
		return strings.Title(field.GoType)
	}
}

func getZapArrayFunc(field *gengo.GenerateField) (f string) {
	switch field.GoType {
	case "[]byte":
		return "String(base64.StdEncoding.EncodeToString(v))"
	default:
		return strings.Title(field.GoType) + "(v)"
	}
}

func getZapMapFunc(field *gengo.GenerateField) (f string) {
	typ := field.WType
	if typ.ElemCustom {
		return fmt.Sprintf(`Object(%s, v)`, typMapKey(typ.KeyBase, "k"))
	}
	return fmt.Sprintf(`%s(%s, %s)`, getZapBaseTypeFunc(field), typMapKey(typ.KeyBase, "k"), getZapMapValue(field))
}

func getZapMapValue(field *gengo.GenerateField) (f string) {
	switch field.GoType {
	case "[]byte":
		return "base64.StdEncoding.EncodeToString(v)"
	default:
		return "v"
	}
}


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

func getZaoImprtConv(field *gengo.GenerateField) (isImport bool) {
	isImport = true
	switch field.Kind {
	case protoreflect.EnumKind:
		isImport = false
	case protoreflect.StringKind:
		isImport = false
	}
	return
}

func getZaoFieldMapKey(field *gengo.GenerateField) (funcName string) {
	switch field.Kind {
	case protoreflect.BoolKind:
		funcName = "strconv.FormatBool(k)"
	case protoreflect.EnumKind:
		funcName = "k.String()"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		funcName = "strconv.FormatInt(int64(k), 10)"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		funcName = "strconv.FormatUint(uint64(k), 10)"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		funcName = "strconv.FormatInt(k, 10)"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		funcName = "strconv.FormatUint(k, 10)"
	case protoreflect.FloatKind:
		funcName = "strconv.FormatFloat32(k,'f', -1, 32)"
	case protoreflect.DoubleKind:
		funcName = "strconv.FormatFloat64(k, 'f', -1, 64)"
	case protoreflect.StringKind:
		funcName = "k"
	case protoreflect.BytesKind:
		funcName = "{Invalid Map Key - []byte}"
		log.Printf("invalid map key type []byte. %#v", field)
	case protoreflect.MessageKind, protoreflect.GroupKind:
		funcName = "{Invalid Map Key - Object}"
		log.Printf("invalid map key type object. %#v", field)
	}
	return
}

func getZapFieldFunc(field *gengo.GenerateField) (funcName string) {
	switch field.Kind {
	case protoreflect.BoolKind:
		funcName = "Bool"
	case protoreflect.EnumKind:
		funcName = "String"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		funcName = "Int32"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		funcName = "Uint32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		funcName = "Int64"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		funcName = "Uint64"
	case protoreflect.FloatKind:
		funcName = "Float32"
	case protoreflect.DoubleKind:
		funcName = "Float64"
	case protoreflect.StringKind:
		funcName = "String"
	case protoreflect.BytesKind:
		funcName = "Binary"
	case protoreflect.MessageKind, protoreflect.GroupKind:
		funcName = "Object"
	}
	return
}

func getZapMethod(field *gengo.GenerateField) (fieldMethod string) {
	switch field.Kind {
	case protoreflect.EnumKind:
		fieldMethod = ".String()"
	}
	return
}

// func genZapMessage(g *protogen.GeneratedFile, m *protogen.Message, msg *gengo.GenerateMessage) {
// 	if m.Desc.IsMapEntry() {
// 		return
// 	}
// 	// object marshal
// 	//g.Annotate(m.GoIdent.GoName, m.Location)
// 	g.Import(protogen.GoImportPath("go.uber.org/zap/zapcore"))
// 	g.QualifiedGoIdent(protogen.GoIdent{GoName: "Abc", GoImportPath: "go.uber.org/zap/zapcore"})
// 	g.P("func (x *", m.GoIdent, ") MarshalLogObject(enc zapcore.ObjectEncoder) error {")
// 	for _, field := range m.Fields {
// 		if field.Desc.IsWeak() {
// 			continue
// 		}

// 		keyName := field.GoName
// 		fieldName := field.GoName
// 		if field.Desc.IsMap() {
// 			g.P(`enc.AddObject("`, keyName, `", zapcore.ObjectMarshalerFunc(func(oe zapcore.ObjectEncoder) error {`)
// 			g.P(`for k,v := range x.`, fieldName, "{")
// 			funcName, fieldMethod := getZapFieldFunc(field.Message.Fields[1])
// 			g.P(fmt.Sprintf(`enc.Add%s(%s, v%s)`, funcName, getZaoFieldMapKey(g, field.Message.Fields[0]), fieldMethod))
// 			g.P("}")
// 			g.P("return nil")
// 			g.P("}))")
// 			continue
// 		}
// 		funcName, fieldMethod := getZapFieldFunc(field)
// 		switch {
// 		case field.Desc.IsList():
// 			g.P(fmt.Sprintf(`enc.AddArray("%s", zapcore.ArrayMarshalerFunc(func(ae zapcore.ArrayEncoder) error {`, keyName))
// 			g.P("for _,v := range x.", fieldName, "{")
// 			if funcName == "Binary" {
// 				g.Import(protogen.GoImportPath("encoding/base64"))
// 				g.QualifiedGoIdent(protogen.GoIdent{GoName: "NewEncodeToString", GoImportPath: "encoding/base64"})

// 				g.P(fmt.Sprintf("ae.AppendString(base64.StdEncoding.EncodeToString(v%s))", fieldMethod))
// 			} else {
// 				g.P(fmt.Sprintf("ae.Append%s(v%s)", funcName, fieldMethod))
// 			}
// 			g.P("}")
// 			g.P("return nil")
// 			g.P("}))")
// 		default:
// 			g.P(fmt.Sprintf(`enc.Add%s("%s", x.%s%s)`, funcName, keyName, fieldName, fieldMethod))
// 		}
// 	}
// 	g.P("return nil")
// 	g.P("}")
// 	g.P()

// 	g.P("type ZapArray", m.GoIdent, " []*", m.GoIdent)
// 	g.P("func (x ZapArray", m.GoIdent, ")  MarshalLogArray(ae zapcore.ArrayEncoder) error {")
// 	g.P("for _, v := range x {")
// 	g.P("ae.AppendObject(v)")
// 	g.P("}")
// 	g.P("return nil")
// 	g.P("}")
// 	g.Import(protogen.GoImportPath("go.uber.org/zap"))
// 	g.QualifiedGoIdent(protogen.GoIdent{GoName: "Abc", GoImportPath: "go.uber.org/zap"})
// 	g.P(`func LogArray`, m.GoIdent, `(name string, v []*`, m.GoIdent, `) zap.Field {`)
// 	g.P(`return zap.Array(name, ZapArray`, m.GoIdent, `(v))`)
// 	g.P("}")

// 	// array marshal
// }
