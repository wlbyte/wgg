package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gitlab.eng.tethrnet.com/liulei/wgg/network"
	"gitlab.eng.tethrnet.com/liulei/wgg/util"
	"gitlab.eng.tethrnet.com/liulei/wgg/wggtypes"
	"gitlab.eng.tethrnet.com/liulei/wgg/wguser"
)

func setConfig(opts *cmdOptions) {
	if opts.Interface == "--help" || opts.Interface == "" || opts.Option == "" {
		showSubCommandUsage("setconf <interface> <configuration filename>", opts)
	}

	fin, err := os.Open(opts.Option)
	util.CheckError(err)
	defer fin.Close()
	cfg, err := loadConfig(fin)
	util.CheckError(err)
	client, err := wguser.New()
	util.CheckError(err)

	//wireguard接口不存在时创建接口
	if !network.DevExist(opts.Interface) {
		if err := network.StartWireguardGo(opts.Interface, cfg.IpAddress); err != nil {
			fmt.Printf("%s started failed, error: %s\n", opts.Interface, err)
			os.Exit(1)
		}
		fmt.Printf("%s started\n", opts.Interface)
	}

	err = client.ConfigureDevice(opts.Interface, *cfg)
	if err != nil {
		fmt.Printf("%s configured failed, error: %s\n", opts.Interface, err)
		os.Exit(1)
	}
	fmt.Printf("%s configured.\n", opts.Interface)

	// 配置路由
	fmt.Printf("%s route configured\n", opts.Interface)
	for _, peer := range cfg.Peers {
		for _, ipNet := range peer.AllowedIPs {
			network.ConfigRouteByIPNet(opts.Interface, "", "add", &ipNet)
		}
	}
	if cfg.Bond != nil {
		for _, ipNet := range cfg.Bond.AllowIPs {
			network.ConfigRouteByIPNet(opts.Interface, "", "add", &ipNet)
		}
	}
}

// Spec: https://git.zx2c4.com/WireGuard/about/src/tools/man/wg.8
// Original code: https://github.com/QuantumGhost/wg-quick-go/blob/master/internal/config/parser.go

type parseError struct {
	message string
	line    int
}

func (p parseError) Error() string {
	return fmt.Sprintf("Parse error: %s, (line %d)", p.message, p.line)
}

const (
	sectionInterface = "Interface"
	sectionPeer      = "Peer"
	sectionBond      = "bond"
	sectionEmpty     = ""
)

var (
	commentPattern = regexp.MustCompile(`#.*$`)
)

func matchSectionHeader(s string) (string, bool) {
	re := regexp.MustCompile(`\[(?P<section>\w+)\]`)
	matched := re.MatchString(s)
	if !matched {
		return "", false
	}
	sec := re.ReplaceAllString(s, "${section}")
	return sec, true
}

type pair struct {
	key   string
	value string
}

func matchKeyValuePair(s string) (pair, bool) {
	re := regexp.MustCompile(`^\s*(?P<key>\w+)\s*=\s*(?P<value>.+)\s*$`)
	matched := re.MatchString(s)
	if !matched {
		return pair{}, false
	}
	key := re.ReplaceAllString(s, "${key}")
	value := re.ReplaceAllString(s, "${value}")
	return pair{key: key, value: value}, true
}

func loadConfig(in io.Reader) (*wggtypes.Config, error) {
	sc := bufio.NewScanner(in)
	var cfg *wggtypes.Config = nil
	peers := make([]wggtypes.PeerConfig, 0, 10)
	currentSec := sectionEmpty
	var currentPeerConfig *wggtypes.PeerConfig = nil
	var currentBondConfig *wggtypes.BondConfig = nil

	for lineNum := 0; sc.Scan(); lineNum++ {
		line := sc.Text()
		line = commentPattern.ReplaceAllString(line, "")
		if strings.TrimSpace(line) == "" {
			// skip comment line
			continue
		}
		if sec, matched := matchSectionHeader(line); matched {
			if sec == sectionInterface {
				if cfg != nil {
					return nil, parseError{message: "duplicated Interface section", line: lineNum}
				}
				cfg = &wggtypes.Config{ReplacePeers: true}
			} else if sec == sectionPeer {
				if currentPeerConfig != nil {
					peers = append(peers, *currentPeerConfig)
				}
				currentPeerConfig = &wggtypes.PeerConfig{ReplaceAllowedIPs: true}
			} else if sec == sectionBond {
				if currentBondConfig != nil {
					return nil, parseError{message: "duplicated bond section", line: lineNum}
				} else {
					currentBondConfig = &wggtypes.BondConfig{SlavePeers: make([]wggtypes.Key, 0, 4)}
				}
			} else {
				return nil, parseError{message: fmt.Sprintf("Unknown section: %s", sec), line: lineNum}
			}
			currentSec = sec
			if currentSec == sectionInterface || currentSec == sectionBond {
				if currentPeerConfig != nil {
					peers = append(peers, *currentPeerConfig)
					currentPeerConfig = nil
				}
			}
			continue
		} else if pair, matched := matchKeyValuePair(line); matched {
			var perr *parseError
			if currentSec == sectionEmpty {
				return nil, parseError{message: "invalid top level key-value pair", line: lineNum}
			}
			if currentSec == sectionInterface {
				perr = parseInterfaceField(cfg, pair)
			} else if currentSec == sectionPeer {
				perr = parsePeerField(currentPeerConfig, pair)
			} else if currentSec == sectionBond {
				perr = parseBondField(currentBondConfig, pair)
			}
			if perr != nil {
				perr.line = lineNum
				return nil, perr
			}
		}
	}
	if currentSec == sectionPeer {
		peers = append(peers, *currentPeerConfig)
	}
	if cfg == nil {
		return nil, parseError{message: "no Interface section found"}
	}
	cfg.Peers = peers
	cfg.Bond = currentBondConfig
	return cfg, nil
}

