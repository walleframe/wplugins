package proto3

import (
	"errors"
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
)

func TestUnmarshal(t *testing.T) {
	v := &FieldTestMessage{
		OptionalBool:     "xxadf",
		OptionalEnum:     100,
		OptionalInt32:    990,
		OptionalSint32:   880,
		OptionalUint32:   550,
		OptionalInt64:    8880,
		OptionalSint64:   11110,
		OptionalUint64:   2220,
		OptionalSfixed32: -3330,
		OptionalFixed32:  4440,
		OptionalFloat:    5.0,
		OptionalSfixed64: -6.0,
		OptionalFixed64:  70,
		OptionalDouble:   888.80,
		OptionalString:   "sdf234",
		OptionalBytes:    []byte{1, 5, 2, 5, 6, 54, 5},
		RepeatedBool:     []bool{true, false, true, false},
		MapInt32Int64:    map[int32]int64{88: 2, 90: 5},
	}
	data, err := proto.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	n := &FieldTestMessage{}
	err = n.UnmarshalObject(data)
	if err != nil {
		t.Fatal(err)
	}

	v.sizeCache = n.sizeCache
	v.state = n.state
	v.unknownFields = n.unknownFields

	assert.EqualValues(t, v, n)

	t.Log(v.MapInt32Int64)

	data, err = v.MarshalObject()
	if err != nil {
		t.Fatal(err)
	}

	err = n.UnmarshalObject(data)
	if err != nil {
		t.Fatal(err)
	}

	err = proto.Unmarshal(data, n)
	if err != nil {
		t.Fatal(err)
	}

	v.sizeCache = n.sizeCache
	v.state = n.state
	v.unknownFields = n.unknownFields

	assert.EqualValues(t, v, n)

	if true {
		return
	}

	n = &FieldTestMessage{}

	err = n.UnmarshalObjectV2(data)
	if err != nil {
		t.Fatal(err)
	}

	v.sizeCache = n.sizeCache
	v.state = n.state
	v.unknownFields = n.unknownFields

	assert.EqualValues(t, v, n)
}

func BenchmarkProtobuf(b *testing.B) {
	v1 := &FieldTestMessage{
		OptionalBool:     "xxadf",
		OptionalEnum:     100,
		OptionalInt32:    990,
		OptionalSint32:   880,
		OptionalUint32:   550,
		OptionalInt64:    8880,
		OptionalSint64:   11110,
		OptionalUint64:   2220,
		OptionalSfixed32: -3330,
		OptionalFixed32:  4440,
		OptionalFloat:    5.0,
		OptionalSfixed64: -6.0,
		OptionalFixed64:  70,
		OptionalDouble:   888.80,
		OptionalString:   "sdf234",
		OptionalBytes:    []byte{1, 5, 2, 5, 6, 54, 5},
		RepeatedBool:     []bool{true, false, true, false},
		MapInt32Int64:    map[int32]int64{88: 2, 90: 5},
	}
	var data []byte
	b.Run("proto-marshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data, _ = proto.Marshal(v1)
		}
	})
	b.Run("custum-marshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data, _ = v1.MarshalObject()
		}
	})
	v1.MapInt32Int64 = nil
	data, _ = proto.Marshal(v1)
	b.Run("proto-unmarshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v1 = &FieldTestMessage{}
			proto.Unmarshal(data, v1)
		}
	})
	b.Run("custom-unmarshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v1 = &FieldTestMessage{}
			v1.UnmarshalObject(data)
		}
	})
	b.Run("custom-unmarshal-map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v1 = &FieldTestMessage{}
			v1.UnmarshalObjectV2(data)
		}
	})
}

