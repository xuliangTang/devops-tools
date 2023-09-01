package config

import (
	"devops-tools/utils/helpers"
	"gopkg.in/yaml.v3"
	"log"
)

var SysConfig *sysConfig

type sysConfig struct {
	Remotes []*Remote `yaml:"remotes"` // 远程主机列表
}

type Remote struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func InitSysConfig() {
	SysConfig = new(sysConfig)
	b := helpers.MustLoadFile("app.yaml")
	if err := yaml.Unmarshal(b, SysConfig); err != nil {
		log.Fatalln("aaa", err)
	}
}
