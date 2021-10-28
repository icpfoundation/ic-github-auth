package main

import (
	"context"
	"encoding/json"
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

func SaveAccessToken(ctx context.Context, state string, auth AuthResponse) error {
	key := datastore.NewKey(state)
	ishas, err := UserInfoDB.Has(ctx, key)
	if err != nil {
		Errorf("entrys: has %s", err)
		return err
	}

	if !ishas {
		in, err := json.Marshal(auth)
		if err != nil {
			return err
		}

		err = UserInfoDB.Put(ctx, key, in)
		if err != nil {
			Infof("entrys: begin %s", err)
			return err
		}
		Infof("write user info for state: %s", key.String())
	}

	return nil
}
