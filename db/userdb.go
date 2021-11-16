package db

import (
	"context"
	"fmt"

	"github.com/ipfs/go-datastore"
	"github.com/lyswifter/ic-auth/util"
)

func (db *AuthDB) SaveAccessToken(ctx context.Context, state string, authResponse []byte) error {
	key := datastore.NewKey(state)
	err := db.UserDb.Put(ctx, key, authResponse)
	if err != nil {
		return err
	}

	util.Infof("write user authorization info for state: %s", key.String())
	return nil
}

func (db *AuthDB) ReadAccessToken(ctx context.Context, state string) ([]byte, error) {
	key := datastore.NewKey(state)
	ishas, err := db.UserDb.Has(ctx, key)
	if err != nil {
		return nil, err
	}

	if ishas {
		token, err := db.UserDb.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		return token, nil
	}
	return nil, nil
}

func (db *AuthDB) SaveInstallCode(ctx context.Context, state string, code string) error {

	key := datastore.NewKey(fmt.Sprintf("%s-installationCode", state))
	err := db.UserDb.Put(ctx, key, []byte(code))
	if err != nil {
		return err
	}
	util.Infof("write installation code for state: %s", key.String())
	return nil
}

func (db *AuthDB) ReadInstallCode(ctx context.Context, state string) ([]byte, error) {
	key := datastore.NewKey(fmt.Sprintf("%s-installationCode", state))
	ishas, err := db.UserDb.Has(ctx, key)
	if err != nil {
		return nil, err
	}

	if ishas {
		code, err := db.UserDb.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		return code, nil
	}

	return nil, nil
}
