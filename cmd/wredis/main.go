package main

import (
	"log"

	"github.com/aggronmagi/wplugins/buildpb"
	"github.com/aggronmagi/wplugins/cmd/wredis/gen"
	"github.com/aggronmagi/wplugins/options"
	"github.com/aggronmagi/wplugins/utils/plugin"
)

func main() {
	plugin.MainOneByOne(generateOneFile)
}

func generateOneFile(prog *buildpb.FileDesc, depend map[string]*buildpb.FileDesc) (out []*buildpb.BuildOutput, err error) {
	for _, msg := range prog.Msgs {
		if !msg.HasOption(options.RedisOpKey) {
			log.Println("ignore message", msg.Name)
			continue
		}
		o, err := gen.GenerateRedisMessage(prog, msg, depend)
		if err != nil {
			return nil, err
		}
		if o == nil {
			continue
		}
		out = append(out, o)
	}
	return
}
