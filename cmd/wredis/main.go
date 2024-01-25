package main

import (
	"log"

	"github.com/walleframe/wplugins/buildpb"
	"github.com/walleframe/wplugins/cmd/wredis/gen"
	"github.com/walleframe/wplugins/options"
	"github.com/walleframe/wplugins/utils/plugin"
)

func main() {
	plugin.MainRangeMessage(nil,
		func(msg *buildpb.MsgDesc) bool {
			if !msg.HasOption(options.RedisOpKey) {
				log.Println("ignore message", msg.Name)
				return false
			}
			return true
		},
		generateOneFile,
	)
}

func generateOneFile(msg *buildpb.MsgDesc, prog *buildpb.FileDesc, depend map[string]*buildpb.FileDesc) (out []*buildpb.BuildOutput, err error) {
	o, err := gen.GenerateRedisMessage(prog, msg, depend)
	if err != nil {
		return nil, err
	}
	if o == nil {
		return
	}
	out = append(out, o)
	return
}
