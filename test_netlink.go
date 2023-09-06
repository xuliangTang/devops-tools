package main

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
	"log"
)

func main() {
	getIp()
	createBridge()
	createVeth()
}

func getIp() {
	// 获取eth0的IP地址
	link, err := netlink.LinkByName("eth0")
	if err != nil {
		log.Fatalln(err)
	}
	addr, err := netlink.AddrList(link, netlink.FAMILY_V4)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(addr)
}

// 创建网桥
func createBridge() {
	br := &netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{
			Name: "mybr",
		},
	}
	if err := netlink.LinkAdd(br); err != nil {
		log.Fatalln(err)
	}

	// 设置网桥ip地址
	addr, err := netlink.ParseAddr("10.16.0.1/16")
	if err != nil {
		log.Fatalln(err)
	}
	if err = netlink.AddrAdd(br, addr); err != nil {
		log.Fatalln(err)
	}

	// 启用网桥
	if err = netlink.LinkSetUp(br); err != nil {
		log.Fatalln(err)
	}
}

// 创建veth设备对，手动设置容器网络
// 手动设置容器网络模式为none: docker run -d --name myngx --net=none nginx:1.18-alpine
func createVeth() {
	// 创建veth设备对
	vethpeer := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{
			Name: "myveth-host", // host端端设备名称
		},
		PeerName: "myveth-docker", // 容器端的设备名称
	}
	// ip link add
	if err := netlink.LinkAdd(vethpeer); err != nil {
		log.Fatalln(err)
	}

	// 启动host这一端的veth
	vethHost, err := netlink.LinkByName("myveth-host")
	if err != nil {
		log.Fatalln(err)
	}
	if err := netlink.LinkSetUp(vethHost); err != nil {
		log.Fatalln(err)
	}
	// 把host这一端的veth绑定到网桥上
	bridge, err := netlink.LinkByName("mybr")
	if err != nil {
		log.Fatalln(err)
	}
	if err := netlink.LinkSetMaster(vethHost, bridge.(*netlink.Bridge)); err != nil {
		log.Fatalln(err)
	}

	// 获取docker容器网络命名空间
	const pid = 310 // 查看docker容器pid: docker inspect xx | grep Pid
	getnetns, err := netns.GetFromPath(fmt.Sprintf("/proc/%d/ns/net", pid))
	if err != nil {
		log.Fatalln(err)
	}
	defer getnetns.Close()
	// 把veth设备其中一端移动到容器命名空间内
	vethDocker, err := netlink.LinkByName("myveth-docker")
	if err != nil {
		log.Fatalln(err)
	}
	if err = netlink.LinkSetNsFd(vethDocker, int(getnetns)); err != nil {
		log.Fatalln(err)
	}
	// 切换当前网络命名空间到容器环境
	if err = netns.Set(getnetns); err != nil {
		log.Fatalln(err)
	}
	// 切换后需要重新获取veth设备(因为要获取移动到容器内的veth，上面获取的是宿主机的veth)
	vethDocker, err = netlink.LinkByName("myveth-docker")
	if err != nil {
		log.Fatalln(err)
	}
	// 设置veth设备ip
	addr, err := netlink.ParseAddr("10.16.0.10/16")
	if err != nil {
		log.Fatalln(err)
	}
	if err = netlink.AddrAdd(vethDocker, addr); err != nil {
		log.Fatalln(err)
	}
	// 重新设置veth设备名称为eth0
	if err = netlink.LinkSetName(vethDocker, "eth0"); err != nil {
		log.Fatalln(err)
	}
	// 启动veth设备
	if err = netlink.LinkSetUp(vethDocker); err != nil {
		log.Fatalln(err)
	}
}
