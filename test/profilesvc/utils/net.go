package utils

import (
	"net"
)

func GetNoLoopbackAddr()string  {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, iface := range ifaces {
		if (iface.Flags & net.FlagUp == net.FlagUp) && (iface.Flags & net.FlagLoopback != net.FlagLoopback) && (iface.Flags & net.FlagPointToPoint != net.FlagPointToPoint) {
			addrs, err := iface.Addrs()
			if err != nil {
				panic(err)
			}
			for _, addr := range addrs {

				if ipAddr, ok := addr.(*net.IPNet); ok {
					ipV4 := ipAddr.IP.To4()
					if ipV4 != nil {
						return ipV4.String()
					}
				}

				if ipAddr, ok := addr.(*net.IPAddr); ok {
					ipV4 := ipAddr.IP.To4()
					if ipV4 != nil {
						return ipV4.String()
					}
				}
			}
		}
	}
	return ""
}
