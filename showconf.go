package main

import (
	"fmt"
	"strings"

	"github.com/wlbyte/wgg/network"
	"github.com/wlbyte/wgg/util"
	"github.com/wlbyte/wgg/wggtypes"
	"github.com/wlbyte/wgg/wguser"
)

func showConfig(opts *cmdOptions) {
	if opts.Interface == "--help" || opts.Interface == "" {
		showSubCommandUsage("showconf <interface>", opts)
	}
	// client, err := wgctrl.New()
	client, err := wguser.New()
	util.CheckError(err)
	dev, err := client.Device(opts.Interface)
	util.CheckError(err)
	address, err := network.GetAddressFirst(opts.Interface)
	util.CheckError(err)
	fmt.Printf("[Interface]\n")
	fmt.Printf("Address =  %s\n", address)
	fmt.Printf("ListenPort =  %d\n", dev.ListenPort)
	fmt.Printf("PrivateKey = %s\n", dev.PrivateKey.String())
	for _, peer := range dev.Peers {
		if dev.Bond != nil {
			showConfigPeers(peer, true)
		} else {
			showConfigPeers(peer, false)
		}

	}
	if dev.Bond != nil {
		showConfigBond(*dev.Bond)
	}
}

func showConfigPeers(peer wggtypes.Peer, bond bool) {
	psk := peer.PresharedKey.String()
	ka := peer.PersistentKeepaliveInterval.Seconds()

	fmt.Printf("\n[Peer]\n")
	fmt.Printf("PublicKey = %s\n", peer.PublicKey.String())
	if psk != "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=" {
		fmt.Printf("PresharedKey = %s\n", peer.PresharedKey.String())
	}
	if bond {
		fmt.Printf("#AllowedIPs = 0.0.0.0/0,::/0\n")
	} else {
		allowedIpStrings := make([]string, 0, len(peer.AllowedIPs))
		for _, v := range peer.AllowedIPs {
			allowedIpStrings = append(allowedIpStrings, v.String())
		}
		fmt.Printf("AllowedIPs = %s\n", strings.Join(allowedIpStrings, ", "))
	}
	fmt.Printf("Endpoint = %s\n", peer.Endpoint.String())
	if ka > 0 {
		fmt.Printf("PersistentKeepalive = %g\n", ka)
	}
}

func showConfigBond(bond wggtypes.BondConfig) {
	fmt.Printf("\n[bond]\n")
	fmt.Printf("bondname = %s\n", bond.BondName)
	fmt.Printf("bondmode = %s\n", bond.BondMode)
	if len(bond.BestPeer) > 0 {
		fmt.Printf("bestslavepeer = %s\n", bond.BestPeer.String())
	}
	for _, peer := range bond.SlavePeers {
		if peer != bond.BestPeer {
			fmt.Printf("slavepeer = %s\n", peer.String())
		}
	}
	allowedIpStrings := make([]string, 0, len(bond.AllowIPs))
	for _, v := range bond.AllowIPs {
		allowedIpStrings = append(allowedIpStrings, v.String())
	}
	fmt.Printf("AllowIPs = %s\n", strings.Join(allowedIpStrings, ", "))
}