func (m *FieldTestMessage) UnmarshalObject(data []byte) (err error) {
	index := 0
	ignoreGroup := 0
	for index < len(data) {
		num, typ, cnt := protowire.ConsumeTag(data[index:])
		if num == 0 {
			err = errors.New("invalid tag")
			return
		}

		index += cnt
		if ignoreGroup > 0 {
			switch typ {
			case protowire.VarintType:
				_, cnt := protowire.ConsumeVarint(data[index:])
				if cnt < 1 {
					err = protowire.ParseError(cnt)
					return
				}
				index += cnt
			case protowire.Fixed32Type:
				index += 4
			case protowire.Fixed64Type:
				index += 8
			case protowire.BytesType:
				v, cnt := protowire.ConsumeBytes(data[index:])
				if v == nil {
					if cnt < 0 {
						err = protowire.ParseError(cnt)
					} else {
						err = errors.New("invalid data")
					}
					return
				}
				index += cnt
			case protowire.StartGroupType:
				ignoreGroup++
			case protowire.EndGroupType:
				ignoreGroup--
			}
			continue
		}
		switch num {
		case 1:
			if typ != protowire.BytesType {
				err = errors.New("invalid field 1, not len type")
				return
			}
			v, cnt := protowire.ConsumeBytes(data[index:])
			if v == nil {
				err = errors.New("invalid field 1, invalid len value")
				return
			}
			m.OptionalBool = string(v)
			index += cnt
		case 2:
			if typ != protowire.VarintType {
				err = errors.New("invlaid field 2. not varint type")
				return
			}
			v, cnt := protowire.ConsumeVarint(data[index:])
			if cnt < 1 {
				err = errors.New("invalid feild 2, invalid varint value")
				return
			}
			index += cnt
			m.OptionalEnum = FieldTestMessage_Enum(v)
		case 3:
			if typ != protowire.VarintType {
				err = errors.New("invlaid field 3. not varint type")
				return
			}
			v, cnt := protowire.ConsumeVarint(data[index:])
			if cnt < 1 {
				err = errors.New("invalid feild 3, invalid varint value")
				return
			}
			index += cnt
			m.OptionalInt32 = int32(v)
		case 4:
			if typ != protowire.VarintType {
				err = errors.New("invlaid field 4. not varint type")
				return
			}
			v, cnt := protowire.ConsumeVarint(data[index:])
			if cnt < 1 {
				err = errors.New("invalid feild 4, invalid varint value")
				return
			}
			index += cnt
			m.OptionalSint32 = int32(protowire.DecodeZigZag(v))
		case 5:
			if typ != protowire.VarintType {
				err = errors.New("invlaid field 5. not varint type")
				return
			}
			v, cnt := protowire.ConsumeVarint(data[index:])
			if cnt < 1 {
				err = errors.New("invalid feild 5, invalid varint value")
				return
			}
			index += cnt
			m.OptionalUint32 = uint32(v)
		case 6:
			if typ != protowire.VarintType {
				err = errors.New("invlaid field 6. not varint type")
				return
			}
			v, cnt := protowire.ConsumeVarint(data[index:])
			if cnt < 1 {
				err = errors.New("invalid feild 6, invalid varint value")
				return
			}
			index += cnt
			m.OptionalInt64 = int64(v)
		case 7:
			if typ != protowire.VarintType {
				err = errors.New("invlaid field 7. not varint type")
				return
			}
			v, cnt := protowire.ConsumeVarint(data[index:])
			if cnt < 1 {
				err = errors.New("invalid feild 7, invalid varint value")
				return
			}
			index += cnt
			m.OptionalSint64 = protowire.DecodeZigZag(v)
		case 8:
			if typ != protowire.VarintType {
				err = errors.New("invlaid field 8. not varint type")
				return
			}
			v, cnt := protowire.ConsumeVarint(data[index:])
			if cnt < 1 {
				err = errors.New("invalid feild 8, invalid varint value")
				return
			}
			index += cnt
			m.OptionalUint64 = v
		case 9:
			if typ != protowire.Fixed32Type {
				err = errors.New("invlaid field 9. not i32 type")
				return
			}
			v, cnt := protowire.ConsumeFixed32(data[index:])
			if cnt < 1 {
				err = errors.New("invalid feild 9, invalid fixed32 value")
				return
			}
			index += cnt
			m.OptionalSfixed32 = int32(v) // int32(protowire.DecodeZigZag(uint64(v)))
		case 10:
			if typ != protowire.Fixed32Type {
				err = errors.New("invlaid field 10. not i32 type")
				return
			}
			v, cnt := protowire.ConsumeFixed32(data[index:])
			if cnt < 1 {
				err = errors.New("invalid feild 10, invalid fixed32 value")
				return
			}
			index += cnt
			m.OptionalFixed32 = v
		case 11:
			if typ != protowire.Fixed32Type {
				err = errors.New("invlaid field 11. not i32 type")
				return
			}
			v, cnt := protowire.ConsumeFixed32(data[index:])
			if cnt < 1 {
				err = errors.New("invalid feild 11, invalid fixed32 value")
				return
			}
			index += cnt
			m.OptionalFloat = math.Float32frombits(v)
		case 12:
			if typ != protowire.Fixed64Type {
				err = errors.New("invlaid field 12. not varint type")
				return
			}
			v, cnt := protowire.ConsumeFixed64(data[index:])
			if cnt < 1 {
				err = errors.New("invalid feild 12, invalid varint value")
				return
			}
			index += cnt
			m.OptionalSfixed64 = int64(v) //  protowire.DecodeZigZag(v)
		case 13:
			if typ != protowire.Fixed64Type {
				err = errors.New("invlaid field 13. not varint type")
				return
			}
			v, cnt := protowire.ConsumeFixed64(data[index:])
			if cnt < 1 {
				err = errors.New("invalid feild 13, invalid varint value")
				return
			}
			index += cnt
			m.OptionalFixed64 = v
		case 14:
			if typ != protowire.Fixed64Type {
				err = errors.New("invlaid field 14. not varint type")
				return
			}
			v, cnt := protowire.ConsumeFixed64(data[index:])
			if cnt < 1 {
				err = errors.New("invalid feild 14, invalid varint value")
				return
			}
			index += cnt
			m.OptionalDouble = math.Float64frombits(v)
		case 15:
			if typ != protowire.BytesType {
				err = errors.New("invlaid field 15. not varint type")
				return
			}
			v, cnt := protowire.ConsumeString(data[index:])
			if cnt < 1 {
				err = errors.New("invalid feild 15, invalid varint value")
				return
			}
			index += cnt
			m.OptionalString = v
		case 16:
			if typ != protowire.BytesType {
				err = errors.New("invlaid field 16. not varint type")
				return
			}
			v, cnt := protowire.ConsumeBytes(data[index:])
			if cnt < 1 {
				err = errors.New("invalid feild 16, invalid varint value")
				return
			}
			index += cnt
			m.OptionalBytes = make([]byte, len(v))
			copy(m.OptionalBytes, v)
		case 17:
			if typ != protowire.BytesType {
				err = errors.New("invlaid field 15. not varint type")
				return
			}
			v, cnt := protowire.ConsumeBytes(data[index:])
			if cnt < 1 {
				err = errors.New("invalid feild 15, invalid varint value")
				return
			}
			if m.Optional_Message == nil {
				m.Optional_Message = &FieldTestMessage_Message{}
			}
			err = m.Optional_Message.UnmarshalObject(v)
			if err != nil {
				return err
			}
			index += cnt
		case 201:
			// packed=false 方式的数据
			if typ == protowire.VarintType {
				v, cnt := protowire.ConsumeVarint(data[index:])
				if cnt < 1 {
					err = errors.New("invalid field 201,packed=false value invliad")
					return
				}
				m.RepeatedBool = append(m.RepeatedBool, protowire.DecodeBool(v))
				index += cnt
				continue
			}
			if typ != protowire.BytesType {
				err = errors.New("invalid field 201,type invliad")
				return
			}
			buf, cnt := protowire.ConsumeBytes(data[index:])
			//fmt.Println(buf, cnt)
			if cnt < 1 {
				// fmt.Println("invalid:", buf, cnt, index, len(data), data[index:])
				err = errors.New("invalid field 201, invalid data")
				return
			}
			index += cnt
			// 只有首次解析才能一次性申请内存. 因为需要支持packed=false的解析的.
			if m.RepeatedBool == nil {
				// 一个bool 一个varint, 又因为值只有 0,1 所以有多少个字节. 就有多少个bool值
				m.RepeatedBool = make([]bool, 0, cnt)
			}
			sub := 0
			for sub < len(buf) {
				v, cnt := protowire.ConsumeVarint(buf[sub:])
				if cnt < 1 {
					err = errors.New("invalid field 201 value.")
					return
				}
				sub += cnt
				m.RepeatedBool = append(m.RepeatedBool, protowire.DecodeBool(v))
			}
		case 202:
			// packed=false 方式的数据
			if typ == protowire.VarintType {
				v, cnt := protowire.ConsumeVarint(data[index:])
				if cnt < 1 {
					err = errors.New("invalid field 202,packed=false value invliad")
					return
				}
				m.RepeatedBool = append(m.RepeatedBool, protowire.DecodeBool(v))
				index += cnt
				continue
			}
			if typ != protowire.BytesType {
				err = errors.New("invalid field 202,type invliad")
				return
			}
			buf, cnt := protowire.ConsumeBytes(data[index:])
			if cnt < 1 {
				err = errors.New("invalid field 202, invalid data")
				return
			}
			index += cnt
			// 只有首次解析才能一次性申请内存. 因为需要支持packed=false的解析的.
			if m.RepeatedBool == nil {
				// 一个bool 一个varint, 又因为值只有 0,1 所以有多少个字节. 就有多少个bool值
				m.RepeatedBool = make([]bool, 0, cnt)
			}
			sub := 0
			for sub < len(buf) {
				v, cnt := protowire.ConsumeVarint(buf[sub:])
				if cnt < 1 {
					err = errors.New("invalid field 202 value.")
					return
				}
				sub += cnt
				m.RepeatedBool = append(m.RepeatedBool, protowire.DecodeBool(v))
			}
		case 500:
			if typ != protowire.BytesType {
				err = errors.New("invalid filed 500, not bytes type")
				return
			}
			buf, cnt := protowire.ConsumeBytes(data[index:])
			if buf == nil {
				err = errors.New("invalid field 500. invalid data")
				return
			}
			index += cnt
			if m.MapInt32Int64 == nil {
				m.MapInt32Int64 = make(map[int32]int64)
			}
			//fmt.Println("map value:", buf)

			for sub := 0; sub < len(buf); {
				// ignore key = 1 ,byte=8
				n, t, scnt := protowire.ConsumeTag(buf)
				if scnt < 1 {
					err = errors.New("invalid field 500 value, key tag invalid")
					return
				}
				if t != protowire.VarintType {
					err = errors.New("invalid field 500 value, key type invalid")
					return
				}
				if n != 1 {
					err = errors.New("invalid field 500 value, key sequence invalid")
					return
				}
				sub += scnt
				k, cnt := protowire.ConsumeVarint(buf[sub:])
				if cnt < 1 {
					err = errors.New("invalid field 500 value, key value invalid")
					return
				}
				// ignore key = 2 ,byte 16
				sub += cnt
				n, t, scnt = protowire.ConsumeTag(buf[sub:])
				if scnt < 1 {
					err = errors.New("invalid field 500 value, value tag invalid")
					return
				}
				if t != protowire.VarintType {
					err = errors.New("invalid field 500 value, value type invalid")
					return
				}
				if n != 2 {
					err = errors.New("invalid field 500 value, value sequence invalid")
					return
				}
				sub++
				v, cnt := protowire.ConsumeVarint(buf[sub:])
				if cnt < 1 {
					err = errors.New("invalid field 500 value, key value invalid")
					return
				}
				sub += cnt
				m.MapInt32Int64[int32(k)] = int64(v)
			}
		}
	}
	return
}

