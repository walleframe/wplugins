package options

// proto 相关插件  使用的 options
const (
	// ProtoGoPkg 生成proto导出的go包
	ProtoGoPkg = "proto.gopkg"
	// ProtoLowerCase 生成的proto消息 使用小写 NOTE: toproto 使用的
	ProtoLowerCase = "proto.lowercase"
	// ProtoSyntax 生成的proto版本
	ProtoSyntax = "proto.syntax"
	// ProtoFieldOptPacked
	ProtoFieldOptPacked = "proto.packed"
	// ProtoFieldOptFixed
	ProtoFieldOptFixed = "proto.fixed"
	// ProtoFieldOptSigned
	ProtoFieldOptSigned = "proto.signed"
)

// redisop 插件使用的option 具体规则看文档
const (
	// RedisOpKey 定义redis key 名字
	RedisOpKey              = "redis.key"
	RedisOpType             = "redis.type"
	RedisOpKeySize          = "redis.keysize"
	RedisOpProtobuf         = "redis.protobuf"
	RedisOpWalleMsg         = "redis.wproto"
	RedisOpMatchField       = "redis.field"
	RedisOpMatchValue       = "redis.value"
	RedisOpMatchMember      = "redis.member"
	RedisScriptPrefix       = "redis.script."
	RedisScriptSuffixScript = ".lua"
	RedisScriptSuffixInput  = ".argv"
	RedisScriptSuffixReply  = ".reply"
)

// mysql key
const (
	SqlDBName        = "sql.db"
	SqlIgnore        = "sql.ignore"
	SqlEngine        = "sql.engine"
	SqlTableName     = "sql.tbl_name"
	SqlPrimaryKey    = "sql.pk"
	SqlPrimaryKey2   = "sql.primary_key"
	SqlAutoIncrement = "sql.auto_incr"
	SqlFieldType     = "sql.type"
	SqlSize          = "sql.size"
	SqlNull          = "sql.null"
	SqlCustomSet     = "sql.custom"
	SqlIndexPrefix   = "sql.index."
	SqlUniquePrefix  = "sql.unique."
)
