package gen

import (
	"strings"
	"text/template"

	"github.com/walleframe/wplugins/buildpb"
	"github.com/walleframe/wplugins/utils"
)

var UseFuncMap = template.FuncMap{}

func init() {
	UseFuncMap["Doc"] = func(doc *buildpb.DocDesc) string {
		if doc == nil {
			return ""
		}
		buf := strings.Builder{}
		for _, v := range doc.Doc {
			buf.WriteString(v)
			buf.WriteByte('\n')
		}
		if len(doc.TailDoc) > 0 {
			buf.WriteString(doc.TailDoc)
			buf.WriteByte('\n')
		}
		return buf.String()
	}
	UseFuncMap["Title"] = utils.Title
}
