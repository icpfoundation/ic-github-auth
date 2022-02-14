package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
	"github.com/lyswifter/ic-auth/types"
	"go.uber.org/multierr"
	"golang.org/x/xerrors"
)

func (db *AuthDB) SaveCanisterInfo(ctx context.Context, canisterInfo types.CanisterInfo) error {
	if canisterInfo.Controller == "" {
		return xerrors.New("canister id is nil")
	}

	key := datastore.NewKey(fmt.Sprintf("%s/%s", canisterInfo.Owner, canisterInfo.CanisterID))
	ishas, err := db.RepoDb.Has(ctx, key)
	if err != nil {
		return err
	}

	if ishas {
		fmt.Printf("record already exist for: %s", key.String())
		return nil
	}

	info, err := json.Marshal(canisterInfo)
	if err != nil {
		return err
	}

	err = db.RepoDb.Put(ctx, key, info)
	if err != nil {
		return err
	}

	fmt.Printf("save cinfo for: %s ok", key.String())
	return nil
}

func (db *AuthDB) ReadCanisterInfo(ctx context.Context, owner string, id string) ([]byte, error) {
	key := datastore.NewKey(fmt.Sprintf("%s/%s", owner, id))
	ishas, err := db.RepoDb.Has(ctx, key)
	if err != nil {
		return nil, err
	}

	if ishas {
		ret, err := db.RepoDb.Get(ctx, key)
		if err != nil {
			return nil, err
		}

		fmt.Printf("read cinfo ok: %s", id)
		return ret, nil
	}

	return nil, nil
}

func (db *AuthDB) ReadCanisterList(ctx context.Context, owner string) ([]string, error) {
	res, err := db.RepoDb.Query(ctx, query.Query{
		Filters: []query.Filter{
			query.FilterKeyPrefix{
				Prefix: datastore.NewKey(owner).String(),
			},
		},
	})
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

	if errs != nil {
		return nil, errs
	}

	fmt.Printf("read canister infos ok, len %d\n", len(canisterids))

	return canisterids, nil
}
