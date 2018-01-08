package inmem

import (
	"testing"

	"bitbucket.org/stack-rox/apollo/apollo/db"
	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
)

func testBenchmarksTriggers(t *testing.T, insertStorage, retrievalStorage db.BenchmarkTriggerStorage) {
	triggerTime1 := ptypes.TimestampNow()
	triggerTime2 := ptypes.TimestampNow()
	triggerTime2.Seconds -= 1000

	triggers := []*v1.BenchmarkTrigger{
		{
			Name: "trigger1",
			Time: triggerTime1,
		},
		{
			Name: "trigger2",
			Time: triggerTime2,
		},
	}

	// Test Add
	for _, trigger := range triggers {
		assert.NoError(t, insertStorage.AddBenchmarkTrigger(trigger))
	}
}

func TestBenchmarkTriggersPersistence(t *testing.T) {
	t.Parallel()
	persistent, err := createBoltDB()
	if err != nil {
		t.Fatal(err)
	}
	storage := newBenchmarkTriggerStore(persistent)
	testBenchmarksTriggers(t, storage, persistent)
}

func TestBenchmarkTriggers(t *testing.T) {
	t.Parallel()
	persistent, err := createBoltDB()
	if err != nil {
		t.Fatal(err)
	}
	storage := newBenchmarkTriggerStore(persistent)
	testBenchmarksTriggers(t, storage, storage)
}

func TestBenchmarkTriggersFiltering(t *testing.T) {
	t.Parallel()
	persistent, err := createBoltDB()
	if err != nil {
		t.Fatal(err)
	}

	triggerTime1 := ptypes.TimestampNow()
	triggerTime2 := ptypes.TimestampNow()
	triggerTime3 := ptypes.TimestampNow()
	triggerTime2.Seconds -= 1000
	triggerTime3.Seconds -= 2000

	storage := newBenchmarkTriggerStore(persistent)
	trigger1 := &v1.BenchmarkTrigger{
		Name:     "trigger1",
		Time:     triggerTime1,
		Clusters: []string{"development"},
	}
	trigger2 := &v1.BenchmarkTrigger{
		Name:     "trigger2",
		Time:     triggerTime2,
		Clusters: []string{"integration", "production"},
	}
	// trigger with no cluster
	trigger3 := &v1.BenchmarkTrigger{
		Name: "trigger3",
		Time: triggerTime3,
	}
	triggers := []*v1.BenchmarkTrigger{
		trigger1,
		trigger2,
		trigger3,
	}

	// Test Add
	for _, trigger := range triggers {
		assert.NoError(t, storage.AddBenchmarkTrigger(trigger))
	}

	actualTriggers, err := storage.GetBenchmarkTriggers(&v1.GetBenchmarkTriggersRequest{})
	assert.NoError(t, err)
	assert.Equal(t, triggers, actualTriggers)

	actualTriggers, err = storage.GetBenchmarkTriggers(&v1.GetBenchmarkTriggersRequest{
		Names: []string{"trigger1"},
	})
	assert.NoError(t, err)
	assert.Equal(t, []*v1.BenchmarkTrigger{trigger1}, actualTriggers)

	actualTriggers, err = storage.GetBenchmarkTriggers(&v1.GetBenchmarkTriggersRequest{
		Clusters: []string{"development"},
	})
	assert.NoError(t, err)
	assert.Equal(t, []*v1.BenchmarkTrigger{trigger1, trigger3}, actualTriggers)
}
