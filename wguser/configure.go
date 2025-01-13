package wguser

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"gitlab.eng.tethrnet.com/liulei/wgg/wggtypes"
)

// configureDevice configures a device specified by its path.
func (c *Client) configureDevice(device string, cfg wggtypes.Config) error {
	conn, err := c.dial(device)
	if err != nil {
		// fmt.Printf("configureDevice() error:%s", err)
		return err
	}
	defer conn.Close()

	// Start with set command.
	var buf bytes.Buffer
	buf.WriteString("set=1\n")

	// Add any necessary configuration from cfg, then finish with an empty line.
	writeConfig(&buf, cfg)
	buf.WriteString("\n")

	// Apply configuration for the device and then check the error number.
	if _, err := io.Copy(conn, &buf); err != nil {
		return err
	}

	res := make([]byte, 32)
	n, err := conn.Read(res)
	if err != nil {
		return err
	}

	// errno=0 indicates success, anything else returns an error number that
	// matches definitions from errno.h.
	str := strings.TrimSpace(string(res[:n]))
	if str != "errno=0" {
		// TODO(mdlayher): return actual errno on Linux?
		return os.NewSyscallError("read", fmt.Errorf("wguser: %s", str))
	}

	return nil
}

// writeConfig writes textual configuration to w as specified by cfg.
func writeConfig(w io.Writer, cfg wggtypes.Config) {
	if cfg.PrivateKey != nil {
		fmt.Fprintf(w, "private_key=%s\n", hexKey(*cfg.PrivateKey))
	}
	
	if cfg.ListenPort != nil {
		fmt.Fprintf(w, "listen_port=%d\n", *cfg.ListenPort)
	}

	if cfg.FirewallMark != nil {
		fmt.Fprintf(w, "fwmark=%d\n", *cfg.FirewallMark)
	}

	if cfg.ReplacePeers {
		fmt.Fprintln(w, "replace_peers=true")
	}

	for _, p := range cfg.Peers {
		fmt.Fprintf(w, "public_key=%s\n", hexKey(p.PublicKey))

		if p.Remove {
			fmt.Fprintln(w, "remove=true")
		}

		if p.UpdateOnly {
			fmt.Fprintln(w, "update_only=true")
		}

		if p.PresharedKey != nil {
			fmt.Fprintf(w, "preshared_key=%s\n", hexKey(*p.PresharedKey))
		}

		if p.Endpoint != nil {
			fmt.Fprintf(w, "endpoint=%s\n", p.Endpoint.String())
		}

		if p.PersistentKeepaliveInterval != nil {
			fmt.Fprintf(w, "persistent_keepalive_interval=%d\n", int(p.PersistentKeepaliveInterval.Seconds()))
		}

		if p.ReplaceAllowedIPs {
			fmt.Fprintln(w, "replace_allowed_ips=true")
		}

		for _, ip := range p.AllowedIPs {
			fmt.Fprintf(w, "allowed_ip=%s\n", ip.String())
		}
	}
	// configure bond
	bondCfg := cfg.Bond
	if bondCfg != nil {
		if bondCfg.BondName != "" {
			fmt.Fprintf(w, "bond_name=%s\n", bondCfg.BondName)
		}
		if bondCfg.BondMode != "" {
			fmt.Fprintf(w, "bond_mode=%s\n", bondCfg.BondMode)
		}
		if len(bondCfg.BestPeer) > 0 {
			fmt.Fprintf(w, "best_slave_peer=%s\n", hexKey(bondCfg.BestPeer))
		}
		for _, slavePeer := range bondCfg.SlavePeers {
			fmt.Fprintf(w, "slave_peer=%s\n", hexKey(slavePeer))
		}
		for _, ip := range bondCfg.AllowIPs {
			fmt.Fprintf(w, "allowed_ip=%s\n", ip.String())
		}
	}
}

// hexKey encodes a wggtypes.Key into a hexadecimal string.
func hexKey(k wggtypes.Key) string {
	return hex.EncodeToString(k[:])
}
