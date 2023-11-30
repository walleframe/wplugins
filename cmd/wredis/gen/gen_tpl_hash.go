package gen

func init() {
	registerTemplate("hash_field_func_arg", `{{ $Name := .Name}} {{$field := .TypeHash.HashDynamic.Field }} {{$value := .TypeHash.HashDynamic.Value }} {{$farg := .TypeHash.HashDynamic.FieldArgs}} {{$varg := .TypeHash.HashDynamic.ValueArgs}} 
{{- if $field -}}
	field {{$field.Type}},
{{- else -}}
	{{- range $i,$arg := $farg -}}
		{{$arg.ArgName}} {{$arg.ArgType }},
	{{- end -}}
{{- end -}}`)
	registerTemplate("hash_value_func_arg", `{{ $Name := .Name}} {{$field := .TypeHash.HashDynamic.Field }} {{$value := .TypeHash.HashDynamic.Value }} {{$farg := .TypeHash.HashDynamic.FieldArgs}} {{$varg := .TypeHash.HashDynamic.ValueArgs}} 
{{- if $value -}}
	value {{$value.Type}},
{{- else -}}
	{{- range $i,$arg := $varg -}}
		{{$arg.ArgName}} {{$arg.ArgType }}, 
	{{- end -}}
{{- end -}}`)
	registerTemplate("hash_filed_str_arg", `{{ $Name := .Name}} {{$field := .TypeHash.HashDynamic.Field }} {{$value := .TypeHash.HashDynamic.Value }} {{$farg := .TypeHash.HashDynamic.FieldArgs}} {{$varg := .TypeHash.HashDynamic.ValueArgs}} 
{{- if $field -}}
	{{- if eq $field.Type "string" -}}
		field
	{{- else -}}
		rdconv.{{$field.RedisFunc}}ToString(field)
	{{- end -}}
{{- else -}}
	Merge{{$Name}}Field({{- range $i,$arg := $farg -}}
		{{$arg.ArgName}},
	{{- end -}})
{{- end -}}`)
	registerTemplate("hash_value_str_arg", `{{ $Name := .Name}} {{$field := .TypeHash.HashDynamic.Field }} {{$value := .TypeHash.HashDynamic.Value }} {{$farg := .TypeHash.HashDynamic.FieldArgs}} {{$varg := .TypeHash.HashDynamic.ValueArgs}} 
{{- if $value -}}
	{{- if eq $value.Type "string" -}}
		value
	{{- else -}}
		rdconv.{{$value.RedisFunc}}ToString(value)
	{{- end -}}
{{- else -}}
	Merge{{$Name}}Value({{- range $i,$arg := $varg -}}
		{{$arg.ArgName}},
	{{- end -}})
{{- end -}}
`)
	registerTemplate("hash_object", `{{ $Name := .Name}}
func (x *x{{$Name}}) Set{{Title .TypeHash.HashObject.Name}}(ctx context.Context, obj *{{.TypeHash.HashObject.Type}}) (err error) {
	n, err := x.rds.HSet(ctx, x.key, {{ range $i,$field := .TypeHash.HashObject.Fields }}
		_{{$Name}}_{{Title $field.Name}}, rdconv.{{$field.RedisFunc}}ToString(obj.{{Title $field.Name}}),{{end}}
	).Result()
	if err != nil {
		return err
	}
	if n != {{len .TypeHash.HashObject.Fields}} {
		return errors.New("set {{Title .TypeHash.HashObject.Name}} failed")
	}
	return
}

{{if .TypeHash.HashObject.HGetAll }}
func (x *x{{$Name}}) Get{{Title .TypeHash.HashObject.Name}}(ctx context.Context) (*{{.TypeHash.HashObject.Type}}, error) {
	ret, err := x.rds.HGetAll(ctx, x.key).Result()
	if err != nil {
		return nil, err
	}
	obj := &{{.TypeHash.HashObject.Type}}{}
{{ range $i,$field := .TypeHash.HashObject.Fields }}
	if val, ok := ret[_{{$Name}}_{{Title $field.Name}}]; ok {
{{if eq $field.Type "string"}}
		obj.{{Title $field.Name}} = val
{{ else -}}
		obj.{{Title $field.Name}}, err = rdconv.StringTo{{$field.RedisFunc}}(val)
		if err != nil {
			return nil, fmt.Errorf("parse {{$Name}}.{{Title $field.Name}} failed,%w", err)
		}
{{ end -}}
	}
{{end}}
	return obj, nil
}
{{end}}

func (x *x{{$Name}}) MGet{{Title .TypeHash.HashObject.Name}}(ctx context.Context) (*{{.TypeHash.HashObject.Type}}, error) {
	ret, err := x.rds.HMGet(ctx, x.key, {{ range $i,$field := .TypeHash.HashObject.Fields }} _{{$Name}}_{{Title $field.Name}},{{end}}).Result()
	if err != nil {
		return nil, err
	}
	obj := &{{.TypeHash.HashObject.Type}}{}
{{ range $i,$field := .TypeHash.HashObject.Fields }}
	if len(ret) > {{$i}} && ret[{{$i}}] != nil {
		obj.{{Title $field.Name}}, err = rdconv.AnyTo{{$field.RedisFunc}}(ret[ {{$i}} ])
		if err != nil {
			return nil, fmt.Errorf("parse {{$Name}}.{{Title $field.Name}} failed,%w", err)
		}
	}
{{end}}
	return obj, nil
}

{{range $i,$field := .TypeHash.HashObject.Fields}}
func (x *x{{$Name}}) Get{{Title $field.Name}}(ctx context.Context)(_ {{$field.Type}},err error) {
{{if eq $field.Type "string" -}}
	return x.rds.HGet(ctx, x.key, _{{$Name}}_{{Title $field.Name}}).Result()
{{ else -}}
	val, err := x.rds.HGet(ctx, x.key, _{{$Name}}_{{Title $field.Name}}).Result()
	if err != nil {
		return
	}
	return rdconv.StringTo{{$field.RedisFunc}}(val)
{{end}}
}
func (x *x{{$Name}}) Set{{Title $field.Name}}(ctx context.Context, val {{$field.Type}}) (err error) {
	n, err := x.rds.HSet(ctx, x.key,  _{{$Name}}_{{Title $field.Name}}, rdconv.{{$field.RedisFunc}}ToString(val)).Result()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("set {{$Name}}.{{Title $field.Name}} failed")
	}
	return nil
}
{{if $field.Number}}
func (x *x{{$Name}}) IncrBy{{Title $field.Name}}(ctx context.Context, incr int) ({{$field.Type}}, error) {
	num, err := x.rds.HIncrBy(ctx, x.key,  _{{$Name}}_{{Title $field.Name}}, int64(incr)).Result()
	if err != nil {
		return 0, err
	}
	return {{$field.Type}}(num), nil
}
{{end}}
{{end}}{{Import "github.com/walleframe/walle/util/rdconv" "Float64ToString"}}
`)
	registerTemplate("hash_dynamic_match", `{{ $Name := .Name}} {{$field := .TypeHash.HashDynamic.Field }} {{$value := .TypeHash.HashDynamic.Value }} {{$farg := .TypeHash.HashDynamic.FieldArgs}} {{$varg := .TypeHash.HashDynamic.ValueArgs}} 

{{if $farg }}
func Merge{{$Name}}Field({{- range $i,$arg := $farg -}}{{$arg.ArgName}} {{$arg.ArgType }}, {{- end -}}) string {
	buf := util.Builder{}
	buf.Grow({{.KeySize}})
{{- range $i,$arg := $farg -}}
{{- if gt $i 0 -}}
	buf.WriteByte(':')
{{- end }}
	{{$arg.FormatCode "buf"}}
{{end -}}
	return buf.String()
}
func Split{{$Name}}Field(val string)({{- range $i,$arg := $farg -}}{{$arg.ArgName}} {{$arg.ArgType }}, {{- end -}}err error) {
	items := strings.Split(val, ":")
	if len(items) != {{len $farg}} {
		err = errors.New("invalid {{$Name}} field value")
		return
	}
{{ range $i,$arg := $farg -}}
{{- if eq $arg.ArgType "string" -}} 
	{{- $arg.ArgName}} = items[{{$i}}]
{{ else -}}
	{{$arg.ArgName}}, err = rdconv.StringTo{{Title $arg.ArgType}}(items[{{$i}}])
	if err != nil {
		return
	}
{{ end -}}
{{end -}}
	return
}
{{end}}

{{if $varg }}
func Merge{{$Name}}Value({{- range $i,$arg := $varg -}}{{$arg.ArgName}} {{$arg.ArgType}},{{- end -}}) string {
	buf := util.Builder{}
	buf.Grow({{.KeySize}})
{{- range $i,$arg := $varg -}}
{{- if gt $i 0 -}}
	buf.WriteByte(':')
{{- end }}
	{{$arg.FormatCode "buf"}}
{{end -}}
	return buf.String()
}
func Split{{$Name}}Value(val string)({{- range $i,$arg := $varg -}}{{$arg.ArgName}} {{$arg.ArgType }}, {{- end -}}err error) {
	items := strings.Split(val, ":")
	if len(items) != {{len $varg}} {
		err = errors.New("invalid {{$Name}} field value")
		return
	}
{{ range $i,$arg := $varg -}}
{{- if eq $arg.ArgType "string" -}} 
	{{$arg.ArgName}} = items[{{$i}}]
{{ else -}}
	{{$arg.ArgName}}, err = rdconv.StringTo{{Title $arg.ArgType}}(items[{{$i}}])
	if err != nil {
		return
	}
{{ end -}}
{{end -}}
	return
}
{{end}}

func (x *x{{.Name}}) GetField(ctx context.Context, {{GenTypeTemplate "hash_field_func_arg" .}}) ({{GenTypeTemplate "hash_value_func_arg" .}} err error) {
{{- if $value }}
{{- if eq $value.Type "string"}}
	return x.rds.HGet(ctx, x.key, {{GenTypeTemplate "hash_filed_str_arg" .}}).Result()
{{- else }}
	v, err := x.rds.HGet(ctx, x.key, {{GenTypeTemplate "hash_filed_str_arg" .}}).Result()
	if err != nil {
		return
	}
	return rdconv.StringTo{{$value.RedisFunc}}(v)
{{- end}}
{{- else }}
	v, err := x.rds.HGet(ctx, x.key, {{GenTypeTemplate "hash_filed_str_arg" .}}).Result()
	if err != nil {
		return
	}
	return Split{{$Name}}Value(v)
{{ end -}}
}
func (x *x{{.Name}}) SetField(ctx context.Context, {{GenTypeTemplate "hash_field_func_arg" .}} {{GenTypeTemplate "hash_value_func_arg" .}}) (err error) {
	num, err := x.rds.HSet(ctx, x.key, {{GenTypeTemplate "hash_filed_str_arg" .}}, {{GenTypeTemplate "hash_value_str_arg" .}}).Result()
	if err != nil {
		return err
	}
	if num != 1 {
		return errors.New("set field failed")
	}
	return nil
}

{{if $field }}
func (x *x{{.Name}}) HKeys(ctx context.Context) (vals []{{$field.Type}}, err error) {
{{if eq $field.Type "string" -}}
	return x.rds.HKeys(ctx, x.key).Result()
{{- else -}}
	ret, err := x.rds.HKeys(ctx, x.key).Result()
	if err != nil {
		return
	}
	for _, v := range ret {
		key,err := rdconv.StringTo{{$field.RedisFunc}}(v)
		if err != nil {
			return nil, err 
		}
		vals = append(vals, key)
	}
	return
{{- end}}
}


func (x *x{{.Name}}) HKeysRange(ctx context.Context, filter func({{$field.Type}}) bool)(error) {
	ret, err := x.rds.HKeys(ctx, x.key).Result()
	if err != nil {
		return err
	}
	for _, v := range ret {
{{if eq $field.Type "string" -}}
		if !filter(v) {
			return nil
		}
{{ else -}}
		key,err := rdconv.StringTo{{$field.RedisFunc}}(v)
		if err != nil {
			return err
		}
		if !filter(key) {
			return nil
		}
{{ end -}}
	}
	return nil
}

{{else}}

func (x *x{{.Name}}) HKeysRange(ctx context.Context, filter func({{GenTypeTemplate "hash_field_func_arg" .}})bool) (err error) {
	ret, err := x.rds.HKeys(ctx, x.key).Result()
	if err != nil {
		return
	}
	for _, v := range ret {
		{{ range $i,$arg := $farg -}} {{$arg.ArgName}}, {{end}}err := Split{{$Name}}Field(v)
		if err != nil {
			return err
		}
		if !filter({{ range $i,$arg := $farg -}} {{$arg.ArgName}}, {{end}}) {
			return nil 
		}
	}
	return nil
}

{{end}}

{{if $value}}
func (x *x{{.Name}}) HVals(ctx context.Context) (vals []{{$value.Type}}, err error) {
{{- if eq $value.Type "string"}}
	return x.rds.HVals(ctx, x.key).Result()
{{- else }}
	ret, err := x.rds.HVals(ctx, x.key).Result()
	if err != nil {
		return
	}
	for _, v := range ret {
		val,err := rdconv.StringTo{{$value.RedisFunc}}(v)
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}
	return
{{- end}}
}

func (x *x{{.Name}}) HValsRange(ctx context.Context, filter func({{$value.Type}}) bool)(error) {
	ret, err := x.rds.HVals(ctx, x.key).Result()
	if err != nil {
		return err
	}
	for _, v := range ret {
{{if eq $value.Type "string" }}
		if !filter(v) {
			return nil
		}
{{ else -}}
		val,err := rdconv.StringTo{{$value.RedisFunc}}(v)
		if err != nil {
			return err
		}
		if !filter(val) {
			return nil
		}
{{ end -}}
	}
	return nil
}
{{else}}
func (x *x{{.Name}}) HValsRange(ctx context.Context, filter func({{GenTypeTemplate "hash_value_func_arg" .}})bool)(err error) {
	ret, err := x.rds.HVals(ctx, x.key).Result()
	if err != nil {
		return
	}
	for _, v := range ret {
		{{ range $i,$arg := $varg -}} {{$arg.ArgName}}, {{end}}err := Split{{$Name}}Value(v)
		if err != nil {
			return err
		}
		if !filter({{ range $i,$arg := $varg -}} {{$arg.ArgName}}, {{end}}) {
			return nil 
		}
	}
	return nil
}

{{end}}

func (x *x{{.Name}}) HExists(ctx context.Context, {{GenTypeTemplate "hash_field_func_arg" .}}) (bool, error) {
	return x.rds.HExists(ctx, x.key, {{GenTypeTemplate "hash_filed_str_arg" .}}).Result()
}

func (x *x{{.Name}}) HDel(ctx context.Context, {{GenTypeTemplate "hash_field_func_arg" .}}) (bool, error) {
	n, err := x.rds.HDel(ctx, x.key, {{GenTypeTemplate "hash_filed_str_arg" .}}).Result()
	if err != nil {
		return false, err
	}
	return n == 1, nil
}

func (x *x{{.Name}}) HLen(ctx context.Context) (count int64, err error) {
	return x.rds.HLen(ctx, x.key).Result()
}
{{if $field }}
	func (x *x{{.Name}}) HRandField(ctx context.Context, count int) (vals []{{$field.Type}}, err error) {
	{{if eq $field.Type "string" -}}
		return x.rds.HRandField(ctx, x.key, count).Result()
	{{- else -}}
		ret, err := x.rds.HRandField(ctx, x.key, count).Result()
		if err != nil {
			return
		}
		for _, v := range ret {
			key, err := rdconv.StringTo{{$field.RedisFunc}}(v)
			if err != nil {
				return nil, err
			}
			vals = append(vals, key)
		}
		return
	{{- end}}
	}
{{else}}
	func (x *x{{.Name}}) HRandFieldRange(ctx context.Context, count int,filter func({{GenTypeTemplate "hash_field_func_arg" .}})bool) (err error) {
		ret, err := x.rds.HRandField(ctx, x.key, count).Result()
		if err != nil {
			return
		}
		for _, v := range ret {
			{{ range $i,$arg := $farg -}} {{$arg.ArgName}}, {{end}}err := Split{{$Name}}Field(v)
			if err != nil {
				return err
			}
			if !filter({{ range $i,$arg := $farg -}} {{$arg.ArgName}}, {{end}}) {
				return nil 
			}
		}
		return
	}
{{end}}

{{if .TypeHash.HashDynamic.GenMap }}
func (x *x{{.Name}}) HRandFieldWithValues(ctx context.Context, count int) (vals map[{{$field.Type}}]{{$value.Type}}, err error) {
	ret, err := x.rds.HRandFieldWithValues(ctx, x.key, count).Result()
	if err != nil {
		return
	}
	vals = make(map[{{$field.Type}}]{{$value.Type}}, len(ret))
	for _, v := range ret {
{{if eq $field.Type "string" -}}
		key := v.Key
{{ else -}}
		key, err := rdconv.StringTo{{$field.RedisFunc}}(v.Key)
		if err != nil {
			return nil, err
		}
{{end -}}
{{ if eq $value.Type "string" -}}
		val := v.Value
{{else -}}
		val, err := rdconv.StringTo{{$value.RedisFunc}}(v.Value)
		if err != nil {
			return nil, err
		}
{{end -}}
		vals[key] = val
	}
	return
}
{{end}}

func (x *x{{.Name}}) HRandFieldWithValuesRange(ctx context.Context, count int, filter func({{GenTypeTemplate "hash_field_func_arg" .}} {{GenTypeTemplate "hash_value_func_arg" .}})bool) (err error) {
	ret, err := x.rds.HRandFieldWithValues(ctx, x.key, count).Result()
	if err != nil {
		return
	}
	for _, v := range ret {
{{if $field}}
{{if eq $field.Type "string" -}}
		key := v.Key
{{ else -}}
		key, err := rdconv.StringTo{{$field.RedisFunc}}(v.Key)
		if err != nil {
			return  err
		}
{{end -}}
{{else}}
	{{ range $i,$arg := $farg -}} {{$arg.ArgName}}, {{end}}err := Split{{$Name}}Field(v.Key)
	if err != nil {
		return err
	}
{{end}}

{{if $value}}
{{if eq $value.Type "string" -}}
		val := v.Value
{{else -}}
		val, err := rdconv.StringTo{{$value.RedisFunc}}(v.Value)
		if err != nil {
			return err
		}
{{end -}}
{{else}}
		{{ range $i,$arg := $varg -}} {{$arg.ArgName}}, {{end}}err := Split{{$Name}}Value(v.Value)
		if err != nil {
			return err
		}
{{end}}

		if !filter({{if $field}} key, {{else}} {{ range $i,$arg := $farg -}} {{$arg.ArgName}}, {{end}}{{end -}}
			{{- if $value}} val {{else}} {{ range $i,$arg := $varg -}} {{$arg.ArgName}}, {{end}} {{end}}) {
			return nil 
		}

	}
	return
}
{{if .TypeHash.HashDynamic.GenMap }}
func (x *x{{.Name}}) HScan(ctx context.Context, match string, count int) (vals map[{{$field.Type}}]{{$value.Type}}, err error) {
	cursor := uint64(0)
	vals = make(map[{{$field.Type}}]{{$value.Type}})
	var kvs []string
	for {
		kvs, cursor, err = x.rds.HScan(ctx, x.key, cursor, match, int64(count)).Result()
		if err != nil {
			return nil, err
		}
		for k := 0; k < len(kvs); k += 2 {
{{if eq $field.Type "string" -}}
			key := kvs[k]
{{ else -}}
			key, err := rdconv.StringTo{{$field.RedisFunc}}(kvs[k])
			if err != nil {
				return nil, err
			}
{{end -}}
{{if eq $value.Type "string" -}}
			val := kvs[k+1]
{{else -}}
			val, err := rdconv.StringTo{{$value.RedisFunc}}(kvs[k+1])
			if err != nil {
				return nil, err
			}
{{end -}}
			vals[key] = val
		}
		if cursor == 0 {
			break
		}
	}

	return
}
{{end}}
func (x *x{{.Name}}) HScanRange(ctx context.Context, match string, count int, filter func({{GenTypeTemplate "hash_field_func_arg" .}} {{GenTypeTemplate "hash_value_func_arg" .}})bool) (err error) {
	cursor := uint64(0)
	var kvs []string
	for {
		kvs, cursor, err = x.rds.HScan(ctx, x.key, cursor, match, int64(count)).Result()
		if err != nil {
			return err
		}
		for k := 0; k < len(kvs); k += 2 {
{{if $field}}
{{if eq $field.Type "string" -}}
			key := kvs[k]
{{ else -}}
			key, err := rdconv.StringTo{{$field.RedisFunc}}(kvs[k])
			if err != nil {
				return err
			}
{{end -}}
{{else -}}
	{{ range $i,$arg := $farg -}} {{$arg.ArgName}}, {{end}}err := Split{{$Name}}Field(kvs[k])
	if err != nil {
		return err
	}
{{- end}}
{{if $value}}
{{if eq $value.Type "string" -}}
		val := kvs[k+1]
{{else -}}
			val, err := rdconv.StringTo{{$value.RedisFunc}}(kvs[k+1])
			if err != nil {
				return err
			}
{{end -}}
{{else}}
			{{ range $i,$arg := $varg -}} {{$arg.ArgName}}, {{end}}err := Split{{$Name}}Value(kvs[k+1])
			if err != nil {
				return err
			}
{{end}}
			if !filter({{if $field}} key, {{else}} {{ range $i,$arg := $farg -}} {{$arg.ArgName}}, {{end}}{{end -}}
				{{- if $value}} val {{else}} {{ range $i,$arg := $varg -}} {{$arg.ArgName}}, {{end}} {{end}}) {
				return nil
			}
		}
		if cursor == 0 {
			break
		}
	}
	return
}
`)
}
