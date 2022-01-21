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
