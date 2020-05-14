package gutils

import (
	"github.com/oschwald/geoip2-golang"
	"net"
)

// IP地址帮助类结构体
type IpAddressHelper struct {
	dbReader *geoip2.Reader
}

// 实例化帮助类
func NewIpAddressHelper(dbPath string) (*IpAddressHelper, error) {
	reader, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, err
	}

	ipAddressHelper := IpAddressHelper{}
	ipAddressHelper.dbReader = reader

	return &ipAddressHelper, nil
}

// 获取国家
func (this *IpAddressHelper) GetCountry(searchIp string) (*geoip2.Country, error) {
	ip := net.ParseIP(searchIp)
	record, err := this.dbReader.Country(ip)
	if err != nil {
		return nil, err
	}
	return record, nil
}
