package gen

func init() {
	registerTemplate("script_return_1", `{{$Name := .Obj.Name}} {{ $Script := .Script}} {{ $count := len $Script.Output}} {{ $ret := index .Script.Output 0}}
var x{{$Name}}{{Title $Script.Name}}Script = svc_redis.NewScript("{{.Script.Script}}")


func (x *x{{$Name}}) {{Title $Script.Name}}(ctx context.Context, {{range $i,$arg := $Script.Args}}{{$arg.ArgName}} {{$arg.ArgType }}, {{- end}}) (_ {{$ret.ArgType}}, err error) {
	cmd := redis.{{$Script.CommandName}}(ctx, "evalsha", x{{$Name}}{{Title $Script.Name}}Script.Hash, "1", x.key, {{range $i,$arg := $Script.Args}} rdconv.{{Title $arg.ArgType}}ToString({{$arg.ArgName}}), {{end}})
	if err != nil {
		if !redis.HasErrorPrefix(err, "NOSCRIPT") {
			return
		}
		cmd = redis.{{$Script.CommandName}}(ctx, "eval", x{{$Name}}{{Title $Script.Name}}Script.Script, "1", x.key, {{range $i,$arg := $Script.Args}} rdconv.{{Title $arg.ArgType}}ToString({{$arg.ArgName}}), {{end}})
		err = x.rds.Process(ctx, cmd)
		if err != nil {
			return
		}
	}
	return {{$ret.ArgType}}(cmd.Val()), nil
}

`)

	registerTemplate("script_return_mul", `{{$Name := .Obj.Name}} {{ $Script := .Script}} {{ $ret := index .Script.Output 0}}
var x{{$Name}}{{Title $Script.Name}}Script = svc_redis.NewScript("{{.Script.Script}}")

func (x *x{{$Name}}) {{Title $Script.Name}}(ctx context.Context, {{range $i,$arg := $Script.Args}}{{$arg.ArgName}} {{$arg.ArgType }}, {{- end}}) ({{range $i,$arg := $Script.Output}}{{$arg.ArgName}} {{$arg.ArgType }}, {{- end}} err error) {
	cmd := redis.NewStringSliceCmd(ctx, "evalsha", x{{$Name}}{{Title $Script.Name}}Script.Hash, "1", x.key, {{range $i,$arg := $Script.Args}} rdconv.{{Title $arg.ArgType}}ToString({{$arg.ArgName}}), {{end}})
	if err != nil {
		if !redis.HasErrorPrefix(err, "NOSCRIPT") {
			return
		}
		cmd = redis.NewStringSliceCmd(ctx, "eval", x{{$Name}}{{Title $Script.Name}}Script.Script, "1", x.key, {{range $i,$arg := $Script.Args}} rdconv.{{Title $arg.ArgType}}ToString({{$arg.ArgName}}), {{end}})
		err = x.rds.Process(ctx, cmd)
		if err != nil {
			return
		}
	}
	vals := cmd.Val()
	if len(vals) != {{len $Script.Output}} {
		err = errors.New("{{$Script.Name}} return value count not equal {{len $Script.Output}}")
		return
	}
	{{range $i,$arg := $Script.Output -}}
{{if eq $arg.ArgType "string" -}}
	{{$arg.ArgName}} = vals[{{$i}}]
{{ else -}}
	{{$arg.ArgName}}, err = rdconv.StringTo{{Title $arg.ArgType }}(vals[{{$i}}]) 
	if err != nil {
		return
	}
{{end -}}
{{ end -}}
	return
}
`)
}
