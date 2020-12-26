package ip

import (
	"errors"
	"math"
	"net"
	"strings"
)

// InternalIP return internal ip.
func InternalIP() string {
	inters, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, inter := range inters {
		if inter.Flags&net.FlagUp != 0 && !strings.HasPrefix(inter.Name, "lo") {
			addrs, err := inter.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}
	return ""
}

// LocalIP 获取机器的IP
func LocalIP() string {
	info, _ := net.InterfaceAddrs()
	for _, addr := range info {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}
	return ""
}

// IPString2Long 把ip字符串转为数值
func IPString2Long(ip string) (uint, error) {
	b := net.ParseIP(ip).To4()
	if b == nil {
		return 0, errors.New("invalid ipv4 format")
	}

	return uint(b[3]) | uint(b[2])<<8 | uint(b[1])<<16 | uint(b[0])<<24, nil
}

// Long2IPString 把数值转为ip字符串
func Long2IPString(i uint) (string, error) {
	if i > math.MaxUint32 {
		return "", errors.New("beyond the scope of ipv4")
	}

	ip := make(net.IP, net.IPv4len)
	ip[0] = byte(i >> 24)
	ip[1] = byte(i >> 16)
	ip[2] = byte(i >> 8)
	ip[3] = byte(i)

	return ip.String(), nil
}
