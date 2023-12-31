package remotes

import (
	"devops-tools/utils/config"
	"log"
)

// GetRemote 从配置文件根据名称寻找远程主机
func GetRemote(name string) (*config.Remote, bool) {
	for _, r := range config.SysConfig.Remotes {
		if r.Name == name {
			return r, true
		}
	}

	return nil, false
}

func MustGetRemote(name string) *config.Remote {
	getRemote, exist := GetRemote(name)
	if !exist {
		log.Fatalln("no such remote:", name)
	}
	return getRemote
}
