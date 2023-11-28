package gengo

var GenProtobufTemplate = map[string]string{
	"check.bool": `
		if typ != protowire.VarintType {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : not varint type")
			return
		}
	`,
	"check.varint": `
		if typ != protowire.VarintType {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : not varint type")
			return
		}
	`,
	"check.sint": `
		if typ != protowire.VarintType {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : not varint zigzag type")
			return
		}
	`,
	"check.fix32": `
		if typ != protowire.Fixed32Type {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : not i32 type")
			return
		}
	`,
	"check.float": `
		if typ != protowire.Fixed32Type {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : not i32 type")
			return
		}
	`,
	"check.fix64": `
		if typ != protowire.Fixed64Type {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : not i64 type")
			return
		}
	`,
	"check.double": `
		if typ != protowire.Fixed64Type {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : not i64 type")
			return
		}
	`,
	"check.string": `
		if typ != protowire.BytesType {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : not len type")
			return
		}
	`,
	"check.bytes": `
		if typ != protowire.BytesType {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : not len type")
			return
		}
	`,
	"check.message": `
		if typ != protowire.BytesType {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : not len type")
			return
		}
	`,

	"encode.bool": `
		// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, {{ .Field.WireType }}) => {{ TagBinary .Field.DescNum .Field.WireType }} 
		{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum .Field.WireType }})
		{{.V.Buffer}} = protowire.AppendVarint({{.V.Buffer}}, protowire.EncodeBool({{.V.VName}}))
	`,
	"encode.varint": `
		// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, {{ .Field.WireType }}) => {{ TagBinary .Field.DescNum .Field.WireType }}
		{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum .Field.WireType }})
		{{.V.Buffer}} = protowire.AppendVarint({{.V.Buffer}}, uint64({{.V.VName}}))
	`,
	"encode.sint": `
		// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, {{ .Field.WireType }}) => {{ TagBinary .Field.DescNum .Field.WireType }}
		{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum .Field.WireType }})
		{{.V.Buffer}} = protowire.AppendVarint({{.V.Buffer}}, protowire.EncodeZigZag(int64({{.V.VName}})))
	`,
	"encode.fix32": `
		// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, {{ .Field.WireType }}) => {{ TagBinary .Field.DescNum .Field.WireType }}
		{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum .Field.WireType }})
		{{.V.Buffer}} = protowire.AppendFixed32({{.V.Buffer}}, uint32({{.V.VName}}))
	`,
	"encode.float": `
		// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, {{ .Field.WireType }}) => {{ TagBinary .Field.DescNum .Field.WireType }}
		{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum .Field.WireType }})
		{{.V.Buffer}} = protowire.AppendFixed32({{.V.Buffer}}, math.Float32bits({{.V.VName}}))
	`,
	"encode.fix64": `
		// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, {{ .Field.WireType }}) => {{ TagBinary .Field.DescNum .Field.WireType }}
		{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum .Field.WireType }})
		{{.V.Buffer}} = protowire.AppendFixed64({{.V.Buffer}}, uint64({{.V.VName}}))
	`,
	"encode.double": `
		// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, {{ .Field.WireType }}) => {{ TagBinary .Field.DescNum .Field.WireType }}
		{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum .Field.WireType }})
		{{.V.Buffer}} = protowire.AppendFixed64({{.V.Buffer}}, math.Float64bits({{.V.VName}}))
	`,
	"encode.string": `
		// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, {{ .Field.WireType }}) => {{ TagBinary .Field.DescNum .Field.WireType }}
		{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum .Field.WireType }})
		{{.V.Buffer}} = protowire.AppendString({{.V.Buffer}}, {{.V.VName}})
	`,
	"encode.bytes": `
		// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, {{ .Field.WireType }}) => {{ TagBinary .Field.DescNum .Field.WireType }}
		{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum .Field.WireType }})
		{{.V.Buffer}} = protowire.AppendBytes({{.V.Buffer}}, {{.V.VName}})
	`,
	"encode.message": `
		// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, {{ .Field.WireType }}) => {{ TagBinary .Field.DescNum .Field.WireType }}
		{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum .Field.WireType }})
		{{.V.Buffer}} = protowire.AppendVarint({{.V.Buffer}}, uint64({{.V.VName}}.MarshalSize()))
		{{.V.Buffer}}, err = {{.V.VName}}.MarshalObjectTo({{.V.Buffer}})
		if err != nil {
			return
		}
	`,

	"encode.packed.bool": `
		// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, protowire.BytesType) => {{ TagBinary .Field.DescNum "protowire.BytesType" }}
		{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum "protowire.BytesType" }})
		{{.V.Buffer}} = protowire.AppendVarint({{.V.Buffer}}, uint64(len({{.V.VName}})))
		for _, v := range {{.V.VName}} {
			{{.V.Buffer}} = protowire.AppendVarint({{.V.Buffer}}, protowire.EncodeBool(v))
		}
	`,
	"encode.packed.varint": `
		// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, protowire.BytesType) => {{ TagBinary .Field.DescNum "protowire.BytesType" }}
		{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum "protowire.BytesType" }})
		size := 0
		for _, v := range {{.V.VName}} {
			size += protowire.SizeVarint(uint64(v))
		}
		{{.V.Buffer}} = protowire.AppendVarint({{.V.Buffer}}, uint64(size))
		for _, v := range {{.V.VName}} {
			{{.V.Buffer}} = protowire.AppendVarint({{.V.Buffer}}, uint64(v))
		}
	`,
	"encode.packed.sint": `
		// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, protowire.BytesType) => {{ TagBinary .Field.DescNum "protowire.BytesType" }}
		{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum "protowire.BytesType" }})
		size := 0
		for _, v := range {{.V.VName}} {
			size += protowire.SizeVarint(protowire.EncodeZigZag(int64(v)))
		}
		{{.V.Buffer}} = protowire.AppendVarint({{.V.Buffer}}, uint64(size))
		for _, v := range {{.V.VName}} {
			{{.V.Buffer}} = protowire.AppendVarint({{.V.Buffer}}, protowire.EncodeZigZag(int64(v)))
		}
	`,
	"encode.packed.fix32": `
		// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, protowire.BytesType) => {{ TagBinary .Field.DescNum "protowire.BytesType" }}
		{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum "protowire.BytesType" }})
		{{.V.Buffer}} = protowire.AppendVarint({{.V.Buffer}}, uint64(4*len({{.V.VName}})))
		for _, v := range {{.V.VName}} {
			{{.V.Buffer}} = protowire.AppendFixed32({{.V.Buffer}}, uint32(v))
		}
	`,
	"encode.packed.float": `
		// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, protowire.BytesType) => {{ TagBinary .Field.DescNum "protowire.BytesType" }}
		{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum "protowire.BytesType" }})
		{{.V.Buffer}} = protowire.AppendVarint({{.V.Buffer}}, uint64(4*len({{.V.VName}})))
		for _, v := range {{.V.VName}} {
			{{.V.Buffer}} = protowire.AppendFixed32({{.V.Buffer}}, math.Float32bits(v))
		}
	`,
	"encode.packed.fix64": `
		// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, protowire.BytesType) => {{ TagBinary .Field.DescNum "protowire.BytesType" }}
		{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum "protowire.BytesType" }})
		{{.V.Buffer}} = protowire.AppendVarint({{.V.Buffer}}, uint64(8*len({{.V.VName}})))
		for _, v := range {{.V.VName}} {
			{{.V.Buffer}} = protowire.AppendFixed64({{.V.Buffer}}, uint64(v))
		}
	`,
	"encode.packed.double": `
		// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, protowire.BytesType) => {{ TagBinary .Field.DescNum "protowire.BytesType" }}
		{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum "protowire.BytesType" }})
		{{.V.Buffer}} = protowire.AppendVarint({{.V.Buffer}}, uint64(8*len({{.V.VName}})))
		for _, v := range {{.V.VName}} {
			{{.V.Buffer}} = protowire.AppendFixed64({{.V.Buffer}}, math.Float64bits(v))
		}
	`,
	"encode.packed.string": `
		for k:=0; k<len({{.V.VName}}); k++ {
			// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, protowire.BytesType) => {{ TagBinary .Field.DescNum "protowire.BytesType" }}
			{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum "protowire.BytesType" }})
			{{.V.Buffer}} = protowire.AppendString({{.V.Buffer}}, {{.V.VName}}[k])
		}
	`,
	"encode.packed.bytes": `
		for _, item := range {{.V.VName}} {
			// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, protowire.BytesType) => {{ TagBinary .Field.DescNum "protowire.BytesType" }}
			{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum "protowire.BytesType" }})
			{{.V.Buffer}} = protowire.AppendBytes({{.V.Buffer}}, item)
		}
	`,
	"encode.packed.message": `
		for _, item := range {{.V.VName}} {
			// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, protowire.BytesType) => {{ TagBinary .Field.DescNum "protowire.BytesType" }}
			{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum "protowire.BytesType" }})
			{{.V.Buffer}} = protowire.AppendVarint({{.V.Buffer}}, uint64(item.MarshalSize()))
			{{.V.Buffer}}, err = item.MarshalObjectTo({{.V.Buffer}})
			if err != nil {
				return
			}
		}
	`,

	"encode.nopack.bool": `
		for _,item := range {{.V.VName}} {
			// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, {{ .Field.WireType }}) => {{ TagBinary .Field.DescNum .Field.WireType }}
			{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum .Field.WireType }})
			{{.V.Buffer}} = protowire.AppendVarint({{.V.Buffer}}, protowire.EncodeBool(item))
		}
	`,
	"encode.nopack.varint": `
		for _,item := range {{.V.VName}} {
			// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, {{ .Field.WireType }}) => {{ TagBinary .Field.DescNum .Field.WireType }}
			{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum .Field.WireType }})
			{{.V.Buffer}} = protowire.AppendVarint({{.V.Buffer}}, uint64(item))
		}
	`,
	"encode.nopack.sint": `
		for _,item := range {{.V.VName}} {
			// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, {{ .Field.WireType }}) => {{ TagBinary .Field.DescNum .Field.WireType }}
			{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum .Field.WireType }})
			{{.V.Buffer}} = protowire.AppendVarint({{.V.Buffer}}, protowire.EncodeZigZag(int64(item)))
		}
	`,
	"encode.nopack.fix32": `
		for _,item := range {{.V.VName}} {
			// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, {{ .Field.WireType }}) => {{ TagBinary .Field.DescNum .Field.WireType }}
			{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum .Field.WireType }})
			{{.V.Buffer}} = protowire.AppendFixed32({{.V.Buffer}}, uint32(item))
		}
	`,
	"encode.nopack.float": `
		for _,item := range {{.V.VName}} {
			// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, {{ .Field.WireType }}) => {{ TagBinary .Field.DescNum .Field.WireType }}
			{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum .Field.WireType }})
			{{.V.Buffer}} = protowire.AppendFixed32({{.V.Buffer}}, math.Float32bits(item))
		}
	`,
	"encode.nopack.fix64": `
		for _,item := range {{.V.VName}} {
			// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, {{ .Field.WireType }}) => {{ TagBinary .Field.DescNum .Field.WireType }}
			{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum .Field.WireType }})
			{{.V.Buffer}} = protowire.AppendFixed64({{.V.Buffer}}, uint64(item))
		}
	`,
	"encode.nopack.double": `
		for _,item := range {{.V.VName}} {
			// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, {{ .Field.WireType }}) => {{ TagBinary .Field.DescNum .Field.WireType }}
			{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum .Field.WireType }})
			{{.V.Buffer}} = protowire.AppendFixed64({{.V.Buffer}}, math.Float64bits(item))
		}
	`,
	"encode.nopack.string": `
		for k:=0; k<len ({{.V.VName}}); k++ {
			// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, {{ .Field.WireType }}) => {{ TagBinary .Field.DescNum .Field.WireType }}
			{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum .Field.WireType }})
			{{.V.Buffer}} = protowire.AppendString({{.V.Buffer}}, {{.V.VName}}[k])
		}
	`,
	"encode.nopack.bytes": `
		for k:=0; k<len ({{.V.VName}}); k++ {
			// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, {{ .Field.WireType }}) => {{ TagBinary .Field.DescNum .Field.WireType }}
			{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum .Field.WireType }})
			{{.V.Buffer}} = protowire.AppendBytes({{.V.Buffer}}, {{.V.VName}}[k])
		}
	`,
	"encode.nopack.message": `
		for _,item := range {{.V.VName}} {
			// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, {{ .Field.WireType }}) => {{ TagBinary .Field.DescNum .Field.WireType }}
			{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum .Field.WireType }})
			{{.V.Buffer}} = protowire.AppendVarint({{.V.Buffer}}, uint64(item.MarshalSize()))
			{{.V.Buffer}}, err = item.MarshalObjectTo({{.V.Buffer}})
			if err != nil {
				return
			}
		}
	`,

	"decode.bool": `
		v, cnt := protowire.ConsumeVarint({{.V.Buffer}})
		if cnt < 1 {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid varint value")
			return
		}
		{{.V.Index}} += cnt
		{{.V.VName}} = protowire.DecodeBool(v)
	`,
	"decode.varint": `
		v, cnt := protowire.ConsumeVarint({{.V.Buffer}})
		if cnt < 1 {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid varint value")
			return
		}
		{{.V.Index}} += cnt
		{{.V.VName}} = {{.Field.GoType}}(v)
	`,
	"decode.sint": `
		v, cnt := protowire.ConsumeVarint({{.V.Buffer}})
		if cnt < 1 {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid varint zigzag value")
			return
		}
		{{.V.Index}} += cnt
		{{.V.VName}} = {{.Field.GoType}}(protowire.DecodeZigZag(v))
	`,
	"decode.fix32": `
		v, cnt := protowire.ConsumeFixed32({{.V.Buffer}})
		if cnt < 1 {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid i32 value")
			return
		}
		{{.V.Index}} += cnt
		{{.V.VName}} = {{.Field.GoType}}(v)
	`,
	"decode.float": `
		v, cnt := protowire.ConsumeFixed32({{.V.Buffer}})
		if cnt < 1 {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid i32 value")
			return
		}
		{{.V.Index}} += cnt
		{{.V.VName}} = math.Float32frombits(v)
	`,
	"decode.fix64": `
		v, cnt := protowire.ConsumeFixed64({{.V.Buffer}})
		if cnt < 1 {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid i64 value")
			return
		}
		{{.V.Index}} += cnt
		{{.V.VName}} = {{.Field.GoType}}(v)
	`,
	"decode.double": `
		v, cnt := protowire.ConsumeFixed64({{.V.Buffer}})
		if cnt < 1 {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid i64 value")
			return
		}
		{{.V.Index}} += cnt
		{{.V.VName}} = math.Float64frombits(v)
	`,
	"decode.string": `
		v, cnt := protowire.ConsumeString({{.V.Buffer}})
		if cnt < 1 {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid len value")
			return
		}
		{{.V.Index}} += cnt
		{{.V.VName}} = v
	`,
	"decode.bytes": `
		v, cnt := protowire.ConsumeBytes({{.V.Buffer}})
		if v == nil {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid len value")
			return
		}
		{{.V.Index}} += cnt
		{{.V.VName}} = make([]byte, len(v))
		copy({{.V.VName}}, v)
	`,
	"decode.message": `
		v, cnt := protowire.ConsumeBytes({{.V.Buffer}})
		if v == nil {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid message value")
			return
		}
		{{.V.Index}} += cnt
		{{.V.VName}} = &{{.Field.GoType}}{}
		err = {{.V.VName}}.UnmarshalObject(v)
		if err != nil {
			return
		}
	`,

	"decode.slice.bool": `
		// packed=false
		if typ == protowire.VarintType {
			v, cnt := protowire.ConsumeVarint({{.V.Buffer}})
			if cnt < 1 {
				err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid varint value")
				return
			}
			{{.V.VName}} = append({{.V.VName}}, protowire.DecodeBool(v))
			{{.V.Index}} += cnt
			continue
		}
		// packed = true
		if typ != protowire.BytesType {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid repeated tag value")
			return
		}
		buf, cnt := protowire.ConsumeBytes({{.V.Buffer}})
		if buf == nil {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid len value")
			return
		}
		{{.V.Index}} += cnt
		if {{.V.VName}} == nil {
			{{.V.VName}} = make([]bool, 0, cnt)
		}
		sub := 0
		for sub < len(buf) {
			v, cnt := protowire.ConsumeVarint(buf[sub:])
			if cnt < 1 {
				err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid item value")
				return
			}
			sub += cnt
			{{.V.VName}} = append({{.V.VName}}, protowire.DecodeBool(v))
		}
	`,
	"decode.slice.varint": `
		// packed=false
		if typ == protowire.VarintType {
			v, cnt := protowire.ConsumeVarint({{.V.Buffer}})
			if cnt < 1 {
				err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid varint value")
				return
			}
			{{.V.VName}} = append({{.V.VName}}, {{.Field.GoType}}(v))
			{{.V.Index}} += cnt
			continue
		}
		// packed = true
		if typ != protowire.BytesType {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid repeated tag value")
			return
		}
		buf, cnt := protowire.ConsumeBytes({{.V.Buffer}})
		if buf == nil {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid len value")
			return
		}
		{{.V.Index}} += cnt
		if {{.V.VName}} == nil {
			{{.V.VName}} = make([]{{.Field.GoType}}, 0, 2)
		}
		sub := 0
		for sub < len(buf) {
			v, cnt := protowire.ConsumeVarint(buf[sub:])
			if cnt < 1 {
				err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid item value")
				return
			}
			sub += cnt
			{{.V.VName}} = append({{.V.VName}}, {{.Field.GoType}}(v))
		}
	`,
	"decode.slice.sint": `
		// packed=false
		if typ == protowire.VarintType {
			v, cnt := protowire.ConsumeVarint({{.V.Buffer}})
			if cnt < 1 {
				err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid varint value")
				return
			}
			{{.V.VName}} = append({{.V.VName}}, {{.Field.GoType}}(protowire.DecodeZigZag(v)))
			{{.V.Index}} += cnt
			continue
		}
		// packed = true
		if typ != protowire.BytesType {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid repeated tag value")
			return
		}
		buf, cnt := protowire.ConsumeBytes({{.V.Buffer}})
		if buf == nil {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid len value")
			return
		}
		{{.V.Index}} += cnt
		if {{.V.VName}} == nil {
			{{.V.VName}} = make([]{{.Field.GoType}}, 0, 2)
		}
		sub := 0
		for sub < len(buf) {
			v, cnt := protowire.ConsumeVarint(buf[sub:])
			if cnt < 1 {
				err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid item value")
				return
			}
			sub += cnt
			{{.V.VName}} = append({{.V.VName}}, {{.Field.GoType}}(protowire.DecodeZigZag(v)))
		}
	`,
	"decode.slice.fix32": `
		// packed=false
		if typ == protowire.Fixed32Type {
			v, cnt := protowire.ConsumeFixed32({{.V.Buffer}})
			if cnt < 1 {
				err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid varint value")
				return
			}
			{{.V.VName}} = append({{.V.VName}}, {{.Field.GoType}}(v))
			{{.V.Index}} += cnt
			continue
		}
		// packed = true
		if typ != protowire.BytesType {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid repeated tag value")
			return
		}
		buf, cnt := protowire.ConsumeBytes({{.V.Buffer}})
		if buf == nil {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid len value")
			return
		}
		{{.V.Index}} += cnt
		if {{.V.VName}} == nil {
			{{.V.VName}} = make([]{{.Field.GoType}}, 0, cnt/4)
		}
		sub := 0
		for sub < len(buf) {
			v, cnt := protowire.ConsumeFixed32(buf[sub:])
			if cnt < 1 {
				err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid item value")
				return
			}
			sub += cnt
			{{.V.VName}} = append({{.V.VName}}, {{.Field.GoType}}(v))
		}
	`,
	"decode.slice.float": `
		// packed=false
		if typ == protowire.Fixed32Type {
			v, cnt := protowire.ConsumeFixed32({{.V.Buffer}})
			if cnt < 1 {
				err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid varint value")
				return
			}
			{{.V.VName}} = append({{.V.VName}}, math.Float32frombits(v))
			{{.V.Index}} += cnt
			continue
		}
		// packed = true
		if typ != protowire.BytesType {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid repeated tag value")
			return
		}
		buf, cnt := protowire.ConsumeBytes({{.V.Buffer}})
		if buf == nil {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid len value")
			return
		}
		{{.V.Index}} += cnt
		if {{.V.VName}} == nil {
			{{.V.VName}} = make([]float32, 0, cnt/4)
		}
		sub := 0
		for sub < len(buf) {
			v, cnt := protowire.ConsumeFixed32(buf[sub:])
			if cnt < 1 {
				err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid item value")
				return
			}
			sub += cnt
			{{.V.VName}} = append({{.V.VName}}, math.Float32frombits(v))
		}
	`,
	"decode.slice.fix64": `
		// packed=false
		if typ == protowire.Fixed64Type {
			v, cnt := protowire.ConsumeFixed64({{.V.Buffer}})
			if cnt < 1 {
				err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid varint value")
				return
			}
			{{.V.VName}} = append({{.V.VName}}, {{.Field.GoType}}(v))
			{{.V.Index}} += cnt
			continue
		}
		// packed = true
		if typ != protowire.BytesType {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid repeated tag value")
			return
		}
		buf, cnt := protowire.ConsumeBytes({{.V.Buffer}})
		if buf == nil {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid len value")
			return
		}
		{{.V.Index}} += cnt
		if {{.V.VName}} == nil {
			{{.V.VName}} = make([]{{.Field.GoType}}, 0, cnt/4)
		}
		sub := 0
		for sub < len(buf) {
			v, cnt := protowire.ConsumeFixed64(buf[sub:])
			if cnt < 1 {
				err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid item value")
				return
			}
			sub += cnt
			{{.V.VName}} = append({{.V.VName}}, {{.Field.GoType}}(v))
		}
	`,
	"decode.slice.double": `
		// packed=false
		if typ == protowire.Fixed64Type {
			v, cnt := protowire.ConsumeFixed64({{.V.Buffer}})
			if cnt < 1 {
				err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid varint value")
				return
			}
			{{.V.VName}} = append({{.V.VName}}, math.Float64frombits(v))
			{{.V.Index}} += cnt
			continue
		}
		// packed = true
		if typ != protowire.BytesType {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid repeated tag value")
			return
		}
		buf, cnt := protowire.ConsumeBytes({{.V.Buffer}})
		if buf == nil {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid len value")
			return
		}
		{{.V.Index}} += cnt
		if {{.V.VName}} == nil {
			{{.V.VName}} = make([]float64, 0, cnt/8)
		}
		sub := 0
		for sub < len(buf) {
			v, cnt := protowire.ConsumeFixed64(buf[sub:])
			if cnt < 1 {
				err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid item value")
				return
			}
			sub += cnt
			{{.V.VName}} = append({{.V.VName}}, math.Float64frombits(v))
		}
	`,
	"decode.slice.string": `
		if typ != protowire.BytesType {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid repeated tag value")
			return
		}
		buf, cnt := protowire.ConsumeBytes({{.V.Buffer}})
		if buf == nil {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid len value")
			return
		}
		{{.V.Index}} += cnt
		if {{.V.VName}} == nil {
			{{.V.VName}} = make([]string, 0, 2)
		}
		{{.V.VName}} = append({{.V.VName}}, string(buf))
	`,
	"decode.slice.bytes": `
		if typ != protowire.BytesType {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid repeated tag value")
			return
		}
		buf, cnt := protowire.ConsumeBytes({{.V.Buffer}})
		if buf == nil {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid len value")
			return
		}
		{{.V.Index}} += cnt
		if {{.V.VName}} == nil {
			{{.V.VName}} = make([][]byte, 0, 2)
		}
		{{.V.VName}} = append({{.V.VName}}, buf)
	`,
	"decode.slice.message": `
		if typ != protowire.BytesType {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid repeated tag value")
			return
		}
		buf, cnt := protowire.ConsumeBytes({{.V.Buffer}})
		if buf == nil {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid len value")
			return
		}
		{{.V.Index}} += cnt
		if {{.V.VName}} == nil {
			{{.V.VName}} = make([]*{{.Field.GoType}}, 0, 2)
		}
		item := &{{.Field.GoType}}{}
		err = item.UnmarshalObject(buf)
		if err != nil {
			return
		}
		{{.V.VName}} = append({{.V.VName}}, item)
	`,

	"size.bool": `
		// {{ TagSize .Field.DescNum }} = protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += 1+{{ TagSize .Field.DescNum }} 
	`,
	"size.varint": `
		// {{ TagSize .Field.DescNum }} = protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += {{ TagSize .Field.DescNum }} + protowire.SizeVarint(uint64({{.V.VName}}))
	`,
	"size.sint": `
		// {{ TagSize .Field.DescNum }} = protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += {{ TagSize .Field.DescNum }} + protowire.SizeVarint(protowire.EncodeZigZag(int64({{.V.VName}})))
	`,
	"size.fix32": `
		// {{ TagSize .Field.DescNum }} = protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += {{ TagSize .Field.DescNum }} + 4
	`,
	"size.float": `
		// {{ TagSize .Field.DescNum }} = protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += {{ TagSize .Field.DescNum }} + 4
	`,
	"size.fix64": `
		// {{ TagSize .Field.DescNum }} = protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += {{ TagSize .Field.DescNum }} + 8
	`,
	"size.double": `
		// {{ TagSize .Field.DescNum }} = protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += {{ TagSize .Field.DescNum }} + 8
	`,
	"size.string": `
		// {{ TagSize .Field.DescNum }} = protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += {{ TagSize .Field.DescNum }} + protowire.SizeBytes(len({{.V.VName}}))
	`,
	"size.bytes": `
		// {{ TagSize .Field.DescNum }} = protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += {{ TagSize .Field.DescNum }} +  protowire.SizeBytes(len({{.V.VName}}))
	`,
	"size.message": `
		// {{ TagSize .Field.DescNum }} = protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += {{ TagSize .Field.DescNum }} +  protowire.SizeBytes({{.V.VName}}.MarshalSize())
	`,

	"size.packed.bool": `
		{{.V.Size}} += {{ TagSize .Field.DescNum }} // {{.V.Size}} += protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += protowire.SizeBytes(len({{.V.VName}}))
	`,
	"size.packed.varint": `
		{{.V.Size}} += {{ TagSize .Field.DescNum }} // {{.V.Size}} += protowire.SizeTag({{.Field.DescNum}})
		if len({{.V.VName}}) > 0 {
			fsize := 0
			for _, item := range {{.V.VName}} {
				fsize += protowire.SizeVarint(uint64(item))
			}
			{{.V.Size}} += protowire.SizeBytes(fsize)
		}
	`,
	"size.packed.sint": `
		{{.V.Size}} += {{ TagSize .Field.DescNum }} // {{.V.Size}} += protowire.SizeTag({{.Field.DescNum}})
		if len({{.V.VName}}) > 0 {
			fsize := 0
			for _, item := range {{.V.VName}} {
				fsize += protowire.SizeVarint(protowire.EncodeZigZag(int64(item)))
			}
			{{.V.Size}} += protowire.SizeBytes(fsize)
		}
	`,
	"size.packed.fix32": `
		{{.V.Size}} += {{ TagSize .Field.DescNum }} // {{.V.Size}} += protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += protowire.SizeBytes(len({{.V.VName}}) * 4)
	`,
	"size.packed.float": `
		{{.V.Size}} += {{ TagSize .Field.DescNum }} // {{.V.Size}} += protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += protowire.SizeBytes(len({{.V.VName}}) * 4)
	`,
	"size.packed.fix64": `
		{{.V.Size}} += {{ TagSize .Field.DescNum }} // {{.V.Size}} += protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += protowire.SizeBytes(len({{.V.VName}}) * 8)
	`,
	"size.packed.double": `
		{{.V.Size}} += {{ TagSize .Field.DescNum }} // {{.V.Size}} += protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += protowire.SizeBytes(len({{.V.VName}}) * 8)
	`,
	"size.packed.string": `
		for _, item := range {{.V.VName}} {
			{{.V.Size}} += {{ TagSize .Field.DescNum }} // {{.V.Size}} += protowire.SizeTag({{.Field.DescNum}})
			{{.V.Size}} += protowire.SizeBytes(len(item))
		}
	`,
	"size.packed.bytes": `
		for _, item := range {{.V.VName}} {
			{{.V.Size}} += {{ TagSize .Field.DescNum }} // {{.V.Size}} += protowire.SizeTag({{.Field.DescNum}})
			{{.V.Size}} += protowire.SizeBytes(len(item))
		}
	`,
	"size.packed.message": `
		for _, item := range {{.V.VName}} {
			{{.V.Size}} += {{ TagSize .Field.DescNum }} // {{.V.Size}} += protowire.SizeTag({{.Field.DescNum}})
			{{.V.Size}} += protowire.SizeBytes(item.MarshalSize())
		}
	`,

	"size.nopack.bool": `
		// {{ TagSize .Field.DescNum }} = protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += ({{ TagSize .Field.DescNum }}+1) * len({{.V.VName}})
	`,
	"size.nopack.varint": `
		// {{ TagSize .Field.DescNum }} = protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += {{ TagSize .Field.DescNum }} * len({{.V.VName}})
		for k:=0; k<len({{.V.VName}}); k++ {
			{{.V.Size}} += protowire.SizeVarint(uint64({{.V.VName}}[k]))
		}
	`,
	"size.nopack.sint": `
		// {{ TagSize .Field.DescNum }} = protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += {{ TagSize .Field.DescNum }} * len({{.V.VName}})
		for k:=0; k<len({{.V.VName}}); k++ {
			{{.V.Size}} += protowire.SizeVarint(protowire.EncodeZigZag(int64({{.V.VName}}[k])))
		}
	`,
	"size.nopack.fix32": `
		// {{ TagSize .Field.DescNum }} = protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += ({{ TagSize .Field.DescNum }}+4) * len({{.V.VName}})
	`,
	"size.nopack.float": `
		// {{ TagSize .Field.DescNum }} = protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += ({{ TagSize .Field.DescNum }}+4) * len({{.V.VName}})
	`,
	"size.nopack.fix64": `
		// {{ TagSize .Field.DescNum }} = protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += ({{ TagSize .Field.DescNum }}+8) * len({{.V.VName}})
	`,
	"size.nopack.double": `
		// {{ TagSize .Field.DescNum }} = protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += ({{ TagSize .Field.DescNum }}+8) * len({{.V.VName}})
	`,
	"size.nopack.string": `
		// {{ TagSize .Field.DescNum }} = protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += {{ TagSize .Field.DescNum }} * len({{.V.VName}})
		for k:=0; k<len({{.V.VName}}); k++ {
			{{.V.Size}} += protowire.SizeBytes(len({{.V.VName}}[k]))
		}
	`,
	"size.nopack.bytes": `
		// {{ TagSize .Field.DescNum }} = protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += {{ TagSize .Field.DescNum }} * len({{.V.VName}})
		for k:=0; k<len({{.V.VName}}); k++ {
			{{.V.Size}} += protowire.SizeBytes(len({{.V.VName}}[k]))
		}
	`,
	"size.nopack.message": `
		// {{ TagSize .Field.DescNum }} = protowire.SizeTag({{.Field.DescNum}})
		{{.V.Size}} += {{ TagSize .Field.DescNum }} * len({{.V.VName}})
		for k:=0; k<len({{.V.VName}}); k++ {
			{{.V.Size}} += protowire.SizeBytes({{.V.VName}}[k].MarshalSize())
		}
	`,

	"encode.map": `
		for mk, mv := range {{ $.V.VName }} {
			// {{.V.Buffer}} = protowire.AppendTag({{.V.Buffer}}, {{ .Field.DescNum }}, protowire.BytesType) => {{ TagBinary .Field.DescNum "protowire.BytesType" }}
			{{.V.Buffer}} = append({{.V.Buffer}}, {{ TagByes .Field.DescNum "protowire.BytesType" }})
			msize := 0 
			{{GenTemplate .Field.MapKey.TemplateSize .Field.MapKey "Size" "msize" "VName" "mk"}}
			{{GenTemplate .Field.MapValue.TemplateSize .Field.MapValue "Size" "msize" "VName" "mv"}}
			{{.V.Buffer}} = protowire.AppendVarint({{.V.Buffer}}, uint64(msize))
			{{GenTemplate .Field.MapKey.TemplateEncode .Field.MapKey "Buffer" .V.Buffer "VName" "mk"}}
			{{GenTemplate .Field.MapValue.TemplateEncode .Field.MapValue "Buffer" .V.Buffer "VName" "mv"}}
		}
	`,
	"decode.map": `
		if typ != protowire.BytesType {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid repeated tag value")
			return
		}
		buf, cnt := protowire.ConsumeBytes({{.V.Buffer}})
		if buf == nil {
			err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid len value")
			return
		}
		{{.V.Index}} += cnt
		if {{.V.VName}}  == nil {
			{{.V.VName}}  = make({{.Field.TypeName}})
		}
		var mk {{.Field.MapKey.TypeName}}
		var mv {{.Field.MapValue.TypeName}}
		for sindex := 0; sindex < len(buf); {
			mi, typ, scnt := protowire.ConsumeTag(buf[sindex:])
			if scnt < 1 {
				err = errors.New("parse {{.Field.Tip}} ID:{{.Field.DescNum}} : invalid varint value")
				return
			}
			_ = typ
			sindex += scnt
			switch mi {
			case 1:
				{{GenTemplate .Field.MapKey.TemplateDecode .Field.MapKey "Buffer" "buf[sindex:]" "VName" "mk" "Index" "sindex"}}
			case 2:
				{{GenTemplate .Field.MapValue.TemplateDecode .Field.MapValue "Buffer" "buf[sindex:]" "VName" "mv" "Index" "sindex"}}
			}
		}
		{{.V.VName}}[mk] = mv
	`,
	"size.map": `
		for mk, mv := range {{.V.VName}} {
			_ = mk
			_ = mv
			// {{ TagSize .Field.DescNum }} = protowire.SizeTag({{.Field.DescNum}})
			{{.V.Size}} += {{ TagSize .Field.DescNum }} 
			msize := 0
			{{GenTemplate .Field.MapKey.TemplateSize .Field.MapKey "Size" "msize" "VName" "mk"}}
			{{GenTemplate .Field.MapValue.TemplateSize .Field.MapValue "Size" "msize" "VName" "mv"}}
			size += protowire.SizeBytes(msize)
		}
	`,
}
