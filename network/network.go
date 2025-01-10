package network

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

// ConfigAddress set ip address for interface
// action must {add|del}
func ConfigAddress(devName, ipMask, action string) error {
	errStr := "ConfigIPAddress() error: "
	link, err := netlink.LinkByName(devName)
	if err != nil {
		return fmt.Errorf("%s%s", errStr, err)
	}
	_, ipNet, err := net.ParseCIDR(ipMask)
	if err != nil {
		return fmt.Errorf("%s%s", errStr, err)
	}
	addr := &netlink.Addr{
		IPNet: ipNet,
	}
	switch action {
	case "add":
		return netlink.AddrReplace(link, addr)
	case "del":
		return netlink.AddrDel(link, addr)
	default:
		return fmt.Errorf("%s%s", errStr, err)
	}
}

// ConfigRoute set route for interface
// action must be {add|del}
func ConfigRouteByStr(devName, gateway, ipMask, action string) error {
	errStr := "ConfigRouteByStr() error: "
	_, ipNet, err := net.ParseCIDR(ipMask)
	if err != nil {
		fmt.Errorf("%s%s", errStr, err)
	}
	return ConfigRouteByIPNet(devName, gateway, "add", ipNet)
}

func ConfigRouteByIPNet(devName, gateway, action string, ipNet *net.IPNet) error {
	errStr := "ConfigRouteByIPNet() error: "
	link, err := netlink.LinkByName(devName)
	if err != nil {
		return fmt.Errorf("%s%s", errStr, err)
	}
	route := &netlink.Route{
		LinkIndex: link.Attrs().Index,
		Dst:       ipNet,
	}
	if gateway != "" {
		route.Gw = net.ParseIP(gateway)
	}

	switch action {
	case "add":
		return netlink.RouteAdd(route)
	case "del":
		return netlink.RouteDel(route)
	default:
		return fmt.Errorf("%s%s", errStr, err)
	}
}

func findRouteByDev(devName string) ([]netlink.Route, error) {
	errStr := "findRouteByDev() error: "
	link, err := netlink.LinkByName(devName)
	if err != nil {
		return nil, fmt.Errorf("%s%s", errStr, err)
	}
	routeList, err := netlink.RouteList(link, 4)
	if err != nil {
		return nil, fmt.Errorf("%s%s", errStr, err)
	}
	return routeList, fmt.Errorf("%s%s", errStr, err)
}

func ConfigInterfaceState(devName, action string) error {
	errStr := "ConfigInterfaceState() error: "
	link, err := netlink.LinkByName(devName)
	if err != nil {
		return fmt.Errorf("%s%s", errStr, err)
	}
	switch action {
	case "up":
		return netlink.LinkSetUp(link)
	case "down":
		return netlink.LinkSetDown(link)
	default:
		return fmt.Errorf("%s%s", errStr, err)
	}
}

func ConfigInterface(devName, ipAddress string) error {
	err := ConfigInterfaceState(devName, "up")
	if err != nil {
		return fmt.Errorf("ConfigInterface() error: %s", err)
	}
	err = ConfigAddress(devName, ipAddress, "add")
	if err != nil {
		return fmt.Errorf("ConfigInterface() error: %s", err)
	}
	return nil
}

// GetAddress get ip address of interface, return first address
func GetAddressFirst(devName string) (string, error) {
	link, err := netlink.LinkByName(devName)
	if err != nil {
		return "", fmt.Errorf("failed to get device %s: %w", devName, err)
	}
	addrs, err := netlink.AddrList(link, netlink.FAMILY_V4)
	if err != nil {
		return "", fmt.Errorf("failed to get address of %s: %w", devName, err)
	}
	if len(addrs) > 0 {
		return addrs[0].IPNet.String(), nil
	}
	return "", nil
}

func GetInterfaceInfo(devName string) (string, error) {
	link, err := netlink.LinkByName(devName)
	if err != nil {
		return "", fmt.Errorf("GetInterfaceInfo() error: %s", err)
	}
	return fmt.Sprintf("interface: %s\n type: %s\n mtu: %d\n state: %d\n mac: %s\n", devName, link.Type(), link.Attrs().MTU, link.Attrs().OperState, link.Attrs().HardwareAddr), nil
}
