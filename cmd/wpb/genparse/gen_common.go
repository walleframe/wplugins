package genparse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aggronmagi/wplugins/buildpb"
	"github.com/aggronmagi/wplugins/gen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Field numbers for google.protobuf.FileDescriptorProto.
const (
	FileDescriptorProto_Name_field_number             protoreflect.FieldNumber = 1
	FileDescriptorProto_Package_field_number          protoreflect.FieldNumber = 2
	FileDescriptorProto_Dependency_field_number       protoreflect.FieldNumber = 3
	FileDescriptorProto_PublicDependency_field_number protoreflect.FieldNumber = 10
	FileDescriptorProto_WeakDependency_field_number   protoreflect.FieldNumber = 11
	FileDescriptorProto_MessageType_field_number      protoreflect.FieldNumber = 4
	FileDescriptorProto_EnumType_field_number         protoreflect.FieldNumber = 5
	FileDescriptorProto_Service_field_number          protoreflect.FieldNumber = 6
	FileDescriptorProto_Extension_field_number        protoreflect.FieldNumber = 7
	FileDescriptorProto_Options_field_number          protoreflect.FieldNumber = 8
	FileDescriptorProto_SourceCodeInfo_field_number   protoreflect.FieldNumber = 9
	FileDescriptorProto_Syntax_field_number           protoreflect.FieldNumber = 12
	FileDescriptorProto_Edition_field_number          protoreflect.FieldNumber = 13
)

// structTags is a data structure for build idiomatic Go struct tags.
// Each [2]string is a key-value pair, where value is the unescaped string.
//
// Example: structTags{{"key", "value"}}.String() -> `key:"value"`
type structTags [][2]string

func (tags structTags) String() string {
	if len(tags) == 0 {
		return ""
	}
	var ss []string
	for _, tag := range tags {
		// NOTE: When quoting the value, we need to make sure the backtick
		// character does not appear. Convert all cases to the escaped hex form.
		key := tag[0]
		val := strings.Replace(strconv.Quote(tag[1]), "`", `\x60`, -1)
		ss = append(ss, fmt.Sprintf("%s:%s", key, val))
	}
	return "`" + strings.Join(ss, " ") + "`"
}

// appendDeprecationSuffix optionally appends a deprecation notice as a suffix.
func appendDeprecationSuffix(prefix protogen.Comments, parentFile protoreflect.FileDescriptor, deprecated bool) protogen.Comments {
	fileDeprecated := parentFile.Options().(*descriptorpb.FileOptions).GetDeprecated()
	if !deprecated && !fileDeprecated {
		return prefix
	}
	if prefix != "" {
		prefix += "\n"
	}
	if fileDeprecated {
		return prefix + " Deprecated: The entire proto file " + protogen.Comments(parentFile.Path()) + " is marked as deprecated.\n"
	}
	return prefix + " Deprecated: Marked as deprecated in " + protogen.Comments(parentFile.Path()) + ".\n"
}

// trailingComment is like protogen.Comments, but lacks a trailing newline.
type trailingComment protogen.Comments

func (c trailingComment) String() string {
	s := strings.TrimSuffix(protogen.Comments(c).String(), "\n")
	if strings.Contains(s, "\n") {
		// We don't support multi-lined trailing comments as it is unclear
		// how to best render them in the generated code.
		return ""
	}
	return s
}

func fieldJSONTagValue(field *buildpb.Field) string {
	return string(field.Name) + ",omitempty"
}

func fieldBasicGoType(typ buildpb.BaseTypeDesc) string {
	switch typ {
	case buildpb.BaseTypeDesc_Int8:
		return "int8"
	case buildpb.BaseTypeDesc_Uint8:
		return "uint8"
	case buildpb.BaseTypeDesc_Int16:
		return "int16"
	case buildpb.BaseTypeDesc_Uint16:
		return "unt16"
	case buildpb.BaseTypeDesc_Int32:
		return "int32"
	case buildpb.BaseTypeDesc_Uint32:
		return "uint32"
	case buildpb.BaseTypeDesc_Int64:
		return "int64"
	case buildpb.BaseTypeDesc_Uint64:
		return "uint64"
	case buildpb.BaseTypeDesc_String:
		return "string"
	case buildpb.BaseTypeDesc_Binary:
		return "[]byte"
	case buildpb.BaseTypeDesc_Bool:
		return "bool"
	case buildpb.BaseTypeDesc_Float32:
		return "float32"
	case buildpb.BaseTypeDesc_Float64:
		return "float64"
	default:
		return ""
	}
}

// fieldGoType returns the Go type used for a field.
//
// If it returns pointer=true, the struct field is a pointer to the type.
func fieldGoType(g *gen.Generator, field *buildpb.Field) (goType string, pointer bool) {
	// if field.Desc.IsWeak() {
	// 	return "struct{}", false
	// }

	pointer = false // field.Desc.HasPresence()
	switch field.GetType().GetType() {
	case buildpb.FieldType_BaseType:
		goType = fieldBasicGoType(field.Type.KeyBase)
	case buildpb.FieldType_CustomType:
		pointer = true
		if strings.Contains(field.Type.Key, ".") {
			typs := strings.SplitN(field.Type.Key, ".", 2)
			goType = typs[0] + "." + g.Key(typs[1])
		} else {
			goType = g.Key(field.Type.Key)
		}
	case buildpb.FieldType_ListType:
		if field.Type.ElemCustom {
			goType = "[]*" + g.Key(field.Type.Key)
		} else {
			goType = "[]" + fieldBasicGoType(field.Type.KeyBase)
		}
	case buildpb.FieldType_MapType:
		goType = "map["
		goType += fieldBasicGoType(field.Type.KeyBase)
		goType += "]"
		if field.GetType().ElemCustom {
			goType += "*" + g.Key(field.Type.Value)
		} else {
			goType += fieldBasicGoType(field.Type.ValueBase)
		}
	}
	return goType, pointer
}

func fieldDefaultValue(g *gen.Generator, m *buildpb.MsgDesc, field *buildpb.Field) string {
	switch {
	case field.IsList(), field.IsCustom(), field.IsMap():
		return "nil"
	default: // basic type
		switch field.Type.KeyBase {
		case buildpb.BaseTypeDesc_Binary:
			return "nil"
		case buildpb.BaseTypeDesc_String:
			return `""`
		case buildpb.BaseTypeDesc_Bool:
			return "false"
		default:
			return "0"
		}
	}
	// TODO: 枚举类型
}
