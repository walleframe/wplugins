package main

import (
	"log"

	"github.com/aggronmagi/wplugins/buildpb"
	"github.com/aggronmagi/wplugins/utils"
	"github.com/aggronmagi/wplugins/utils/plugin"
)

//
func main() {
	plugin.MainRoot(func(rq *buildpb.BuildRQ) (rs *buildpb.BuildRS, err error) {
		log.Println(utils.Sdump(rq, "rq"))
		return
	})
}
