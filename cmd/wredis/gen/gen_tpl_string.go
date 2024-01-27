package gen

func init() {
	registerTemplate("string_number", `
func (x *x{{.Name}}) Incr(ctx context.Context) ({{.TypeString.Type}}, error) {
	n,err := x.rds.Incr(ctx, x.key).Result()
	return {{.TypeString.Type}}(n), err
}

func (x *x{{.Name}}) IncrBy(ctx context.Context, val int) (_ {{.TypeString.Type}},err error) {
	cmd := redis.NewIntCmd(ctx, "incrby", x.key, strconv.FormatInt(int64(val), 10))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	return {{.TypeString.Type}}(cmd.Val()), nil
}

func (x *x{{.Name}}) Decr(ctx context.Context) ({{.TypeString.Type}}, error) {
	n,err := x.rds.Decr(ctx, x.key).Result()
	return {{.TypeString.Type}}(n), err
}

func (x *x{{.Name}}) DecrBy(ctx context.Context, val int) (_ {{.TypeString.Type}}, err error) {
	cmd := redis.NewIntCmd(ctx, "decrby", x.key, strconv.FormatInt(int64(val), 10))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	return {{.TypeString.Type}}(cmd.Val()), nil
}

func (x *x{{.Name}}) Get(ctx context.Context) ({{.TypeString.Type}}, error) {
	data,err := x.rds.Get(ctx, x.key).Result()
	if err != nil {
		return 0,err
	}
	val,err := strconv.Parse{{if .TypeString.Signed}}Int{{else}}Uint{{end}}(data, 10, 64)
	if err != nil {
		return 0,err
	}
	return {{.TypeString.Type}}(val), nil
}

func (x *x{{.Name}}) Set(ctx context.Context, val {{.TypeString.Type}}, expire time.Duration) error {
	return x.rds.Set(ctx, x.key, strconv.Format{{if .TypeString.Signed}}Int(int64{{else}}Uint(uint64{{end}}(val), 10), expire).Err()
}

func (x *x{{.Name}}) SetNX(ctx context.Context, val {{.TypeString.Type}}, expire time.Duration) (bool, error) {
	return x.rds.SetNX(ctx, x.key, strconv.Format{{if .TypeString.Signed}}Int(int64{{else}}Uint(uint64{{end}}(val), 10), expire).Result()
}

func (x *x{{.Name}}) SetEx(ctx context.Context, val {{.TypeString.Type}}, expire time.Duration) error {
	return x.rds.SetEx(ctx, x.key, strconv.Format{{if .TypeString.Signed}}Int(int64{{else}}Uint(uint64{{end}}(val), 10), expire).Err()
}
`)
	registerTemplate("string_float", `
func (x *x{{.Name}}) Get(ctx context.Context) ({{.TypeString.Type}}, error) {
	data,err := x.rds.Get(ctx, x.key).Result()
	if err != nil {
		return 0,err
	}
	val,err := strconv.ParseFloat(data, 64)
	if err != nil {
		return 0,err
	}
	return {{.TypeString.Type}}(val), nil
}

func (x *x{{.Name}}) IncrBy(ctx context.Context, val int) (_ {{.TypeString.Type}},err error) {
	cmd := redis.NewFloatCmd(ctx, "incrbyfloat", x.key, strconv.FormatInt(int64(val), 10))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	return {{.TypeString.Type}}(cmd.Val()), nil
}

func (x *x{{.Name}}) Set(ctx context.Context, val {{.TypeString.Type}}, expire time.Duration) error {
	return x.rds.Set(ctx, x.key, rdconv.Float64ToString(float64(val)), expire).Err()
}

func (x *x{{.Name}}) SetNX(ctx context.Context, val {{.TypeString.Type}}, expire time.Duration) (bool, error) {
	return x.rds.SetNX(ctx, x.key, rdconv.Float64ToString(float64(val)), expire).Result()
}

func (x *x{{.Name}}) SetEx(ctx context.Context, val {{.TypeString.Type}}, expire time.Duration) error {
	return x.rds.SetEx(ctx, x.key, rdconv.Float64ToString(float64(val)), expire).Err()
} {{Import "github.com/walleframe/walle/util/rdconv" "Float64ToString"}}
`)
	registerTemplate("string_string", `
func (x *x{{.Name}}) GetRange(ctx context.Context, start, end int64) (_ string, err error) {
	cmd := redis.NewStringCmd(ctx, "getrange", x.key, strconv.FormatInt(start, 10), strconv.FormatInt(end, 10))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	return cmd.Val(), nil
}

func (x *x{{.Name}}) SetRange(ctx context.Context, offset int64, value string) (_ int64, err error) {
		cmd := redis.NewIntCmd(ctx, "setrange", x.key, strconv.FormatInt(offset, 10), value)
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	return cmd.Val(), nil
}

func (x *x{{.Name}}) Append(ctx context.Context, val string) (int64, error) {
	return x.rds.Append(ctx, x.key, val).Result()
}

func (x *x{{.Name}}) StrLen(ctx context.Context) (int64, error) {
	return x.rds.StrLen(ctx, x.key).Result()
}

func (x *x{{.Name}}) Get(ctx context.Context) (string, error) {
	return x.rds.Get(ctx, x.key).Result()
}

func (x *x{{.Name}}) Set(ctx context.Context, data string, expire time.Duration) error {
	return x.rds.Set(ctx, x.key, data, expire).Err()
}

func (x *x{{.Name}}) SetNX(ctx context.Context, data string, expire time.Duration) (bool, error) {
	return x.rds.SetNX(ctx, x.key, data, expire).Result()
}

func (x *x{{.Name}}) SetEx(ctx context.Context, data string, expire time.Duration) error {
	return x.rds.SetEx(ctx, x.key, data, expire).Err()
}
`)
	registerTemplate("string_protobuf", `
func (x *x{{.Name}}) Set(ctx context.Context, pb proto.Message, expire time.Duration) error {
	data, err := proto.Marshal(pb)
	if err != nil {
		return err
	}
	return x.rds.Set(ctx, x.key, util.BytesToString(data), expire).Err()
}

func (x *x{{.Name}}) SetNX(ctx context.Context, pb proto.Message, expire time.Duration) error {
	data, err := proto.Marshal(pb)
	if err != nil {
		return err
	}
	return x.rds.SetNX(ctx, x.key, util.BytesToString(data), expire).Err()
}

func (x *x{{.Name}}) SetEx(ctx context.Context, pb proto.Message, expire time.Duration) error {
	data, err := proto.Marshal(pb)
	if err != nil {
		return err
	}
	return x.rds.SetEx(ctx, x.key, util.BytesToString(data), expire).Err()
}

func (x *x{{.Name}}) Get(ctx context.Context, pb proto.Message) error {
	data, err := x.rds.Get(ctx, x.key).Result()
	if err != nil {
		return err
	}
	err = proto.Unmarshal(util.StringToBytes(data), pb)
	if err != nil {
		return err
	}
	return nil
} {{Import "github.com/walleframe/walle/util" "StringToBytes"}}
`)
	registerTemplate("string_walle", `
func (x *x{{.Name}}) Set(ctx context.Context, pb {{.WPbPkg}}.Message, expire time.Duration) error {
	data, err := pb.MarshalObject()
	if err != nil {
		return err
	}
	return x.rds.Set(ctx, x.key, util.BytesToString(data), expire).Err()
}

func (x *x{{.Name}}) SetNX(ctx context.Context, pb {{.WPbPkg}}.Message, expire time.Duration) error {
	data, err := pb.MarshalObject()
	if err != nil {
		return err
	}
	return x.rds.SetNX(ctx, x.key, util.BytesToString(data), expire).Err()
}

func (x *x{{.Name}}) SetEx(ctx context.Context, pb {{.WPbPkg}}.Message, expire time.Duration) error {
	data, err := pb.MarshalObject()
	if err != nil {
		return err
	}
	return x.rds.SetEx(ctx, x.key, util.BytesToString(data), expire).Err()
}

func (x *x{{.Name}}) Get(ctx context.Context, pb {{.WPbPkg}}.Message) error {
	data, err := x.rds.Get(ctx, x.key).Result()
	if err != nil {
		return err
	}
	err = pb.UnmarshalObject(util.StringToBytes(data))
	if err != nil {
		return err
	}
	return nil
} {{Import "github.com/walleframe/walle/util" "StringToBytes"}}
`)
	registerTemplate("string_custom_walle", `
func (x *x{{.Name}}) Set(ctx context.Context, pb *{{.TypeString.Type}}, expire time.Duration) error {
	data, err := pb.MarshalObject()
	if err != nil {
		return err
	}
	return x.rds.Set(ctx, x.key, util.BytesToString(data), expire).Err()
}

func (x *x{{.Name}}) SetNX(ctx context.Context, pb *{{.TypeString.Type}}, expire time.Duration) error {
	data, err := pb.MarshalObject()
	if err != nil {
		return err
	}
	return x.rds.SetNX(ctx, x.key, util.BytesToString(data), expire).Err()
}

func (x *x{{.Name}}) SetEx(ctx context.Context, pb *{{.TypeString.Type}}, expire time.Duration) error {
	data, err := pb.MarshalObject()
	if err != nil {
		return err
	}
	return x.rds.SetEx(ctx, x.key, util.BytesToString(data), expire).Err()
}

func (x *x{{.Name}}) Get(ctx context.Context, pb *{{.TypeString.Type}}) error {
	data, err := x.rds.Get(ctx, x.key).Result()
	if err != nil {
		return err
	}
	err = pb.UnmarshalObject(util.StringToBytes(data))
	if err != nil {
		return err
	}
	return nil
} {{Import "github.com/walleframe/walle/util" "StringToBytes"}}
`)
	registerTemplate("string_custom_protobuf", `
func (x *x{{.Name}}) Set(ctx context.Context, pb *{{.TypeString.Type}}, expire time.Duration) error {
	data, err := proto.Marshal(pb)
	if err != nil {
		return err
	}
	return x.rds.Set(ctx, x.key, util.BytesToString(data), expire).Err()
}

func (x *x{{.Name}}) SetNX(ctx context.Context, pb *{{.TypeString.Type}}, expire time.Duration) error {
	data, err := proto.Marshal(pb)
	if err != nil {
		return err
	}
	return x.rds.SetNX(ctx, x.key, util.BytesToString(data), expire).Err()
}

func (x *x{{.Name}}) SetEx(ctx context.Context, pb *{{.TypeString.Type}}, expire time.Duration) error {
	data, err := proto.Marshal(pb)
	if err != nil {
		return err
	}
	return x.rds.SetEx(ctx, x.key, util.BytesToString(data), expire).Err()
}

func (x *x{{.Name}}) Get(ctx context.Context, pb *{{.TypeString.Type}}) error {
	data, err := x.rds.Get(ctx, x.key).Result()
	if err != nil {
		return err
	}
	err = proto.Unmarshal(util.StringToBytes(data), pb)
	if err != nil {
		return err
	}
	return nil
} {{Import "github.com/walleframe/walle/util" "StringToBytes"}} 
`)
}
