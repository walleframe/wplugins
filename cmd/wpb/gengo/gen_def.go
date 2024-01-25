package gengo

import (
	"strconv"
	"strings"

	"github.com/walleframe/wplugins/buildpb"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type GenerateStruct struct {
	// 版本信息
	VersionInfo func() string
	// 包名
	Package string
	// 所有待生成的枚举
	Enums []*GenerateEnums
	// 所有待生成的消息
	Messages []*GenerateMessage
	// 导入
	ImportPackages []string
	//
	WirePkg string
}

func (g *GenerateStruct) Import(pkg, name string) string {
	for _, v := range g.ImportPackages {
		if v == pkg {
			return pkg
		}
	}
	g.ImportPackages = append(g.ImportPackages, pkg)
	return pkg
}

func (g *GenerateStruct) customImport() string {
	buf := strings.Builder{}
	buf.Grow(len(g.ImportPackages) * 30)
	for _, v := range g.ImportPackages {
		buf.WriteByte('\n')
		buf.WriteByte('\t')
		buf.WriteString(`"`)
		buf.WriteString(v)
		buf.WriteByte('"')

	}
	return buf.String()
}

type GenerateDoc struct {
	// 前置注释
	LeadingComments string
	// 尾注释
	TrailingComment string
}

type GenerateEnums struct {
	GenerateDoc
	// 类型名
	TypeName string
	// go里面的名字
	GoName string
	// 枚举值
	Values []*GenerateEnumValue
}

type GenerateEnumValue struct {
	GenerateDoc
	// enum值的名字
	ValueName string
	// 描述
	Desc string
	// 编号
	Num int32
	//
	Duplicate string
}

type GenerateMessage struct {
	GenerateDoc
	// 类型名
	TypeName string
	// go里面的名字
	GoName string
	// 字段
	Fields []*GenerateField
	// 生成get方法
	GenGetter bool
	// 自定义模板列表
	CustomTemplates []string
}

type GenerateField struct {
	GenerateDoc
	// 类型名
	TypeName string
	// go里面的名字
	GoName string
	// 元素类型. 原始类型
	GoType string
	// 错误提示. 消息.字段
	Tip string

	// tags
	tags string
	// getter 辅助
	GetNilCheck  bool
	DefaultValue string

	// protobuf field num
	DescNum int
	// PB 中的名字
	DescName string
	// proto-wire 类型 -string
	WireType string // protowire.VarintType
	// proto-wire 类型
	DescType int
	// 类型属性
	IsMap  bool
	IsList bool
	Kind   protoreflect.Kind
	WType *buildpb.TypeDesc

	// marshal 辅助
	CheckNotEmpty func(vname string) string // 检测是否为空的条件. 是否需要序列化
	// 模板名称
	TemplateEncode string
	TemplateSize   string
	TemplateDecode string
	//
	MapKey   *GenerateField
	MapValue *GenerateField
}

func (f *GenerateField) AddTag(tag, usefun string, vals ...string) (_ string) {
	val := f.DescName
	if len(vals) > 0 {
		val = strings.Join(vals, ",")
	}
	switch usefun {
	case "":
	case "go":
		val = GoCamelCase(val)
	case "snake":
		val = JSONSnakeCase(val)
	case "json":
		val = JSONCamelCase(val)
	}
	switch tag {
	case "json":
		val += ",omitempty"
	}
	if len(f.tags) > 0 {
		f.tags += " "
	}
	val = strings.Replace(strconv.Quote(val), "`", `\x60`, -1)
	f.tags += tag + ":" + val
	return f.tags
}

func (f *GenerateField) Tags() string {
	if len(f.tags) == 0 {
		return ""
	}
	return "`" + f.tags + "`"
}
