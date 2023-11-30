package gen

import (
	"testing"

	"github.com/aggronmagi/wplugins/buildpb"
	"github.com/aggronmagi/wplugins/options"
	"github.com/stretchr/testify/assert"
)

func TestGenerateStringType(t *testing.T) {
	envFlag.ProtobufPackage = "github.com/gogo/protobuf/proto"

	runFunc := func(name string, msg *buildpb.MsgDesc, depend map[string]*buildpb.FileDesc) {
		t.Run(name, func(t *testing.T) {
			out, err := generateRedisMessage(&buildpb.FileDesc{
				Pkg: &buildpb.PackageDesc{
					Package: "pkg",
				},
			}, msg, depend)
			assert.Nil(t, err, "generate key error")
			if out != nil {
				t.Log(string(out.Data))
			}
		})
	}

	runFunc("basic keys", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
			},
		},
		Name: "XX",
	}, nil)

	runFunc("basic string empty string", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!string",
				},
			},
		},
		Name: "XX",
	}, nil)

	runFunc("basic string empty string-withkey", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "string",
				},
			},
		},
		Name: "XX",
	}, nil)

	runFunc("basic string empty protobuf", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!string",
				},
				options.RedisOpProtobuf: {
					IntValue: 1,
				},
			},
		},
		Name: "XX",
	}, nil)

	runFunc("basic string empty wallemsg", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!string",
				},
				options.RedisOpWalleMsg: {
					IntValue: 1,
				},
			},
		},
		Name: "XX",
	}, nil)

	runFunc("basic string 1-field int32", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!string",
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "x",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					KeyBase: buildpb.BaseTypeDesc_Int32,
					Key:     "int32",
				},
			},
		},
	}, nil)

	runFunc("basic string 1-field uint32", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!string",
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "x",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					KeyBase: buildpb.BaseTypeDesc_Uint32,
					Key:     "uint32",
				},
			},
		},
	}, nil)

	runFunc("basic string 1-field float32", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!string",
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "x",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					KeyBase: buildpb.BaseTypeDesc_Float32,
					Key:     "float32",
				},
			},
		},
	}, nil)

	runFunc("basic string 1-field custom-walle", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!string",
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "x",
				Type: &buildpb.TypeDesc{
					Type: buildpb.FieldType_CustomType,
					Key:  "abc.ABC",
				},
			},
		},
	}, nil)

	runFunc("basic string 1-field custom-protobuf", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!string",
				},
				options.RedisOpProtobuf: {
					IntValue: 1,
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "x",
				Type: &buildpb.TypeDesc{
					Type: buildpb.FieldType_CustomType,
					Key:  "abc.ABC",
				},
			},
		},
	}, nil)
}