func (m *FieldTestMessage) MarshalObject() (data []byte, err error) {
	data = make([]byte, 0, m.ObjectSize())
	if len(m.OptionalBool) > 0 {
		data = protowire.AppendTag(data, 1, protowire.BytesType)
		//data = protowire.AppendBytes(data, utils.StringToBytes(m.OptionalBool))
		data = protowire.AppendBytes(data, *(*[]byte)(unsafe.Pointer(&m.OptionalBool)))
	}
	if m.OptionalEnum > 0 {
		data = protowire.AppendTag(data, 2, protowire.VarintType)
		data = protowire.AppendVarint(data, uint64(m.OptionalEnum))
	}
	if m.OptionalInt32 != 0 {
		data = protowire.AppendTag(data, 3, protowire.VarintType)
		data = protowire.AppendVarint(data, uint64(m.OptionalInt32))
	}
	if m.OptionalSint32 != 0 {
		data = protowire.AppendTag(data, 4, protowire.VarintType)
		data = protowire.AppendVarint(data, protowire.EncodeZigZag(int64(m.OptionalSint32)))
	}
	if m.OptionalUint32 != 0 {
		data = protowire.AppendTag(data, 5, protowire.VarintType)
		data = protowire.AppendVarint(data, uint64(m.OptionalUint32))
	}
	if m.OptionalInt64 != 0 {
		data = protowire.AppendTag(data, 6, protowire.VarintType)
		data = protowire.AppendVarint(data, uint64(m.OptionalInt64))
	}
	if m.OptionalSint64 != 0 {
		data = protowire.AppendTag(data, 7, protowire.VarintType)
		data = protowire.AppendVarint(data, protowire.EncodeZigZag(m.OptionalSint64))
	}
	if m.OptionalUint64 != 0 {
		data = protowire.AppendTag(data, 8, protowire.VarintType)
		data = protowire.AppendVarint(data, m.OptionalUint64)
	}

	if m.OptionalSfixed32 != 0 {
		data = protowire.AppendTag(data, 9, protowire.Fixed32Type)
		data = protowire.AppendFixed32(data, uint32(m.OptionalSfixed32))
	}
	if m.OptionalFixed32 != 0 {
		data = protowire.AppendTag(data, 10, protowire.Fixed32Type)
		data = protowire.AppendFixed32(data, m.OptionalFixed32)
	}
	if m.OptionalFloat != 0 {
		data = protowire.AppendTag(data, 11, protowire.Fixed32Type)
		data = protowire.AppendFixed32(data, math.Float32bits(m.OptionalFloat))
	}
	if m.OptionalSfixed64 != 0 {
		data = protowire.AppendTag(data, 12, protowire.Fixed64Type)
		data = protowire.AppendFixed64(data, uint64(m.OptionalSfixed64))
	}
	if m.OptionalFixed64 != 0 {
		data = protowire.AppendTag(data, 13, protowire.Fixed64Type)
		data = protowire.AppendFixed64(data, m.OptionalFixed64)
	}
	if m.OptionalDouble != 0 {
		data = protowire.AppendTag(data, 14, protowire.Fixed64Type)
		data = protowire.AppendFixed64(data, math.Float64bits(m.OptionalDouble))
	}
	if len(m.OptionalString) > 0 {
		data = protowire.AppendTag(data, 15, protowire.BytesType)
		data = protowire.AppendString(data, m.OptionalString)
	}
	if len(m.OptionalBytes) > 0 {
		data = protowire.AppendTag(data, 16, protowire.BytesType)
		data = protowire.AppendBytes(data, m.OptionalBytes)
	}
	if m.Optional_Message != nil {
		data = protowire.AppendTag(data, 17, protowire.BytesType)
		//
		tmp, err := m.Optional_Message.MarshalObject()
		if err != nil {
			return nil, err
		}
		data = protowire.AppendBytes(data, tmp)
	}
	//
	if len(m.RepeatedBool) > 0 {
		// fmt.Println("pre:", len(data))
		data = protowire.AppendTag(data, 201, protowire.BytesType)
		data = protowire.AppendVarint(data, uint64(len(m.RepeatedBool)))
		// fmt.Println("size:", len(m.RepeatedBool), protowire.SizeVarint(uint64(len(m.RepeatedBool))))
		for _, v := range m.RepeatedBool {
			data = protowire.AppendVarint(data, protowire.EncodeBool(v))
			// fmt.Println("index:", protowire.SizeVarint(protowire.EncodeBool(v)))
		}
	}
	if len(m.RepeatedInt32) > 0 {
		data = protowire.AppendTag(data, 202, protowire.BytesType)
		size := 0
		for _, v := range m.RepeatedInt32 {
			size += protowire.SizeVarint(uint64(v))
		}
		data = protowire.AppendVarint(data, uint64(size))
	}
	// //
	// if len(m.MapInt32Int64) > 0 {
	// 	//
	// 	data = protowire.AppendTag(data, 500, protowire.BytesType)
	// 	size := 0
	// 	for k, v := range m.MapInt32Int64 {
	// 		size += protowire.SizeTag(1)
	// 		size += protowire.SizeVarint(uint64(k))
	// 		size += protowire.SizeTag(2)
	// 		size += protowire.SizeVarint(uint64(v))
	// 	}
	// 	data = protowire.AppendVarint(data, uint64(size))
	// 	//
	// 	for k, v := range m.MapInt32Int64 {
	// 		data = protowire.AppendTag(data, 1, protowire.VarintType)
	// 		data = protowire.AppendVarint(data, uint64(k))
	// 		data = protowire.AppendTag(data, 2, protowire.VarintType)
	// 		data = protowire.AppendVarint(data, uint64(v))
	// 	}
	// }

	//
	if len(m.MapInt32Int64) > 0 {
		for k, v := range m.MapInt32Int64 {
			size := 0
			data = protowire.AppendTag(data, 500, protowire.BytesType)
			size += protowire.SizeTag(1)
			size += protowire.SizeVarint(uint64(k))
			size += protowire.SizeTag(2)
			size += protowire.SizeVarint(uint64(v))
			
			data = protowire.AppendVarint(data, uint64(size))

			data = protowire.AppendTag(data, 2, protowire.VarintType)
			data = protowire.AppendVarint(data, uint64(v))
			data = protowire.AppendTag(data, 1, protowire.VarintType)
			data = protowire.AppendVarint(data, uint64(k))
		}
	}

	//fmt.Println("all:", len(data), data)

	return
}

