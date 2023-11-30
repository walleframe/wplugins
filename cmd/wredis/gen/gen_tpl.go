package gen

var GenRedisTemplate = `
package {{.Package}} {{$obj := .}}

$Import-Packages$

type x{{.Name}} struct{
	key string
	rds redis.UniversalClient
}

func init() {
	{{.SvcPkg}}.RegisterDBName({{.SvcPkg}}.DBType, "{{.Package}}.{{.Name}}")
}

{{Doc .Doc}}
func {{.Name}}({{range $i,$arg := .Args}}{{if $arg.ConstructArg }}{{$arg.ArgName}} {{$arg.ArgType }}, {{end}}{{- end}}) *x{{.Name}} {
	buf := util.Builder{}
	buf.Grow({{.KeySize}})
{{- range $i,$arg := .Args -}}
{{- if gt $i 0 -}}
	buf.WriteByte(':')
{{- end }}
	{{$arg.FormatCode "buf"}}
{{end -}}
	return &x{{.Name}}{
		key: buf.String(),
		rds: {{.SvcPkg}}.GetDBLink({{.SvcPkg}}.DBType, "{{.Package}}.{{.Name}}"),
	}
}

// With reset redis client
func (x *x{{.Name}}) With(rds redis.UniversalClient) *x{{.Name}} {
	x.rds = rds
	return x
}

func (x *x{{.Name}}) Key() string {
	return x.key
}

{{if .Lock }}
////////////////////////////////////////////////////////////
// redis lock operation
{{GenTypeTemplate "redis_lock" .}}
{{end}}

{{if .TypeKeys }}
////////////////////////////////////////////////////////////
// redis keys operation
{{GenTypeTemplate "key" .}} {{end}}

{{- if .TypeString }}
////////////////////////////////////////////////////////////
// redis string operation
{{if .TypeString.Custom}}
	 {{if .TypeString.Protobuf}}
{{GenTypeTemplate "string_custom_protobuf" .}}
	{{else}}
{{GenTypeTemplate "string_custom_walle" .}}
	{{end}}
{{else}}
{{if .TypeString.Number }} {{GenTypeTemplate "string_number" .}} {{end}}
{{if .TypeString.Float }} {{GenTypeTemplate "string_float" .}} {{end}}
{{if .TypeString.String }} {{GenTypeTemplate "string_string" .}} {{end}}
{{if .TypeString.Protobuf }} {{GenTypeTemplate "string_protobuf" .}} {{end}}
{{if .TypeString.WProto }} {{GenTypeTemplate "string_walle" .}} {{end}}
{{end}}
{{- end -}}

{{- if .TypeHash -}}
////////////////////////////////////////////////////////////
// redis hash operation
{{if .TypeHash.HashObject }}
const ({{ $Name := .Name}}
{{ range $i,$field := .TypeHash.HashObject.Fields }}
	 _{{$Name}}_{{Title $field.Name}} = "{{$field.Name}}"{{ end }}
)
{{GenTypeTemplate "hash_object" .}} {{end}}
{{if .TypeHash.HashDynamic }} {{GenTypeTemplate "hash_dynamic_match" .}} {{end}}
{{ end -}}

{{- if .TypeSet -}}
////////////////////////////////////////////////////////////
// redis set operation
{{if .TypeSet.BaseType }} {{GenTypeTemplate "set_type" .}} {{end}}
{{if .TypeSet.Message }} {{GenTypeTemplate "set_message" .}} {{end}}
{{end}}

{{- if .TypeZSet -}}
////////////////////////////////////////////////////////////
// redis zset operation
{{if .TypeZSet.Message -}}
{{GenTypeTemplate "zset_message" .}}
{{else -}}
{{GenTypeTemplate "zset_basic" .}}
{{end}}
{{end}}

{{- range $i,$Script := .Scripts }}
{{GenScriptTemplate $obj $Script}}
{{end}}
`

func registerTemplate(name, template string) {
	redisTypeTemplate[name] = template
}

var redisTypeTemplate = map[string]string{}
