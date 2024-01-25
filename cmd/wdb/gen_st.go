package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/walleframe/wplugins/buildpb"
	"github.com/walleframe/wplugins/gen"
	"github.com/walleframe/wplugins/options"
	"github.com/walleframe/wplugins/utils"
	"go.uber.org/multierr"
)

type SqlTable struct {
	gen.GoObject
	SvcDB      string
	Charset    string
	Collate    string
	Engine     string
	DB         string
	SqlTable   string
	Name       string //
	Struct     string // 结构体名
	allColumns []*SqlColumn
	AutoIncr   *SqlColumn
	PrimaryKey []*SqlColumn
	Columns    []*SqlColumn
	Index      []*SqlIndex
}

func (tbl *SqlTable) AllColumns(ignoreAutoID bool) []*SqlColumn {
	if !ignoreAutoID || tbl.AutoIncr == nil {
		return tbl.allColumns
	}
	tmp := make([]*SqlColumn, 0, len(tbl.allColumns)-1)
	for _, col := range tbl.allColumns {
		if col.Name == tbl.AutoIncr.Name {
			continue
		}
		tmp = append(tmp, col)
	}
	return tmp
}

func (tbl *SqlTable) Placeholder(ignoreAutoID bool) string {
	count := len(tbl.allColumns)
	if ignoreAutoID && tbl.AutoIncr != nil {
		count--
	}

	buf := strings.Builder{}
	buf.Grow(count * 2)
	buf.WriteByte('?')
	for k := 1; k < count; k++ {
		buf.Write([]byte(",?"))
	}

	return buf.String()
}

type SqlColumn struct {
	Name      string
	GoType    string
	SqlType   string
	Marshal   string
	Unmarshal string
	Doc       *buildpb.DocDesc
}

func (col *SqlColumn) SqlName() string {
	return "`" + col.Name + "`"
}

func (col *SqlColumn) Cond(name string) string {
	if col.GoType == "string" || strings.HasPrefix(col.GoType, "map") || strings.HasPrefix(col.GoType, "[]") || strings.HasPrefix(col.GoType, "*") {
		return fmt.Sprintf("StringCondition[%sWhereStmt]", utils.Title(name))
	}
	if strings.HasPrefix(col.GoType, "uint") {
		return fmt.Sprintf("IntUnSignedCondition[%sWhereStmt, %s]", utils.Title(name), col.GoType)
	}
	return fmt.Sprintf("IntSignedCondition[%sWhereStmt, %s]", utils.Title(name), col.GoType)
}

type SqlIndex struct {
	Name    string
	Columns []*SqlColumn
	unique  bool
}

func (idx *SqlIndex) Unique() string {
	if idx.unique {
		return "unique "
	}
	return ""
}

type sqlIdxTmp struct {
	name    string
	columns []string
	unique  bool
}

func (x *sqlIdxTmp) hasCol(col string) bool {
	if x == nil {
		return false
	}
	for _, v := range x.columns {
		if v == col {
			return true
		}
	}
	return false
}

