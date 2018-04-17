package main

import (
	"net"
	"fmt"
)

func main() {
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
						fmt.Printf("ip: %s\n", ipV4.String())
					}
				}

				if ipAddr, ok := addr.(*net.IPAddr); ok {
					ipV4 := ipAddr.IP.To4()
					if ipV4 != nil {
						fmt.Printf("ip: %s\n", ipV4.String())
					}
				}
			}
		}
	}
}
