package gen

func init() {
	registerTemplate("zset_func_arg", `{{ $Name := .Name}} {{$mem := .TypeZSet.Member }} {{$farg := .TypeZSet.Args}}
	{{- if $mem -}}
		mem {{$mem.Type}},
	{{- else if $farg -}}
		{{- range $i,$arg := $farg -}}
			{{$arg.ArgName}} {{$arg.ArgType}},
		{{- end -}}
	{{- end -}}
	`)
	registerTemplate("zset_str_arg", `{{ $Name := .Name}} {{$mem := .TypeZSet.Member }} {{$farg := .TypeZSet.Args}}
	{{- if $mem -}}
		{{- if eq $mem.Type "string" -}}
			mem
		{{- else -}}
			rdconv.{{$mem.RedisFunc}}ToString(mem)
		{{- end -}}
	{{- else if $farg -}}
		Merge{{$Name}}Member({{- range $i,$arg := $farg -}}
			{{$arg.ArgName}},
		{{- end -}})
	{{- end -}}
	`)
	registerTemplate("zset_basic", `{{ $Name := .Name}} {{$mem := .TypeZSet.Member }} {{$score := .TypeZSet.Score }} {{$farg := .TypeZSet.Args}}
{{if $farg }}
func Merge{{$Name}}Member({{- range $i,$arg := $farg -}}{{$arg.ArgName}} {{$arg.ArgType }}, {{- end -}}) string {
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
func Split{{$Name}}Member(val string)({{- range $i,$arg := $farg -}}{{$arg.ArgName}} {{$arg.ArgType }}, {{- end -}}err error) {
	items := strings.Split(val, ":")
	if len(items) != {{len $farg}} {
		err = errors.New("invalid {{$Name}} mem value")
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
	
func (x *x{{$Name}}) ZCard(ctx context.Context) (int64, error) {
	cmd := redis.NewIntCmd(ctx, "zcard", x.key)
	x.rds.Process(ctx, cmd)
	return cmd.Result()
}

func (x *x{{$Name}}) ZAdd(ctx context.Context, {{GenTypeTemplate "zset_func_arg" .}} score {{$score.Type}}) error {
	cmd := redis.NewIntCmd(ctx, "zadd", x.key, {{GenTypeTemplate "zset_str_arg" .}}, rdconv.{{$score.RedisFunc}}ToString(score))
	return x.rds.Process(ctx, cmd)
}

func (x *x{{$Name}}) ZAddNX(ctx context.Context, {{GenTypeTemplate "zset_func_arg" .}} score {{$score.Type}}) error {
	cmd := redis.NewIntCmd(ctx, "zadd", x.key, "nx", {{GenTypeTemplate "zset_str_arg" .}}, rdconv.{{$score.RedisFunc}}ToString(score))
	return x.rds.Process(ctx, cmd)
}

func (x *x{{$Name}}) ZAddXX(ctx context.Context, {{GenTypeTemplate "zset_func_arg" .}} score {{$score.Type}}) error {
	cmd := redis.NewIntCmd(ctx, "zadd", x.key, "xx", {{GenTypeTemplate "zset_str_arg" .}}, rdconv.{{$score.RedisFunc}}ToString(score))
	return x.rds.Process(ctx, cmd)
}

func (x *x{{$Name}}) ZAddLT(ctx context.Context, {{GenTypeTemplate "zset_func_arg" .}} score {{$score.Type}}) error {
	cmd := redis.NewIntCmd(ctx, "zadd", x.key, "lt", {{GenTypeTemplate "zset_str_arg" .}}, rdconv.{{$score.RedisFunc}}ToString(score))
	return x.rds.Process(ctx, cmd)
}

func (x *x{{$Name}}) ZAddGT(ctx context.Context, {{GenTypeTemplate "zset_func_arg" .}} score {{$score.Type}}) error {
	cmd := redis.NewIntCmd(ctx, "zadd", x.key, "gt", {{GenTypeTemplate "zset_str_arg" .}}, rdconv.{{$score.RedisFunc}}ToString(score))
	return x.rds.Process(ctx, cmd)
}

{{ if $mem -}}
func (x *x{{$Name}}) ZAdds(ctx context.Context, vals map[{{$mem.Type}}]{{$score.Type}}) error {
	args := make([]interface{}, 2, 2+len(vals)*2)
	args[0] = "zadd"
	args[1] = x.key
	for k, v := range vals {
		{{if eq $mem.Type "string"}} args = append(args, k) {{else}} args = append(args, rdconv.{{$mem.RedisFunc}}ToString(k)) {{end}}
		args = append(args, rdconv.{{$score.RedisFunc}}ToString(v))
	}
	cmd := redis.NewIntCmd(ctx, args...)
	return x.rds.Process(ctx, cmd)
}
{{ end -}}

func (x *x{{$Name}}) ZRem(ctx context.Context, {{GenTypeTemplate "zset_func_arg" .}})error {
	cmd := redis.NewIntCmd(ctx, "zrem", x.key, {{GenTypeTemplate "zset_str_arg" .}})
	return x.rds.Process(ctx, cmd)
}

func (x *x{{$Name}}) ZIncrBy(ctx context.Context, increment {{$score.Type}}, {{GenTypeTemplate "zset_func_arg" .}})(_ {{$score.Type}}, err error) {
	cmd := redis.{{if $score.IsFloat}}NewFloatCmd{{else}}NewIntCmd{{end}}(ctx, "zincrby", x.key, rdconv.{{$score.RedisFunc}}ToString(increment), {{GenTypeTemplate "zset_str_arg" .}})
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	return {{$score.Type}}(cmd.Val()), nil
}

func (x *x{{$Name}}) ZScore(ctx context.Context, {{GenTypeTemplate "zset_func_arg" .}})(_ {{$score.Type}}, err error) {
	cmd := redis.{{if $score.IsFloat}}NewFloatCmd{{else}}NewIntCmd{{end}}(ctx, "zscore", x.key, {{GenTypeTemplate "zset_str_arg" .}})
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	return {{$score.Type}}(cmd.Val()), nil
}

func (x *x{{$Name}}) ZRank(ctx context.Context, {{GenTypeTemplate "zset_func_arg" .}})(_ int64, err error) {
	cmd := redis.NewIntCmd(ctx, "zrank", x.key, {{GenTypeTemplate "zset_str_arg" .}})
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	return cmd.Val(), nil
}

func (x *x{{$Name}}) ZRankWithScore(ctx context.Context, {{GenTypeTemplate "zset_func_arg" .}})(rank int64, score {{$score.Type}}, err error) {
	cmd := redis.NewRankWithScoreCmd(ctx, "zrank", x.key, {{GenTypeTemplate "zset_str_arg" .}}, "withscore")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	rank = cmd.Val().Rank
	score = {{$score.Type}}(cmd.Val().Score)
	return
}

func (x *x{{$Name}}) ZRevRank(ctx context.Context, {{GenTypeTemplate "zset_func_arg" .}})(_ int64, err error) {
	cmd := redis.NewIntCmd(ctx, "zrevrank", x.key, {{GenTypeTemplate "zset_str_arg" .}})
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	return cmd.Val(), nil
}

func (x *x{{$Name}}) ZRevRankWithScore(ctx context.Context, {{GenTypeTemplate "zset_func_arg" .}})(rank int64, score {{$score.Type}}, err error) {
	cmd := redis.NewRankWithScoreCmd(ctx, "zrevrank", x.key, {{GenTypeTemplate "zset_str_arg" .}}, "withscore")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	rank = cmd.Val().Rank
	score = {{$score.Type}}(cmd.Val().Score)
	return
}
{{if $mem }}
{{if ne $mem.Type "string" -}}
func (x *x{{$Name}}) parseMemberSliceCmd(cmd *redis.StringSliceCmd) (vals []{{$mem.Type}}, err error) {
	for _, v := range cmd.Val() {
		val, err := rdconv.StringTo{{$mem.RedisFunc}}(v)
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}
	return
}
{{end }}

func (x *x{{$Name}}) ZRange(ctx context.Context, start, stop int64) (vals []{{$mem.Type}}, err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(start), rdconv.Int64ToString(stop))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	{{if eq $mem.Type "string" -}} return cmd.Result() {{ else -}} return x.parseMemberSliceCmd(cmd) {{end}}
}

func (x *x{{$Name}}) ZRangeByScore(ctx context.Context, start, stop {{$score.Type}}) (vals []{{$mem.Type}}, err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(start), rdconv.{{$score.RedisFunc}}ToString(stop), "byscore")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	{{if eq $mem.Type "string" -}} return cmd.Result() {{ else -}} return x.parseMemberSliceCmd(cmd) {{end}}
}

func (x *x{{$Name}}) ZRevRange(ctx context.Context, start, stop int64) (vals []{{$mem.Type}}, err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(stop), rdconv.Int64ToString(start), "rev")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	{{if eq $mem.Type "string" -}} return cmd.Result() {{ else -}} return x.parseMemberSliceCmd(cmd) {{end}}
}

func (x *x{{$Name}}) ZRevRangeByScore(ctx context.Context, start, stop {{$score.Type}}) (vals []{{$mem.Type}}, err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(stop), rdconv.{{$score.RedisFunc}}ToString(start), "byscore", "rev")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	{{if eq $mem.Type "string" -}} return cmd.Result() {{ else -}} return x.parseMemberSliceCmd(cmd) {{end}}
}
{{ end }}

func (x *x{{$Name}}) rangeMemberSliceCmd(cmd *redis.StringSliceCmd, f func({{GenTypeTemplate "zset_func_arg" .}})bool) (err error) {
	{{if $farg }}var ({{ range $i,$arg := $farg}}
		{{$arg.ArgName}} {{$arg.ArgType}} {{end}}
	){{else if ne $mem.Type "string" -}}
		var mem {{$mem.Type}}
	{{end }}
	for _, v := range cmd.Val() {
		{{- if $mem -}}
			{{if eq $mem.Type "string" -}}
				if !f(v) {
					return 
				}
			{{else -}}
				mem,err = rdconv.StringTo{{$mem.RedisFunc}}(v)
				if err != nil {
					return 
				}
				if !f(mem) {
					return 
				}
			{{end}}
		{{- else -}}
			{{ range $i,$arg := $farg -}} {{$arg.ArgName}}, {{end}}err = Split{{$Name}}Member(v)
			if err != nil {
				return 
			}
			if !f({{ range $i,$arg := $farg -}} {{$arg.ArgName}},{{end}}) {
				return  
			}
		{{- end -}}
	}
	return
}

func (x *x{{$Name}}) ZRangeF(ctx context.Context, start, stop int64, f func({{GenTypeTemplate "zset_func_arg" .}}) bool) (err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(start), rdconv.Int64ToString(stop))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeMemberSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZRangeByScoreF(ctx context.Context, start, stop {{$score.Type}}, f func({{GenTypeTemplate "zset_func_arg" .}}) bool) (err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(start), rdconv.{{$score.RedisFunc}}ToString(stop), "byscore")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeMemberSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZRevRangeF(ctx context.Context, start, stop int64, f func({{GenTypeTemplate "zset_func_arg" .}}) bool) (err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(stop), rdconv.Int64ToString(start), "rev")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeMemberSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZRevRangeByScoreF(ctx context.Context, start, stop {{$score.Type}}, f func({{GenTypeTemplate "zset_func_arg" .}}) bool) (err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(stop), rdconv.{{$score.RedisFunc}}ToString(start), "byscore", "rev")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeMemberSliceCmd(cmd, f)
}

{{if $mem }}
func (x *x{{$Name}}) parseZSliceCmd(cmd *redis.ZSliceCmd) (vals map[{{$mem.Type}}]{{$score.Type}}, err error) {
	vals = make(map[{{$mem.Type}}]{{$score.Type}})
	for _, v := range cmd.Val() {
		val, err := rdconv.AnyTo{{$mem.RedisFunc}}(v.Member)
		if err != nil {
			return nil, err
		}
		vals[val] = {{$score.Type}}(v.Score)
	}
	return
}

func (x *x{{$Name}}) ZRangeWithScores(ctx context.Context, start, stop int64) (vals map[{{$mem.Type}}]{{$score.Type}}, err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(start), rdconv.Int64ToString(stop), "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.parseZSliceCmd(cmd)
}

func (x *x{{$Name}}) ZRangeByScoreWithScores(ctx context.Context, start, stop {{$score.Type}}) (vals map[{{$mem.Type}}]{{$score.Type}}, err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(start), rdconv.{{$score.RedisFunc}}ToString(stop), "byscore", "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.parseZSliceCmd(cmd)
}
func (x *x{{$Name}}) ZRevRangeWithScores(ctx context.Context, start, stop int64) (vals map[{{$mem.Type}}]{{$score.Type}}, err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(stop), rdconv.Int64ToString(start), "rev", "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.parseZSliceCmd(cmd)
}

func (x *x{{$Name}}) ZRevRangeByScoreWithScores(ctx context.Context, start, stop {{$score.Type}}) (vals map[{{$mem.Type}}]{{$score.Type}}, err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(stop), rdconv.{{$score.RedisFunc}}ToString(start), "byscore", "rev", "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.parseZSliceCmd(cmd)
}
{{end}}
func (x *x{{$Name}}) rangeZSliceCmd(cmd *redis.ZSliceCmd, f func({{GenTypeTemplate "zset_func_arg" .}}  score {{$score.Type}}) bool) (err error) {
	{{if $farg }}var ({{ range $i,$arg := $farg }}
		{{$arg.ArgName}} {{$arg.ArgType}} {{end}}
	){{else if ne $mem.Type "string" -}}
		var mem {{$mem.Type}}
	{{end }}
	for _, v := range cmd.Val() {
		{{- if $mem -}}
			{{if eq $mem.Type "string" -}}
				if !f(v.Member.(string), {{$score.Type}}(v.Score)) {
					return
				}
			{{else -}}
				mem,err = rdconv.AnyTo{{$mem.RedisFunc}}(v.Member)
				if err != nil {
					return 
				}
				if !f(mem, {{$score.Type}}(v.Score)) {
					return
				}
			{{end}}
		{{- else -}}
			{{ range $i,$arg := $farg -}} {{$arg.ArgName}}, {{end}}err = Split{{$Name}}Member(v.Member.(string))
			if err != nil {
				return 
			}
			if !f({{ range $i,$arg := $farg -}} {{$arg.ArgName}},{{end}} {{$score.Type}}(v.Score)) {
				return 
			}
		{{- end -}}
	}
	return
}

func (x *x{{$Name}}) ZRangeWithScoresF(ctx context.Context, start, stop int64, f func({{GenTypeTemplate "zset_func_arg" .}} score {{$score.Type}}) bool) (err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(start), rdconv.Int64ToString(stop), "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeZSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZRangeByScoreWithScoresF(ctx context.Context, start, stop {{$score.Type}}, f func({{GenTypeTemplate "zset_func_arg" .}} score {{$score.Type}}) bool) (err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(start), rdconv.{{$score.RedisFunc}}ToString(stop), "byscore", "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeZSliceCmd(cmd, f)
}
func (x *x{{$Name}}) ZRevRangeWithScoresF(ctx context.Context, start, stop int64, f func({{GenTypeTemplate "zset_func_arg" .}} score {{$score.Type}}) bool) (err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(stop), rdconv.Int64ToString(start), "rev", "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeZSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZRevRangeByScoreWithScoresF(ctx context.Context, start, stop {{$score.Type}}, f func({{GenTypeTemplate "zset_func_arg" .}} score {{$score.Type}}) bool) (err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(stop), rdconv.{{$score.RedisFunc}}ToString(start), "byscore", "rev", "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeZSliceCmd(cmd, f)
}

{{if $mem }}
func (x *x{{$Name}}) ZPopMin(ctx context.Context, count int64) (_ map[{{$mem.Type}}]{{$score.Type}}, err error) {
	cmd := x.rds.ZPopMin(ctx, x.key, int64(count))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return x.parseZSliceCmd(cmd)
}

func (x *x{{$Name}}) ZPopMax(ctx context.Context, count int64) (_ map[{{$mem.Type}}]{{$score.Type}}, err error) {
	cmd := x.rds.ZPopMax(ctx, x.key, int64(count))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return x.parseZSliceCmd(cmd)
}
{{end}}

func (x *x{{$Name}}) ZPopMinF(ctx context.Context, count int64, f func({{GenTypeTemplate "zset_func_arg" .}} score {{$score.Type}}) bool) (err error) {
	cmd := x.rds.ZPopMin(ctx, x.key, int64(count))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return x.rangeZSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZPopMaxF(ctx context.Context, count int64, f func({{GenTypeTemplate "zset_func_arg" .}} score {{$score.Type}}) bool) (err error) {
	cmd := x.rds.ZPopMax(ctx, x.key, int64(count))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return x.rangeZSliceCmd(cmd, f)
}

{{if $mem }}
func (x *x{{$Name}}) ZPopGTScore(ctx context.Context, limitScore {{$score.Type}}, count int64) (vals []{{$mem.Type}}, err error) {
	cmd := redis.NewStringSliceCmd(ctx, "evalsha", {{.SvcPkg}}.ZPopMaxValue.Hash, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		if !redis.HasErrorPrefix(err, "NOSCRIPT") {
			return
		}
		cmd = redis.NewStringSliceCmd(ctx, "eval", {{.SvcPkg}}.ZPopMaxValue.Script, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
		err = x.rds.Process(ctx, cmd)
		if err != nil {
			return
		}
	}
	{{if eq $mem.Type "string" -}} return cmd.Result() {{ else -}} return x.parseMemberSliceCmd(cmd) {{end}}
}
{{end}}

func (x *x{{$Name}}) ZPopGTScoreF(ctx context.Context, limitScore {{$score.Type}}, count int64, f func({{GenTypeTemplate "zset_func_arg" .}}) bool) (err error) {
	cmd := redis.NewStringSliceCmd(ctx, "evalsha", {{.SvcPkg}}.ZPopMaxValue.Hash, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		if !redis.HasErrorPrefix(err, "NOSCRIPT") {
			return
		}
		cmd = redis.NewStringSliceCmd(ctx, "eval", {{.SvcPkg}}.ZPopMaxValue.Script, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
		err = x.rds.Process(ctx, cmd)
		if err != nil {
			return
		}
	}
	//return cmd.Val(), nil
	return x.rangeMemberSliceCmd(cmd, f)
}

{{if $mem }}
func (x *x{{$Name}}) ZPopGTScoreWithScores(ctx context.Context, limitScore {{$score.Type}}, count int64) (vals map[{{$mem.Type}}]{{$score.Type}}, err error) {
	cmd := redis.NewZSliceCmd(ctx, "evalsha", {{.SvcPkg}}.ZPopMaxValueWithScore.Hash, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		if !redis.HasErrorPrefix(err, "NOSCRIPT") {
			return
		}
		cmd = redis.NewZSliceCmd(ctx, "eval", {{.SvcPkg}}.ZPopMaxValueWithScore.Script, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
		err = x.rds.Process(ctx, cmd)
		if err != nil {
			return
		}
	}
	//return cmd.Val(), nil
	return x.parseZSliceCmd(cmd)
}
{{end}}

func (x *x{{$Name}}) ZPopGTScoreWithScoresF(ctx context.Context, limitScore {{$score.Type}}, count int64, f func({{GenTypeTemplate "zset_func_arg" .}} score {{$score.Type}}) bool) (err error) {
	cmd := redis.NewZSliceCmd(ctx, "evalsha", {{.SvcPkg}}.ZPopMaxValueWithScore.Hash, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		if !redis.HasErrorPrefix(err, "NOSCRIPT") {
			return
		}
		cmd = redis.NewZSliceCmd(ctx, "eval", {{.SvcPkg}}.ZPopMaxValueWithScore.Script, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
		err = x.rds.Process(ctx, cmd)
		if err != nil {
			return
		}
	}
	return x.rangeZSliceCmd(cmd, f)
}
`)
	registerTemplate("zset_message", ` {{ $Name := .Name}} {{$msg := .TypeZSet.Message }} {{$score := .TypeZSet.Score }} 
func (x *x{{$Name}}) ZCard(ctx context.Context) (int64, error) {
	cmd := redis.NewIntCmd(ctx, "zcard", x.key)
	x.rds.Process(ctx, cmd)
	return cmd.Result()
}

func (x *x{{$Name}}) ZAdd(ctx context.Context, mem {{$msg.Type}}, score {{$score.Type}}) error {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return err
	}
	cmd := redis.NewIntCmd(ctx, "zadd", x.key, util.BytesToString(data), rdconv.{{$score.RedisFunc}}ToString(score))
	return x.rds.Process(ctx, cmd)
}

func (x *x{{$Name}}) ZAddNX(ctx context.Context, mem {{$msg.Type}}, score {{$score.Type}}) error {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return err
	}
	cmd := redis.NewIntCmd(ctx, "zadd", x.key, "nx", util.BytesToString(data), rdconv.{{$score.RedisFunc}}ToString(score))
	return x.rds.Process(ctx, cmd)
}

func (x *x{{$Name}}) ZAddXX(ctx context.Context, mem {{$msg.Type}}, score {{$score.Type}}) error {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return err
	}
	cmd := redis.NewIntCmd(ctx, "zadd", x.key, "xx", util.BytesToString(data), rdconv.{{$score.RedisFunc}}ToString(score))
	return x.rds.Process(ctx, cmd)
}

func (x *x{{$Name}}) ZAddLT(ctx context.Context, mem {{$msg.Type}}, score {{$score.Type}}) error {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return err
	}
	cmd := redis.NewIntCmd(ctx, "zadd", x.key, "lt", util.BytesToString(data), rdconv.{{$score.RedisFunc}}ToString(score))
	return x.rds.Process(ctx, cmd)
}

func (x *x{{$Name}}) ZAddGT(ctx context.Context, mem {{$msg.Type}}, score {{$score.Type}}) error {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return err
	}
	cmd := redis.NewIntCmd(ctx, "zadd", x.key, "gt", util.BytesToString(data), rdconv.{{$score.RedisFunc}}ToString(score))
	return x.rds.Process(ctx, cmd)
}

func (x *x{{$Name}}) ZAdds(ctx context.Context, vals map[{{$msg.Type}}]{{$score.Type}}) error {
	args := make([]interface{}, 2, 2+len(vals)*2)
	args[0] = "zadd"
	args[1] = x.key
	for k, v := range vals {
		data, err := {{call $msg.Marshal "k"}}
		if err != nil {
			return err
		}
		args = append(args, util.BytesToString(data))
		args = append(args, rdconv.{{$score.RedisFunc}}ToString(v))
	}
	cmd := redis.NewIntCmd(ctx, args...)
	return x.rds.Process(ctx, cmd)
}

func (x *x{{$Name}}) ZRem(ctx context.Context, mem {{$msg.Type}}) error {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return err
	}
	cmd := redis.NewIntCmd(ctx, "zrem", x.key, util.BytesToString(data))
	return x.rds.Process(ctx, cmd)
}

func (x *x{{$Name}}) ZIncrBy(ctx context.Context, increment {{$score.Type}}, mem {{$msg.Type}}) (_ {{$score.Type}}, err error) {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return
	}
	cmd := redis.{{if $score.IsFloat}}NewFloatCmd{{else}}NewIntCmd{{end}}(ctx, "zincrby", x.key, rdconv.{{$score.RedisFunc}}ToString(increment), util.BytesToString(data))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	return {{$score.Type}}(cmd.Val()), nil
}

func (x *x{{$Name}}) ZScore(ctx context.Context, mem {{$msg.Type}}) (_ {{$score.Type}}, err error) {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return
	}
	cmd := redis.{{if $score.IsFloat}}NewFloatCmd{{else}}NewIntCmd{{end}}(ctx, "zscore", x.key, util.BytesToString(data))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	return {{$score.Type}}(cmd.Val()), nil
}

func (x *x{{$Name}}) ZRank(ctx context.Context, mem {{$msg.Type}}) (_ int64, err error) {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return
	}
	cmd := redis.NewIntCmd(ctx, "zrank", x.key, util.BytesToString(data))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	return cmd.Val(), nil
}

func (x *x{{$Name}}) ZRankWithScore(ctx context.Context, mem {{$msg.Type}}) (rank int64, score {{$score.Type}}, err error) {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return
	}
	cmd := redis.NewRankWithScoreCmd(ctx, "zrank", x.key, util.BytesToString(data), "withscore")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	rank = cmd.Val().Rank
	score = {{$score.Type}}(cmd.Val().Score)
	return
}

func (x *x{{$Name}}) ZRevRank(ctx context.Context, mem {{$msg.Type}}) (_ int64, err error) {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return
	}
	cmd := redis.NewIntCmd(ctx, "zrevrank", x.key, util.BytesToString(data))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	return cmd.Val(), nil
}

func (x *x{{$Name}}) ZRevRankWithScore(ctx context.Context, mem {{$msg.Type}}) (rank int64, score {{$score.Type}}, err error) {
	data, err := {{call $msg.Marshal "mem"}}
	if err != nil {
		return
	}
	cmd := redis.NewRankWithScoreCmd(ctx, "zrevrank", x.key, util.BytesToString(data), "withscore")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	rank = cmd.Val().Rank
	score = {{$score.Type}}(cmd.Val().Score)
	return
}

func (x *x{{$Name}}) parseMemberSliceCmd(cmd *redis.StringSliceCmd) (vals []{{$msg.Type}}, err error) {
	for _, v := range cmd.Val() {
		val := {{$msg.New}}
		err = {{call $msg.Unmarshal "val" "util.StringToBytes(v)"}}
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}
	return
}

func (x *x{{$Name}}) ZRange(ctx context.Context, start, stop int64) (vals []{{$msg.Type}}, err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(start), rdconv.Int64ToString(stop))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.parseMemberSliceCmd(cmd)
}

func (x *x{{$Name}}) ZRangeByScore(ctx context.Context, start, stop {{$score.Type}}) (vals []{{$msg.Type}}, err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(start), rdconv.{{$score.RedisFunc}}ToString(stop), "byscore")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.parseMemberSliceCmd(cmd)
}

func (x *x{{$Name}}) ZRevRange(ctx context.Context, start, stop int64) (vals []{{$msg.Type}}, err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(stop), rdconv.Int64ToString(start), "rev")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.parseMemberSliceCmd(cmd)
}

func (x *x{{$Name}}) ZRevRangeByScore(ctx context.Context, start, stop {{$score.Type}}) (vals []{{$msg.Type}}, err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(stop), rdconv.{{$score.RedisFunc}}ToString(start), "byscore", "rev")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.parseMemberSliceCmd(cmd)
}

func (x *x{{$Name}}) rangeMemberSliceCmd(cmd *redis.StringSliceCmd, f func({{$msg.Type}}) bool) (err error) {
	for _, v := range cmd.Val() {
		val := {{$msg.New}}
		err = {{call $msg.Unmarshal "val" "util.StringToBytes(v)"}}
		if err != nil {
			return err
		}
		if !f(val) {
			return nil
		}
	}
	return
}

func (x *x{{$Name}}) ZRangeF(ctx context.Context, start, stop int64, f func({{$msg.Type}}) bool) (err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(start), rdconv.Int64ToString(stop))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeMemberSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZRangeByScoreF(ctx context.Context, start, stop {{$score.Type}}, f func({{$msg.Type}}) bool) (err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(start), rdconv.{{$score.RedisFunc}}ToString(stop), "byscore")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeMemberSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZRevRangeF(ctx context.Context, start, stop int64, f func({{$msg.Type}}) bool) (err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(stop), rdconv.Int64ToString(start), "rev")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeMemberSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZRevRangeByScoreF(ctx context.Context, start, stop {{$score.Type}}, f func({{$msg.Type}}) bool) (err error) {
	cmd := redis.NewStringSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(stop), rdconv.{{$score.RedisFunc}}ToString(start), "byscore", "rev")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeMemberSliceCmd(cmd, f)
}

func (x *x{{$Name}}) parseZSliceCmd(cmd *redis.ZSliceCmd) (vals map[{{$msg.Type}}]{{$score.Type}}, err error) {
	vals = make(map[{{$msg.Type}}]{{$score.Type}})
	for _, v := range cmd.Val() {
		str, err := rdconv.AnyToString(v.Member)
		if err != nil {
			return nil, err
		}
		val := {{$msg.New}}
		err = {{call $msg.Unmarshal "val" "util.StringToBytes(str)"}}
		if err != nil {
			return nil, err
		}
		vals[val] = {{$score.Type}}(v.Score)
	}
	return
}

func (x *x{{$Name}}) ZRangeWithScores(ctx context.Context, start, stop int64) (vals map[{{$msg.Type}}]{{$score.Type}}, err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(start), rdconv.Int64ToString(stop), "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.parseZSliceCmd(cmd)
}

func (x *x{{$Name}}) ZRangeByScoreWithScores(ctx context.Context, start, stop {{$score.Type}}) (vals map[{{$msg.Type}}]{{$score.Type}}, err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(start), rdconv.{{$score.RedisFunc}}ToString(stop), "byscore", "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.parseZSliceCmd(cmd)
}
func (x *x{{$Name}}) ZRevRangeWithScores(ctx context.Context, start, stop int64) (vals map[{{$msg.Type}}]{{$score.Type}}, err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(stop), rdconv.Int64ToString(start), "rev", "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.parseZSliceCmd(cmd)
}

func (x *x{{$Name}}) ZRevRangeByScoreWithScores(ctx context.Context, start, stop {{$score.Type}}) (vals map[{{$msg.Type}}]{{$score.Type}}, err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(stop), rdconv.{{$score.RedisFunc}}ToString(start), "byscore", "rev", "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.parseZSliceCmd(cmd)
}

func (x *x{{$Name}}) rangeZSliceCmd(cmd *redis.ZSliceCmd, f func({{$msg.Type}}, {{$score.Type}}) bool) (err error) {
	for _, v := range cmd.Val() {
		str, err := rdconv.AnyToString(v.Member)
		if err != nil {
			return err
		}
		val := {{$msg.New}}
		err = {{call $msg.Unmarshal "val" "util.StringToBytes(str)"}}
		if err != nil {
			return err
		}
		if !f(val, {{$score.Type}}(v.Score)) {
			return nil
		}
	}
	return
}

func (x *x{{$Name}}) ZRangeWithScoresF(ctx context.Context, start, stop int64, f func({{$msg.Type}}, {{$score.Type}}) bool) (err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(start), rdconv.Int64ToString(stop), "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeZSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZRangeByScoreWithScoresF(ctx context.Context, start, stop {{$score.Type}}, f func({{$msg.Type}}, {{$score.Type}}) bool) (err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(start), rdconv.{{$score.RedisFunc}}ToString(stop), "byscore", "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeZSliceCmd(cmd, f)
}
func (x *x{{$Name}}) ZRevRangeWithScoresF(ctx context.Context, start, stop int64, f func({{$msg.Type}}, {{$score.Type}}) bool) (err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.Int64ToString(stop), rdconv.Int64ToString(start), "rev", "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeZSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZRevRangeByScoreWithScoresF(ctx context.Context, start, stop {{$score.Type}}, f func({{$msg.Type}}, {{$score.Type}}) bool) (err error) {
	cmd := redis.NewZSliceCmd(ctx, "zrange", x.key, rdconv.{{$score.RedisFunc}}ToString(stop), rdconv.{{$score.RedisFunc}}ToString(start), "byscore", "rev", "withscores")
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}

	return x.rangeZSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZPopMin(ctx context.Context, count int64) (_ map[{{$msg.Type}}]{{$score.Type}}, err error) {
	cmd := x.rds.ZPopMin(ctx, x.key, int64(count))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return x.parseZSliceCmd(cmd)
}

func (x *x{{$Name}}) ZPopMinF(ctx context.Context, count int64, f func({{$msg.Type}}, {{$score.Type}}) bool) (err error) {
	cmd := x.rds.ZPopMin(ctx, x.key, int64(count))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return x.rangeZSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZPopMax(ctx context.Context, count int64) (_ map[{{$msg.Type}}]{{$score.Type}}, err error) {
	cmd := x.rds.ZPopMax(ctx, x.key, int64(count))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return x.parseZSliceCmd(cmd)
}

func (x *x{{$Name}}) ZPopMaxF(ctx context.Context, count int64, f func({{$msg.Type}}, {{$score.Type}}) bool) (err error) {
	cmd := x.rds.ZPopMax(ctx, x.key, int64(count))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return x.rangeZSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZPopGTScore(ctx context.Context, limitScore {{$score.Type}}, count int64) (vals []{{$msg.Type}}, err error) {
	cmd := redis.NewStringSliceCmd(ctx, "evalsha", {{.SvcPkg}}.ZPopMaxValue.Hash, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		if !redis.HasErrorPrefix(err, "NOSCRIPT") {
			return
		}
		cmd = redis.NewStringSliceCmd(ctx, "eval", {{.SvcPkg}}.ZPopMaxValue.Script, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
		err = x.rds.Process(ctx, cmd)
		if err != nil {
			return
		}
	}
	return x.parseMemberSliceCmd(cmd)
}

func (x *x{{$Name}}) ZPopGTScoreF(ctx context.Context, limitScore {{$score.Type}}, count int64, f func({{$msg.Type}}) bool) (err error) {
	cmd := redis.NewStringSliceCmd(ctx, "evalsha", {{.SvcPkg}}.ZPopMaxValue.Hash, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		if !redis.HasErrorPrefix(err, "NOSCRIPT") {
			return
		}
		cmd = redis.NewStringSliceCmd(ctx, "eval", {{.SvcPkg}}.ZPopMaxValue.Script, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
		err = x.rds.Process(ctx, cmd)
		if err != nil {
			return
		}
	}
	//return cmd.Val(), nil
	return x.rangeMemberSliceCmd(cmd, f)
}

func (x *x{{$Name}}) ZPopGTScoreWithScores(ctx context.Context, limitScore {{$score.Type}}, count int64) (vals map[{{$msg.Type}}]{{$score.Type}}, err error) {
	cmd := redis.NewZSliceCmd(ctx, "evalsha", {{.SvcPkg}}.ZPopMaxValueWithScore.Hash, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		if !redis.HasErrorPrefix(err, "NOSCRIPT") {
			return
		}
		cmd = redis.NewZSliceCmd(ctx, "eval", {{.SvcPkg}}.ZPopMaxValueWithScore.Script, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
		err = x.rds.Process(ctx, cmd)
		if err != nil {
			return
		}
	}
	//return cmd.Val(), nil
	return x.parseZSliceCmd(cmd)
}

func (x *x{{$Name}}) ZPopGTScoreWithScoresF(ctx context.Context, limitScore {{$score.Type}}, count int64, f func({{$msg.Type}}, {{$score.Type}}) bool) (err error) {
	cmd := redis.NewZSliceCmd(ctx, "evalsha", {{.SvcPkg}}.ZPopMaxValueWithScore.Hash, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		if !redis.HasErrorPrefix(err, "NOSCRIPT") {
			return
		}
		cmd = redis.NewZSliceCmd(ctx, "eval", {{.SvcPkg}}.ZPopMaxValueWithScore.Script, "1", x.key, rdconv.{{$score.RedisFunc}}ToString(limitScore), rdconv.Int64ToString(count))
		err = x.rds.Process(ctx, cmd)
		if err != nil {
			return
		}
	}
	//return cmd.Val(), nil
	return x.rangeZSliceCmd(cmd, f)
}
`)
}
