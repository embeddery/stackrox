package datastore

import (
	"context"

	"github.com/stackrox/rox/central/activecomponent/datastore/internal/store"
	"github.com/stackrox/rox/central/activecomponent/datastore/search"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/dackbox/graph"
	pkgSearch "github.com/stackrox/rox/pkg/search"
)

// DataStore is an intermediary to ActiveComponent storage.
//go:generate mockgen-wrapper DataStore
type DataStore interface {
	Search(ctx context.Context, query *v1.Query) ([]pkgSearch.Result, error)
	SearchRawActiveComponents(ctx context.Context, q *v1.Query) ([]*storage.ActiveComponent, error)

	Exists(ctx context.Context, id string) (bool, error)
	Get(ctx context.Context, id string) (*storage.ActiveComponent, bool, error)
	GetBatch(ctx context.Context, id []string) ([]*storage.ActiveComponent, error)
}

// New returns a new instance of a DataStore.
func New(graphProvider graph.Provider, storage store.Store, searcher search.Searcher) DataStore {
	ds := &datastoreImpl{
		storage:       storage,
		graphProvider: graphProvider,
		searcher:      searcher,
	}
	return ds
}