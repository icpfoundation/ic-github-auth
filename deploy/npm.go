package deploy

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"text/template"

	"github.com/lyswifter/ic-auth/util"
)

type Dfxjson struct {
	CanisterName string
	ResourcePath string
}

const templateText = `{
	"canisters": {
	  "{{ .CanisterName }}": {
		"type": "assets",
		"source": ["{{ .ResourcePath }}"]
	  }
	}
  }`

func NewDfxjson(targetpath string, source string, canistername string) error {
	//generate dfxjson
	dfxjson := Dfxjson{
		CanisterName: canistername,
		ResourcePath: source,
	}

	var bytew bytes.Buffer
	tpl := template.Must(template.New("anyname").Parse(templateText))

	err := tpl.Execute(&bytew, dfxjson)
	if err != nil {
		return err
	}

	fmt.Printf("dfxjson: %s %+v\n", bytew.String(), dfxjson)
	f, err := os.Create(path.Join(targetpath, "dfx.json"))
	if err != nil {
		return err
	}

	_, err = f.Write(bytew.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func NpmInstall(targetpath string, f *os.File) error {
	//npm install
	installCmd := exec.Command("npm", "install")
	installCmd.Dir = targetpath

	stderr, err := installCmd.StderrPipe()
	if err != nil {
		return err
	}

	stdout, err := installCmd.StdoutPipe()
	if err != nil {
		return err
	}

	installCmd.Start()

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

		// write local
		_, err = f.WriteString(util.Format(line))
		if err != nil {
			break
		}
	}

	installCmd.Wait()

	return nil
}

func NpmRunBuild(targetpath string, f *os.File) error {
	//npm install
	npmBuildCmd := exec.Command("npm", "run", "build")
	npmBuildCmd.Dir = targetpath

	stderr, err := npmBuildCmd.StderrPipe()
	if err != nil {
		return err
	}

	stdout, err := npmBuildCmd.StdoutPipe()
	if err != nil {
		return err
	}

	npmBuildCmd.Start()

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

		// write local
		_, err = f.WriteString(util.Format(line))
		if err != nil {
			break
		}
	}

	npmBuildCmd.Wait()

	return nil
}
