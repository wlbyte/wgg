# wgg

A Golang implementation of the WireGuard [wg(8)](https://git.zx2c4.com/wireguard-tools/about/src/man/wg.8) utility.

This tool could be used to get and set the configuration of WireGuard tunnel interfaces.

It can be used in conjunction with [wireguard-go](https://git.zx2c4.com/wireguard-go/about/) for an almost complete userspace implementation of WireGuard on platforms which can be targeted by Go but do not have an implementation of WireGuard available.

`wgg` can also control a kernel-based WireGuard configuration.

For more information on WireGuard, please see https://www.wireguard.com/.

## Supported Sub-commands

This implementation supports the following sub-commands as specified in [wg(8)](https://git.zx2c4.com/wireguard-tools/about/src/man/wg.8):
```
  show:     Shows the current configuration and device information
  showconf: Shows the current configuration of a given WireGuard interface, for use with 'setconf'
  setconf:  Applies a configuration file to a WireGuard interface
  genkey:   Generates a new private key and writes it to stdout
  genpsk:   Generates a new preshared key and writes it to stdout
  pubkey:   Reads a private key from stdin and writes a public key to stdout
```

The `--version` command line option is also supported to show the release version.

### Script Wrapper

The `wg` script provides a convenient wrapper around `wgg` to provide a level of compatibility with the [wg(8)](https://git.zx2c4.com/wireguard-tools/about/src/man/wg.8) utility.

## How does this work?

This tool uses [wgctrl-go](https://github.com/WireGuard/wgctrl-go/) to enable control of WireGuard devices on multiple platforms.

## Building

This requires an installation of go ≥ 1.16.
```
git clone https://github.com/wlbyte/wgg.git
cd wgg
make
```

You can build the executable for different architectures and operating systems by setting the `GOOS`, `GOARCH`, and, if necessary, `GOARM` environment variables before running `make`, as specified in https://golang.org/doc/install/source#environment.

## Original Code

This project was inspired by and based upon [QuantumGhost/wg-quick-go](https://github.com/QuantumGhost/wg-quick-go).

