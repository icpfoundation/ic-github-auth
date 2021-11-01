package main

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/ipfs/go-datastore"
	levelds "github.com/ipfs/go-ds-leveldb"
	"github.com/mitchellh/go-homedir"
	ldbopts "github.com/syndtr/goleveldb/leveldb/opt"
)

var InfoDB datastore.Batching
var UserInfoDB datastore.Batching

func setupLevelDs(path string, readonly bool) (datastore.Batching, error) {
	if _, err := os.ReadDir(path); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(path, 0777)
			if err != nil {
				return nil, err
			}
		}
	}

	db, err := levelds.NewDatastore(path, &levelds.Options{
		Compression: ldbopts.NoCompression,
		NoSync:      false,
		Strict:      ldbopts.StrictAll,
		ReadOnly:    readonly,
	})
	if err != nil {
		Errorf("NewDatastore: %s", err)
		return nil, err
	}
	return db, err
}

func DataStores() {
	repodir, err := homedir.Expand(repoPath)
	if err != nil {
		return
	}

	ldb, err := setupLevelDs(repodir, false)
	if err != nil {
		Errorf("setup beacondb: err %s", err)
		return
	}
	InfoDB = ldb
	Infof("InfoDB: %+v", InfoDB)

	idb, err := setupLevelDs(path.Join(repodir, "uinfo"), false)
	if err != nil {
		Errorf("setup infodb: err %s", err)
		return
	}
	UserInfoDB = idb
	Infof("UserInfoDB: %+v", UserInfoDB)
}

func SaveAccessToken(ctx context.Context, state string, authResponse []byte) error {
	key := datastore.NewKey(state)
	ishas, err := UserInfoDB.Has(ctx, key)
	if err != nil {
		return err
	}

	if !ishas {
		err = UserInfoDB.Put(ctx, key, authResponse)
		if err != nil {
			return err
		}
		Infof("write user info for state: %s", key.String())
	} else {
		Infof("token is already exist for; %s", key.String())
	}

	return nil
}

func ReadAccessToken(ctx context.Context, state string) ([]byte, error) {
	key := datastore.NewKey(state)
	ishas, err := UserInfoDB.Has(ctx, key)
	if err != nil {
		return nil, err
	}

	if ishas {
		token, err := UserInfoDB.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		return token, nil
	}

	return nil, nil
}

func SaveInstallCode(ctx context.Context, state string, code string) error {
	key := datastore.NewKey(fmt.Sprintf("%s-installationCode", state))

	err := UserInfoDB.Put(ctx, key, []byte(code))
	if err != nil {
		return err
	}
	Infof("write installation code for state: %s", key.String())
	return nil
}

func ReadInstallCode(ctx context.Context, state string) ([]byte, error) {
	key := datastore.NewKey(fmt.Sprintf("%s-installationCode", state))
	ishas, err := UserInfoDB.Has(ctx, key)
	if err != nil {
		return nil, err
	}

	if ishas {
		code, err := UserInfoDB.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		return code, nil
	}

	return nil, nil
}
