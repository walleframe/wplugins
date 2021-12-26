package options

// toproto 插件  使用的 options
const (
	// ProtoGoPkg 生成proto导出的go包
	ProtoGoPkg = "proto.gopkg"
	// ProtoLowerCase 生成的proto消息 使用小写
	ProtoLowerCase = "proto.lowercase"
	// ProtoSyntax 生成的proto版本
	ProtoSyntax = "proto.syntax"
)

// dblog 插件 使用的options
const (
	// DBLogDB 写入db名称选项. 选项值必须是 字符串
	//
	// 如果在文件级别设置. 为本文件消息默认生成的db名.
	// 如果即在文件设置了,也在msg上设置了.那么将使用msg上设置的db名.
	DBLogDB = "dblog.db"
	// DBLogTable 写入数据库的表名称
	//
	// 文件级别 值忽略. 代表当前文件所有msg都会生成.如果msg上没有设置此属性. 表名 = t_[msg名]
	// 如果msg上设置了该属性. 那么消息使用设置的值导出(值必须是 字符串)
	DBLogTable = "dblog.table"
	// DBLogMode 定义写入数据库及表的格式(月库天表)
	//
	// 值为 month => 月库天表
	// 既可以设置在文件级别.也可以设置消息级别.
	DBLogDTMode = "dblog.mode"
	// DBLogStringSize 用于指定string类型存储长度
	DBLogStringSize = "dblog.size"
	// DBLogTableIndex 数据库表索引. 使用,分隔的字符串. 字符串必须是字段名或者默认字段(auto_id,auto_create)
	//
	// 应用在msg级别. 可以直接设置多个字段.
	// 应用在字段级别. 值忽略. 如果已经msg级别设置. 会合并
	DBLogTableIndex = "dblog.index"
	// DBLogAbsPath 生成的配置和sql语句是否使用原始文件的路径. 仅在文件级别有效
	// 未设置. 默认生成目录 为 dblog/sql 和 dblog/configs
	DBLogAbsPath = "dblog.abs_path"
)

// redisop 插件使用的option
const (
	// RedisOpKey 定义redis key 名字
	// key名格式
	//  key 以 "$" 开头
	//  "key:$1:name" 结构体字段. 序号1
	//  "xx:$mod(1%1024):xxx" 使用结构体序号为1的字段进行取模1024.
	//  "xx:$(def):xx" $(def)  def 为定义
	// "xx:$mod(def%1024):xx" $(def)  def 必须为数值类型 取模1024
	// 支持的定义如下
	//  daystamp 当日0点0分时间戳
	//  weekstamp 本周周一0点0分时间戳
	//  monthstamp 本月1号0点0分时间戳
	//  int8,int16,int32,int64,uint8,uint16,uint32,uint64,int,uint golang类型名
	//  string 自定义符串
	RedisOpKey = "redis.key"
	// RedisOpType 生成redis操作类型
	// 支持的redis数据类型 string,hash
	RedisOpType = "redis.type"
	// RedisOpLang 生成的语言.默认golang
	RedisOpLang = "redis.language"
	// RedisOpKeySep key切分字符 默认":"
	RedisOpKeySep = "redis.keysep"
	// RedisOpKeyPrefix key切分.定义前缀 默认 "$"
	RedisOpKeyPrefix = "redis.splitprefix"
	// RedisOpFeildIncr 字段是否生成incr方法
	RedisOpFeildIncr = "redis.incr"
	// 脚本可读性名字. 用于生成函数名
	RedisOpLuaName = "redis.scriptname"
	// 脚本参数
	RedisOpLuaArgs = "redis.scriptargs"
	// 脚本返回值
	RedisOpLuaRets = "redis.scriptrets"
	// 脚本
	RedisOpLuaScript = "redis.script"
)

// GMT 名称
const (
	// gmt视图生成标记
	GMTSchema = "gmt.schema"
	// gmt显示名称
	GMTName = "gmt.name"
	// 权限前缀
	GMTPermission = "group."
	////////////////////////////////////////////////////////////
	// 字段可设置
	// 字段类型
	GMTFieldType = "gmt.type"
	// 默认值
	GMTFieldDefault = "gmt.default"
	// 不在列表
	GMTNotInList = "gmt.not_in_list"
	// 不在详情列表
	GMTNotInDetail = "gmt.not_in_detail"
	// 排序
	GMTSortable = "gmt.sort"
	// 搜索
	GMTSearch = "gmt.search"
	// 只读
	GMTReadOnly = "gmt.readonly"
	//
	GMTCreateOnly = "gmt.createonly"
	// 列可以直接编辑
	GMTEditColumn = "gmt.columnedit"
	// Radio 选项
	GMTValues = "gmt.values"
	// GMT key 默认第一字段是key
	GMTKey = "gmt.key"
	//
	GMTLFSelect    = "lf.select"
	GMTLFGet       = "lf.get"
	GMTInput       = "lf.input"
	GMTLFAct       = "lfact."
	GMTLFActAll    = "all"
	GMTLFActAdd    = "add"
	GMTLFActDel    = "del"
	GMTLFActImport = "import"
	////////////////////////////////////////////////////////////
	// TGA 支持
	//
	// TGAEvent 标记生成. 事件名
	TGAEvent = "tga.ev"
	// 数据库表名
	TGATable = "tga.tbl"
	// TGAFieldName 字段名
	TGAFieldName = "tga.fname"
	// TGA类型
	TGAFieldType = "tga.ftype"
)
