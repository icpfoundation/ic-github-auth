package deploy

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/lyswifter/ic-auth/types"
	"github.com/lyswifter/ic-auth/util"
	"github.com/shirou/gopsutil/v3/process"
	"golang.org/x/xerrors"
)

var noProcessErr = xerrors.New("process not found")

func KillProcess(name string) error {
	processes, err := process.Processes()
	if err != nil {
		return err
	}

	var noProcess = true
	for _, p := range processes {
		n, err := p.Name()
		if err != nil {
			return err
		}

		if n == name || strings.Contains(n, name) {
			noProcess = false
			err = p.Kill()
			if err != nil {
				continue
			}
		}
	}

	if noProcess {
		return noProcessErr
	}

	return nil
}

// dfx wallet --network ic addresses
func getController(targetpath string, islocal bool) (string, error) {
	var idcmd *exec.Cmd
	if islocal {
		idcmd = exec.Command("dfx", "identity", "get-wallet")
	} else {
		idcmd = exec.Command("dfx", "identity", "--network", "ic", "get-wallet")
	}

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

func restartDfx(targetpath string, islocal bool) error {
	if islocal {
		_ = KillProcess("dfx")

		startCmd := exec.Command("dfx", "start", "--background", "--clean")
		startCmd.Dir = targetpath
		startCmd.Stderr = os.Stderr
		startCmd.Stdout = os.Stdout
		err := startCmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func DeployWithDfx(targetpath string, f *os.File, repo string, islocal bool, framework string, buildcmd string) ([]types.CanisterInfo, error) {

	var deploycmd *exec.Cmd
	if islocal {
		err := restartDfx(targetpath, islocal)
		if err != nil {
			return nil, err
		}

		if buildcmd != "" {
			deploycmd = exec.Command("bash", "-c", buildcmd)
		} else {
			deploycmd = exec.Command("dfx", "deploy")
		}
	} else {
		if buildcmd != "" {
			idx := 10
			buildcmdIC := buildcmd[:idx] + " --network ic" + buildcmd[idx:]
			deploycmd = exec.Command("bash", "-c", buildcmdIC)
		} else {
			deploycmd = exec.Command("dfx", "deploy", "--network", "ic", "--with-cycles", "10000")
		}
	}

	fmt.Printf("deploycmd: %v\n", deploycmd)

	deploycmd.Dir = targetpath

	stderr, err := deploycmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := deploycmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	deploycmd.Start()

	errReader := bufio.NewReader(stderr)
	outReader := bufio.NewReader(stdout)

	canisterName := []string{}
	canisterId := []string{}

	var hasAlreadyCreated bool = false

	for {
		line, err := errReader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}

		if strings.Contains(line, "canister_id") {
			name, id, err := extractCanisterInfo(line)
			if err != nil {
				break
			}

			canisterName = append(canisterName, name)
			canisterId = append(canisterId, id)
		}

		if strings.Contains(line, "All canisters have already been created.") {
			hasAlreadyCreated = true
		}

		// write local
		_, err = f.WriteString(util.Format(line))
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

		if strings.Contains(line, "canister_id") {
			name, id, err := extractCanisterInfo(line)
			if err != nil {
				break
			}

			canisterName = append(canisterName, name)
			canisterId = append(canisterId, id)
		}

		if strings.Contains(line, "All canisters have already been created.") {
			hasAlreadyCreated = true
		}

		// write local
		_, err = f.WriteString(util.Format(line))
		if err != nil {
			break
		}
	}

	type CanisterID struct {
		IC string `json:"ic"`
	}

	controller, err := getController(targetpath, islocal)
	if err != nil {
		return nil, err
	}

	var network string = "ic"
	if islocal {
		network = "local"
	}

	//read canister id
	cinfofile := filepath.Join(targetpath, "canister_ids.json")

	fmt.Printf("canister ids file path: %s\n", cinfofile)

	cinfos := []types.CanisterInfo{}

	if util.Exists(cinfofile) {
		ret, err := os.ReadFile(cinfofile)
		if err != nil {
			return nil, err
		}

		var infos map[string]CanisterID = make(map[string]CanisterID)
		err = json.Unmarshal(ret, &infos)
		if err != nil {
			return nil, err
		}

		fmt.Printf("canister info map: %+v", infos)

		for k, v := range infos {
			var ctype = "asssets"
			if framework == "dfx" && !strings.Contains(k, "assets") {
				ctype = "other"
			}

			cinfo := types.CanisterInfo{
				Repo:         repo,
				Controller:   controller,
				CanisterName: k,
				CanisterID:   v.IC,
				CanisterType: ctype,
				Framework:    framework,
				Network:      network,
			}

			if hasAlreadyCreated {
				cinfo.UpdateTimestamp = time.Now().Unix()
			} else {
				cinfo.CreateTimestamp = time.Now().Unix()
				cinfo.UpdateTimestamp = time.Now().Unix()
			}
			cinfos = append(cinfos, cinfo)
		}
	} else {
		for i, v := range canisterName {
			id := canisterId[i]
			name := v

			cinfo := types.CanisterInfo{
				Repo:         repo,
				Controller:   controller,
				CanisterName: name,
				CanisterID:   id,
				CanisterType: "",
				Framework:    framework,
				Network:      network,
			}

			if hasAlreadyCreated {
				cinfo.UpdateTimestamp = time.Now().Unix()
			} else {
				cinfo.CreateTimestamp = time.Now().Unix()
				cinfo.UpdateTimestamp = time.Now().Unix()
			}
			cinfos = append(cinfos, cinfo)
		}
	}

	deploycmd.Wait()

	return cinfos, nil
}

func extractCanisterInfo(input string) (string, string, error) {
	var split = "with canister_id"
	arr := strings.Split(input, split)

	first := strings.TrimSpace(arr[0])
	last := strings.TrimSpace(arr[1])

	first = strings.TrimSuffix(first, ",")
	firstarr := strings.Split(first, " ")
	canisterName := firstarr[len(firstarr)-1]

	fmt.Printf("input: %s canister name: %s canister id: %s", input, canisterName, last)
	return canisterName, last, nil
}
