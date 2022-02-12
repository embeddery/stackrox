// Code generated by pg-bindings generator. DO NOT EDIT.

package postgres

import (
	"context"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	ops "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
)

const (
	countStmt  = "SELECT COUNT(*) FROM singlekey"
	existsStmt = "SELECT EXISTS(SELECT 1 FROM singlekey WHERE Key = $1)"

	getStmt     = "SELECT serialized FROM singlekey WHERE Key = $1"
	deleteStmt  = "DELETE FROM singlekey WHERE Key = $1"
	walkStmt    = "SELECT serialized FROM singlekey"
	getIDsStmt  = "SELECT Key FROM singlekey"
	getManyStmt = "SELECT serialized FROM singlekey WHERE Key = ANY($1::text[])"

	deleteManyStmt = "DELETE FROM singlekey WHERE Key = ANY($1::text[])"
)

var (
	log = logging.LoggerForModule()

	table = "singlekey"
)

type Store interface {
	Count() (int, error)
	Exists(key string) (bool, error)
	Get(key string) (*storage.TestSingleKeyStruct, bool, error)
	Upsert(obj *storage.TestSingleKeyStruct) error
	UpsertMany(objs []*storage.TestSingleKeyStruct) error
	Delete(key string) error
	GetIDs() ([]string, error)
	GetMany(ids []string) ([]*storage.TestSingleKeyStruct, []int, error)
	DeleteMany(ids []string) error

	Walk(fn func(obj *storage.TestSingleKeyStruct) error) error
	AckKeysIndexed(keys ...string) error
	GetKeysToIndex() ([]string, error)
}

type storeImpl struct {
	db *pgxpool.Pool
}

func createTableSinglekey(db *pgxpool.Pool) {
	// hack for testing, remove
	db.Exec(context.Background(), "DROP TABLE singlekey CASCADE")

	table := `
create table if not exists singlekey (
    Key varchar,
    Name varchar UNIQUE,
    StringSlice text[],
    Bool bool,
    Uint64 numeric,
    Int64 numeric,
    Float numeric,
    Labels jsonb,
    Timestamp timestamp,
    Enum integer,
    Enums int[],
    Embedded_Embedded varchar,
    Oneofstring varchar,
    Oneofnested_Nested varchar,
    Oneofnested_Nested2_Nested2 varchar,
    serialized bytea,
    PRIMARY KEY(Key)
)
`

	_, err := db.Exec(context.Background(), table)
	if err != nil {
		panic("error creating table: " + table)
	}

	indexes := []string{

		"create index if not exists singlekey_Key on singlekey using hash(Key)",
	}
	for _, index := range indexes {
		if _, err := db.Exec(context.Background(), index); err != nil {
			panic(err)
		}
	}

	createTableSinglekeyNested(db)
}

func createTableSinglekeyNested(db *pgxpool.Pool) {
	// hack for testing, remove
	db.Exec(context.Background(), "DROP TABLE singlekey_Nested CASCADE")

	table := `
create table if not exists singlekey_Nested (
    parent_Key varchar,
    idx numeric,
    Nested varchar,
    Nested2_Nested2 varchar,
    PRIMARY KEY(parent_Key, idx),
    CONSTRAINT fk_parent_table FOREIGN KEY (parent_Key) REFERENCES singlekey(Key) ON DELETE CASCADE
)
`

	_, err := db.Exec(context.Background(), table)
	if err != nil {
		panic("error creating table: " + table)
	}

	indexes := []string{

		"create index if not exists singlekeyNested_idx on singlekey_Nested using btree(idx)",
	}
	for _, index := range indexes {
		if _, err := db.Exec(context.Background(), index); err != nil {
			panic(err)
		}
	}

}

