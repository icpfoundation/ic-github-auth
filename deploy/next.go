package deploy

import (
	"fmt"
	"os"

	"github.com/lyswifter/ic-auth/types"
)

// npm install
// npm run build
// npm run export
func DeployWithNext(path string, f *os.File, canister string, resource string, repo string, islocal bool, framework string) ([]types.CanisterInfo, error) {
	err := NpmInstall(path, f)
	if err != nil {
		return nil, err
	}

	err = NpmRunBuild(path, f)
	if err != nil {
		return nil, err
	}

	err = NpmRunExport(path, f)
	if err != nil {
		return nil, err
	}

	err = NewDfxjson(path, resource, canister)
	if err != nil {
		return nil, err
	}

	cinfos, err := DeployWithDfx(path, f, repo, islocal, framework, "", canister)
	if err != nil {
		return nil, err
	}

	fmt.Printf("canister infos: %+v", cinfos)

	return cinfos, nil
}
