package genparse

import (
	"strings"

	"github.com/walleframe/wplugins/buildpb"
	"github.com/walleframe/wplugins/cmd/wpb/gengo"
	"github.com/walleframe/wplugins/gen"
	"github.com/walleframe/wplugins/options"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// 外部配置
var (
	Getter bool
	// "google.golang.org/protobuf/encoding/protowire"
	WirePkg string = "github.com/walleframe/walle/util/protowire"
	Zap     bool   = true
)

// 版本信息
var (
	Version = "0.0.3"
)

func ParseEnum(t *gengo.GenerateStruct, g *gen.Generator, e *buildpb.EnumDesc) (err error) {
	t.WirePkg = WirePkg
	enum := &gengo.GenerateEnums{}
	// enum.LeadingComments = appendDeprecationSuffix(e.Comments.Leading,
	// 	e.Desc.ParentFile(),
	// 	e.Desc.Options().(*descriptorpb.EnumOptions).GetDeprecated()).String()
	enum.LeadingComments = e.Doc.ToDoc() //e.GetDoc().String()
	enum.TrailingComment = e.Doc.GetTailDoc()
	enum.TypeName = e.Name      //g.QualifiedGoIdent(e.GoIdent)
	enum.GoName = g.Key(e.Name) // e.GoIdent.GoName

	for _, value := range e.Values {
		val := &gengo.GenerateEnumValue{}
		// val.LeadingComments = appendDeprecationSuffix(value.Comments.Leading,
		// 	value.Desc.ParentFile(),
		// 	value.Desc.Options().(*descriptorpb.EnumValueOptions).GetDeprecated()).String()
		// val.TrailingComment = trailingComment(value.Comments.Trailing).String()
		val.LeadingComments = value.Doc.ToDoc() //value.GetDoc().String()
		val.TrailingComment = value.Doc.GetTailDoc()
		// val.Desc = string(value.Desc.Name())
		// if value.Desc != e.Desc.Values().ByNumber(value.Desc.Number()) {
		// 	val.Duplicate = "// Duplicate value: "
		// }
		val.Desc = value.Name
		val.Num = int32(value.Value)
		val.ValueName = g.Key(value.Name) // g.QualifiedGoIdent(value.GoIdent)
		enum.Values = append(enum.Values, val)
	}

	t.Enums = append(t.Enums, enum)

	// g.Import(protogen.GoImportPath("strconv"))
	// g.QualifiedGoIdent(protogen.GoIdent{GoName: "Abc", GoImportPath: "strconv"})
	return
}

func ParseMessage(t *gengo.GenerateStruct, g *gen.Generator, m *buildpb.MsgDesc, ignoreCheck func(m *buildpb.MsgDesc) bool) (err error) {

	if !ignoreCheck(m) {
		return
	}

	t.WirePkg = WirePkg
	msg := &gengo.GenerateMessage{}
	// msg.LeadingComments = appendDeprecationSuffix(m.Comments.Leading,
	// 	m.Desc.ParentFile(),
	// 	m.Desc.Options().(*descriptorpb.MessageOptions).GetDeprecated()).String()
	msg.LeadingComments = m.Doc.ToDoc() // m.GetDoc().String()
	msg.TrailingComment = m.Doc.GetTailDoc()
	msg.TypeName = g.Key(m.Name)
	msg.GoName = g.Key(m.Name)
	msg.GenGetter = Getter

	for _, field := range m.Fields {
		gf, ne := parseMessageField(msg, g, m, field)
		err = multierr.Append(err, ne)
		if gf != nil {
			msg.Fields = append(msg.Fields, gf)
		}
	}
	t.Messages = append(t.Messages, msg)

	if Zap {
		msg.CustomTemplates = append(msg.CustomTemplates, "genzap")
	}

	// // sub enum
	// for _, en := range m.Enums {
	// 	ParseEnum(t, g, f, en)
	// }
	// sub message
	for _, msg := range m.SubMsgs {
		ParseMessage(t, g, msg, ignoreCheck)
	}
	return
}

func parseMessageField(msg *gengo.GenerateMessage, g *gen.Generator, m *buildpb.MsgDesc, field *buildpb.Field) (genField *gengo.GenerateField, err error) {

	goType, pointer := fieldGoType(g, field)
	if pointer {
		goType = "*" + goType
	}

	genField = &gengo.GenerateField{}

	// 类型相关
	genField.LeadingComments = field.Doc.ToDoc() //field.GetDoc().String()
	genField.TrailingComment = field.Doc.GetTailDoc()
	//genField.TrailingComment = trailingComment(field.Comments.Trailing).String()
	genField.TypeName = goType
	genField.GoName = g.Key(field.Name)
	if goType == "[]byte" {
		genField.GoType = goType
	} else {
		//genField.GoType = strings.TrimPrefix(goType, "[]")
		genField.GoType = strings.TrimPrefix(strings.TrimPrefix(goType, "[]"), "*")
	}
	genField.Tip = msg.GoName + "." + g.Key(field.Name)
	// getter 相关
	defaultValue := fieldDefaultValue(g, m, field)
	genField.GetNilCheck = defaultValue == "nil"
	genField.DefaultValue = defaultValue
	// tag 相关
	genField.DescNum = int(field.No)
	genField.DescName = string(field.Name)
	genField.DescType, genField.WireType = switchProtoType(field)
	genField.IsList = field.IsList()
	genField.IsMap = field.IsMap()
	genField.Kind = switchProtoKind(field)

	// 序列化
	switch {
	case field.IsMap():
		genField.CheckNotEmpty = func(vname string) string {
			return "len(" + vname + ") > 0"
		}

		genField.MapKey, err = parseMessageField(msg, g, m, &buildpb.Field{
			Name:    "Key",
			Options: field.Options,
			No:      1,
			Type: &buildpb.TypeDesc{
				Type:       buildpb.FieldType_BaseType,
				Key:        field.Type.Key,
				ElemCustom: false,
				KeyBase:    field.Type.KeyBase,
			},
		})
		if err != nil {
			return
		}

		nf := &buildpb.Field{
			Name:    "Value",
			Options: field.Options,
			No:      2,
			Type: &buildpb.TypeDesc{
				Type:       buildpb.FieldType_BaseType,
				Key:        field.Type.Value,
				ElemCustom: false,
				KeyBase:    field.Type.ValueBase,
				Msg:        field.Type.Msg,
			},
		}
		if field.Type.ElemCustom {
			nf.Type.Type = buildpb.FieldType_CustomType
			// log.Println("----", goType)
			// log.Println(utils.Sdump(field, "origin"))
			// log.Println(utils.Sdump(nf, "map value"))
		}

		genField.MapValue, err = parseMessageField(msg, g, m, nf)
		if err != nil {
			return
		}
		genField.TemplateDecode = "decode.map"
		genField.TemplateSize = "size.map"
		genField.TemplateEncode = "encode.map"

	case field.IsList():
		genField.CheckNotEmpty = func(vname string) string {
			return "len(" + vname + ") > 0"
		}

		if field.Options.GetOptionBool(options.ProtoFieldOptPacked) {
			parseFillListPackedFiled(genField, field)
		} else {
			parseFillListNoPackedFiled(genField, field)
		}

	default:
		parseFillBasicFiled(genField, field)
	}

	// // import
	// g.Import(protogen.GoImportPath(WirePkg))
	// g.QualifiedGoIdent(protogen.GoIdent{GoName: "VarintType", GoImportPath: protogen.GoImportPath(WirePkg)})
	// g.Import(protogen.GoImportPath("errors"))
	// g.QualifiedGoIdent(protogen.GoIdent{GoName: "New", GoImportPath: "errors"})

	return
}

func switchProtoType(field *buildpb.Field) (typ int, desc string) {
	switch {
	case field.IsBasicType():
		switch field.Type.KeyBase {
		case buildpb.BaseTypeDesc_Binary, buildpb.BaseTypeDesc_String:
			desc = "protowire.BytesType"
			typ = int(protowire.BytesType)
			return
		case buildpb.BaseTypeDesc_Float32:
			desc = "protowire.Fixed32Type"
			typ = int(protowire.Fixed32Type)
			return
		case buildpb.BaseTypeDesc_Float64:
			desc = "protowire.Fixed64Type"
			typ = int(protowire.Fixed64Type)
			return
		}
		if field.Options.GetOptionBool(options.ProtoFieldOptFixed) {
			switch field.Type.KeyBase {
			case buildpb.BaseTypeDesc_Uint64, buildpb.BaseTypeDesc_Int64:
				desc = "protowire.Fixed64Type"
				typ = int(protowire.Fixed64Type)
				return
			case buildpb.BaseTypeDesc_Bool:
				desc = "protowire.VarintType"
				typ = int(protowire.VarintType)
				return
			case buildpb.BaseTypeDesc_Int32, buildpb.BaseTypeDesc_Uint32:
				desc = "protowire.Fixed32Type"
				typ = int(protowire.Fixed32Type)
				return
			}
		}
		desc = "protowire.VarintType"
		typ = int(protowire.VarintType)
	case field.IsList():
		if field.Options.GetOptionBool(options.ProtoFieldOptPacked) {
			desc = "protowire.BytesType"
			typ = int(protowire.BytesType)
			return
		}
		if field.Type.ElemCustom {
			desc = "protowire.BytesType"
			typ = int(protowire.BytesType)
			return
		}
		return switchProtoType(&buildpb.Field{
			Options: field.Options,
			Type: &buildpb.TypeDesc{
				Type:       buildpb.FieldType_BaseType,
				ElemCustom: false,
				KeyBase:    field.Type.KeyBase,
			},
		})
	default:
		// case field.IsMap():

		// case field.IsCustom():
		desc = "protowire.BytesType"
		typ = int(protowire.BytesType)
	}
	return
}

func switchProtoKind(field *buildpb.Field) protoreflect.Kind {
	switch {
	case field.IsBasicType():
		switch field.Type.KeyBase {
		case buildpb.BaseTypeDesc_String:
			return protoreflect.StringKind
		case buildpb.BaseTypeDesc_Binary:
			return protoreflect.BytesKind
		case buildpb.BaseTypeDesc_Bool:
			return protoreflect.BoolKind
		case buildpb.BaseTypeDesc_Float32:
			return protoreflect.FloatKind
		case buildpb.BaseTypeDesc_Float64:
			return protoreflect.DoubleKind
		}

		if field.Options.GetOptionBool(options.ProtoFieldOptFixed) {
			if field.Options.GetOptionBool(options.ProtoFieldOptSigned) {
				switch field.Type.KeyBase {
				case buildpb.BaseTypeDesc_Int32, buildpb.BaseTypeDesc_Uint32:
					return protoreflect.Sfixed32Kind
				case buildpb.BaseTypeDesc_Int64, buildpb.BaseTypeDesc_Uint64:
					return protoreflect.Sfixed64Kind
				}
			}
			switch field.Type.KeyBase {
			case buildpb.BaseTypeDesc_Int32, buildpb.BaseTypeDesc_Uint32:
				return protoreflect.Fixed32Kind
			case buildpb.BaseTypeDesc_Int64, buildpb.BaseTypeDesc_Uint64:
				return protoreflect.Fixed64Kind
			}
		}
		if field.Options.GetOptionBool(options.ProtoFieldOptSigned) {
			switch field.Type.KeyBase {
			case buildpb.BaseTypeDesc_Int32, buildpb.BaseTypeDesc_Uint32:
				return protoreflect.Sint32Kind
			case buildpb.BaseTypeDesc_Int64, buildpb.BaseTypeDesc_Uint64:
				return protoreflect.Sint64Kind
			}
		}

		switch field.Type.KeyBase {
		case buildpb.BaseTypeDesc_Int8, buildpb.BaseTypeDesc_Int16, buildpb.BaseTypeDesc_Int32:
			return protoreflect.Int32Kind
		case buildpb.BaseTypeDesc_Uint8, buildpb.BaseTypeDesc_Uint16, buildpb.BaseTypeDesc_Uint32:
			return protoreflect.Uint32Kind
		case buildpb.BaseTypeDesc_Int64:
			return protoreflect.Int64Kind
		case buildpb.BaseTypeDesc_Uint64:
			return protoreflect.Uint64Kind
		}

		return protoreflect.Int32Kind
	case field.IsMap():
		return protoreflect.MessageKind
	case field.IsCustom():
		return protoreflect.MessageKind
	case field.IsList():
		if field.Type.ElemCustom {
			return protoreflect.MessageKind
		}
		return switchProtoKind(&buildpb.Field{
			Options: field.Options,
			Type: &buildpb.TypeDesc{
				Type:       buildpb.FieldType_BaseType,
				Key:        field.Type.Key,
				ElemCustom: false,
				KeyBase:    field.Type.KeyBase,
			},
		})
	default:
		panic("invalid type")
	}
}

func parseFillBasicFiled(genField *gengo.GenerateField, field *buildpb.Field) {
	if field.IsCustom() {
		genField.CheckNotEmpty = func(x string) string {
			return x + " != nil"
		}
		genField.TemplateSize = "size.message"
		genField.TemplateEncode = "encode.message"
		genField.TemplateDecode = "decode.message"
		return
	}

	// if utils.Title(field.Name) == "OptionalSint32" {
	// 	log.Println(utils.Sdump(field, "OptionalSint32"))
	// 	log.Println(field.Options.GetOptionBool(options.ProtoFieldOptSigned))
	// }

	genField.CheckNotEmpty = func(x string) string {
		return x + " != 0"
	}
	switch field.Type.KeyBase {
	case buildpb.BaseTypeDesc_Bool:
		genField.CheckNotEmpty = func(x string) string {
			return x
		}
		genField.TemplateSize = "size.bool"
		genField.TemplateEncode = "encode.bool"
		genField.TemplateDecode = "decode.bool"

	case buildpb.BaseTypeDesc_Float32:
		genField.TemplateSize = "size.float"
		genField.TemplateEncode = "encode.float"
		genField.TemplateDecode = "decode.float"

		// g.Import(protogen.GoImportPath("math"))
		// g.QualifiedGoIdent(protogen.GoIdent{GoName: "Float32bits", GoImportPath: "math"})

	case buildpb.BaseTypeDesc_Float64:
		genField.TemplateSize = "size.double"
		genField.TemplateEncode = "encode.double"
		genField.TemplateDecode = "decode.double"

		// g.Import(protogen.GoImportPath("math"))
		// g.QualifiedGoIdent(protogen.GoIdent{GoName: "Float64bits", GoImportPath: "math"})
	case buildpb.BaseTypeDesc_String:
		genField.CheckNotEmpty = func(x string) string {
			return "len(" + x + ") > 0"
		}
		genField.TemplateSize = "size.string"
		genField.TemplateEncode = "encode.string"
		genField.TemplateDecode = "decode.string"

	case buildpb.BaseTypeDesc_Binary:
		genField.CheckNotEmpty = func(x string) string {
			return "len(" + x + ") > 0"
		}
		genField.TemplateSize = "size.bytes"
		genField.TemplateEncode = "encode.bytes"
		genField.TemplateDecode = "decode.bytes"

	default:
		if field.Options.GetOptionBool(options.ProtoFieldOptFixed) {
			switch field.Type.KeyBase {
			case buildpb.BaseTypeDesc_Int32, buildpb.BaseTypeDesc_Uint32:
				genField.TemplateSize = "size.fix32"
				genField.TemplateEncode = "encode.fix32"
				genField.TemplateDecode = "decode.fix32"
				return
			case buildpb.BaseTypeDesc_Int64, buildpb.BaseTypeDesc_Uint64:
				genField.TemplateSize = "size.fix64"
				genField.TemplateEncode = "encode.fix64"
				genField.TemplateDecode = "decode.fix64"
				return
			}
		}
		if field.Options.GetOptionBool(options.ProtoFieldOptSigned) {
			switch field.Type.KeyBase {
			case buildpb.BaseTypeDesc_Int32, buildpb.BaseTypeDesc_Int64:
				// log.Println("use sint")
				genField.TemplateSize = "size.sint"
				genField.TemplateEncode = "encode.sint"
				genField.TemplateDecode = "decode.sint"
				return
			}
		}

		genField.TemplateSize = "size.varint"
		genField.TemplateEncode = "encode.varint"
		genField.TemplateDecode = "decode.varint"
	}
}

func parseFillListPackedFiled(genField *gengo.GenerateField, field *buildpb.Field) {
	if field.Type.ElemCustom {
		genField.CheckNotEmpty = func(x string) string {
			return x + " != nil"
		}
		genField.TemplateSize = "size.packed.message"
		genField.TemplateEncode = "encode.packed.message"
		genField.TemplateDecode = "decode.slice.message"
		return
	}
	switch field.Type.KeyBase {
	case buildpb.BaseTypeDesc_Bool:
		genField.TemplateSize = "size.packed.bool"
		genField.TemplateEncode = "encode.packed.bool"
		genField.TemplateDecode = "decode.slice.bool"

	case buildpb.BaseTypeDesc_Float32:
		genField.TemplateSize = "size.packed.float"
		genField.TemplateEncode = "encode.packed.float"
		genField.TemplateDecode = "decode.slice.float"

		// g.Import(protogen.GoImportPath("math"))
		// g.QualifiedGoIdent(protogen.GoIdent{GoName: "Float32bits", GoImportPath: "math"})

	case buildpb.BaseTypeDesc_Float64:
		genField.TemplateSize = "size.packed.double"
		genField.TemplateEncode = "encode.packed.double"
		genField.TemplateDecode = "decode.slice.double"

		// g.Import(protogen.GoImportPath("math"))
		// g.QualifiedGoIdent(protogen.GoIdent{GoName: "Float64bits", GoImportPath: "math"})
	case buildpb.BaseTypeDesc_String:
		genField.CheckNotEmpty = func(x string) string {
			return "len(" + x + ") > 0"
		}
		genField.TemplateSize = "size.packed.string"
		genField.TemplateEncode = "encode.packed.string"
		genField.TemplateDecode = "decode.slice.string"

	case buildpb.BaseTypeDesc_Binary:
		genField.CheckNotEmpty = func(x string) string {
			return "len(" + x + ") > 0"
		}
		genField.TemplateSize = "size.packed.bytes"
		genField.TemplateEncode = "encode.packed.bytes"
		genField.TemplateDecode = "decode.slice.bytes"
	default:
		if field.Options.GetOptionBool(options.ProtoFieldOptFixed) {
			switch field.Type.KeyBase {
			case buildpb.BaseTypeDesc_Int32, buildpb.BaseTypeDesc_Uint32:
				genField.TemplateSize = "size.packed.fix32"
				genField.TemplateEncode = "encode.packed.fix32"
				genField.TemplateDecode = "decode.slice.fix32"
				return
			case buildpb.BaseTypeDesc_Int64, buildpb.BaseTypeDesc_Uint64:
				genField.TemplateSize = "size.packed.fix64"
				genField.TemplateEncode = "encode.packed.fix64"
				genField.TemplateDecode = "decode.slice.fix64"
				return
			}
		}
		if field.Options.GetOptionBool(options.ProtoFieldOptSigned) {
			switch field.Type.KeyBase {
			case buildpb.BaseTypeDesc_Int32, buildpb.BaseTypeDesc_Int64:
				genField.TemplateSize = "size.packed.sint"
				genField.TemplateEncode = "encode.packed.sint"
				genField.TemplateDecode = "decode.slice.sint"
				return
			}
		}

		genField.TemplateSize = "size.packed.varint"
		genField.TemplateEncode = "encode.packed.varint"
		genField.TemplateDecode = "decode.slice.varint"
	}
}

func parseFillListNoPackedFiled(genField *gengo.GenerateField, field *buildpb.Field) {
	if field.Type.ElemCustom {
		genField.CheckNotEmpty = func(x string) string {
			return x + " != nil"
		}
		genField.TemplateSize = "size.nopack.message"
		genField.TemplateEncode = "encode.nopack.message"
		genField.TemplateDecode = "decode.slice.message"
		return
	}
	switch field.Type.KeyBase {
	case buildpb.BaseTypeDesc_Bool:
		genField.TemplateSize = "size.nopack.bool"
		genField.TemplateEncode = "encode.nopack.bool"
		genField.TemplateDecode = "decode.slice.bool"

	case buildpb.BaseTypeDesc_Float32:
		genField.TemplateSize = "size.nopack.float"
		genField.TemplateEncode = "encode.nopack.float"
		genField.TemplateDecode = "decode.slice.float"

		// g.Import(protogen.GoImportPath("math"))
		// g.QualifiedGoIdent(protogen.GoIdent{GoName: "Float32bits", GoImportPath: "math"})

	case buildpb.BaseTypeDesc_Float64:
		genField.TemplateSize = "size.nopack.double"
		genField.TemplateEncode = "encode.nopack.double"
		genField.TemplateDecode = "decode.slice.double"

		// g.Import(protogen.GoImportPath("math"))
		// g.QualifiedGoIdent(protogen.GoIdent{GoName: "Float64bits", GoImportPath: "math"})
	case buildpb.BaseTypeDesc_String:
		genField.CheckNotEmpty = func(x string) string {
			return "len(" + x + ") > 0"
		}
		genField.TemplateSize = "size.nopack.string"
		genField.TemplateEncode = "encode.nopack.string"
		genField.TemplateDecode = "decode.slice.string"

	case buildpb.BaseTypeDesc_Binary:
		genField.CheckNotEmpty = func(x string) string {
			return "len(" + x + ") > 0"
		}
		genField.TemplateSize = "size.nopack.bytes"
		genField.TemplateEncode = "encode.nopack.bytes"
		genField.TemplateDecode = "decode.slice.bytes"
	default:
		if field.Options.GetOptionBool(options.ProtoFieldOptFixed) {
			switch field.Type.KeyBase {
			case buildpb.BaseTypeDesc_Int32, buildpb.BaseTypeDesc_Uint32:
				genField.TemplateSize = "size.nopack.fix32"
				genField.TemplateEncode = "encode.nopack.fix32"
				genField.TemplateDecode = "decode.slice.fix32"
				return
			case buildpb.BaseTypeDesc_Int64, buildpb.BaseTypeDesc_Uint64:
				genField.TemplateSize = "size.nopack.fix64"
				genField.TemplateEncode = "encode.nopack.fix64"
				genField.TemplateDecode = "decode.slice.fix64"
				return
			}
		}
		if field.Options.GetOptionBool(options.ProtoFieldOptSigned) {
			switch field.Type.KeyBase {
			case buildpb.BaseTypeDesc_Int32, buildpb.BaseTypeDesc_Int64:
				genField.TemplateSize = "size.nopack.sint"
				genField.TemplateEncode = "encode.nopack.sint"
				genField.TemplateDecode = "decode.slice.sint"
				return
			}
		}

		genField.TemplateSize = "size.nopack.varint"
		genField.TemplateEncode = "encode.nopack.varint"
		genField.TemplateDecode = "decode.slice.varint"
	}
}