func insertIntoSinglekey(db *pgxpool.Pool, obj *storage.TestSingleKeyStruct) error {

	serialized, marshalErr := obj.Marshal()
	if marshalErr != nil {
		return marshalErr
	}

	values := []interface{}{
		// parent primary keys start

		obj.GetKey(),
		obj.GetName(),
		obj.GetStringSlice(),
		obj.GetBool(),
		obj.GetUint64(),
		obj.GetInt64(),
		obj.GetFloat(),
		obj.GetLabels(),
		obj.GetTimestamp(),
		obj.GetEnum(),
		obj.GetEnums(),
		obj.GetEmbedded().GetEmbedded(),
		obj.GetOneofstring(),
		obj.GetOneofnested().GetNested(),
		obj.GetOneofnested().GetNested2().GetNested2(),
		serialized,
	}

	finalStr := "INSERT INTO singlekey (Key, Name, StringSlice, Bool, Uint64, Int64, Float, Labels, Timestamp, Enum, Enums, Embedded_Embedded, Oneofstring, Oneofnested_Nested, Oneofnested_Nested2_Nested2, serialized) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) ON CONFLICT(Key) DO UPDATE SET Key = EXCLUDED.Key, Name = EXCLUDED.Name, StringSlice = EXCLUDED.StringSlice, Bool = EXCLUDED.Bool, Uint64 = EXCLUDED.Uint64, Int64 = EXCLUDED.Int64, Float = EXCLUDED.Float, Labels = EXCLUDED.Labels, Timestamp = EXCLUDED.Timestamp, Enum = EXCLUDED.Enum, Enums = EXCLUDED.Enums, Embedded_Embedded = EXCLUDED.Embedded_Embedded, Oneofstring = EXCLUDED.Oneofstring, Oneofnested_Nested = EXCLUDED.Oneofnested_Nested, Oneofnested_Nested2_Nested2 = EXCLUDED.Oneofnested_Nested2_Nested2, serialized = EXCLUDED.serialized"
	_, err := db.Exec(context.Background(), finalStr, values...)
	if err != nil {
		return err
	}

	var query string

	for childIdx, child := range obj.GetNested() {
		if err := insertIntoSinglekeyNested(db, child, obj.GetKey(), childIdx); err != nil {
			return err
		}
	}

	query = "delete from singlekey_Nested where parent_Key = $1 AND idx >= $2"
	_, err = db.Exec(context.Background(), query, obj.GetKey(), len(obj.GetNested()))
	if err != nil {
		return err
	}
	return nil
}

func insertIntoSinglekeyNested(db *pgxpool.Pool, obj *storage.TestSingleKeyStruct_Nested, parent_Key string, idx int) error {

	values := []interface{}{
		// parent primary keys start

		parent_Key,
		idx,
		obj.GetNested(),
		obj.GetNested2().GetNested2(),
	}

	finalStr := "INSERT INTO singlekey_Nested (parent_Key, idx, Nested, Nested2_Nested2) VALUES($1, $2, $3, $4) ON CONFLICT(parent_Key, idx) DO UPDATE SET parent_Key = EXCLUDED.parent_Key, idx = EXCLUDED.idx, Nested = EXCLUDED.Nested, Nested2_Nested2 = EXCLUDED.Nested2_Nested2"
	_, err := db.Exec(context.Background(), finalStr, values...)
	if err != nil {
		return err
	}

	return nil
}

// New returns a new Store instance using the provided sql instance.
func New(db *pgxpool.Pool) Store {
	globaldb.RegisterTable(table, "storage.TestSingleKeyStruct")

	createTableSinglekey(db)

	return &storeImpl{
		db: db,
	}
}

func (s *storeImpl) Upsert(obj *storage.TestSingleKeyStruct) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Upsert, "storage.TestSingleKeyStruct")

	return insertIntoSinglekey(s.db, obj)
}

func (s *storeImpl) UpsertMany(objs []*storage.TestSingleKeyStruct) error {
	for _, obj := range objs {
		if err := insertIntoSinglekey(s.db, obj); err != nil {
			return err
		}
	}
	return nil
}

// Count returns the number of objects in the store
func (s *storeImpl) Count() (int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Count, "storage.TestSingleKeyStruct")

	row := s.db.QueryRow(context.Background(), countStmt)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// Exists returns if the id exists in the store
func (s *storeImpl) Exists(key string) (bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Exists, "storage.TestSingleKeyStruct")

	row := s.db.QueryRow(context.Background(), existsStmt, key)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, pgutils.ErrNilIfNoRows(err)
	}
	return exists, nil
}

