package util

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

// LocalIP1 获取本机IP地址，多网卡时不要用这个方法，用下面的
func LocalIP1() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("not found")
}

func LocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println("cannot get current host ip")
		return "", err
	}
	addr := conn.LocalAddr().(*net.UDPAddr)
	ip := strings.Split(addr.String(), ":")[0]
	return ip, nil
}
