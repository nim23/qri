package fsrepo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
	"github.com/qri-io/cafs"
	"github.com/qri-io/dataset"
	"github.com/qri-io/dataset/dsfs"
	"github.com/qri-io/qri/repo"
)

// Datasets is a file-based implementation of the repo.Datasets interface
type Datasets struct {
	basepath
	file  File
	store cafs.Filestore
}

// NewDatasets allocates a new Datasets instance
func NewDatasets(base string, file File, store cafs.Filestore) Datasets {
	return Datasets{basepath: basepath(base), file: file, store: store}
}

// CreateDataset initializes a dataset from a dataset pointer and data file
func (r Datasets) CreateDataset(ds *dataset.Dataset, data cafs.File) (path datastore.Key, err error) {

	return
}

// PutDataset adds a dataset to the store
func (r Datasets) PutDataset(path datastore.Key, ds *dataset.Dataset) error {
	d, err := r.datasets()
	if err != nil {
		return err
	}
	d[path.String()] = ds
	return r.saveFile(d, r.file)
}

// PutDatasets adds a number of datasets to the store
func (r Datasets) PutDatasets(datasets []*repo.DatasetRef) error {
	ds, err := r.datasets()
	if err != nil {
		return err
	}
	for _, dr := range datasets {
		if dr.Path != "" && dr.Dataset != nil {
			ds[dr.Path] = dr.Dataset
		}
	}
	return r.saveFile(ds, r.file)
}

// GetDataset grabs a dataset from the store
func (r Datasets) GetDataset(path datastore.Key) (*dataset.Dataset, error) {
	ds, err := r.datasets()
	if err != nil {
		return nil, err
	}
	ps := path.String()
	for p, d := range ds {
		if ps == p {
			return d, nil
		}
	}
	if r.store != nil {
		return dsfs.LoadDataset(r.store, path)
	}

	return nil, datastore.ErrNotFound
}

// DeleteDataset removes a dataset from the store
func (r Datasets) DeleteDataset(path datastore.Key) error {
	ds, err := r.datasets()
	if err != nil {
		return err
	}
	delete(ds, path.String())
	return r.saveFile(ds, r.file)
}

// Query fetches a set of Dataset References from the store based on a set of query params
func (r Datasets) Query(q query.Query) (query.Results, error) {
	ds, err := r.datasets()
	if err != nil {
		return nil, err
	}

	re := make([]query.Entry, 0, len(ds))
	for path, d := range ds {
		re = append(re, query.Entry{Key: path, Value: d})
	}
	res := query.ResultsWithEntries(q, re)
	res = query.NaiveQueryApply(q, res)
	return res, nil
}

func (r *Datasets) datasets() (map[string]*dataset.Dataset, error) {
	ds := map[string]*dataset.Dataset{}
	data, err := ioutil.ReadFile(r.filepath(r.file))
	if err != nil {
		if os.IsNotExist(err) {
			return ds, nil
		}
		return ds, fmt.Errorf("error loading datasets: %s", err.Error())
	}

	if err := json.Unmarshal(data, &ds); err != nil {
		return ds, fmt.Errorf("error unmarshaling datasets: %s", err.Error())
	}
	return ds, nil
}