func (m *FieldTestMessage) MarshalPackedMap() (data []byte, err error) {
	if len(m.MapInt32Int64) > 0 {
		data = protowire.AppendTag(data, 500, protowire.BytesType)
		size := 0
		for k, v := range m.MapInt32Int64 {

			size += protowire.SizeTag(1)
			size += protowire.SizeVarint(uint64(k))
			size += protowire.SizeTag(2)
			size += protowire.SizeVarint(uint64(v))
		}

		data = protowire.AppendVarint(data, uint64(size))
		for k, v := range m.MapInt32Int64 {
			data = protowire.AppendTag(data, 1, protowire.VarintType)
			data = protowire.AppendVarint(data, uint64(k))
			data = protowire.AppendTag(data, 2, protowire.VarintType)
			data = protowire.AppendVarint(data, uint64(v))
		}
	}
	return
}

func (m *FieldTestMessage) ObjectSize() (size int) {

	size = 0
	for k, v := range m.MapInt32Int64 {
		one := 0
		one += protowire.SizeTag(1)
		one += protowire.SizeVarint(uint64(k))
		one += protowire.SizeTag(2)
		one += protowire.SizeVarint(uint64(v))
		one += protowire.SizeVarint(uint64(one))
		one += protowire.SizeTag(500)
		size += one
	}
	return 96 + size
}

