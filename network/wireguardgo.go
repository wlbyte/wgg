package network

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"gitlab.eng.tethrnet.com/liulei/wgg/wguser"
)

func StartWireguardGo(devName, ipAddress string) error {
	logDir := "/var/log/wireguard/"
	logFile := logDir + "wireguard.log"
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		return fmt.Errorf("startWireguardGo() create dir %s failed, error: %s", logDir, err)
	}
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("startWireguardGo() open log error: %s", err)
	}
	defer f.Close()

	cmd := exec.Command("wireguard-go", devName)
	cmd.Stdout = f
	cmd.Stderr = f
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true, // 分离子进程
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("startWireguardGo() error: %s", err)
	}
	for {
		if !DevExist(devName) {
			fmt.Println(devName, devName+".sock not ready, waiting...")
			time.Sleep(time.Second)
			continue
		}
		break
	}
	
	return nil
}

func DevExist(devName string) bool {
	client, err := wguser.New()
	if err != nil {
		fmt.Printf("FindDev() new client failed, error: %s", err)
	}
	if _, err = client.Device(devName); err != nil {
		return false
	}
	return true
}
