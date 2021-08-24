package badgermigration

import (
	"strconv"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/skl"
	"github.com/stackrox/rox/migrator/badgerhelpers"
	"github.com/stackrox/rox/migrator/bolthelpers"
	"github.com/stretchr/testify/suite"
	"go.etcd.io/bbolt"
)

var (
	bucketName = []byte("bucket")

	maxTxnSize  = 15 * badger.DefaultOptions("").MaxTableSize / 100
	maxTxnCount = maxTxnSize / int64(skl.MaxNodeSize)
)

func TestRewrite(t *testing.T) {
	suite.Run(t, new(MigrationTestSuite))
}

type MigrationTestSuite struct {
	suite.Suite

	boltDB   *bbolt.DB
	badgerDB *badger.DB
}

func (suite *MigrationTestSuite) SetupTest() {
	var err error
	suite.badgerDB, err = badgerhelpers.NewTemp("single")
	suite.Require().NoError(err)

	suite.boltDB, err = bolthelpers.NewTemp("single")
	suite.Require().NoError(err)
}

func (suite *MigrationTestSuite) TearDownTest() {
	_ = suite.badgerDB.Close()
	_ = suite.boltDB.Close()
}

type kv struct {
	key, value string
}

func (suite *MigrationTestSuite) checkBadger(keypairs ...kv) {
	suite.NoError(suite.badgerDB.View(func(txn *badger.Txn) error {
		for _, kp := range keypairs {
			item, err := txn.Get([]byte(kp.key))
			suite.NoError(err)
			if item != nil {
				dst, err := item.ValueCopy(nil)
				suite.NoError(err)
				suite.Equal([]byte(kp.value), dst)
			}
		}
		return nil
	}))
}

func (suite *MigrationTestSuite) TestSmall() {
	// Test the case where both the size and the count are under the limits

	err := suite.boltDB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}
		return bucket.Put([]byte("key"), []byte("value"))
	})
	suite.NoError(err)

	// Rewrite into badger
	suite.NoError(rewrite(suite.boltDB, suite.badgerDB, bucketName))

	// check badger for result
	suite.checkBadger(kv{key: "bucket\x00key", value: "value"})
}

func (suite *MigrationTestSuite) TestLargerThanCount() {
	// Test the case where the count is greater than the max count
	var keypairs []kv
	err := suite.boltDB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}

		for i := 0; i < int(maxTxnCount)+1; i++ {
			v := strconv.Itoa(i)
			keypairs = append(keypairs, kv{key: "bucket\x00" + v, value: "1"})
			if err := bucket.Put([]byte(v), []byte("1")); err != nil {
				return err
			}
		}
		return nil
	})
	suite.NoError(err)

	// Rewrite into badger
	suite.NoError(rewrite(suite.boltDB, suite.badgerDB, bucketName))

	// check badger for result
	suite.checkBadger(keypairs...)
}