func (m *FieldTestMessage_Message) UnmarshalObject(data []byte) (err error) {
	return
}

func (m *FieldTestMessage_Message) MarshalObject() (data []byte, err error) {
	return
}

func unmarshalOptionBool(typ protowire.Type, data []byte, m *FieldTestMessage) (cnt int, err error) {
	if typ != protowire.BytesType {
		err = errors.New("invalid field 1, not len type")
		return
	}
	buf, cnt := protowire.ConsumeBytes(data)
	if buf == nil {
		err = errors.New("invalid field 1, invalid len value")
		return
	}
	m.OptionalBool = string(buf)
	return
}

func unmarshalOptionEnum(typ protowire.Type, data []byte, m *FieldTestMessage) (cnt int, err error) {
	if typ != protowire.VarintType {
		err = errors.New("invlaid field 2. not varint type")
		return
	}
	v, cnt := protowire.ConsumeVarint(data)
	if cnt < 1 {
		err = errors.New("invalid feild 2, invalid varint value")
		return
	}
	m.OptionalEnum = FieldTestMessage_Enum(v)
	return
}

func unmarshalOptionInt32(typ protowire.Type, data []byte, m *FieldTestMessage) (cnt int, err error) {
	if typ != protowire.VarintType {
		err = errors.New("invlaid field 2. not varint type")
		return
	}
	v, cnt := protowire.ConsumeVarint(data)
	if cnt < 1 {
		err = errors.New("invalid feild 2, invalid varint value")
		return
	}
	m.OptionalInt32 = int32(v)
	return
}

