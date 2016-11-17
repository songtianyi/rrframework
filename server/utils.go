package rrserver

import (
	"errors"
	"fmt"
	"net"
)

var (
	IP_PROTOCOL = "ipv4"
)

func getIpAddrByInterface(inf string) (error, string) {
	if len(inf) < 1 {
		return errors.New("Interface name is an empty string!"), inf
	}
	if net.ParseIP(inf) != nil {
		return nil, inf
	}
	infHandle, err := net.InterfaceByName(inf)
	if err != nil {
		return err, inf
	}
	infAddrs, err := infHandle.Addrs()
	if err != nil {
		return err, inf
	}
	for _, addr := range infAddrs {
		ipHandle, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			return err, inf
		}
		if IP_PROTOCOL == "ipv4" {
			if ipHandle.To4() != nil {
				return nil, ipHandle.To4().String()
			}
		} else if IP_PROTOCOL == "ipv6" {
			if ipHandle.To16() != nil {
				return nil, ipHandle.To16().String()
			}
		} else {
			return errors.New(fmt.Sprintf("Ip protocol [%s] not support", IP_PROTOCOL)), inf
		}
	}
	return errors.New(fmt.Sprintf("Failed when try to get ip address, [%s]", inf)), inf
}
