package main

import (
	"github.com/mitchellh/packer/builder/cloudstack"
	"github.com/mitchellh/packer/packer/plugin"
)

func main() {
	plugin.ServeBuilder(new(cloudstack.Builder))
}
