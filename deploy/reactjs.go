package deploy

import (
	"fmt"
	"os"
)

func DeployWithReactjs(path string, f *os.File, canister string, resource string, repo string, islocal bool, framework string) error {
	err := NpmInstall(path, f)
	if err != nil {
		return err
	}

	err = NpmRunBuild(path, f)
	if err != nil {
		return err
	}

	err = NewDfxjson(path, resource, canister)
	if err != nil {
		return err
	}

	cinfos, err := DeployWithDfx(path, f, repo, islocal, framework)
	if err != nil {
		return err
	}

	fmt.Printf("canister infos: %+v", cinfos)

	return nil
}
