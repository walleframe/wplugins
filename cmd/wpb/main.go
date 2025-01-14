package main

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/walleframe/wplugins/buildpb"
	"github.com/walleframe/wplugins/cmd/wpb/gengo"
	"github.com/walleframe/wplugins/cmd/wpb/genparse"
	"github.com/walleframe/wplugins/gen"
	"github.com/walleframe/wplugins/options"
	"github.com/walleframe/wplugins/utils"
	"github.com/walleframe/wplugins/utils/plugin"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"
)

var GenerateEmptyMessage = false

func init() {
	// 如果环境变量设置了值, 读取作为为默认值. 优先使用传递的参数
	utils.GetEnvString("GOPB_WIRE_PACKAGE", &genparse.WirePkg)
	utils.GetEnvBool("GOPB_GEN_GET", &genparse.Getter)
	utils.GetEnvBool("GOPB_GEN_ZAP", &genparse.Zap)
	utils.GetEnvBool("GOPB_GEN_EMPTY_MSG", &GenerateEmptyMessage)
}

func main() {
	plugin.MainRangeFile(nil, generateWalleProrobuf)
}

func generateWalleProrobuf(prog *buildpb.FileDesc, depend map[string]*buildpb.FileDesc) (out []*buildpb.BuildOutput, err error) {

	g := gen.New(
		gen.WithFormat(gen.GoFormat2),
		gen.WithIndent("    "),
		gen.WithKeyTitle(true),
	)

	// parse
	data := &gengo.GenerateStruct{}
	data.Package = string(prog.Pkg.Package)
	// log.Println("package:[", data.Package, "]")
	// 打印版本信息
	data.VersionInfo = func() string {
		buf := bytes.Buffer{}

		buf.WriteString("// Code generated by wpb. DO NOT EDIT.\n")

		// buf.WriteString("// versions:\n")
		// protocGenGoVersion := genparse.Version
		// protocVersion := "(unknown)"
		// if v := gen.Request.GetCompilerVersion(); v != nil {
		// 	protocVersion = fmt.Sprintf("v%v.%v.%v", v.GetMajor(), v.GetMinor(), v.GetPatch())
		// 	if s := v.GetSuffix(); s != "" {
		// 		protocVersion += "-" + s
		// 	}
		// }
		// buf.WriteString("// \tprotoc-gen-gopb ")
		// buf.WriteString(protocGenGoVersion)
		// buf.WriteByte('\n')
		// buf.WriteString("// \tprotoc          ")
		// buf.WriteString(protocVersion)
		// buf.WriteByte('\n')

		// if f.Proto.GetOptions().GetDeprecated() {
		// 	buf.WriteString(fmt.Sprintln("// ", f.Desc.Path(), " is a deprecated file."))
		// } else {
		// 	buf.WriteString(fmt.Sprintln("// source: ", f.Desc.Path()))
		// }
		// buf.WriteByte('\n')

		return buf.String()
	}
	genCount := 0
	for _, e := range prog.Enums {
		// 默认不生成空枚举
		if GenerateEmptyMessage == false && len(e.Values) < 1 {
			continue
		}
		err = multierr.Append(err, genparse.ParseEnum(data, g, e))
		genCount++
	}

	genDB := prog.HasOption(options.SqlDBName)

	for _, m := range prog.Msgs {
		err = multierr.Append(err, genparse.ParseMessage(data, g, m, func(m *buildpb.MsgDesc) bool {
			// 默认不生成空消息
			if GenerateEmptyMessage == false && len(m.Fields) < 1 {
				return false
			}
			// 过滤不生成redis定义消息 - 无意义
			if m.HasOption(options.RedisOpKey) {
				return false
			}
			// 手动指定忽略消息生成
			if m.Options.GetOptionBool(options.WPBIngore) {
				return false
			}
			genCount++
			return true
		}))
		// sql 额外生成一个Ex消息,附带 modify_stamp,create_stamp
		if genDB && !m.HasOption(options.SqlIgnore) {
			// 设置sql.ex == false,不生成ex
			if m.HasOption(options.SqlExSwitch) && !m.Options.GetOptionBool(options.SqlExSwitch) {
				continue
			}
			No := int32(len(m.Fields))
			sqlEx := proto.Clone(m).(*buildpb.MsgDesc)
			sqlEx.Name += "_ex"
			sqlEx.Doc = &buildpb.DocDesc{
				Doc: []string{"// sql row extern data"},
			}
			sqlEx.Fields = append(sqlEx.Fields,
				&buildpb.Field{
					Type: &buildpb.TypeDesc{
						Type:    buildpb.FieldType_BaseType,
						Key:     "int64",
						KeyBase: buildpb.BaseTypeDesc_Int64,
					},
					Name: "modify_stamp",
					No:   No + 1,
				},
				&buildpb.Field{
					Type: &buildpb.TypeDesc{
						Type:    buildpb.FieldType_BaseType,
						Key:     "int64",
						KeyBase: buildpb.BaseTypeDesc_Int64,
					},
					Name: "create_stamp",
					No:   No + 2,
				},
			)
			err = multierr.Append(err, genparse.ParseMessage(data, g, sqlEx, func(m *buildpb.MsgDesc) bool {
				return true
			}))
		}
	}
	if err != nil {
		return
	}
	// 没有生成任何消息,直接返回.不生成文件
	if genCount == 0 {
		return
	}
	// import others
	for _, imp := range prog.Imports {
		dep, ok := depend[imp.File]
		if !ok {
			err = multierr.Append(err, fmt.Errorf("%s import %s, but not found %s", prog.File, imp.File, imp.File))
			continue
		}
		pkg, ok := dep.Options.GetStringCheck(options.ProtoGoPkg)
		if !ok {
			err = multierr.Append(err, fmt.Errorf("%s import %s, but %s not set '%s' option", prog.File, imp.File, imp.File, options.ProtoGoPkg))
			continue
		}
		data.Import(pkg, "import others")
	}
	//
	if err != nil {
		return
	}
	// generate
	buf, err := gengo.GenExec(data)
	if err != nil {
		return
	}
	// output
	g.P(string(buf))

	bdata, err := g.Bytes()
	if err != nil {
		log.Println("format code failed.", err)
		log.Println(string(g.Buffer.Bytes()))
		//err = nil
		return nil, fmt.Errorf("format failed %w", err)
		// data = g.Buffer.Bytes()
		// err = nil
	}
	out = append(out, &buildpb.BuildOutput{
		File: strings.TrimSuffix(prog.File, filepath.Ext(prog.File)) + ".wpb.go",
		Data: bdata,
	})

	return
}