func ParseSqlTableColumns(msg *buildpb.MsgDesc, tbl *SqlTable) (err error) {
	// 所有列, 方便后续检测索引等是否设置了无效的列名
	colCheck := make(map[string]*SqlColumn)
	for _, field := range msg.Fields {
		colCheck[field.Name] = &SqlColumn{}
	}
	// primary key
	pk := &sqlIdxTmp{}
	// allindex
	idx := make([]*sqlIdxTmp, 0, 4)
	// 分析消息级primary key 设置
	optPkStr := msg.GetString(options.SqlPrimaryKey2, "")
	optPkStr = msg.GetString(options.SqlPrimaryKey, optPkStr)
	if optPkStr != "" {
		pk.unique = true // 这里用unique字段表示已经使用消息级选项设置.
		for _, col := range strings.Split(optPkStr, ",") {
			col = strings.TrimSpace(col)
			if _, ok := colCheck[col]; !ok {
				err = multierr.Append(err, fmt.Errorf("%s primary key has invalid columns[%s]", msg.Name, col))
			}
			pk.columns = append(pk.columns, col)
		}
	}

	// 分析消息级的索引设置
	for name, v := range msg.Options.Options {
		if strings.HasPrefix(name, options.SqlIndexPrefix) {
			name = strings.TrimPrefix(name, options.SqlIndexPrefix)
			x := &sqlIdxTmp{name: name}
			for _, col := range strings.Split(v.Value, ",") {
				col = strings.TrimSpace(col)
				x.columns = append(x.columns, col)
				if _, ok := colCheck[col]; !ok {
					err = multierr.Append(err, fmt.Errorf("%s index [%s] has invalid columns[%s]", msg.Name, name, col))
				}
			}
			idx = append(idx, x)
		}
		if strings.HasPrefix(name, options.SqlUniquePrefix) {
			name = strings.TrimPrefix(name, options.SqlUniquePrefix)
			x := &sqlIdxTmp{name: name, unique: true}
			for _, col := range strings.Split(v.Value, ",") {
				col = strings.TrimSpace(col)
				x.columns = append(x.columns, col)
				if _, ok := colCheck[col]; !ok {
					err = multierr.Append(err, fmt.Errorf("%s unique-index [%s] has invalid columns[%s]", msg.Name, name, col))
				}
			}
			idx = append(idx, x)
		}
	}

	// 分析字段
	for _, field := range msg.Fields {
		col := &SqlColumn{
			Name: field.Name,
			Doc:  field.Doc,
		}
		setCustom := field.Options.GetOptionBool(options.SqlCustomSet)
		null := field.Options.GetOptionBool(options.SqlNull)
		autoIncr := field.Options.GetOptionBool(options.SqlAutoIncrement)
		sqlType := field.GetString(options.SqlFieldType, "")
		if field.Options.GetOptionBool(options.SqlPrimaryKey) || field.Options.GetOptionBool(options.SqlPrimaryKey2) {
			if pk.unique {
				log.Printf("WARN: %s.%s set option 'sql.pk' both message level and field level\n", msg.Name, field.Name)
			} else {
				pk.columns = append(pk.columns, field.Name)
			}
		}
		defVal := ""
		switch field.Type.Type {
		case buildpb.FieldType_BaseType:
			col.GoType = field.Type.Key
			col.Unmarshal = utils.Title(col.GoType)
			if setCustom {
				col.Marshal = fmt.Sprintf("Any[ %s ]", col.GoType)
			}
			switch field.Type.KeyBase {
			case buildpb.BaseTypeDesc_Int8, buildpb.BaseTypeDesc_Bool:
				col.SqlType = "tinyint"
				defVal = "default 0"
			case buildpb.BaseTypeDesc_Int16, buildpb.BaseTypeDesc_Int32:
				col.SqlType = "int"
				defVal = "default 0"
			case buildpb.BaseTypeDesc_Int64:
				col.SqlType = "bigint"
				defVal = "default 0"
			case buildpb.BaseTypeDesc_Uint8:
				col.SqlType = "tinyint unsigned"
				defVal = "default 0"
			case buildpb.BaseTypeDesc_Uint16, buildpb.BaseTypeDesc_Uint32:
				col.SqlType = "int unsigned"
				defVal = "default 0"
			case buildpb.BaseTypeDesc_Uint64:
				col.SqlType = "bigint unsigned"
				defVal = "default 0"
			case buildpb.BaseTypeDesc_String:
				col.SqlType = fmt.Sprintf("varchar(%d)", field.Options.GetInt64(options.SqlSize, 64))
				defVal = "default ''"
			case buildpb.BaseTypeDesc_Binary:
				col.SqlType = "blob"
			default:
				err = multierr.Append(err, fmt.Errorf("%s.%s type not support now.[%+v]", msg.Name, field.Name, field.Type))
			}
		case buildpb.FieldType_CustomType:
			col.GoType = field.Type.Key
			col.Unmarshal = fmt.Sprintf("Object[%s]", strings.TrimPrefix(col.GoType, "*"))
			if setCustom {
				col.Marshal = col.Unmarshal
			}
			col.SqlType = fmt.Sprintf("varchar(%d)", field.Options.GetInt64(options.SqlSize, 256))
			defVal = "default ''"
		case buildpb.FieldType_ListType:
			col.GoType = "[]" + field.Type.Key
			col.Unmarshal = fmt.Sprintf("Slice[%s]", col.GoType)
			if setCustom {
				col.Marshal = col.Unmarshal
			}
			col.SqlType = fmt.Sprintf("varchar(%d)", field.Options.GetInt64(options.SqlSize, 256))
			defVal = "default ''"
		case buildpb.FieldType_MapType:
			col.GoType = fmt.Sprintf("map[%s]%s", field.Type.Key, field.Type.Value)
			col.Unmarshal = fmt.Sprintf("Map[%s,%s]", field.Type.Key, field.Type.Value)
			if setCustom {
				col.Marshal = col.Unmarshal
			}
			col.SqlType = fmt.Sprintf("varchar(%d)", field.Options.GetInt64(options.SqlSize, 256))
			defVal = "default ''"
		}

		// 设置了NULL, 检测是否是主键
		if null && pk.hasCol(col.Name) {
			log.Printf("WARN: %s.%s set primary key NULL,invalid\n", msg.Name, field.Name)
			null = false
		}

		if !null && col.SqlType != "blob" {
			col.SqlType += " not null"
		}

		if autoIncr {
			// auto_increment
			col.SqlType += " auto_increment"
			tbl.AutoIncr = col
		}

		if !null && col.SqlType != "blob" && defVal != "" && !autoIncr {
			col.SqlType += " " + defVal
		}

		// 优先使用sql.type 指定的类型
		if sqlType != "" {
			col.SqlType = sqlType
			sqlType = strings.ToLower(sqlType)
			sqlType = strings.TrimSpace(sqlType)
			if strings.HasPrefix(sqlType, "timestamp") {
				switch col.GoType {
				case "int64":
					col.Unmarshal = "StampInt64"
					col.Marshal = "StampInt64"
				case "string":

				default:
					err = multierr.Append(err, fmt.Errorf("%s.%s type[%s] invalid,timestamp not int64 or string", msg.Name, field.Name, col.SqlType))
				}

			}
		}

		// auto incr
		if typ := strings.ToLower(col.SqlType); strings.Contains(typ, "auto_incr") && !strings.Contains(typ, "int") {
			err = multierr.Append(err, fmt.Errorf("%s.%s set auto_increment but not integer type", msg.Name, field.Name))
		}

		if pk.hasCol(field.Name) {
			tbl.PrimaryKey = append(tbl.PrimaryKey, col)
		} else {
			tbl.Columns = append(tbl.Columns, col)
		}

		tbl.allColumns = append(tbl.allColumns, col)

		// replace columns
		colCheck[field.Name] = col
	}
	// auto incr 必须是primary key
	if tbl.AutoIncr != nil && !pk.hasCol(tbl.AutoIncr.Name) {
		err = multierr.Append(err, fmt.Errorf("%s.%s set auto_increment but in primary key", msg.Name, tbl.AutoIncr.Name))
	}
	if err != nil {
		return
	}
	// index
	for _, v := range idx {
		x := &SqlIndex{
			Name:    v.name,
			Columns: make([]*SqlColumn, 0, len(v.columns)),
			unique:  v.unique,
		}
		for _, col := range v.columns {
			x.Columns = append(x.Columns, colCheck[col])
		}

		tbl.Index = append(tbl.Index, x)
	}
	return
}
