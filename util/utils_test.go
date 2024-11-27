package util

import "testing"

func TestRunCmd(t *testing.T) {
	type args struct {
		cmd string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"testRunCmd", args{"lsb_release -a"}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RunCmd(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				// t.Errorf("RunCmd() = %v, want %v", got, tt.want)
				t.Log(got)
			}
		})
	}
}
