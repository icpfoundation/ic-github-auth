package db

import (
	"os"
	"path"

	"github.com/ipfs/go-datastore"
	levelds "github.com/ipfs/go-ds-leveldb"
	"github.com/mitchellh/go-homedir"
	ldbopts "github.com/syndtr/goleveldb/leveldb/opt"
	"golang.org/x/xerrors"
)

const RepoDB = "rinfo"
const UserDB = "uinfo"

type AuthDB struct {
	RepoDb datastore.Batching
	UserDb datastore.Batching
}

func SetupAuth(repopath string) (*AuthDB, error) {
	rdb, idb, err := DataStores(repopath)
	if err != nil {
		return nil, err
	}

	return &AuthDB{
		RepoDb: rdb,
		UserDb: idb,
	}, nil
}

func DataStores(repopath string) (datastore.Batching, datastore.Batching, error) {
	repodir, err := homedir.Expand(repopath)
	if err != nil {
		return nil, nil, err
	}

	rdb, err := setupLevelDs(path.Join(repodir, RepoDB), false)
	if err != nil {
		return nil, nil, err
	}

	idb, err := setupLevelDs(path.Join(repodir, UserDB), false)
	if err != nil {
		return nil, nil, err
	}

	return rdb, idb, nil
}

func setupLevelDs(path string, readonly bool) (datastore.Batching, error) {
	if _, err := os.ReadDir(path); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0777)
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

func (db *AuthDB) LoadRepoInfoDB(dbName string) (datastore.Batching, error) {
	switch dbName {
	case RepoDB:
		return db.RepoDb, nil
	case UserDB:
		return db.UserDb, nil
	default:
		return nil, xerrors.New("no partern")
	}
}
