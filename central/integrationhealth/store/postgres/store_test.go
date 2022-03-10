// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package postgres

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/testutils/envisolator"
	"github.com/stretchr/testify/suite"
)

type IntegrationhealthStoreSuite struct {
	suite.Suite
	envIsolator *envisolator.EnvIsolator
}

func TestIntegrationhealthStore(t *testing.T) {
	suite.Run(t, new(IntegrationhealthStoreSuite))
}

func (s *IntegrationhealthStoreSuite) SetupTest() {
	s.envIsolator = envisolator.NewEnvIsolator(s.T())
	s.envIsolator.Setenv(features.PostgresDatastore.EnvVar(), "true")

	if !features.PostgresDatastore.Enabled() {
		s.T().Skip("Skip postgres store tests")
		s.T().SkipNow()
	}
}

func (s *IntegrationhealthStoreSuite) TearDownTest() {
	s.envIsolator.RestoreAll()
}

func (s *IntegrationhealthStoreSuite) TestStore() {
	ctx := context.Background()

	source := pgtest.GetConnectionString(s.T())
	config, err := pgxpool.ParseConfig(source)
	s.Require().NoError(err)
	pool, err := pgxpool.ConnectConfig(ctx, config)
	s.NoError(err)
	defer pool.Close()

	Destroy(ctx, pool)
	store := New(ctx, pool)

	integrationHealth := &storage.IntegrationHealth{}
	s.NoError(testutils.FullInit(integrationHealth, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundIntegrationHealth, exists, err := store.Get(ctx, integrationHealth.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundIntegrationHealth)

	s.NoError(store.Upsert(ctx, integrationHealth))
	foundIntegrationHealth, exists, err = store.Get(ctx, integrationHealth.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(integrationHealth, foundIntegrationHealth)

	integrationHealthCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(integrationHealthCount, 1)

	integrationHealthExists, err := store.Exists(ctx, integrationHealth.GetId())
	s.NoError(err)
	s.True(integrationHealthExists)
	s.NoError(store.Upsert(ctx, integrationHealth))

	foundIntegrationHealth, exists, err = store.Get(ctx, integrationHealth.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(integrationHealth, foundIntegrationHealth)

	s.NoError(store.Delete(ctx, integrationHealth.GetId()))
	foundIntegrationHealth, exists, err = store.Get(ctx, integrationHealth.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundIntegrationHealth)
}