func unmarshalOptionSint32(typ protowire.Type, data []byte, m *FieldTestMessage) (cnt int, err error) {
	if typ != protowire.VarintType {
		err = errors.New("invlaid field 2. not varint type")
		return
	}
	v, cnt := protowire.ConsumeVarint(data)
	if cnt < 1 {
		err = errors.New("invalid feild 2, invalid varint value")
		return
	}
	m.OptionalSint32 = int32(protowire.DecodeZigZag(v))
	return
}

func unmarshalOptionUint32(typ protowire.Type, data []byte, m *FieldTestMessage) (cnt int, err error) {
	if typ != protowire.VarintType {
		err = errors.New("invlaid field 2. not varint type")
		return
	}
	v, cnt := protowire.ConsumeVarint(data)
	if cnt < 1 {
		err = errors.New("invalid feild 2, invalid varint value")
		return
	}
	m.OptionalUint32 = uint32(v)
	return
}

func unmarshalOptionInt64(typ protowire.Type, data []byte, m *FieldTestMessage) (cnt int, err error) {
	if typ != protowire.VarintType {
		err = errors.New("invlaid field 2. not varint type")
		return
	}
	v, cnt := protowire.ConsumeVarint(data)
	if cnt < 1 {
		err = errors.New("invalid feild 2, invalid varint value")
		return
	}
	m.OptionalInt64 = int64(v)
	return
}

func unmarshalOptionSint64(typ protowire.Type, data []byte, m *FieldTestMessage) (cnt int, err error) {
	if typ != protowire.VarintType {
		err = errors.New("invlaid field 2. not varint type")
		return
	}
	v, cnt := protowire.ConsumeVarint(data)
	if cnt < 1 {
		err = errors.New("invalid feild 2, invalid varint value")
		return
	}

	m.OptionalSint64 = protowire.DecodeZigZag(v)
	return
}

func unmarshalOptionUint64(typ protowire.Type, data []byte, m *FieldTestMessage) (cnt int, err error) {
	if typ != protowire.VarintType {
		err = errors.New("invlaid field 2. not varint type")
		return
	}
	v, cnt := protowire.ConsumeVarint(data)
	if cnt < 1 {
		err = errors.New("invalid feild 2, invalid varint value")
		return
	}

	m.OptionalUint64 = v

	return
}

func unmarshalOptionSfix32(typ protowire.Type, data []byte, m *FieldTestMessage) (cnt int, err error) {
	if typ != protowire.Fixed32Type {
		err = errors.New("invlaid field 9. not i32 type")
		return
	}
	v, cnt := protowire.ConsumeFixed32(data)
	if cnt < 1 {
		err = errors.New("invalid feild 9, invalid fixed32 value")
		return
	}

	m.OptionalSfixed32 = int32(v)

	return
}

