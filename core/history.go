package core

import (
	"fmt"

	"github.com/ipfs/go-datastore"
	"github.com/qri-io/dataset"
	"github.com/qri-io/dataset/dsfs"
	"github.com/qri-io/qri/repo"
)

type HistoryRequests struct {
	repo repo.Repo
}

func NewHistoryRequests(r repo.Repo) *HistoryRequests {
	return &HistoryRequests{
		repo: r,
	}
}

type LogParams struct {
	Path datastore.Key
	ListParams
}

func (d *HistoryRequests) Log(params *LogParams, log *[]*dataset.Dataset) (err error) {
	dss := []*dataset.Dataset{}
	limit := params.Limit
	ds := &dataset.Dataset{Previous: params.Path}

	if params.Path.String() == "" {
		return fmt.Errorf("path is required")
	}

	for {
		if ds.Previous.String() == "" {
			break
		}

		ds, err = dsfs.LoadDataset(d.repo.Store(), ds.Previous)
		if err != nil {
			return err
		}
		dss = append(dss, ds)

		limit--
		if limit == 0 {
			break
		}
	}

	*log = dss
	return nil
}
