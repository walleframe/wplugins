package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/walleframe/wplugins/buildpb"
	"github.com/walleframe/wplugins/utils"
	"github.com/walleframe/wplugins/utils/plugin"
)

func main() {
	plugin.MainRoot(func(rq *buildpb.BuildRQ) (rs *buildpb.BuildRS, err error) {
		if os.Getenv("PRINT_JSON") != "" {
			datas := make([]*buildpb.FileDesc, 0, len(rq.Files))
			for _, f := range rq.Files {
				file := rq.Programs[f]
				for _, msg := range file.Msgs {
					resetMsg(msg)
				}
				datas = append(datas, file)
			}
			buf, err2 := json.MarshalIndent(&datas, "", "  ")
			if err2 != nil {
				log.Println(err2)
				return
			}
			log.Println(string(buf))
			return
		}
		log.Println(utils.Sdump(rq, "rq"))
		return
	})
}

func resetMsg(msg *buildpb.MsgDesc) {
	for _, f := range msg.Fields {
		if f.Type.Msg != nil {
			f.Type.Msg = nil
		}
	}
	for _, sub := range msg.SubMsgs {
		resetMsg(sub)
	}
}
