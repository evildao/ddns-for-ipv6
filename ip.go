package main

import (
	"errors"
	"net"
)

// GetAllIPv6 返回所有可用的公网 IPv6 地址
func GetAllIPv6() ([]string, error) {
	ipv6s := make([]string, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	// https://en.wikipedia.org/wiki/IPv6_address#General_allocation
	_, ipv6Unicast, _ := net.ParseCIDR("2000::/3")
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok &&
				len(ipNet.IP) == net.IPv6len &&
				ipNet.IP.IsGlobalUnicast() &&
				ipv6Unicast.Contains(ipNet.IP) {
				ipv6s = append(ipv6s, ipNet.IP.String())
			}
		}
	}
	if len(ipv6s) == 0 {
		return nil, errors.New("未找到可用的公网 IPv6 地址")
	}
	return ipv6s, nil
}
