package gen

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/walleframe/wplugins/buildpb"
	"github.com/walleframe/wplugins/cmd/wredis/keyarg"
	"github.com/walleframe/wplugins/gen"
	"github.com/walleframe/wplugins/options"
	"github.com/walleframe/wplugins/utils"
)

var envFlag = struct {
	ProtobufPackage string
	WProtoPackage   string
	ServicePackage  string
}{
	ProtobufPackage: "", // github.com/gogo/protobuf/proto
	WProtoPackage:   "github.com/walleframe/walle/process/message",
	ServicePackage:  "github.com/walleframe/svc_redis",
}

func init() {
	utils.GetEnvString("WREDIS_PB_PKG", &envFlag.ProtobufPackage)
	utils.GetEnvString("WREDIS_WPB_PKG", &envFlag.WProtoPackage)
	utils.GetEnvString("WREDIS_SVC_PKG", &envFlag.ServicePackage)
}

func GenerateRedisMessage(prog *buildpb.FileDesc, msg *buildpb.MsgDesc, depend map[string]*buildpb.FileDesc) (out *buildpb.BuildOutput, err error) {
	// 解析redis-key
	redisKeyArgs, err := keyarg.MatchKey(msg.GetString(options.RedisOpKey, ""), nil)
	if err != nil {
		return nil, fmt.Errorf("parse %s option failed. MessageName:%s err:[%w]", options.RedisOpKey, msg.Name, err)
	}

	redisType := strings.ToLower(strings.TrimSpace(msg.GetString(options.RedisOpType, "")))

	// 构造redis对象
	obj := &RedisObject{
		Package: prog.GetPkg().Package,
		Name:    utils.Title(msg.Name),
		Doc:     msg.Doc,
		Args:    redisKeyArgs,
		SvcPkg:  filepath.Base(envFlag.ServicePackage),
		WPbPkg:  filepath.Base(envFlag.WProtoPackage),
		KeySize: int(msg.GetInt64(options.RedisOpKeySize, 64)),
	}

	obj.TypeKeys = true
	if strings.HasPrefix(redisType, "!") {
		obj.TypeKeys = false
		redisType = strings.TrimSpace(strings.TrimPrefix(redisType, "!"))
	}

	switch redisType {
	case "string":
		err = analyseTypeString(msg, obj)
		if err != nil {
			return nil, fmt.Errorf("MessageName:%s invalid string type define. err:[%w]", msg.Name, err)
		}
	case "hash":
		err = analyseTypeHash(msg, obj)
		if err != nil {
			return nil, fmt.Errorf("MessageName:%s invalid hash type define. err:[%w]", msg.Name, err)
		}
	case "set":
		err = analyseTypeSet(msg, obj)
		if err != nil {
			return nil, fmt.Errorf("MessageName:%s invalid set type define. err:[%w]", msg.Name, err)
		}
	case "zset":
		err = analyseTypeZSet(msg, obj)
		if err != nil {
			return nil, fmt.Errorf("MessageName:%s invalid zset type define. err:[%w]", msg.Name, err)
		}
	case "lock":
		obj.Lock = true
	case "": // 只生成key操作
		if !obj.TypeKeys {
			return nil, fmt.Errorf("MessageName:%s type config invalid. err:[%w]", msg.Name, err)
		}
	default:
		return nil, fmt.Errorf("MessageName:%s not support redis type[%s]", msg.Name, redisType)
	}
	// range for redis script
	err = analyseScript(msg, obj)
	if err != nil {
		return
	}

	if obj.Lock {
		if !strings.HasSuffix(obj.Name, "Lock") {
			obj.Name += "Lock"
		}
	} else {
		if !strings.HasSuffix(obj.Name, "RedisOpt") {
			obj.Name += "RedisOpt"
		}
	}

	obj.Import("context", "Context")
	obj.Import("github.com/redis/go-redis/v9", "UniversalClient")
	obj.Import("github.com/walleframe/walle/util", "Builder")
	obj.Import("github.com/walleframe/walle/util/rdconv", "AnyToString/Int64ToString/...")
	obj.Import(envFlag.ServicePackage, "RegisterDBName/GetDBLink")
	for _, arg := range redisKeyArgs {
		for _, pkg := range arg.Imports() {
			if pkg == "" {
				continue
			}
			obj.Import(pkg, "keyargv")
		}
	}

	// 基础结构
	tpl := template.New("redis.generate").Funcs(UseFuncMap)

	// 基础导入go包函数
	tpl.Funcs(template.FuncMap{
		"Import": obj.Import,
	})
	// 循环调用模板函数
	tpl.Funcs(template.FuncMap{
		"GenTypeTemplate": func(typeTplName string, obj *RedisObject) (string, error) {

			buf := &bytes.Buffer{}
			dst := tpl.Lookup(typeTplName)
			if dst == nil {
				return "", fmt.Errorf("%v not found", typeTplName)
			}

			err = dst.Execute(buf, obj)
			if err != nil {
				log.Println(err)
				return "", err
			}

			//		log.Println(buf.String())

			return buf.String(), nil
		},
		"GenScriptTemplate": func(obj *RedisObject, script *RedisScript) (string, error) {
			buf := &bytes.Buffer{}
			dst := tpl.Lookup(script.TemplateName)
			if dst == nil {
				return "", fmt.Errorf("%v not found", script.TemplateName)
			}

			err = dst.Execute(buf, map[string]interface{}{
				"Obj":    obj,
				"Script": script,
			})
			if err != nil {
				log.Println(err)
				return "", err
			}

			return buf.String(), nil
		},
	})

	// 解析前置的全部模板
	for k, v := range redisTypeTemplate {
		ts := fmt.Sprintf(`{{define "%s"}} %s {{end}}`, k, strings.TrimSpace(v))
		tpl, err = tpl.Parse(ts)
		if err != nil {
			err = fmt.Errorf("parse template %s failed:%+v", k, err.Error())
			return
		}
		// log.Println(ts)
	}

	tpl, err = tpl.Parse(GenRedisTemplate)
	if err != nil {
		return nil, fmt.Errorf("MessageName:%s parse basic template failed,[%w]", msg.Name, err)
	}

	buf := bytes.NewBuffer(make([]byte, 0, 4096))

	err = tpl.Execute(buf, obj)
	if err != nil {
		return nil, fmt.Errorf("MessageName:%s exec basic template failed,[%w]", msg.Name, err)
	}
	// import
	for _, imp := range prog.Imports {
		dep, ok := depend[imp.File]
		if !ok {
			return nil, fmt.Errorf("import %s file not found", imp.File)
		}
		pkg, ok := dep.Options.GetStringCheck(options.ProtoGoPkg)
		if !ok {
			return nil, fmt.Errorf("import %s not set '%s' option", imp.File, options.ProtoGoPkg)
		}
		obj.Import(pkg, "import depend file")
	}

	bdata := bytes.Replace(buf.Bytes(), []byte("$Import-Packages$"), []byte(fmt.Sprintf(`
import (
	%s
)`, obj.customImport())), 1)

	data, err := gen.GoFormat2(bdata)
	if err != nil {
		printWithLine(bdata)
		return nil, fmt.Errorf("MessageName:%s format code failed,[%w]", msg.Name, err)
	}

	out = &buildpb.BuildOutput{
		File: filepath.Join(filepath.Dir(prog.File), msg.Name+".go"),
		Data: data,
	}

	return
}

