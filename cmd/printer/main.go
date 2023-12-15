package main

import (
	"log"

	"github.com/walleframe/wplugins/buildpb"
	"github.com/walleframe/wplugins/utils"
	"github.com/walleframe/wplugins/utils/plugin"
)

func main() {
	plugin.MainRoot(func(rq *buildpb.BuildRQ) (rs *buildpb.BuildRS, err error) {
		log.Println(utils.Sdump(rq, "rq"))
		return
	})
}
