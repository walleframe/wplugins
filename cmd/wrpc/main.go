package main

import (
	"github.com/aggronmagi/wplugins/buildpb"
	"github.com/aggronmagi/wplugins/utils/plugin"
)

func main() {
	plugin.MainOneByOne(generateWalleRPC)
}

func generateWalleRPC(prog *buildpb.FileDesc, depend map[string]*buildpb.FileDesc) (out []*buildpb.BuildOutput, err error) {
	return
}
