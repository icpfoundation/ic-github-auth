package db

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
	levelds "github.com/ipfs/go-ds-leveldb"
	"github.com/lyswifter/ic-auth/params"
	"github.com/lyswifter/ic-auth/types"
	"github.com/lyswifter/ic-auth/util"
	"github.com/mitchellh/go-homedir"
	ldbopts "github.com/syndtr/goleveldb/leveldb/opt"
	"go.uber.org/multierr"
	"golang.org/x/xerrors"
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
		util.Errorf("NewDatastore: %s", err)
		return nil, err
	}
	return db, err
}

func DataStores() {
	repodir, err := homedir.Expand(params.RepoPath)
	if err != nil {
		return
	}

	ldb, err := setupLevelDs(repodir, false)
	if err != nil {
		util.Errorf("setup beacondb: err %s", err)
		return
	}
	InfoDB = ldb
	util.Infof("InfoDB: %+v", InfoDB)

	idb, err := setupLevelDs(path.Join(repodir, "uinfo"), false)
	if err != nil {
		util.Errorf("setup infodb: err %s", err)
		return
	}
	UserInfoDB = idb
	util.Infof("UserInfoDB: %+v", UserInfoDB)
}

func SaveAccessToken(ctx context.Context, state string, authResponse []byte) error {
	key := datastore.NewKey(state)
	err := UserInfoDB.Put(ctx, key, authResponse)
	if err != nil {
		return err
	}

	util.Infof("write user authorization info for state: %s", key.String())
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
	util.Infof("write installation code for state: %s", key.String())
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

func SaveCanisterInfo(ctx context.Context, canisterInfo types.CanisterInfo) error {
	if canisterInfo.CanisterID == "" {
		return xerrors.New("canister id is nil")
	}

	key := datastore.NewKey(canisterInfo.CanisterID)

	info, err := json.Marshal(canisterInfo)
	if err != nil {
		return err
	}

	err = InfoDB.Put(ctx, key, info)
	if err != nil {
		return err
	}

	fmt.Printf("save cinfo ok: %s", canisterInfo.CanisterID)

	return nil
}

func ReadCanisterInfo(ctx context.Context, id string) ([]byte, error) {
	///
	key := datastore.NewKey(id)

	ishas, err := InfoDB.Has(ctx, key)
	if err != nil {
		return nil, err
	}

	if ishas {
		ret, err := InfoDB.Get(ctx, key)
		if err != nil {
			return nil, err
		}

		fmt.Printf("read cinfo ok: %s", id)
		return ret, nil
	}

	return nil, nil
}

func ReadCanisterList(ctx context.Context) ([]string, error) {
	res, err := InfoDB.Query(ctx, query.Query{})
	if err != nil {
		return nil, err
	}

	defer res.Close()

	canisterids := []string{}

	var errs error
	for {
		res, ok := res.NextSync()
		if !ok {
			break
		}

		if res.Error != nil {
			return nil, res.Error
		}

		cinfo := &types.CanisterInfo{}
		err := json.Unmarshal(res.Value, cinfo)
		if err != nil {
			errs = multierr.Append(errs, xerrors.Errorf("decoding state for key '%s': %w", res.Key, err))
			continue
		}

		if cinfo.CanisterID == "" {
			continue
		}

		canisterids = append(canisterids, cinfo.CanisterID)
	}

	fmt.Printf("read canister infos ok, len %d\n", len(canisterids))
	return canisterids, nil
}
