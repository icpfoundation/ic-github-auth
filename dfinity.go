package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func getController(targetpath string) (string, error) {
	// dfx wallet --network ic addresses

	idcmd := exec.Command("dfx", "wallet", "--network", "ic", "address")
	idcmd.Dir = targetpath

	var b bytes.Buffer
	idcmd.Stderr = &b
	idcmd.Stdout = &b

	err := idcmd.Run()
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func deployWithDfx(targetpath string, f *os.File, repo string, islocal bool, framework string) error {

	var deploycmd *exec.Cmd
	if islocal {
		deploycmd = exec.Command("dfx", "deploy")
	} else {
		deploycmd = exec.Command("dfx", "deploy", "--network", "ic")
	}

	deploycmd.Dir = targetpath

	stderr, err := deploycmd.StderrPipe()
	if err != nil {
		return err
	}

	stdout, err := deploycmd.StdoutPipe()
	if err != nil {
		return err
	}

	deploycmd.Start()

	errReader := bufio.NewReader(stderr)
	outReader := bufio.NewReader(stdout)

	for {
		line, err := errReader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}

		// write local
		_, err = f.WriteString(Format(line))
		if err != nil {
			break
		}
	}

	for {
		line, err := outReader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}

		// write local
		_, err = f.WriteString(Format(line))
		if err != nil {
			break
		}
	}

	deploycmd.Wait()

	type CanisterID struct {
		ID string `json:"id"`
	}

	controller, err := getController(targetpath)
	if err != nil {
		return err
	}

	//read canister id
	cinfofile := filepath.Join(targetpath, "canister_ids.json")

	if Exists(cinfofile) {
		ret, err := os.ReadFile(cinfofile)
		if err != nil {
			return err
		}

		var infos map[string]CanisterID = make(map[string]CanisterID)
		err = json.Unmarshal(ret, &infos)
		if err != nil {
			return err
		}

		fmt.Printf("canister info map: %+v", infos)

		var network string = "ic"
		if islocal {
			network = "local"
		}

		for k, v := range infos {

			var ctype = "asssets"
			if framework == "dfx" && !strings.Contains(k, "assets") {
				ctype = "other"
			}

			cinfo := CanisterInfo{
				Repository:   repo,
				Controller:   controller,
				CanisterName: k,
				CanisterID:   v.ID,
				CanisterType: ctype,
				Framework:    framework,
				Network:      network,
			}

			err = SaveCanisterInfo(context.TODO(), cinfo)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
