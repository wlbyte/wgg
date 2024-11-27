package util

import (
	"fmt"
	"log"
	"os/exec"
)

func CheckError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func CheckErrorPrintln(err error) {
	if err != nil {
		log.Println("error: ", err)
	}
}

func RunCmd(cmd string) (string, error) {
	execCmd := exec.Command("sh", "-c", cmd)
	result, err := execCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("RunCmd() error: %s", err)
	}
	return string(result), nil
}