func unmarshalOptionFix32(typ protowire.Type, data []byte, m *FieldTestMessage) (cnt int, err error) {
	if typ != protowire.Fixed32Type {
		err = errors.New("invlaid field 9. not i32 type")
		return
	}
	v, cnt := protowire.ConsumeFixed32(data)
	if cnt < 1 {
		err = errors.New("invalid feild 9, invalid fixed32 value")
		return
	}

	m.OptionalFixed32 = v

	return
}

func unmarshalOptionFloat(typ protowire.Type, data []byte, m *FieldTestMessage) (cnt int, err error) {
	if typ != protowire.Fixed32Type {
		err = errors.New("invlaid field 9. not i32 type")
		return
	}
	v, cnt := protowire.ConsumeFixed32(data)
	if cnt < 1 {
		err = errors.New("invalid feild 9, invalid fixed32 value")
		return
	}

	m.OptionalFloat = math.Float32frombits(v)
	return
}

func unmarshalOptionSfix64(typ protowire.Type, data []byte, m *FieldTestMessage) (cnt int, err error) {
	if typ != protowire.Fixed64Type {
		err = errors.New("invlaid field 9. not i32 type")
		return
	}
	v, cnt := protowire.ConsumeFixed64(data)
	if cnt < 1 {
		err = errors.New("invalid feild 9, invalid fixed32 value")
		return
	}

	m.OptionalSfixed64 = int64(v)

	return
}

func unmarshalOptionFix64(typ protowire.Type, data []byte, m *FieldTestMessage) (cnt int, err error) {
	if typ != protowire.Fixed64Type {
		err = errors.New("invlaid field 9. not i32 type")
		return
	}
	v, cnt := protowire.ConsumeFixed64(data)
	if cnt < 1 {
		err = errors.New("invalid feild 9, invalid fixed32 value")
		return
	}

	m.OptionalFixed64 = v

	return
}

func unmarshalOptionDouble(typ protowire.Type, data []byte, m *FieldTestMessage) (cnt int, err error) {
	if typ != protowire.Fixed64Type {
		err = errors.New("invlaid field 9. not i32 type")
		return
	}
	v, cnt := protowire.ConsumeFixed64(data)
	if cnt < 1 {
		err = errors.New("invalid feild 9, invalid fixed32 value")
		return
	}

	m.OptionalDouble = math.Float64frombits(v)

	return
}

func unmarshalOptionString(typ protowire.Type, data []byte, m *FieldTestMessage) (cnt int, err error) {
	if typ != protowire.BytesType {
		err = errors.New("invalid field 1, not len type")
		return
	}
	buf, cnt := protowire.ConsumeBytes(data)
	if buf == nil {
		err = errors.New("invalid field 1, invalid len value")
		return
	}
	m.OptionalString = string(buf)
	return
}

func unmarshalOptionBytes(typ protowire.Type, data []byte, m *FieldTestMessage) (cnt int, err error) {
	if typ != protowire.BytesType {
		err = errors.New("invalid field 1, not len type")
		return
	}
	buf, cnt := protowire.ConsumeBytes(data)
	if buf == nil {
		err = errors.New("invalid field 1, invalid len value")
		return
	}
	m.OptionalBytes = make([]byte, len(buf))
	copy(m.OptionalBytes, buf)
	return
}

func unmarshalOptionMessage(typ protowire.Type, data []byte, m *FieldTestMessage) (cnt int, err error) {
	if typ != protowire.BytesType {
		err = errors.New("invlaid field 15. not varint type")
		return
	}
	v, cnt := protowire.ConsumeBytes(data)
	if cnt < 1 {
		err = errors.New("invalid feild 15, invalid varint value")
		return
	}
	if m.Optional_Message == nil {
		m.Optional_Message = &FieldTestMessage_Message{}
	}
	err = m.Optional_Message.UnmarshalObject(v)
	if err != nil {
		return
	}
	return
}

