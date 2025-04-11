//go:build !linux

package network

import (
	"net"
)

// ConfigAddress set ip address for interface
// action must {add|del}
func ConfigAddress(devName, ipMask, action string) error {
	return nil
}

// ConfigRoute set route for interface
// action must be {add|del}
func ConfigRouteByStr(devName, gateway, ipMask, action string) error {
	return nil
}

func ConfigRouteByIPNet(devName, gateway, action string, ipNet *net.IPNet) error {
	return nil
}

func ConfigInterfaceState(devName, action string) error {
	return nil
}

func ConfigInterface(devName, ipAddress string) error {
	return nil
}

// GetAddress get ip address of interface, return first address
func GetAddressFirst(devName string) (string, error) {

	return "", nil
}

func GetInterfaceInfo(devName string) (string, error) {
	return "", nil
}
