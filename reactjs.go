package main

import "os"

func deployWithReactjs(path string, f *os.File, canister string, resource string) error {
	err := npmInstall(path, f)
	if err != nil {
		return err
	}

	err = npmRunBuild(path, f)
	if err != nil {
		return err
	}

	err = newDfxjson(path, resource, canister)
	if err != nil {
		return err
	}

	err = deployWithDfx(path, f)
	if err != nil {
		return err
	}

	return nil
}
