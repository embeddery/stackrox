// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package postgres

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	storage "github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/testutils/envisolator"
	"github.com/stretchr/testify/suite"
)

type ImageComponentRelationStoreSuite struct {
	suite.Suite
	envIsolator *envisolator.EnvIsolator
}

func TestImageComponentRelationStore(t *testing.T) {
	suite.Run(t, new(ImageComponentRelationStoreSuite))
}

func (s *ImageComponentRelationStoreSuite) SetupTest() {
	s.envIsolator = envisolator.NewEnvIsolator(s.T())
	s.envIsolator.Setenv(features.PostgresDatastore.EnvVar(), "true")

	if !features.PostgresDatastore.Enabled() {
		s.T().Skip("Skip postgres store tests")
		s.T().SkipNow()
	}
}

func (s *ImageComponentRelationStoreSuite) TearDownTest() {
	s.envIsolator.RestoreAll()
}

func (s *ImageComponentRelationStoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	source := pgtest.GetConnectionString(s.T())
	config, err := pgxpool.ParseConfig(source)
	s.Require().NoError(err)
	pool, err := pgxpool.ConnectConfig(ctx, config)
	s.NoError(err)
	defer pool.Close()

	Destroy(ctx, pool)
	store := New(ctx, pool)

	imageComponentEdge := &storage.ImageComponentEdge{}
	s.NoError(testutils.FullInit(imageComponentEdge, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundImageComponentEdge, exists, err := store.Get(ctx, imageComponentEdge.GetImageId(), imageComponentEdge.GetImageComponentId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundImageComponentEdge)

}