func unmarshalRepeatedBool(typ protowire.Type, data []byte, m *FieldTestMessage) (cnt int, err error) {
	// packed=false 方式的数据
	if typ == protowire.VarintType {
		var v uint64
		v, cnt = protowire.ConsumeVarint(data)
		if cnt < 1 {
			err = errors.New("invalid field 201,packed=false value invliad")
			return
		}
		m.RepeatedBool = append(m.RepeatedBool, protowire.DecodeBool(v))
		return
	}
	if typ != protowire.BytesType {
		err = errors.New("invalid field 201,type invliad")
		return
	}
	buf, cnt := protowire.ConsumeBytes(data)
	//fmt.Println(buf, cnt)
	if cnt < 1 {
		// fmt.Println("invalid:", buf, cnt, index, len(data), data[index:])
		err = errors.New("invalid field 201, invalid data")
		return
	}
	// 只有首次解析才能一次性申请内存. 因为需要支持packed=false的解析的.
	if m.RepeatedBool == nil {
		// 一个bool 一个varint, 又因为值只有 0,1 所以有多少个字节. 就有多少个bool值
		m.RepeatedBool = make([]bool, 0, cnt)
	}
	sub := 0
	for sub < len(buf) {
		v, scnt := protowire.ConsumeVarint(buf[sub:])
		if scnt < 1 {
			err = errors.New("invalid field 201 value.")
			return
		}
		sub += scnt
		m.RepeatedBool = append(m.RepeatedBool, protowire.DecodeBool(v))
	}
	return
}

func unmarshalMapInt32Int64(typ protowire.Type, data []byte, m *FieldTestMessage) (cnt int, err error) {
	if typ != protowire.BytesType {
		err = errors.New("invalid filed 500, not bytes type")
		return
	}
	buf, cnt := protowire.ConsumeBytes(data)
	if buf == nil {
		err = errors.New("invalid field 500. invalid data")
		return
	}
	if m.MapInt32Int64 == nil {
		m.MapInt32Int64 = make(map[int32]int64)
	}
	//fmt.Println("map value:", buf)

	for sub := 0; sub < len(buf); {
		//fmt.Println(protowire.ConsumeTag(buf))
		// ignore key = 1 ,byte=8
		sub += 1
		k, scnt := protowire.ConsumeVarint(buf[sub:])
		// fmt.Println(k, cnt, m.MapInt32Int64)
		if scnt < 1 {
			err = errors.New("invalid field 500 value, key type invalid")
			return
		}
		// ignore key = 2 ,byte 16
		//fmt.Println(protowire.ConsumeTag(buf[sub+cnt:]))
		sub += scnt + 1
		v, scnt := protowire.ConsumeVarint(buf[sub:])
		if scnt < 1 {
			err = errors.New("invalid field 500 value, key value invalid")
			return
		}
		sub += scnt
		m.MapInt32Int64[int32(k)] = int64(v)
	}
	return
}

var (
	fieldTestMessageMap = map[protowire.Number]func(typ protowire.Type, data []byte, m *FieldTestMessage) (cnt int, err error){
		1:   unmarshalOptionBool,
		2:   unmarshalOptionEnum,
		3:   unmarshalOptionInt32,
		4:   unmarshalOptionSint32,
		5:   unmarshalOptionUint32,
		6:   unmarshalOptionInt64,
		7:   unmarshalOptionSint64,
		8:   unmarshalOptionUint64,
		9:   unmarshalOptionSfix32,
		10:  unmarshalOptionFix32,
		11:  unmarshalOptionFloat,
		12:  unmarshalOptionSfix64,
		13:  unmarshalOptionFix64,
		14:  unmarshalOptionDouble,
		15:  unmarshalOptionString,
		16:  unmarshalOptionBytes,
		17:  unmarshalOptionMessage,
		201: unmarshalRepeatedBool,
		500: unmarshalMapInt32Int64,
	}
)

func (m *FieldTestMessage) UnmarshalObjectV2(data []byte) (err error) {
	index := 0
	for index < len(data) {
		num, typ, cnt := protowire.ConsumeTag(data[index:])
		if num == 0 {
			err = errors.New("invalid tag")
			return
		}
		index += cnt
		f, ok := fieldTestMessageMap[num]
		if !ok {
			// TODO skip field
			continue
		}
		cnt, err = f(typ, data[index:], m)
		if err != nil {
			return err
		}
		index += cnt
	}
	return
}

func TestMarshalPackedMap(t *testing.T) {
	v1 := &FieldTestMessage{
		MapInt32Int64: map[int32]int64{1: 100, 3: 100, 4: 4, 100: 9},
	}

	// data, err := v1.MarshalPackedMap()
	// if err != nil {
	// 	t.Fatal(err)
	// }

	n1 := &FieldTestMessage{}
	// err = proto.Unmarshal(data, n1)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// failed. map 使用使用packed=false 方式 .
	//assert.EqualValues(t, v1.MapInt32Int64, n1.MapInt32Int64, "compare value")

	data, err := v1.MarshalObject()
	if err != nil {
		t.Fatal(err)
	}
	err = proto.Unmarshal(data, n1)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, v1.MapInt32Int64, n1.MapInt32Int64, "compare value 12")

}
