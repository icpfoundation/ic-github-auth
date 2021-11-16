package util

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Directory and File

func MkDir(path string) error {
	if _, err := os.ReadDir(path); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0777)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}

	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}

	return s.IsDir()
}

func IsFile(path string) bool {
	return !IsDir(path)
}

// Process

func RunCommand(cmd string) (string, error) {
	fmt.Println("Running Linux cmd:" + cmd)
	result, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(result)), err
}

func CheckProRunning(serverName string) (bool, error) {
	a := `ps ux | awk '/` + serverName + `/ && !/awk/ {print $2}'`
	pid, err := RunCommand(a)
	if err != nil {
		return false, err
	}
	return pid != "", nil
}

func GetPid(serverName string) ([]string, error) {
	a := `ps ux | awk '/` + serverName + `/ && !/awk/ {print $2}'`
	pid, err := RunCommand(a)
	pids := strings.Split(pid, "\n")
	for _, v := range pids {
		if v == "" {
			continue
		}
		pids = append(pids, v)
	}

	Infof("pids: %+v", pids)
	return pids, err
}

//////////////
func Format(v string) string {
	return fmt.Sprintf("[%s]	%s", time.Now().Format("2006-01-02 15:04:05.999"), v)
}
