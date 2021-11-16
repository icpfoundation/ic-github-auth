package db

import (
	"os"
	"path"

	"github.com/ipfs/go-datastore"
	levelds "github.com/ipfs/go-ds-leveldb"
	"github.com/lyswifter/ic-auth/params"
	"github.com/mitchellh/go-homedir"
	ldbopts "github.com/syndtr/goleveldb/leveldb/opt"
)

const RepoDB = "rinfo"
const UserDB = "uinfo"

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
		return nil, err
	}
	return db, err
}

func DataStores(repopath string) (datastore.Batching, datastore.Batching) {
	repodir, err := homedir.Expand(repopath)
	if err != nil {
		return nil, nil
	}

	rdb, err := setupLevelDs(path.Join(repodir, RepoDB), false)
	if err != nil {
		return nil, nil
	}

	idb, err := setupLevelDs(path.Join(repodir, UserDB), false)
	if err != nil {
		return nil, nil
	}

	return rdb, idb
}

func LoadRepoInfoDB(dbName string) (datastore.Batching, error) {
	repodir, err := homedir.Expand(params.RepoPath)
	if err != nil {
		return nil, err
	}

	db, err := setupLevelDs(path.Join(repodir, dbName), false)
	if err != nil {
		return nil, err
	}

	return db, nil
}
