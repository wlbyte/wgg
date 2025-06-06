package main

import (
	"fmt"
	"os"
)

const (
	ENV_WG_COMMAND   = "WG_COMMAND"
	ENV_WG_HIDE_KEYS = "WG_HIDE_KEYS"
)

var (
	appVersion    = "dev"
	wgctrlVersion = "unknown"
)

func main() {
	opts := getOptions()

	switch opts.SubCommand {
	case "show":
		show(opts)
	case "showconf":
		showConfig(opts)
	case "setconf":
		setConfig(opts)
	case "genkey":
		genKey(opts)
	case "genpsk":
		genPSK(opts)
	case "pubkey":
		pubKey(opts)
	case "--version":
		showVersion(opts)
	default:
		fmt.Printf("Invalid subcommand: '%s'\n", opts.Command)
		showCommandUsage(1, opts)
	}
}

type cmdOptions struct {
	Command    string
	SubCommand string
	Interface  string
	Option     string
	ShowKeys   bool
}

func getOptions() *cmdOptions {
	args := len(os.Args[1:])
	base := 0
	opts := cmdOptions{}

	opts.ShowKeys = os.Getenv(ENV_WG_HIDE_KEYS) == "never"
	opts.Command = os.Getenv(ENV_WG_COMMAND)
	if opts.Command == "" {
		opts.Command = "wgg"
	}

	if args == 0 {
		opts.SubCommand = "show"
		opts.Interface = "all"
	} else if args == 1 && os.Args[base+1] == "--help" {
		showCommandUsage(0, &opts)
	} else if args > 3 {
		showCommandUsage(1, &opts)
	} else {
		opts.SubCommand = os.Args[base+1]
		if args >= 2 {
			opts.Interface = os.Args[base+2]
			if args == 3 {
				opts.Option = os.Args[base+3]
			}
		} else if opts.SubCommand == "show" {
			opts.Interface = "all"
		}
	}

	return &opts
}

func showCommandUsage(code int, opts *cmdOptions) {
	subcommands := `Available subcommands:
  show:     Shows the current configuration and device information
  showconf: Shows the current configuration of a given WireGuard interface, for use with 'setconf'
  setconf:  Applies a configuration file to a WireGuard interface
  genkey:   Generates a new private key and writes it to stdout
  genpsk:   Generates a new preshared key and writes it to stdout
  pubkey:   Reads a private key from stdin and writes a public key to stdout`

	fmt.Printf("Usage: %s <cmd> [<args>]\n\n", opts.Command)
	fmt.Printf("%s\n\n", subcommands)
	fmt.Println("You may pass '--help' to any of these subcommands to view showCommandUsage.")
	os.Exit(code)
}

func showSubCommandUsage(parameters string, opts *cmdOptions) {
	fmt.Printf("Usage: %s %s\n", opts.Command, parameters)
	if opts.Interface == "--help" {
		os.Exit(0)
	} else {
		os.Exit(2)
	}
}

func showVersion(opts *cmdOptions) {
	if opts.Interface == "--help" || opts.Interface != "" || opts.Option != "" {
		showSubCommandUsage("version", opts)
	}
	fmt.Printf("wgg v%s https://github.com/wlbyte/wgg (wgctrl %s)\n", appVersion, wgctrlVersion)
	os.Exit(0)
}