func printWithLine(data []byte) {
	for k, v := range bytes.Split(data, []byte{'\n'}) {
		log.Printf("%d\t%s\n", k, string(v))
	}
}

func analyseTypeString(msg *buildpb.MsgDesc, obj *RedisObject) (err error) {
	// 最多只能有一个字段
	if len(msg.Fields) > 1 {
		return errors.New("redis-string type fields too many. max 1 fields")
	}
	opt := &RedisTypeString{}
	obj.TypeString = opt
	// 无数据设置,查看是否生成通用接口
	if len(msg.Fields) < 1 {
		// 生成protobuf接口
		if msg.Options.GetOptionBool(options.RedisOpProtobuf) && len(envFlag.ProtobufPackage) > 0 {
			obj.Import(envFlag.ProtobufPackage, "Marshal/Unmarshal")
			opt.Protobuf = true
			return
		}
		// 生成walle message 接口
		if msg.Options.GetOptionBool(options.RedisOpWalleMsg) && len(envFlag.WProtoPackage) > 0 {
			obj.Import(envFlag.WProtoPackage, "MarshalObject/UnmarshalObject")
			opt.WProto = true
			return
		}
		// 无任何设置,直接生成string类型接口
		opt.String = true
		return
	}
	fieldType := msg.Fields[0].Type

	switch fieldType.Type {
	case buildpb.FieldType_BaseType:
		switch fieldType.KeyBase {
		case buildpb.BaseTypeDesc_Binary:
			return errors.New("redis-string type generation not support binary basic type.")
		case buildpb.BaseTypeDesc_Bool:
			return errors.New("redis-string type generation not support bool basic type.")
		case buildpb.BaseTypeDesc_String:
			opt.String = true
		case buildpb.BaseTypeDesc_Float32, buildpb.BaseTypeDesc_Float64:
			opt.Float = true
			opt.Type = fieldType.Key
		default:
			opt.Number = true
			opt.Type = fieldType.Key
			switch fieldType.KeyBase {
			case buildpb.BaseTypeDesc_Int8, buildpb.BaseTypeDesc_Int16, buildpb.BaseTypeDesc_Int32, buildpb.BaseTypeDesc_Int64:
				opt.Signed = true
			}
		}
	case buildpb.FieldType_CustomType:
		opt.Custom = true
		opt.Type = keyType(fieldType.Key)
		if msg.Options.GetOptionBool(options.RedisOpProtobuf) && len(envFlag.ProtobufPackage) > 0 {
			opt.Protobuf = true
			obj.Import(envFlag.ProtobufPackage, "Marshal/Unmarshal")
		}
	default:
		return errors.New("redis-string type generation not support array or map type.")
	}

	return
}

