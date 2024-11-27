package network

import (
	"testing"
)

func TestConfigAddress(t *testing.T) {
	type args struct {
		devName string
		ipMask  string
		action  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"configIPAddress", args{"wg0", "100.71.192.5/32", "add"}, false},
		{"configIPAddress", args{"wg0", "100.71.192.5/32", "del"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ConfigAddress(tt.args.devName, tt.args.ipMask, tt.args.action); (err != nil) != tt.wantErr {
				t.Errorf("ConfigIPAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigRoute(t *testing.T) {
	type args struct {
		devName string
		gateway string
		ipMask  string
		action  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"configRoute", args{"wg0", "", "172.21.0.4/32", "add"}, false},
		{"configRoute", args{"wg0", "", "172.21.0.4/32", "del"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ConfigRoute(tt.args.devName, tt.args.gateway, tt.args.ipMask, tt.args.action); (err != nil) != tt.wantErr {
				t.Errorf("ConfigRoute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigInterfaceState(t *testing.T) {
	type args struct {
		devName string
		action  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"configInterfaceState", args{"wg0", "up"}, false},
		{"configInterfaceState", args{"wg0", "down"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ConfigInterfaceState(tt.args.devName, tt.args.action); (err != nil) != tt.wantErr {
				t.Errorf("ConfigInterfaceState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
