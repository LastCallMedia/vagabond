package util

import "net"

func IpSliceContains(ips []net.IP, ip net.IP) bool {
	for _, v := range ips {
		if v.String() == ip.String() {
			return true
		}
	}
	return false
}