func analyseTypeHash(msg *buildpb.MsgDesc, obj *RedisObject) (err error) {
	obj.TypeHash = &RedisTypeHash{}
	// 分析hash object
	hashObject := func(ft *buildpb.TypeDesc, gen *RedisTypeHash) (err error) {
		if ft.Type != buildpb.FieldType_CustomType {
			return errors.New("redis-hash type 1 field must be custom struct.")
		}
		gen.HashObject = &RedisHashObject{
			Name:    utils.Title(ft.Msg.Name),
			Type:    keyType(ft.Key),
			Fields:  make([]*RedisGenType, 0, len(ft.Msg.GetFields())),
			HGetAll: true,
		}
		for _, v := range ft.Msg.Fields {
			if v.Type.Type != buildpb.FieldType_BaseType {
				return fmt.Errorf("redis-hash type message field [%s.%s] is not basic type.", ft.Msg.Name, v.Name)
			}
			field := &RedisGenType{
				Name:      v.Name,
				Type:      v.Type.Key,
				Number:    isNumber(v.Type.KeyBase),
				RedisFunc: v.Type.KeyBase.String(),
			}
			gen.HashObject.Fields = append(gen.HashObject.Fields, field)
		}
		return nil
	}
	hashDynamic := func(key *buildpb.TypeDesc, value *buildpb.TypeDesc, gen *RedisTypeHash) (err error) {
		defer func() {
			if err != nil {
				return
			}
			err = checkHashMatchModeArg(gen.HashDynamic)
		}()
		gen.HashDynamic = &RedisHashDynamic{GenMap: true}
		switch key.Type {
		case buildpb.FieldType_BaseType:
			switch key.KeyBase {
			case buildpb.BaseTypeDesc_Float32, buildpb.BaseTypeDesc_Float64, buildpb.BaseTypeDesc_Binary:
				gen.HashDynamic.GenMap = false
			}
			if key.KeyBase == buildpb.BaseTypeDesc_String && msg.HasOption(options.RedisOpMatchField) {
				// 拼接string做field
				gen.HashDynamic.FieldArgs, err = keyarg.MatchGoTypes(msg.Options.GetString(options.RedisOpMatchField, ""), nil)
				if err != nil {
					return fmt.Errorf("redis-hash analyse redis.field failed.%v", err)
				}
				gen.HashDynamic.GenMap = false
			} else {
				gen.HashDynamic.Field = &RedisGenType{
					Name:      "field",
					Type:      key.Key,
					Number:    isNumber(key.KeyBase),
					RedisFunc: key.KeyBase.String(),
				}
			}
		default:
			return fmt.Errorf("redis-hash type field not support type")
		}
		switch value.Type {
		case buildpb.FieldType_BaseType:
			if value.KeyBase == buildpb.BaseTypeDesc_String && msg.HasOption(options.RedisOpMatchValue) {
				// 拼接string做value
				gen.HashDynamic.ValueArgs, err = keyarg.MatchGoTypes(msg.Options.GetString(options.RedisOpMatchValue, ""), nil)
				if err != nil {
					return fmt.Errorf("redis-hash analyse redis.field failed.%v", err)
				}
				gen.HashDynamic.GenMap = false
			} else {
				gen.HashDynamic.Value = &RedisGenType{
					Name:      "value",
					Type:      value.Key,
					Number:    isNumber(value.KeyBase),
					RedisFunc: value.KeyBase.String(),
				}
			}
		default:
			return fmt.Errorf("redis-hash type value not support type")
		}
		return
	}
	switch len(msg.Fields) {
	case 1:
		err = hashObject(msg.Fields[0].Type, obj.TypeHash)
		if err != nil {
			return err
		}
	case 2:
		err = hashDynamic(msg.Fields[0].Type, msg.Fields[1].Type, obj.TypeHash)
		if err != nil {
			return err
		}

	case 3: // NOTE: hval,hfields 生成需要过滤,有需要再改吧. 先禁掉功能了.
		return fmt.Errorf("redis-hash type fields count 3 not support now,need modify hvals/hfields/range functions")
	// 	err = hashObject(msg.Fields[0].Type, obj.TypeHash)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	err = hashDynamic(msg.Fields[1].Type, msg.Fields[2].Type, obj.TypeHash)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	// 同时存在object和动态字段,禁止使用hgetall
	// 	obj.TypeHash.HashObject.HGetAll = false
	default:
		return fmt.Errorf("redis-hash type fields count %d not support", len(msg.Fields))
	}

	return
}

