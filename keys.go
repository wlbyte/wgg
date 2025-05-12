package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/wlbyte/wgg/util"
	"github.com/wlbyte/wgg/wggtypes"
)

func genKey(opts *cmdOptions) {
	if opts.Interface == "--help" || opts.Interface != "" || opts.Option != "" {
		showSubCommandUsage("genkey", opts)
	}

	key, err := wggtypes.GeneratePrivateKey()
	util.CheckError(err)
	fmt.Println(key.String())
}

func genPSK(opts *cmdOptions) {
	if opts.Interface == "--help" || opts.Interface != "" || opts.Option != "" {
		showSubCommandUsage("genpsk", opts)
	}

	key, err := wggtypes.GenerateKey()
	util.CheckError(err)
	fmt.Println(key.String())
}

func pubKey(opts *cmdOptions) {
	if opts.Interface == "--help" || opts.Interface != "" || opts.Option != "" {
		showSubCommandUsage("pubkey", opts)
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	util.CheckError(err)
	input = strings.TrimSpace(input)
	private, err := wggtypes.ParseKey(input)
	util.CheckError(err)
	public := private.PublicKey()
	fmt.Println(public.String())
}
