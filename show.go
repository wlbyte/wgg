package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/wlbyte/wgg/util"
	"github.com/wlbyte/wgg/wggtypes"
	"github.com/wlbyte/wgg/wguser"
)

func show(opts *cmdOptions) {
	if opts.Interface == "--help" || (opts.Interface == "interfaces" && opts.Option != "") || !(opts.Option == "" || opts.Option == "public-key" || opts.Option == "private-key" || opts.Option == "listen-port" || opts.Option == "fwmark" || opts.Option == "peers" || opts.Option == "preshared-keys" || opts.Option == "endpoints" || opts.Option == "allowed-ips" || opts.Option == "latest-handshakes" || opts.Option == "transfer" || opts.Option == "persistent-keepalive" || opts.Option == "dump") {
		showSubCommandUsage("show { <interface> | all | interfaces } [public-key | private-key | listen-port | fwmark | peers | preshared-keys | endpoints | allowed-ips | latest-handshakes | transfer | persistent-keepalive | dump]", opts)
	}

	// client, err := wgctrl.New()
	// util.CheckError(err)
	client, err := wguser.New()
	if err != nil {
		fmt.Println("show() new client failed, error:", err)
	}
	switch opts.Interface {
	case "interfaces":
		devices, err := client.Devices()
		if err != nil {
			fmt.Println("show() get Devices failed, error:", err)
			os.Exit(1)
		}
		for i := 0; i < len(devices); i++ {
			fmt.Println(devices[i].Name)
		}
	case "all":
		devices, err := client.Devices()
		if err != nil {
			fmt.Println("show() get Devices failed, error:", err)
			os.Exit(1)
		}
		for _, dev := range devices {
			showDevice(*dev, opts)
		}
	default:
		dev, err := client.Device(opts.Interface)
		if err != nil {
			fmt.Println("show() get Device failed, error:", err)
			os.Exit(1)
		}
		showDevice(*dev, opts)
	}
	client.Close()
}

func showDevice(dev wggtypes.Device, opts *cmdOptions) {
	if opts.Option == "" {
		showKeys := opts.ShowKeys
		fmt.Printf("interface: %s (%s)\n", dev.Name, dev.Type.String())
		fmt.Printf("  public key: %s\n", dev.PublicKey.String())
		fmt.Printf("  private key: %s\n", formatKey(dev.PrivateKey, showKeys))
		fmt.Printf("  listening port: %d\n", dev.ListenPort)
		fmt.Println()
		for _, peer := range dev.Peers {
			if dev.Bond != nil {
				showPeers(peer, showKeys)
			} else {
				showPeers(peer, showKeys)
			}

		}
		if dev.Bond != nil {
			showBond(*dev.Bond)
		}
	} else {
		deviceName := ""
		if opts.Interface == "all" {
			deviceName = dev.Name + "\t"
		}
		switch opts.Option {
		case "public-key":
			fmt.Printf("%s%s\n", deviceName, dev.PublicKey.String())
		case "private-key":
			fmt.Printf("%s%s\n", deviceName, dev.PrivateKey.String())
		case "listen-port":
			fmt.Printf("%s%d\n", deviceName, dev.ListenPort)
		case "fwmark":
			fmt.Printf("%s%d\n", deviceName, dev.FirewallMark)
		case "peers":
			for _, peer := range dev.Peers {
				fmt.Printf("%s%s\n", deviceName, peer.PublicKey.String())
			}
		case "preshared-keys":
			for _, peer := range dev.Peers {
				fmt.Printf("%s%s\t%s\n", deviceName, peer.PublicKey.String(), formatPSK(peer.PresharedKey, "(none)"))
			}
		case "endpoints":
			for _, peer := range dev.Peers {
				fmt.Printf("%s%s\t%s\n", deviceName, peer.PublicKey.String(), formatEndpoint(peer.Endpoint))
			}
		case "allowed-ips":
			for _, peer := range dev.Peers {
				fmt.Printf("%s%s\t%s\n", deviceName, peer.PublicKey.String(), joinIPs(peer.AllowedIPs))
			}
		case "latest-handshakes":
			for _, peer := range dev.Peers {
				fmt.Printf("%s%s\t%d\n", deviceName, peer.PublicKey.String(), peer.LastHandshakeTime.Unix())
			}
		case "transfer":
			for _, peer := range dev.Peers {
				fmt.Printf("%s%s\t%d\t%d\n", deviceName, peer.PublicKey.String(), peer.ReceiveBytes, peer.TransmitBytes)
			}
		case "persistent-keepalive":
			for _, peer := range dev.Peers {
				fmt.Printf("%s%s\t%s\n", deviceName, peer.PublicKey.String(), zeroToOff(strconv.FormatFloat(peer.PersistentKeepaliveInterval.Seconds(), 'g', 0, 64)))
			}
		case "dump":
			fmt.Printf("%s%s\t%s\t%d\t%s\n", deviceName, dev.PrivateKey.String(), dev.PublicKey.String(), dev.ListenPort, zeroToOff(strconv.FormatInt(int64(dev.FirewallMark), 10)))
			for _, peer := range dev.Peers {
				fmt.Printf("%s%s\t%s\t%s\t%s\t%d\t%d\t%d\t%s\n",
					deviceName,
					peer.PublicKey.String(),
					formatPSK(peer.PresharedKey, "(none)"),
					formatEndpoint(peer.Endpoint),
					joinIPs(peer.AllowedIPs),
					peer.LastHandshakeTime.Unix(),
					peer.ReceiveBytes,
					peer.TransmitBytes,
					zeroToOff(strconv.FormatFloat(peer.PersistentKeepaliveInterval.Seconds(), 'g', 0, 64)))
			}
		}
	}
}