func analyseTypeSet(msg *buildpb.MsgDesc, obj *RedisObject) (err error) {
	// 最多只能有一个字段
	if len(msg.Fields) > 1 {
		return errors.New("redis-set type fields too many. max 1 fields")
	}
	opt := &RedisTypeSet{}
	obj.TypeSet = opt

	// 无数据设置,查看是否生成通用接口
	if len(msg.Fields) < 1 {
		// 生成protobuf接口
		if msg.Options.GetOptionBool(options.RedisOpProtobuf) && len(envFlag.ProtobufPackage) > 0 {
			obj.Import(envFlag.ProtobufPackage, "Marshal/Unmarshal")
			opt.Message = &RedisGenMsg{
				Type: "proto.Message",
				Marshal: func(objName string) string {
					return fmt.Sprintf("proto.Marshal(%s)", objName)
				},
				Unmarshal: func(objName string, paramName string) string {
					return fmt.Sprintf("proto.Unmarshal(%s,%s)", objName, paramName)
				},
				New: "",
			}
			return
		}
		// 生成walle message 接口
		if msg.Options.GetOptionBool(options.RedisOpWalleMsg) && len(envFlag.WProtoPackage) > 0 {
			obj.Import(envFlag.WProtoPackage, "MarshalObject/UnmarshalObject")
			opt.Message = &RedisGenMsg{
				Type: "message.Message",
				Marshal: func(objName string) string {
					return fmt.Sprintf("%s.MarshalObject()", objName)
				},
				Unmarshal: func(objName string, paramName string) string {
					return fmt.Sprintf("%s.UnmarshalObject(%s)", objName, paramName)
				},
				New: "",
			}
			return
		}
		// 无任何设置,直接生成string类型接口
		opt.BaseType = &RedisGenType{
			Name:      "",
			Type:      "string",
			Number:    false,
			RedisFunc: "String",
		}
		return
	}

	fieldType := msg.Fields[0].Type

	switch fieldType.Type {
	case buildpb.FieldType_BaseType:
		switch fieldType.KeyBase {
		// case buildpb.BaseTypeDesc_Binary:
		// 	return errors.New("redis-string type generation not support binary basic type.")
		case buildpb.BaseTypeDesc_Bool:
			return errors.New("redis-set type generation not support bool basic type")
		default:
			opt.BaseType = &RedisGenType{
				Name:      fieldType.Key,
				Type:      fieldType.Key,
				Number:    false,
				RedisFunc: fieldType.KeyBase.String(),
			}
		}
	case buildpb.FieldType_CustomType:
		opt.Message = &RedisGenMsg{
			Type: "*" + keyType(fieldType.Key),
			Marshal: func(objName string) string {
				return fmt.Sprintf("%s.MarshalObject()", objName)
			},
			Unmarshal: func(objName string, paramName string) string {
				return fmt.Sprintf("%s.UnmarshalObject(%s)", objName, paramName)
			},
			New: "&" + keyType(fieldType.Key) + "{}",
		}
		if msg.Options.GetOptionBool(options.RedisOpProtobuf) && len(envFlag.ProtobufPackage) > 0 {
			obj.Import(envFlag.ProtobufPackage, "Marshal/Unmarshal")
			opt.Message.Marshal = func(objName string) string {
				return fmt.Sprintf("proto.Marshal(%s)", objName)
			}
			opt.Message.Unmarshal = func(objName string, paramName string) string {
				return fmt.Sprintf("proto.Unmarshal(%s,%s)", objName, paramName)
			}
		}
	default:
		return errors.New("redis-set type generation not support array or map type")
	}

	return
}

