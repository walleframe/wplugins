/*
Copyright © 2020 aggronmagi <czy463@163.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package gen

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/walleframe/wplugins/utils"
	"golang.org/x/tools/imports"
)

//go:generate gogen option -n Option -o options.go
func generateOptions() interface{} {
	return map[string]interface{}{
		// 缩进
		"Indent": "\t",
		// 代码格式化
		"Format": func(in []byte) (out []byte, err error) { return in, nil },
		// key 是否大写
		"KeyTitle": true,
	}
}

// Generator 生成器
type Generator struct {
	*bytes.Buffer
	indent int
	buf    []*bytes.Buffer
	inbuf  []int
	tab    string
	cc     *Options
}

// New 新建生成器
func New(opts ...Option) *Generator {
	gen := new(Generator)
	gen.cc = NewOptions(opts...)
	gen.buf = make([]*bytes.Buffer, 0, 20)
	gen.inbuf = make([]int, 0, 20)
	gen.Buffer = &bytes.Buffer{}
	return gen
}

func (gen *Generator) Key(k string) string {
	if gen.cc.KeyTitle {
		return utils.Title(k)
	}
	return k
}

// Error reports a problem, including an error, and exits the program.
func (gen *Generator) Error(err error, msgs ...string) {
	s := strings.Join(msgs, " ") + ":" + err.Error()
	log.Println("Error: ", s)
	panic(err)
}

// Fail reports a problem and exits the program.
func (gen *Generator) Fail(msgs ...string) {
	s := strings.Join(msgs, " ")
	log.Println("Error: ", s)
	panic(s)
}

// Print 打印数据
func (gen *Generator) Print(str ...interface{}) {
	for _, v := range str {
		switch s := v.(type) {
		case string:
			gen.WriteString(s)
		case *string:
			gen.WriteString(*s)
		case bool:
			fmt.Fprintf(gen, "%t", s)
		case *bool:
			fmt.Fprintf(gen, "%t", *s)
		case int, int32, uint, uint32, int64, uint64, int8, uint8, int16, uint16:
			fmt.Fprintf(gen, "%d", s)
		case *int32:
			fmt.Fprintf(gen, "%d", *s)
		case *int64:
			fmt.Fprintf(gen, "%d", *s)
		case float64:
			fmt.Fprintf(gen, "%g", s)
		case *float64:
			fmt.Fprintf(gen, "%g", *s)
		default:
			gen.Fail(fmt.Sprintf("unknown type in printer: %T", v))
			//panic("here")

		}
	}

}

// Pf 同 Printf
func (gen *Generator) Pf(fmts string, v ...interface{}) {
	gen.Printf(fmts, v...)
}

// P 同 Println
func (gen *Generator) P(str ...interface{}) {
	gen.Println(str...)
}

// Printf
func (gen *Generator) Printf(fmts string, v ...interface{}) {
	if len(v) < 1 {
		if len(fmts) > 0 {
			gen.BeginLine()
			gen.Print(fmts)
			gen.EndLine()
			return
		}
		gen.EndLine()
		return
	}
	gen.BeginLine()
	fmt.Fprintf(gen, fmts, v...)
	gen.EndLine()
}

// Println 打印数据
func (gen *Generator) Println(str ...interface{}) {
	if len(str) < 1 {
		gen.EndLine()
		return
	}
	gen.BeginLine()
	gen.Print(str...)
	gen.EndLine()
}

// BeginLine 行开始
func (gen *Generator) BeginLine() {
	// fmt.Println("gen.indent", gen.indent, "[", gen.tab, "]")
	for i := 0; i < gen.indent; i++ {
		gen.WriteString(gen.cc.Indent)
	}
}

// EndLine 行结束
func (gen *Generator) EndLine() {
	gen.WriteString("\n")
}

// In Indents the output one tab stop.
func (gen *Generator) In() { gen.indent++ }

// Out unindents the output one tab stop.
func (gen *Generator) Out() {
	gen.indent--
}

// WriteString 写string
func (gen *Generator) WriteString(str string) {
	if len(gen.buf) > 0 {
		gen.buf[len(gen.buf)-1].WriteString(str)
		return
	}
	gen.Write([]byte(str))
}

// GoFmtByest go文件格式化
func (gen *Generator) Bytes() (data []byte, err error) {
	data = gen.Buffer.Bytes()
	if gen.cc.Format != nil {
		return gen.cc.Format(data)
	}
	return
}

// PushObj 入栈
func (gen *Generator) PushObj() {
	gen.buf = append(gen.buf, &bytes.Buffer{})
	gen.inbuf = append(gen.inbuf, gen.indent)
	gen.indent = 0
}

// PopAndWrite 出栈
func (gen *Generator) PopAndWrite() {
	buf := gen.buf[len(gen.buf)-1]
	gen.buf = gen.buf[:len(gen.buf)-1]
	lines := strings.Split(buf.String(), "\n")
	indent := gen.inbuf[len(gen.inbuf)-1]
	gen.inbuf = gen.inbuf[:len(gen.inbuf)-1]
	gen.indent = indent
	num := -1
	for _, line := range lines {
		if num == -1 {
			num = strings.Count(line, gen.tab)
			//logFile.Println("---------------------", num, gen.indent)
		}
		line = strings.Replace(line, gen.cc.Indent, "", num)
		for i := 0; i < gen.indent; i++ {
			gen.Write([]byte(gen.cc.Indent))
		}
		gen.Write([]byte(line))
		gen.Write([]byte("\n"))
	}
}

func GoFormat(in []byte) (out []byte, err error) {
	return format.Source(in)
}

func GoimportsCmdFormat(in []byte) (out []byte, err error) {
	ib := bytes.NewBuffer(in)
	ob := bytes.NewBuffer(nil)
	cmd := exec.Command("goimports")
	cmd.Stdin = ib
	cmd.Stdout = ob

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	return ob.Bytes(), err
}

func GoFormat2(in []byte) ([]byte, error) {
	opts := &imports.Options{
		TabIndent: false,
		TabWidth:  4,
		Fragment:  true,
		Comments:  true,
	}
	data, err := imports.Process("", in, opts)
	if err != nil {
		return nil, err
	}
	return format.Source(data)
}

// OptionGoimportsFormtat use the command line `goimports` command to
// format first. If an error is reported, use the default formatting method.
func OptionGoimportsFormtat(in []byte) ([]byte, error) {
	out := bytes.NewBuffer(nil)
	cmd := exec.Command("goimports")
	cmd.Stdin = bytes.NewBuffer(in)
	cmd.Stdout = out
	err := cmd.Run()
	if err == nil {
		return out.Bytes(), nil
	}
	opts := &imports.Options{
		TabIndent: false,
		TabWidth:  4,
		Fragment:  true,
		Comments:  true,
	}
	data, err := imports.Process("", in, opts)
	if err != nil {
		return nil, err
	}
	return format.Source(data)
}

// Content returns the contents of the generated file.
func Content(filename string, original []byte, packageNames map[string]string) ([]byte, error) {
	if !strings.HasSuffix(filename, ".go") {
		return original, nil
	}

	// Reformat generated code.
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", original, parser.ParseComments)
	if err != nil {
		// Print out the bad code with line numbers.
		// This should never happen in practice, but it can while changing generated code
		// so consider this a debugging aid.
		var src bytes.Buffer
		s := bufio.NewScanner(bytes.NewReader(original))
		for line := 1; s.Scan(); line++ {
			fmt.Fprintf(&src, "%5d\t%s\n", line, s.Bytes())
		}
		return nil, fmt.Errorf("%v: unparsable Go source: %v\n%v", filename, err, src.String())
	}

	// Collect a sorted list of all imports.
	var importPaths [][2]string
	rewriteImport := func(importPath string) string {
		// if f := g.gen.opts.ImportRewriteFunc; f != nil {
		// 	return string(f(GoImportPath(importPath)))
		// }
		return importPath
	}
	for importPath := range packageNames {
		pkgName := string(packageNames[importPath])
		pkgPath := rewriteImport(string(importPath))
		importPaths = append(importPaths, [2]string{pkgName, pkgPath})
	}
	// for importPath := range g.manualImports {
	// 	if _, ok := g.packageNames[importPath]; !ok {
	// 		pkgPath := rewriteImport(string(importPath))
	// 		importPaths = append(importPaths, [2]string{"_", pkgPath})
	// 	}
	// }
	sort.Slice(importPaths, func(i, j int) bool {
		return importPaths[i][1] < importPaths[j][1]
	})

	// Modify the AST to include a new import block.
	if len(importPaths) > 0 {
		// Insert block after package statement or
		// possible comment attached to the end of the package statement.
		pos := file.Package
		tokFile := fset.File(file.Package)
		pkgLine := tokFile.Line(file.Package)
		for _, c := range file.Comments {
			if tokFile.Line(c.Pos()) > pkgLine {
				break
			}
			pos = c.End()
		}

		// Construct the import block.
		impDecl := &ast.GenDecl{
			Tok:    token.IMPORT,
			TokPos: pos,
			Lparen: pos,
			Rparen: pos,
		}
		for _, importPath := range importPaths {
			impDecl.Specs = append(impDecl.Specs, &ast.ImportSpec{
				Name: &ast.Ident{
					Name:    importPath[0],
					NamePos: pos,
				},
				Path: &ast.BasicLit{
					Kind:     token.STRING,
					Value:    strconv.Quote(importPath[1]),
					ValuePos: pos,
				},
				EndPos: pos,
			})
		}
		file.Decls = append([]ast.Decl{impDecl}, file.Decls...)
	}

	var out bytes.Buffer
	if err = (&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(&out, fset, file); err != nil {
		return nil, fmt.Errorf("%v: can not reformat Go source: %v", filename, err)
	}
	return out.Bytes(), nil
}
