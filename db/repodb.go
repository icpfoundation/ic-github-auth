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

func SaveCanisterInfo(ctx context.Context, canisterInfo types.CanisterInfo) error {
	if canisterInfo.CanisterID == "" {
		return xerrors.New("canister id is nil")
	}

	db, err := LoadRepoInfoDB(RepoDB)
	if err != nil {
		return err
	}

	key := datastore.NewKey(canisterInfo.CanisterID)

	info, err := json.Marshal(canisterInfo)
	if err != nil {
		return err
	}

	err = db.Put(ctx, key, info)
	if err != nil {
		return err
	}

	fmt.Printf("save cinfo ok: %s", canisterInfo.CanisterID)

	return nil
}

func ReadCanisterInfo(ctx context.Context, id string) ([]byte, error) {
	db, err := LoadRepoInfoDB(RepoDB)
	if err != nil {
		return nil, err
	}

	key := datastore.NewKey(id)

	ishas, err := db.Has(ctx, key)
	if err != nil {
		return nil, err
	}

	if ishas {
		ret, err := db.Get(ctx, key)
		if err != nil {
			return nil, err
		}

		fmt.Printf("read cinfo ok: %s", id)
		return ret, nil
	}

	return nil, nil
}

func ReadCanisterList(ctx context.Context) ([]string, error) {
	db, err := LoadRepoInfoDB(RepoDB)
	if err != nil {
		return nil, err
	}

	res, err := db.Query(ctx, query.Query{})
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
