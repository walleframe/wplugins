package gengo

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"
)

func GenExec(data *GenerateStruct) (_ []byte, err error) {
	tpl := template.New("").Funcs(UseFuncMap)
	// 基础导入go包函数
	tpl.Funcs(template.FuncMap{
		"Import": data.Import,
	})
	// 循环调用模板函数
	tpl.Funcs(template.FuncMap{
		"GenTemplate": func(tplName string, field *GenerateField, vals ...string) (string, error) {
			if len(vals)%2 != 0 {
				return "", fmt.Errorf("%v recv vals invalid. %v. need 2n", tplName, vals)
			}
			buf := &bytes.Buffer{}
			dst := tpl.Lookup(tplName)
			if dst == nil {
				return "", fmt.Errorf("%v not found", tplName)
			}

			Vals := make(map[string]string)

			for i := 0; i < len(vals); i += 2 {
				Vals[vals[i]] = vals[i+1]
			}

			err = dst.Execute(buf, &map[string]interface{}{
				"Field": field,
				"V":     Vals,
			})
			if err != nil {
				log.Println(err)
				return "", err
			}

			//		log.Println(buf.String())

			return buf.String(), nil
		},
		"GenCustomTemplate": func(tplName string, msg *GenerateMessage) (string, error) {
			buf := &bytes.Buffer{}
			dst := tpl.Lookup(tplName)
			if dst == nil {
				return "", fmt.Errorf("%v not found", tplName)
			}

			err = dst.Execute(buf, msg)
			if err != nil {
				log.Println(err)
				return "", err
			}

			return buf.String(), nil
		},
	})

	// 解析前置的全部魔板
	for k, v := range GenProtobufTemplate {
		ts := fmt.Sprintf(`{{define "%s"}} %s {{end}}`, k, strings.TrimSpace(v))
		tpl, err = tpl.Parse(ts)
		if err != nil {
			err = fmt.Errorf("parse template %s failed:%w", k, err)
			return
		}
		// log.Println(ts)
	}
	// 嵌入模板
	for _, model := range globalModules {
		tpl.Funcs(model.Funcs)
		for _, nt := range model.Templates {
			ts := fmt.Sprintf(`{{define "%s"}} %s {{end}}`, nt[0], strings.TrimSpace(nt[1]))
			tpl.Parse(ts)
			if err != nil {
				err = fmt.Errorf("parse template %s failed:%w", nt[0], err)
				return
			}
		}
	}
	// 完整模板
	tpl, err = tpl.Parse(GenerateTemplate)
	if err != nil {
		err = fmt.Errorf("parse full template failed:%w", err)
		return
	}
	// basic import 
	data.Import("strconv", "")
	data.Import("errors", "")
	data.Import("math", "")
	// 执行生成
	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, data)
	if err != nil {
		err = fmt.Errorf("parse full template failed:%w", err)
		return
	}
	// return buf.Bytes(), nil
	return bytes.Replace(buf.Bytes(), []byte("$Import-Packages$"), []byte(fmt.Sprintf(`
import (
	"%s"%s
)`, data.WirePkg, data.customImport())), 1), nil
}

type CustomModule struct {
	Templates [][2]string // 模板名称和模板内容
	Funcs     template.FuncMap
}

var globalModules []*CustomModule

func RegisterCustomModule(m *CustomModule) {
	globalModules = append(globalModules, m)
}