func analyseTypeZSet(msg *buildpb.MsgDesc, obj *RedisObject) (err error) {
	// 只支持1个或者2个字段
	fieldCount := len(msg.Fields)
	if fieldCount < 1 || fieldCount > 2 {
		return errors.New("redis-zset type fields invalid. only support 1 or 2 fields")
	}
	opt := &RedisTypeZSet{}
	obj.TypeZSet = opt

	fieldType := msg.Fields[0].Type

	switch fieldType.Type {
	case buildpb.FieldType_BaseType:
		switch fieldType.KeyBase {
		case buildpb.BaseTypeDesc_Bool:
			return errors.New("redis-zset type member not support bool basic type")
		default:
			if fieldType.KeyBase == buildpb.BaseTypeDesc_String && msg.HasOption(options.RedisOpMatchMember) {
				// 拼接string做value
				opt.Args, err = keyarg.MatchGoTypes(msg.Options.GetString(options.RedisOpMatchMember, ""), nil)
				if err != nil {
					return fmt.Errorf("redis-zset analyse redis.member failed.%v", err)
				}
			} else {
				opt.Member = &RedisGenType{
					Name:      fieldType.Key,
					Type:      fieldType.Key,
					Number:    false,
					RedisFunc: fieldType.KeyBase.String(),
				}
			}
		}
	case buildpb.FieldType_CustomType:
		opt.Message = &RedisGenMsg{
			Type: "*" + keyType(fieldType.Key),
			Marshal: func(objName string) string {
				return fmt.Sprintf("%s.MarshalObject()", objName)
			},
			Unmarshal: func(objName string, paramName string) string {
				return fmt.Sprintf("%s.UnmarshalObject(%s)", objName, paramName)
			},
			New: "&" + keyType(fieldType.Key) + "{}",
		}
		if msg.Options.GetOptionBool(options.RedisOpProtobuf) && len(envFlag.ProtobufPackage) > 0 {
			obj.Import(envFlag.ProtobufPackage, "Marshal/Unmarshal")
			opt.Message.Marshal = func(objName string) string {
				return fmt.Sprintf("proto.Marshal(%s)", objName)
			}
			opt.Message.Unmarshal = func(objName string, paramName string) string {
				return fmt.Sprintf("proto.Unmarshal(%s,%s)", objName, paramName)
			}
		}
	default:
		return errors.New("redis-zset type generation not support array or map type")
	}

	if fieldCount < 2 {
		opt.Score = &RedisGenType{
			Name:      "score",
			Type:      "float64",
			Number:    true,
			RedisFunc: "Float64",
		}
		return
	}

	scoreType := msg.Fields[1].Type
	if scoreType.Type != buildpb.FieldType_BaseType {
		return errors.New("redis-zset type score only support signed int or float type")
	}

	if !strings.HasPrefix(scoreType.Key, "int") && !strings.HasPrefix(scoreType.Key, "float") {
		return errors.New("redis-zset type score only support signed int or float type")
	}

	opt.Score = &RedisGenType{
		Name:      "score",
		Type:      scoreType.Key,
		Number:    false,
		RedisFunc: scoreType.KeyBase.String(),
	}

	return
}