func TestGenerateHashType(t *testing.T) {
	envFlag.ProtobufPackage = "github.com/gogo/protobuf/proto"

	runFunc := func(name string, msg *buildpb.MsgDesc, depend map[string]*buildpb.FileDesc) {
		t.Run(name, func(t *testing.T) {
			out, err := generateRedisMessage(&buildpb.FileDesc{
				Pkg: &buildpb.PackageDesc{
					Package: "pkg",
				},
			}, msg, depend)
			assert.Nil(t, err, "generate key error")
			if out != nil {
				t.Log(string(out.Data))
			}
		})
	}

	runFunc("hash object 1-field", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!hash",
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "x",
				Type: &buildpb.TypeDesc{
					Type: buildpb.FieldType_CustomType,
					Key:  "abc.ABC",
					Msg: &buildpb.MsgDesc{
						Name: "ABC",
						Fields: []*buildpb.Field{
							{
								Name: "i8",
								Type: &buildpb.TypeDesc{
									Type:    buildpb.FieldType_BaseType,
									Key:     "int8",
									KeyBase: buildpb.BaseTypeDesc_Int8,
								},
							},
							{
								Name: "i64",
								Type: &buildpb.TypeDesc{
									Type:    buildpb.FieldType_BaseType,
									Key:     "int64",
									KeyBase: buildpb.BaseTypeDesc_Int64,
								},
							},
							{
								Name: "u64",
								Type: &buildpb.TypeDesc{
									Type:    buildpb.FieldType_BaseType,
									Key:     "uint64",
									KeyBase: buildpb.BaseTypeDesc_Uint64,
								},
							},
							{
								Name: "str",
								Type: &buildpb.TypeDesc{
									Type:    buildpb.FieldType_BaseType,
									Key:     "string",
									KeyBase: buildpb.BaseTypeDesc_String,
								},
							},
							{
								Name: "bytes",
								Type: &buildpb.TypeDesc{
									Type:    buildpb.FieldType_BaseType,
									Key:     "[]byte",
									KeyBase: buildpb.BaseTypeDesc_Binary,
								},
							},
							{
								Name: "boolean",
								Type: &buildpb.TypeDesc{
									Type:    buildpb.FieldType_BaseType,
									Key:     "bool",
									KeyBase: buildpb.BaseTypeDesc_Bool,
								},
							},
							{
								Name: "f32",
								Type: &buildpb.TypeDesc{
									Type:    buildpb.FieldType_BaseType,
									Key:     "float32",
									KeyBase: buildpb.BaseTypeDesc_Float32,
								},
							},
						},
					},
				},
			},
		},
	}, nil)

	runFunc("hash 2-field basic int64:float64", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!hash",
				},
				options.RedisOpProtobuf: {
					IntValue: 1,
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "xkey",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "int64",
					KeyBase: buildpb.BaseTypeDesc_Int64,
				},
			},
			{
				Name: "xval",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "float64",
					KeyBase: buildpb.BaseTypeDesc_Float64,
				},
			},
		},
	}, nil)

	runFunc("hash 2-field basic int64:string", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!hash",
				},
				options.RedisOpProtobuf: {
					IntValue: 1,
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "xkey",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "int64",
					KeyBase: buildpb.BaseTypeDesc_Int64,
				},
			},
			{
				Name: "xval",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "string",
					KeyBase: buildpb.BaseTypeDesc_String,
				},
			},
		},
	}, nil)

	runFunc("hash 2-field basic string:int64", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!hash",
				},
				options.RedisOpProtobuf: {
					IntValue: 1,
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "xval",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "string",
					KeyBase: buildpb.BaseTypeDesc_String,
				},
			},
			{
				Name: "xkey",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "int64",
					KeyBase: buildpb.BaseTypeDesc_Int64,
				},
			},
		},
	}, nil)

	runFunc("hash 2-field basic string:string", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!hash",
				},
				options.RedisOpProtobuf: {
					IntValue: 1,
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "xval",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "string",
					KeyBase: buildpb.BaseTypeDesc_String,
				},
			},
			{
				Name: "xkey",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "string",
					KeyBase: buildpb.BaseTypeDesc_String,
				},
			},
		},
	}, nil)

	runFunc("hash 2-field basic float64:int64 - nomap", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!hash",
				},
				options.RedisOpProtobuf: {
					IntValue: 1,
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "xval",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "float64",
					KeyBase: buildpb.BaseTypeDesc_Float64,
				},
			},
			{
				Name: "xkey",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "int64",
					KeyBase: buildpb.BaseTypeDesc_Int64,
				},
			},
		},
	}, nil)

	runFunc("hash 2-field match string:int64", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!hash",
				},
				options.RedisOpMatchField: {
					Value: "$x1=int32:$x2=int8:$string",
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "xval1",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "string",
					KeyBase: buildpb.BaseTypeDesc_String,
				},
			},
			{
				Name: "xkey2",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "int64",
					KeyBase: buildpb.BaseTypeDesc_Int64,
				},
			},
		},
	}, nil)

	runFunc("hash 2-field match int64:string", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!hash",
				},
				options.RedisOpMatchValue: {
					Value: "$x1=int32:$x2=int8:$string",
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "xkey1",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "int64",
					KeyBase: buildpb.BaseTypeDesc_Int64,
				},
			},
			{
				Name: "xval2",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "string",
					KeyBase: buildpb.BaseTypeDesc_String,
				},
			},
		},
	}, nil)

	runFunc("hash 2-field match string:string", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!hash",
				},
				options.RedisOpMatchValue: {
					Value: "$x1=int32:$x2=int8:$string",
				},
				options.RedisOpMatchField: {
					Value: "$a1=int32:$a2=int8:$a4=string",
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "xkey1",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "string",
					KeyBase: buildpb.BaseTypeDesc_String,
				},
			},
			{
				Name: "xval2",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "string",
					KeyBase: buildpb.BaseTypeDesc_String,
				},
			},
		},
	}, nil)

	runFunc("hash 3-field mix message float64:int64 - nomap", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!hash",
				},
				options.RedisOpProtobuf: {
					IntValue: 1,
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "x",
				Type: &buildpb.TypeDesc{
					Type: buildpb.FieldType_CustomType,
					Key:  "abc.ABC",
					Msg: &buildpb.MsgDesc{
						Name: "ABC",
						Fields: []*buildpb.Field{
							{
								Name: "i8",
								Type: &buildpb.TypeDesc{
									Type:    buildpb.FieldType_BaseType,
									Key:     "int8",
									KeyBase: buildpb.BaseTypeDesc_Int8,
								},
							},
							{
								Name: "i64",
								Type: &buildpb.TypeDesc{
									Type:    buildpb.FieldType_BaseType,
									Key:     "int64",
									KeyBase: buildpb.BaseTypeDesc_Int64,
								},
							},
							{
								Name: "u64",
								Type: &buildpb.TypeDesc{
									Type:    buildpb.FieldType_BaseType,
									Key:     "uint64",
									KeyBase: buildpb.BaseTypeDesc_Uint64,
								},
							},
							{
								Name: "str",
								Type: &buildpb.TypeDesc{
									Type:    buildpb.FieldType_BaseType,
									Key:     "string",
									KeyBase: buildpb.BaseTypeDesc_String,
								},
							},
							{
								Name: "bytes",
								Type: &buildpb.TypeDesc{
									Type:    buildpb.FieldType_BaseType,
									Key:     "[]byte",
									KeyBase: buildpb.BaseTypeDesc_Binary,
								},
							},
							{
								Name: "boolean",
								Type: &buildpb.TypeDesc{
									Type:    buildpb.FieldType_BaseType,
									Key:     "bool",
									KeyBase: buildpb.BaseTypeDesc_Bool,
								},
							},
							{
								Name: "f32",
								Type: &buildpb.TypeDesc{
									Type:    buildpb.FieldType_BaseType,
									Key:     "float32",
									KeyBase: buildpb.BaseTypeDesc_Float32,
								},
							},
						},
					},
				},
			},
			{
				Name: "xval",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "float64",
					KeyBase: buildpb.BaseTypeDesc_Float64,
				},
			},
			{
				Name: "xkey",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "int64",
					KeyBase: buildpb.BaseTypeDesc_Int64,
				},
			},
		},
	}, nil)
}

