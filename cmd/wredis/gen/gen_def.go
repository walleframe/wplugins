package gen

import (
	"strings"

	"github.com/walleframe/wplugins/buildpb"
	"github.com/walleframe/wplugins/cmd/wredis/keyarg"
)

// RedisObject 生成redis消息的对象
type RedisObject struct {
	// innerImportPkgs 导入包
	innerImportPkgs []string

	// 生成包名
	Package string
	// 操作名,消息名
	Name string
	// 消息体注释
	Doc *buildpb.DocDesc
	// key 参数
	Args []keyarg.KeyArg
	// service package
	SvcPkg string
	// walle package
	WPbPkg string
	//
	KeySize int

	// redis 脚本
	Scripts []*RedisScript

	// key 相关接口
	TypeKeys bool
	// string 生成
	TypeString *RedisTypeString
	// hash 生成
	TypeHash *RedisTypeHash
	// set 生成
	TypeSet *RedisTypeSet
	// zset 生成
	TypeZSet *RedisTypeZSet
	// redis 分布式锁
	Lock bool
}

func (obj *RedisObject) Import(pkg, fun string) (_ string) {
	for _, v := range obj.innerImportPkgs {
		if v == pkg {
			return
		}
	}
	obj.innerImportPkgs = append(obj.innerImportPkgs, pkg)
	return
}

func (g *RedisObject) customImport() string {
	buf := strings.Builder{}
	buf.Grow(len(g.innerImportPkgs) * 30)
	for _, v := range g.innerImportPkgs {
		buf.WriteByte('\n')
		buf.WriteByte('\t')
		buf.WriteString(`"`)
		buf.WriteString(v)
		buf.WriteByte('"')

	}
	return buf.String()
}

type RedisTypeString struct {
	Type     string
	Signed   bool
	Number   bool
	String   bool
	Protobuf bool // github.com/gogo/protobuf/proto
	WProto   bool // github.com/walleframe/walle/process/message
	Custom   bool
	Float    bool
}

type RedisTypeHash struct {
	HashObject  *RedisHashObject
	HashDynamic *RedisHashDynamic
}

type RedisGenType struct {
	// 字段名
	Name string
	// 字段类型
	Type string
	// Number
	Number bool
	//
	RedisFunc string
}

func (x *RedisGenType) IsFloat() bool {
	return strings.HasPrefix(x.Type, "float")
}

func (x *RedisGenType) IsInt() bool {
	return strings.Contains(x.Type, "int")
}

type RedisHashObject struct {
	Fields  []*RedisGenType
	Name    string
	Type    string
	HGetAll bool
}

type RedisHashDynamic struct {
	Field     *RedisGenType
	Value     *RedisGenType
	GenMap    bool
	FieldArgs []keyarg.KeyArg
	ValueArgs []keyarg.KeyArg
}

type RedisGenMsg struct {
	Type      string
	Marshal   func(objName string) string
	Unmarshal func(objName, paramName string) string
	New       string
}

type RedisTypeSet struct {
	BaseType *RedisGenType
	Message  *RedisGenMsg
}

type RedisTypeZSet struct {
	// score
	Score *RedisGenType
	// mem
	Member *RedisGenType
	// 拼接string做field
	Args []keyarg.KeyArg
	// 使用消息作为field
	Message *RedisGenMsg
}

// RedisScript redis脚本
type RedisScript struct {
	// 脚本操作名称
	Name string
	// 脚本数据
	Script string
	// 脚本参数
	Args []keyarg.KeyArg
	// 脚本输出
	Output []keyarg.KeyArg
	//
	TemplateName string
	//
	CommandName string
}