// Get returns the object, if it exists from the store
func (s *storeImpl) Get(key string) (*storage.TestSingleKeyStruct, bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Get, "storage.TestSingleKeyStruct")

	conn, release := s.acquireConn(ops.Get, "storage.TestSingleKeyStruct")
	defer release()

	row := conn.QueryRow(context.Background(), getStmt, key)
	var data []byte
	if err := row.Scan(&data); err != nil {
		return nil, false, pgutils.ErrNilIfNoRows(err)
	}

	var msg storage.TestSingleKeyStruct
	if err := proto.Unmarshal(data, &msg); err != nil {
		return nil, false, err
	}
	return &msg, true, nil
}

func (s *storeImpl) acquireConn(op ops.Op, typ string) (*pgxpool.Conn, func()) {
	defer metrics.SetAcquireDuration(time.Now(), op, typ)
	conn, err := s.db.Acquire(context.Background())
	if err != nil {
		panic(err)
	}
	return conn, conn.Release
}

// Delete removes the specified ID from the store
func (s *storeImpl) Delete(key string) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Remove, "storage.TestSingleKeyStruct")

	conn, release := s.acquireConn(ops.Remove, "storage.TestSingleKeyStruct")
	defer release()

	if _, err := conn.Exec(context.Background(), deleteStmt, key); err != nil {
		return err
	}
	return nil
}

// GetIDs returns all the IDs for the store
func (s *storeImpl) GetIDs() ([]string, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.GetAll, "storage.TestSingleKeyStructIDs")

	rows, err := s.db.Query(context.Background(), getIDsStmt)
	if err != nil {
		return nil, pgutils.ErrNilIfNoRows(err)
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// GetMany returns the objects specified by the IDs or the index in the missing indices slice
func (s *storeImpl) GetMany(ids []string) ([]*storage.TestSingleKeyStruct, []int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.GetMany, "storage.TestSingleKeyStruct")

	conn, release := s.acquireConn(ops.GetMany, "storage.TestSingleKeyStruct")
	defer release()

	rows, err := conn.Query(context.Background(), getManyStmt, ids)
	if err != nil {
		if err == pgx.ErrNoRows {
			missingIndices := make([]int, 0, len(ids))
			for i := range ids {
				missingIndices = append(missingIndices, i)
			}
			return nil, missingIndices, nil
		}
		return nil, nil, err
	}
	defer rows.Close()
	elems := make([]*storage.TestSingleKeyStruct, 0, len(ids))
	foundSet := make(map[string]struct{})
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, nil, err
		}
		var msg storage.TestSingleKeyStruct
		if err := proto.Unmarshal(data, &msg); err != nil {
			return nil, nil, err
		}
		foundSet[msg.GetKey()] = struct{}{}
		elems = append(elems, &msg)
	}
	missingIndices := make([]int, 0, len(ids)-len(foundSet))
	for i, id := range ids {
		if _, ok := foundSet[id]; !ok {
			missingIndices = append(missingIndices, i)
		}
	}
	return elems, missingIndices, nil
}

// Delete removes the specified IDs from the store
func (s *storeImpl) DeleteMany(ids []string) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.RemoveMany, "storage.TestSingleKeyStruct")

	conn, release := s.acquireConn(ops.RemoveMany, "storage.TestSingleKeyStruct")
	defer release()
	if _, err := conn.Exec(context.Background(), deleteManyStmt, ids); err != nil {
		return err
	}
	return nil
}

// Walk iterates over all of the objects in the store and applies the closure
func (s *storeImpl) Walk(fn func(obj *storage.TestSingleKeyStruct) error) error {
	rows, err := s.db.Query(context.Background(), walkStmt)
	if err != nil {
		return pgutils.ErrNilIfNoRows(err)
	}
	defer rows.Close()
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return err
		}
		var msg storage.TestSingleKeyStruct
		if err := proto.Unmarshal(data, &msg); err != nil {
			return err
		}
		return fn(&msg)
	}
	return nil
}

//// Stubs for satisfying legacy interfaces

// AckKeysIndexed acknowledges the passed keys were indexed
func (s *storeImpl) AckKeysIndexed(keys ...string) error {
	return nil
}

// GetKeysToIndex returns the keys that need to be indexed
func (s *storeImpl) GetKeysToIndex() ([]string, error) {
	return nil, nil
}