func TestGenerateSetType(t *testing.T) {
	envFlag.ProtobufPackage = "github.com/gogo/protobuf/proto"

	runFunc := func(name string, msg *buildpb.MsgDesc, depend map[string]*buildpb.FileDesc) {
		t.Run(name, func(t *testing.T) {
			out, err := generateRedisMessage(&buildpb.FileDesc{
				Pkg: &buildpb.PackageDesc{
					Package: "pkg",
				},
			}, msg, depend)
			assert.Nil(t, err, "generate key error")
			if out != nil {
				t.Log(string(out.Data))
			}
		})
	}

	runFunc("basic int64", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!set",
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "xkey",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "int64",
					KeyBase: buildpb.BaseTypeDesc_Int64,
				},
			},
		},
	}, nil)

	runFunc("basic string", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!set",
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "xkey",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "string",
					KeyBase: buildpb.BaseTypeDesc_String,
				},
			},
		},
	}, nil)

	runFunc("basic default string", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!set",
				},
			},
		},
		Name: "XX",
	}, nil)

	runFunc("basic default protobuf", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!set",
				},
				options.RedisOpProtobuf: {
					IntValue: 1,
				},
			},
		},
		Name: "XX",
	}, nil)

	runFunc("basic default wallemsg", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!set",
				},
				options.RedisOpWalleMsg: {
					IntValue: 1,
				},
			},
		},
		Name: "XX",
	}, nil)

	runFunc("basic custom", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!set",
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "xkey",
				Type: &buildpb.TypeDesc{
					Type: buildpb.FieldType_CustomType,
					Key:  "abc.ABC",
				},
			},
		},
	}, nil)

	runFunc("basic custom-pb", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!set",
				},
				options.RedisOpProtobuf: {
					IntValue: 1,
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "xkey",
				Type: &buildpb.TypeDesc{
					Type: buildpb.FieldType_CustomType,
					Key:  "abc.ABC",
				},
			},
		},
	}, nil)

}

