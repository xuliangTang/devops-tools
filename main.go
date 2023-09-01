package main

import (
	"devops-tools/cmd"
	"devops-tools/utils/config"
)

func main() {
	config.InitSysConfig()
	cmd.RunCmd()
}