func showPeers(peer wggtypes.Peer, showKeys bool) {

	tmpl := `peer: {{ .PublicKey }}
  endpoint: {{ .Endpoint }}
  allowed ips: {{ .AllowedIPs }}
  {{- if .PresharedKey}}
  preshared key: {{ .PresharedKey }}
  {{- end}}
  keep alive interval: {{ .KeepAliveInterval }}s
  last handshake time: {{ .LastHandshakeTime }}
  transfer: {{ .ReceiveBytes }} bytes received, {{ .TransmitBytes }} bytes sent
  protocol version: {{ .ProtocolVersion }} 

`

	type tmplContent struct {
		PublicKey         string
		PresharedKey      string
		Endpoint          string
		KeepAliveInterval float64
		LastHandshakeTime string
		ReceiveBytes      int64
		TransmitBytes     int64
		AllowedIPs        string
		ProtocolVersion   int
	}

	t := template.Must(template.New("peer_tmpl").Parse(tmpl))
	c := tmplContent{
		PublicKey:         peer.PublicKey.String(),
		PresharedKey:      formatPSK(peer.PresharedKey, ""),
		Endpoint:          formatEndpoint(peer.Endpoint),
		KeepAliveInterval: peer.PersistentKeepaliveInterval.Seconds(),
		LastHandshakeTime: peer.LastHandshakeTime.Format(time.RFC3339),
		ReceiveBytes:      peer.ReceiveBytes,
		TransmitBytes:     peer.TransmitBytes,
		AllowedIPs:        joinIPs(peer.AllowedIPs),
		ProtocolVersion:   peer.ProtocolVersion,
	}

	err := t.Execute(os.Stdout, c)
	util.CheckError(err)
}

func showBond(bond wggtypes.BondConfig) {
	fmt.Printf("bond: %s\n", bond.BondName)
	fmt.Printf("  bond mode: %s\n", bond.BondMode)
	if len(bond.BestPeer) > 0 {
		fmt.Printf("  best peer: %s\n", bond.BestPeer.String())
	}
	for _, peer := range bond.SlavePeers {
		if peer != bond.BestPeer {
			fmt.Printf("  slave peer: %s\n", peer.String())
		}
	}
	fmt.Printf("  active peer: %s\n", bond.ActivePeer)
	fmt.Printf("  transfer: %.2f KiB received, %.2f KiB sent\n", float64(bond.ReceiveBytes/1024), float64(bond.TransmitBytes/1024))
	allowedIpStrings := make([]string, 0, len(bond.AllowIPs))
	for _, v := range bond.AllowIPs {
		allowedIpStrings = append(allowedIpStrings, v.String())
	}
	fmt.Printf("  allowed ips: %s\n", strings.Join(allowedIpStrings, ","))
}

func formatEndpoint(endpoint *net.UDPAddr) string {
	ip := endpoint.String()
	if ip == "<nil>" {
		ip = "(none)"
	}
	return ip
}

func formatKey(key wggtypes.Key, showKeys bool) string {
	k := "(hidden)"
	if showKeys {
		k = key.String()
	}
	return k
}

func formatPSK(key wggtypes.Key, none string) string {
	psk := key.String()
	if psk == "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=" {
		return none
	}
	return psk
}

func joinIPs(ips []net.IPNet) string {
	ipStrings := make([]string, 0, len(ips))
	for _, v := range ips {
		ipStrings = append(ipStrings, v.String())
	}
	return strings.Join(ipStrings, ", ")
}

func zeroToOff(value string) string {
	if value == "0" {
		return "off"
	}
	return value
}