func TestGenerateZsetType(t *testing.T) {
	envFlag.ProtobufPackage = "github.com/gogo/protobuf/proto"

	runFunc := func(name string, msg *buildpb.MsgDesc, depend map[string]*buildpb.FileDesc) {
		t.Run(name, func(t *testing.T) {
			out, err := generateRedisMessage(&buildpb.FileDesc{
				Pkg: &buildpb.PackageDesc{
					Package: "pkg",
				},
			}, msg, depend)
			assert.Nil(t, err, "generate key error")
			if out != nil {
				t.Log(string(out.Data))
			}
		})
	}

	runFunc("basic int64", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!zset",
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "xkey",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "int64",
					KeyBase: buildpb.BaseTypeDesc_Int64,
				},
			},
		},
	}, nil)

	runFunc("basic string", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!zset",
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "xkey",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "string",
					KeyBase: buildpb.BaseTypeDesc_String,
				},
			},
		},
	}, nil)

	runFunc("basic string:int64", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!zset",
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "xkey",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "string",
					KeyBase: buildpb.BaseTypeDesc_String,
				},
			},
			{
				Name: "score",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "int64",
					KeyBase: buildpb.BaseTypeDesc_Int64,
				},
			},
		},
	}, nil)

	runFunc("basic merge-string:int64", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!zset",
				},
				options.RedisOpMatchMember: {
					Value: "$x1=int32:$x2=int8:$string",
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "xkey",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "string",
					KeyBase: buildpb.BaseTypeDesc_String,
				},
			},
			{
				Name: "score",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "int64",
					KeyBase: buildpb.BaseTypeDesc_Int64,
				},
			},
		},
	}, nil)

	runFunc("basic obj:int64", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!zset",
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "xkey",
				Type: &buildpb.TypeDesc{
					Type: buildpb.FieldType_CustomType,
					Key:  "abc.ABC",
				},
			},
			{
				Name: "score",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "int64",
					KeyBase: buildpb.BaseTypeDesc_Int64,
				},
			},
		},
	}, nil)

	runFunc("basic obj:int64 protobuf", &buildpb.MsgDesc{
		Options: &buildpb.OptionDesc{
			Options: map[string]*buildpb.OptionValue{
				options.RedisOpKey: {
					Value: "u:data:$uid=int64:@daystamp+8",
				},
				options.RedisOpType: {
					Value: "!zset",
				},
				options.RedisOpProtobuf: {
					IntValue: 1,
				},
			},
		},
		Name: "XX",
		Fields: []*buildpb.Field{
			{
				Name: "xkey",
				Type: &buildpb.TypeDesc{
					Type: buildpb.FieldType_CustomType,
					Key:  "abc.ABC",
				},
			},
			{
				Name: "score",
				Type: &buildpb.TypeDesc{
					Type:    buildpb.FieldType_BaseType,
					Key:     "int64",
					KeyBase: buildpb.BaseTypeDesc_Int64,
				},
			},
		},
	}, nil)
}
