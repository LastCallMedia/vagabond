package util

import(
	"net"
	"os"
	"errors"
	"fmt"
)

func IpSliceContains(ips []net.IP, ip net.IP) bool {
	for _, v := range ips {
		if v.String() == ip.String() {
			return true
		}
	}
	return false
}

func DirExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err == nil {
		if stat.IsDir() {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("Exists, but not a directory %s", path))
	}
	if os.IsNotExist(err) { return false, nil }
	return true, nil
}