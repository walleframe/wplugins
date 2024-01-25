package main

import (
	"log"
	"path/filepath"

	"github.com/walleframe/wplugins/buildpb"
	"github.com/walleframe/wplugins/gen"
	"github.com/walleframe/wplugins/options"
	"github.com/walleframe/wplugins/utils"
	"github.com/walleframe/wplugins/utils/plugin"
)

var GenerateEmptyMessage = false
var cfg = struct {
	SvcPkg   string
	CodePath string
	Charset  string
	Collate  string
}{
	SvcPkg:   "github.com/walleframe/svc_db",
	CodePath: "pkg/dbop/",
	Charset:  "utf8mb4",
	Collate:  "utf8mb4_general_ci",
}

func init() {
	// 如果环境变量设置了值, 读取作为为默认值. 优先使用传递的参数
	utils.GetEnvString("WDB_COLLATE", &cfg.Collate)
	utils.GetEnvString("WDB_CHARSET", &cfg.Charset)
	utils.GetEnvString("WDB_OPCODE_PATH", &cfg.CodePath)
	utils.GetEnvString("WDB_SVC_PKG", &cfg.SvcPkg)
}

func main() {
	plugin.MainRangeMessage(
		func(file *buildpb.FileDesc) bool {
			return file.GetString(options.SqlDBName, "") != ""
		},
		func(msg *buildpb.MsgDesc) bool {
			if len(msg.Fields) < 1 {
				return false
			}
			if msg.HasOption(options.SqlIgnore) {
				return false 
			}
			return true
		},
		generateWalleDB,
	)
}

func generateWalleDB(msg *buildpb.MsgDesc, prog *buildpb.FileDesc, depend map[string]*buildpb.FileDesc) (out []*buildpb.BuildOutput, err error) {
	dbName := prog.GetString(options.SqlDBName, "")
	tblName := msg.GetString(options.SqlTableName, msg.Name)
	engine := msg.GetString(options.SqlEngine, "InnoDB")

	tbl := &SqlTable{
		DB:       dbName,
		SqlTable: tblName,
		Name:     msg.Name,
		Struct:   prog.Pkg.Package + "." + msg.Name,
		SvcDB:    filepath.Base(cfg.SvcPkg),
		Charset:  cfg.Charset,
		Collate:  cfg.Collate,
		Engine:   engine,
	}

	// 分析表定义
	err = ParseSqlTableColumns(msg, tbl)
	if err != nil {
		return nil, err
	}

	tpl := gen.NewTemplate("wdb")
	tpl.AddImportFunc(tbl)
	err = tpl.Parse(wdbTpl)
	if err != nil {
		return nil, err
	}

	for _, pkg := range []string{
		"context",
		"database/sql",
		"errors",
		"strings",
		"fmt",
		"sync/atomic",
		"github.com/jmoiron/sqlx",
		"github.com/walleframe/walle/util",
	} {
		tbl.Import(pkg, "pkg")
	}
	tbl.Import(cfg.SvcPkg, "svc_db")
	if pkg, ok := prog.Options.GetStringCheck(options.ProtoGoPkg); ok {
		tbl.Import(pkg, "")
	}

	data, err := tpl.Exec(tbl)
	if err != nil {
		return nil, err
	}

	out = append(out, &buildpb.BuildOutput{
		File: filepath.Join(cfg.CodePath, dbName, msg.Name+".dbop.go"),
		Data: data,
	})

	log.Println(dbName, msg.Name)

	return
}
