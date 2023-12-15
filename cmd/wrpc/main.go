package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/walleframe/wplugins/buildpb"
	"github.com/walleframe/wplugins/gen"
	"github.com/walleframe/wplugins/utils/plugin"
)

func main() {
	plugin.MainOneByOne(generateWalleRPC)
}

func generateWalleRPC(prog *buildpb.FileDesc, depend map[string]*buildpb.FileDesc) (out []*buildpb.BuildOutput, err error) {

	g := gen.New(
		gen.WithFormat(gen.GoFormat2),
		gen.WithIndent("    "),
		gen.WithKeyTitle(true),
	)

	g.P("// Generate by wctl plugin(wrpc). DO NOT EDIT.")
	g.P("package ", prog.Pkg.Package, ";")
	g.P()
	g.P("import (")
	g.In()
	g.P(`"github.com/walleframe/walle/network"`)
	g.P(`"github.com/walleframe/walle/network/rpc"`)
	g.P(`"github.com/walleframe/walle/process"`)
	g.P()
	g.P(`"go.uber.org/zap/zapcore"`)
	g.P()
	g.P(`"context"`)
	// 用于binary数据打印
	g.P(`"encoding/base64"`)
	// 用于map key 打印
	g.P(`"strconv"`)
	g.Out()
	g.P(")")

	g.P()
	g.P()

	for _, svc := range prog.Services {
		g.Doc(svc.Doc)
		// 消息路由定义
		g.P("// ", svc.Name, " method uri define")
		g.P(`const (`)
		g.In()
		for _, method := range svc.Methods {
			g.P("__", g.Key(svc.Name), g.Key(method.Name), ` = "/`, method.Name, `"`)
		}
		g.Out()
		g.P(`)`)
		g.P()
		// service 定义
		g.P("type ", g.Key(svc.Name), "Service interface {")
		g.In()
		for _, method := range svc.Methods {
			g.Doc(method.Doc)
			if method.Reply == nil {
				g.Pf("%s(ctx network.SessionContext, rq* %s)(err error)",
					g.Key(method.Name), g.Key(method.Request.Name),
				)
				continue
			}
			g.Pf("%s(ctx network.SessionContext, rq* %s,rs *%s)(err error)",
				g.Key(method.Name), g.Key(method.Request.Name),
				g.Key(method.Reply.Name),
			)

		}
		g.Out()
		g.P("}")
		g.P()
		// 注册函数
		g.P("func Register", g.Key(svc.Name), "Service(router process.Router, s ", g.Key(svc.Name), "Service) {")
		g.In()
		g.P("svc := &w", g.Key(svc.Name), "Service{svc: s}")
		for _, method := range svc.Methods {
			g.P("router.Register(__", g.Key(svc.Name), g.Key(method.Name), `, svc.`, g.Key(method.Name), `)`)
		}
		g.Out()
		g.P("}")
		g.P()

		// client 定义
		g.P("type ", g.Key(svc.Name), "Client interface {")
		g.In()
		for _, method := range svc.Methods {
			g.Doc(method.Doc)
			if method.IsNotify() {
				g.Pf("%s(ctx context.Context, rq *%s, opts ...rpc.NoticeOption)(err error)",
					g.Key(method.Name), g.Key(method.Request.Name),
				)

				continue
			}
			if method.Reply == nil {
				g.Pf("%s(ctx context.Context, rq *%s, opts ...rpc.CallOption)(err error)",
					g.Key(method.Name), g.Key(method.Request.Name),
				)
				continue
			}

			g.Pf("%s(ctx context.Context, rq *%s, opts ...rpc.CallOption)(rs *%s,err error)",
				g.Key(method.Name), g.Key(method.Request.Name),
				g.Key(method.Reply.Name),
			)

			g.Pf("%sAsync(ctx context.Context, rq *%s, resp func(ctx process.Context, rs *%s, err error), opts ...rpc.AsyncCallOption)(err error)",
				g.Key(method.Name), g.Key(method.Request.Name),
				g.Key(method.Reply.Name),
			)
		}
		g.Out()
		g.P("}")
		g.P()

		// service 实现
		g.P("type w", g.Key(svc.Name), "Service struct {")
		g.In()
		g.P("svc ", g.Key(svc.Name), "Service")
		g.Out()
		g.P("}")
		g.P()
		for _, method := range svc.Methods {
			if method.IsNotify() {
				g.Pf(`
func (s *w%[1]sService) %[2]s(c process.Context) {
	ctx := c.(network.SessionContext)
	rq := %[3]s{}
	err := ctx.Bind(&rq)
	if err != nil {
		return
	}
	err = s.svc.%[2]s(ctx, &rq)
	if err != nil {
		return
	}
	return
}`, g.Key(svc.Name), g.Key(method.Name), g.Key(method.Request.Name))
				continue
			}
			if method.Reply == nil {
				g.Pf(`
func (s *w%[1]sService) %[2]s(c process.Context) {
	ctx := c.(network.SessionContext)
	rq := %[3]s{}
	err := ctx.Bind(&rq)
	if err != nil {
		ctx.Respond(ctx, err, nil)
		return
	}
	err = s.svc.%[2]s(ctx, &rq)
	if err != nil {
		ctx.Respond(ctx, err, nil)
		return
	}
	ctx.Respond(ctx, nil, nil)
	return
}`, g.Key(svc.Name), g.Key(method.Name), g.Key(method.Request.Name))
				continue
			}

			g.Pf(`
func (s *w%[1]sService) %[2]s(c process.Context) {
	ctx := c.(network.SessionContext)
	rq := %[3]s{}
	rs := %[4]s{}
	err := ctx.Bind(&rq)
	if err != nil {
		ctx.Respond(ctx, err, nil)
		return
	}
	err = s.svc.%[2]s(ctx, &rq, &rs)
	if err != nil {
		ctx.Respond(ctx, err, nil)
		return
	}
	ctx.Respond(ctx, &rs, nil)
	return
}`, g.Key(svc.Name), g.Key(method.Name),
				g.Key(method.Request.Name),
				g.Key(method.Reply.Name),
			)
		}
		// 客户端实现
		g.Pf(`
type w%[1]sClient struct {
	cli network.Client
}

func New%[1]sClient(cli network.Client) %[1]sClient {
	return &w%[1]sClient{
		cli: cli,
	}
}
`, g.Key(svc.Name))
		for _, method := range svc.Methods {
			if method.IsNotify() {
				g.Pf(`
func (c *w%[1]sClient) %[2]s(ctx context.Context, rq *%[3]s, opts ...rpc.NoticeOption) (err error) {
	cc := rpc.NewNoticeOptions(opts...)
	err = c.cli.Notify(ctx, __%[1]s%[2]s, rq, cc)
	if err != nil {
		return err
	}
	return
}
`, g.Key(svc.Name), g.Key(method.Name), g.Key(method.Request.Name))
				continue
			}
			if method.Reply == nil {
				g.Pf(`
func (c *w%[1]sClient) %[2]s(ctx context.Context, rq *%[3]s, opts ...rpc.CallOption) (err error) {
	cc := rpc.NewCallOptions(opts...)
	err = c.cli.Call(ctx, __%[1]s%[2]s, rq, nil, cc)
	if err != nil {
		return err
	}
	return
}
`, g.Key(svc.Name), g.Key(method.Name), g.Key(method.Request.Name))
				continue
			}
			g.Pf(`
func (c *w%[1]sClient) %[2]s(ctx context.Context, rq *%[3]s, opts ...rpc.CallOption)(rs *%[4]s, err error) {
	cc := rpc.NewCallOptions(opts...)
	rs = &%[4]s{}
	err = c.cli.Call(ctx, __%[1]s%[2]s, rq, rs, cc)
	if err != nil {
		return nil, err
	}
	return rs, nil
}
`, g.Key(svc.Name), g.Key(method.Name), g.Key(method.Request.Name), g.Key(method.Reply.Name))

			g.Pf(`
func (c *w%[1]sClient) %[2]sAsync(ctx context.Context, rq *%[3]s, 
	rf func(ctx process.Context, rs *%[4]s, err error),
	opts ...rpc.AsyncCallOption) (err error) {
	cc := rpc.NewAsyncCallOptions(opts...)
	err = c.cli.AsyncCall(ctx, __%[1]s%[2]s, rq, func(c process.Context) {
		rs := &%[4]s{}
		err := c.Bind(rs)
		if err != nil {
			rf(c, nil, err)
			return
		}
		rf(c, rs, nil)
		return
	}, cc)
	if err != nil {
		return err
	}
	return nil
}
`, g.Key(svc.Name), g.Key(method.Name), g.Key(method.Request.Name), g.Key(method.Reply.Name))

		}
	}
	data, err := g.Bytes()
	if err != nil {
		log.Println("format code failed.", err)
		log.Println(string(g.Buffer.Bytes()))
		//err = nil
		return nil, fmt.Errorf("format failed %w", err)
		// data = g.Buffer.Bytes()
		// err = nil
	}
	out = append(out, &buildpb.BuildOutput{
		File: strings.TrimSuffix(prog.File, filepath.Ext(prog.File)) + ".rpc.go",
		Data: data,
	})
	return
}
