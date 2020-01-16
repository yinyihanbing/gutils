package gutils

import (
	"net"
	"strings"
)

// 获取域名对应IP
func GetIPByDomainName(domainName string) (ips []string, err error) {
	ns, err := net.LookupHost(domainName)
	if err != nil {
		return nil, err
	}
	return ns, nil
}

// 获取没有前缀的主机地址, 如: 参数=wss://127.0.0.1:80, 结果=127.0.0.1:80
func GetHostWithoutPrefix(addr string) string {
	idx := strings.Index(addr, "//")
	if idx != -1 {
		addr = addr[idx+2:]
	}

	idx = strings.Index(addr, ":")
	if idx != -1 {
		addr = addr[:idx]
	}

	idx = strings.Index(addr, "/")
	if idx != -1 {
		addr = addr[:idx]
	}

	return addr
}

// 获取本机IPV4地址
func GetLocalIPv4s() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}
	return ips, nil
}
