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

	key := datastore.NewKey(canisterInfo.Controller)
	info, err := json.Marshal(canisterInfo)
	if err != nil {
		return err
	}

	err = db.UserDb.Put(ctx, key, info)
	if err != nil {
		return err
	}

	fmt.Printf("save cinfo ok: %s", canisterInfo.CanisterID)

	return nil
}

func (db *AuthDB) ReadCanisterInfo(ctx context.Context, id string) ([]byte, error) {
	key := datastore.NewKey(id)
	ishas, err := db.UserDb.Has(ctx, key)
	if err != nil {
		return nil, err
	}

	if ishas {
		ret, err := db.UserDb.Get(ctx, key)
		if err != nil {
			return nil, err
		}

		fmt.Printf("read cinfo ok: %s", id)
		return ret, nil
	}

	return nil, nil
}

func (db *AuthDB) ReadCanisterList(ctx context.Context, controller string) ([]string, error) {
	res, err := db.UserDb.Query(ctx, query.Query{
		Filters: []query.Filter{
			query.FilterKeyCompare{
				Op:  query.Equal,
				Key: controller,
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

	fmt.Printf("read canister infos ok, len %d\n", len(canisterids))
	return canisterids, nil
}