func parseInterfaceField(cfg *wggtypes.Config, p pair) *parseError {
	switch p.key {
	case "PrivateKey":
		key, err := decodeKey(p.value)
		if err != nil {
			return err
		}
		cfg.PrivateKey = &key
	case "ListenPort":
		port, err := strconv.Atoi(p.value)
		if err != nil {
			return &parseError{message: err.Error()}
		}
		cfg.ListenPort = &port
	case "FwMark":
		return &parseError{message: "FwMark is not supported"}
	case "Address":
		cfg.IpAddress = strings.TrimSpace(p.value)
	case "DNS":
	default:
		return &parseError{message: fmt.Sprintf("invalid key %s for Interface section", p.key)}
	}
	return nil
}

func parsePeerField(cfg *wggtypes.PeerConfig, p pair) *parseError {
	switch p.key {
	case "PublicKey":
		key, err := decodeKey(p.value)
		if err != nil {
			return err
		}
		cfg.PublicKey = key
	case "PresharedKey":
		key, err := decodeKey(p.value)
		if err != nil {
			return err
		}
		cfg.PresharedKey = &key
	case "AllowedIPs":
		allowedIPs := make([]net.IPNet, 0, 10)
		splitted := strings.Split(p.value, ",")
		for _, seg := range splitted {
			seg = strings.TrimSpace(seg)
			ip, err := parseIPNet(seg)
			if err != nil {
				return err
			}
			allowedIPs = append(allowedIPs, *ip)
		}
		cfg.AllowedIPs = allowedIPs
	case "Endpoint":
		addr, err := net.ResolveUDPAddr("udp", p.value)
		if err != nil {
			return &parseError{message: err.Error()}
		}
		cfg.Endpoint = addr
	case "PersistentKeepalive":
		if p.value == "off" {
			cfg.PersistentKeepaliveInterval = nil
			return nil
		}
		sec, err := strconv.Atoi(p.value)
		if err != nil {
			return &parseError{message: err.Error()}
		}
		duration := time.Second * time.Duration(sec)
		cfg.PersistentKeepaliveInterval = &duration
	default:
		return &parseError{message: fmt.Sprintf("invalid key %s for Peer section", p.key)}
	}
	return nil
}

func parseBondField(cfg *wggtypes.BondConfig, p pair) *parseError {
	switch p.key {
	case "bondname":
		cfg.BondName = p.value
	case "bondmode":
		if ok, err := wggtypes.ValidBondMode(p.value); !ok {
			return &parseError{message: err.Error()}
		}
		cfg.BondMode = p.value
	case "bestslavepeer":
		key, err := decodeKey(p.value)
		if err != nil {
			return err
		}
		cfg.BestPeer = key
	case "slavepeer":
		key, err := decodeKey(p.value)
		if err != nil {
			return err
		}
		cfg.SlavePeers = append(cfg.SlavePeers, key)
	case "AllowIPs":
		allowedIPs := make([]net.IPNet, 0, 10)
		splitted := strings.Split(p.value, ",")
		for _, seg := range splitted {
			seg = strings.TrimSpace(seg)
			ip, err := parseIPNet(seg)
			if err != nil {
				return err
			}
			allowedIPs = append(allowedIPs, *ip)
		}
		cfg.AllowIPs = allowedIPs
	default:
		return &parseError{message: fmt.Sprintf("invalid key %s for Bond section", p.key)}
	}
	return nil
}

func decodeKey(s string) (wggtypes.Key, *parseError) {
	key, err := wggtypes.ParseKey(s)
	if err != nil {
		return wggtypes.Key{}, &parseError{message: err.Error()}
	}
	return key, nil
}

func parseIPNet(s string) (*net.IPNet, *parseError) {
	_, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, &parseError{message: err.Error()}
	}
	if ipnet == nil {
		return nil, &parseError{message: "invalid cidr string"}
	}
	return ipnet, nil
}