func analyseScript(msg *buildpb.MsgDesc, obj *RedisObject) (err error) {
	for optKey := range msg.Options.Options {
		if !strings.HasPrefix(optKey, options.RedisScriptPrefix) {
			continue
		}
		if !strings.HasSuffix(optKey, options.RedisScriptSuffixScript) {
			continue
		}
		scriptName := strings.TrimSuffix(strings.TrimPrefix(optKey, options.RedisScriptPrefix), options.RedisScriptSuffixScript)
		if strings.Contains(scriptName, ".") {
			return fmt.Errorf("MessageName:%s define redis script failed. script name [%s] invalid", msg.Name, scriptName)
		}
		scriptData := msg.GetString(options.RedisScriptPrefix+scriptName+options.RedisScriptSuffixScript, "")
		scriptArgv := msg.GetString(options.RedisScriptPrefix+scriptName+options.RedisScriptSuffixInput, "")
		scriptReply := msg.GetString(options.RedisScriptPrefix+scriptName+options.RedisScriptSuffixReply, "")
		// log.Println(options.RedisScriptPrefix+scriptName+options.RedisScriptSuffixScript,
		// 	options.RedisScriptPrefix+scriptName+options.RedisScriptSuffixInput,
		// 	options.RedisScriptPrefix+scriptName+options.RedisScriptSuffixReply)

		if scriptData == "" {
			return fmt.Errorf("MessageName:%s redis script [%s] data empty", msg.Name, scriptName)
		}
		if scriptArgv == "" {
			return fmt.Errorf("MessageName:%s redis script [%s] argv empty", msg.Name, scriptName)
		}
		if scriptReply == "" {
			return fmt.Errorf("MessageName:%s redis script [%s] reply empty", msg.Name, scriptName)
		}
		argv, err := keyarg.MatchGoTypes(scriptArgv, nil)
		if err != nil {
			return fmt.Errorf("MessageName:%s redis script [%s] argv invalid. %+v", msg.Name, scriptName, err)
		}

		reply, err := keyarg.MatchGoTypes(scriptReply, nil)
		if err != nil {
			return fmt.Errorf("MessageName:%s redis script [%s] reply invalid. %+v", msg.Name, scriptName, err)
		}
		if len(reply) < 1 {
			return fmt.Errorf("MessageName:%s redis script [%s] reply must >= 1", msg.Name, scriptName)
		}

		script := &RedisScript{
			Name:         scriptName,
			Script:       scriptData,
			Args:         argv,
			Output:       reply,
			TemplateName: "script_return_mul",
			CommandName:  "",
		}
		if len(reply) == 1 {
			script.TemplateName = "script_return_1"
			switch reply[0].ArgType() {
			case "bool":
				script.CommandName = "NewBoolCmd"
			case "float32", "float64":
				script.CommandName = "NewFloatCmd"
			case "string":
				script.CommandName = "NewStringCmd"
			default:
				script.CommandName = "NewIntCmd"
			}
		}

		obj.Scripts = append(obj.Scripts, script)

	}

	return
}

func isNumber(typ buildpb.BaseTypeDesc) bool {
	switch typ {
	case buildpb.BaseTypeDesc_String, buildpb.BaseTypeDesc_Binary, buildpb.BaseTypeDesc_Bool, buildpb.BaseTypeDesc_Float32, buildpb.BaseTypeDesc_Float64:
		return false
	default:
		return true
	}
}

func checkHashMatchModeArg(obj *RedisHashDynamic) (err error) {
	if obj == nil {
		return
	}
	if obj.FieldArgs == nil || obj.ValueArgs == nil {
		return
	}
	checks := make(map[string]struct{})
	for _, v := range obj.FieldArgs {
		checks[v.ArgName()] = struct{}{}
	}
	for _, v := range obj.ValueArgs {
		if _, ok := checks[v.ArgName()]; ok {
			return fmt.Errorf("redis-hash match field named repeated[%s]", v.ArgName())
		}
	}
	return
}

func keyType(key string) string {
	if strings.Contains(key, ".") {
		typs := strings.SplitN(key, ".", 2)
		return typs[0] + "." + utils.Title(typs[1])
	}
	return utils.Title(key)
}
