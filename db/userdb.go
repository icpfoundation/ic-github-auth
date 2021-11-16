package db

import (
	"context"
	"fmt"

	"github.com/ipfs/go-datastore"
	"github.com/lyswifter/ic-auth/util"
)

func SaveAccessToken(ctx context.Context, state string, authResponse []byte) error {
	db, err := LoadRepoInfoDB(UserDB)
	if err != nil {
		return err
	}

	key := datastore.NewKey(state)
	err = db.Put(ctx, key, authResponse)
	if err != nil {
		return err
	}

	util.Infof("write user authorization info for state: %s", key.String())
	return nil
}

func ReadAccessToken(ctx context.Context, state string) ([]byte, error) {
	db, err := LoadRepoInfoDB(UserDB)
	if err != nil {
		return nil, err
	}

	key := datastore.NewKey(state)
	ishas, err := db.Has(ctx, key)
	if err != nil {
		return nil, err
	}

	if ishas {
		token, err := db.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		return token, nil
	}
	return nil, nil
}

func SaveInstallCode(ctx context.Context, state string, code string) error {
	db, err := LoadRepoInfoDB(UserDB)
	if err != nil {
		return err
	}

	key := datastore.NewKey(fmt.Sprintf("%s-installationCode", state))
	err = db.Put(ctx, key, []byte(code))
	if err != nil {
		return err
	}
	util.Infof("write installation code for state: %s", key.String())
	return nil
}

func ReadInstallCode(ctx context.Context, state string) ([]byte, error) {
	db, err := LoadRepoInfoDB(UserDB)
	if err != nil {
		return nil, err
	}

	key := datastore.NewKey(fmt.Sprintf("%s-installationCode", state))
	ishas, err := db.Has(ctx, key)
	if err != nil {
		return nil, err
	}

	if ishas {
		code, err := db.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		return code, nil
	}

	return nil, nil
}
