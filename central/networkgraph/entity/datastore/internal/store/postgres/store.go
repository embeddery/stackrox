// Code generated by pg-bindings generator. DO NOT EDIT.

package postgres

import (
	"bytes"
	"context"
	"reflect"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/types"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	ops "github.com/stackrox/rox/pkg/metrics"
)

const (
	countStmt  = "SELECT COUNT(*) FROM NetworkEntity"
	existsStmt = "SELECT EXISTS(SELECT 1 FROM NetworkEntity WHERE )"

	getStmt    = "SELECT serialized FROM NetworkEntity WHERE "
	deleteStmt = "DELETE FROM NetworkEntity WHERE "
	walkStmt   = "SELECT serialized FROM NetworkEntity"
)

var (
	log = logging.LoggerForModule()

	table = "NetworkEntity"

	marshaler = &jsonpb.Marshaler{EnumsAsInts: true, EmitDefaults: true}
)

type Store interface {
	Count() (int, error)
	Exists() (bool, error)
	Get() (*storage.NetworkEntity, bool, error)
	Upsert(obj *storage.NetworkEntity) error
	UpsertMany(objs []*storage.NetworkEntity) error
	Delete() error

	Walk(fn func(obj *storage.NetworkEntity) error) error
	AckKeysIndexed(keys ...string) error
	GetKeysToIndex() ([]string, error)
}

type storeImpl struct {
	db *pgxpool.Pool
}

const (
	batchInsertTemplate = "<no value>"
)

// New returns a new Store instance using the provided sql instance.
func New(db *pgxpool.Pool) Store {
	globaldb.RegisterTable(table, "NetworkEntity")

	for _, table := range []string{
		"create table if not exists NetworkEntity(serialized jsonb not null, Info_Desc_ExternalSource_Default bool, PRIMARY KEY ());",
	} {
		_, err := db.Exec(context.Background(), table)
		if err != nil {
			panic("error creating table: " + table)
		}
	}

	//
	return &storeImpl{
		db: db,
	}
	//
}

// Count returns the number of objects in the store
func (s *storeImpl) Count() (int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Count, "NetworkEntity")

	row := s.db.QueryRow(context.Background(), countStmt)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// Exists returns if the id exists in the store
func (s *storeImpl) Exists() (bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Exists, "NetworkEntity")

	row := s.db.QueryRow(context.Background(), existsStmt)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, nilNoRows(err)
	}
	return exists, nil
}

func nilNoRows(err error) error {
	if err == pgx.ErrNoRows {
		return nil
	}
	return err
}

// Get returns the object, if it exists from the store
func (s *storeImpl) Get() (*storage.NetworkEntity, bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Get, "NetworkEntity")

	conn, release := s.acquireConn(ops.Get, "NetworkEntity")
	defer release()

	row := conn.QueryRow(context.Background(), getStmt)
	var data []byte
	if err := row.Scan(&data); err != nil {
		return nil, false, nilNoRows(err)
	}

	var msg storage.NetworkEntity
	buf := bytes.NewBuffer(data)
	defer metrics.SetJSONPBOperationDurationTime(time.Now(), "Unmarshal", "NetworkEntity")
	if err := jsonpb.Unmarshal(buf, &msg); err != nil {
		return nil, false, err
	}
	return &msg, true, nil
}

func convertEnumSliceToIntArray(i interface{}) []int32 {
	enumSlice := reflect.ValueOf(i)
	enumSliceLen := enumSlice.Len()
	resultSlice := make([]int32, 0, enumSliceLen)
	for i := 0; i < enumSlice.Len(); i++ {
		resultSlice = append(resultSlice, int32(enumSlice.Index(i).Int()))
	}
	return resultSlice
}

func nilOrStringTimestamp(t *types.Timestamp) *string {
	if t == nil {
		return nil
	}
	s := t.String()
	return &s
}

// Upsert inserts the object into the DB
func (s *storeImpl) Upsert(obj0 *storage.NetworkEntity) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Add, "NetworkEntity")

	t := time.Now()
	serialized, err := marshaler.MarshalToString(obj0)
	if err != nil {
		return err
	}
	metrics.SetJSONPBOperationDurationTime(t, "Marshal", "NetworkEntity")
	conn, release := s.acquireConn(ops.Add, "NetworkEntity")
	defer release()

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return err
	}
	doRollback := true
	defer func() {
		if doRollback {
			if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil {
				log.Errorf("error rolling backing: %v", err)
			}
		}
	}()

	localQuery := "insert into NetworkEntity(serialized, Info_Desc_ExternalSource_Default) values($1, $2) on conflict() do update set serialized = EXCLUDED.serialized, Info_Desc_ExternalSource_Default = EXCLUDED.Info_Desc_ExternalSource_Default"
	_, err = tx.Exec(context.Background(), localQuery, serialized, obj0.GetInfo().GetExternalSource().GetDefault())
	if err != nil {
		return err
	}

	doRollback = false
	return tx.Commit(context.Background())
}

func (s *storeImpl) acquireConn(op ops.Op, typ string) (*pgxpool.Conn, func()) {
	defer metrics.SetAcquireDuration(time.Now(), op, typ)
	conn, err := s.db.Acquire(context.Background())
	if err != nil {
		panic(err)
	}
	return conn, conn.Release
}

// UpsertMany batches objects into the DB
func (s *storeImpl) UpsertMany(objs []*storage.NetworkEntity) error {
	if len(objs) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.AddMany, "NetworkEntity")
	for _, obj0 := range objs {
		t := time.Now()
		serialized, err := marshaler.MarshalToString(obj0)
		if err != nil {
			return err
		}
		metrics.SetJSONPBOperationDurationTime(t, "Marshal", "NetworkEntity")
		localQuery := "insert into NetworkEntity(serialized, Info_Desc_ExternalSource_Default) values($1, $2) on conflict() do update set serialized = EXCLUDED.serialized, Info_Desc_ExternalSource_Default = EXCLUDED.Info_Desc_ExternalSource_Default"
		batch.Queue(localQuery, serialized, obj0.GetInfo().GetExternalSource().GetDefault())

	}

	conn, release := s.acquireConn(ops.AddMany, "NetworkEntity")
	defer release()

	results := conn.SendBatch(context.Background(), batch)
	if err := results.Close(); err != nil {
		return err
	}
	return nil
}

// Delete removes the specified ID from the store
func (s *storeImpl) Delete() error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Remove, "NetworkEntity")

	conn, release := s.acquireConn(ops.Remove, "NetworkEntity")
	defer release()

	if _, err := conn.Exec(context.Background(), deleteStmt); err != nil {
		return err
	}
	return nil
}

// Walk iterates over all of the objects in the store and applies the closure
func (s *storeImpl) Walk(fn func(obj *storage.NetworkEntity) error) error {
	rows, err := s.db.Query(context.Background(), walkStmt)
	if err != nil {
		return nilNoRows(err)
	}
	defer rows.Close()
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return err
		}
		var msg storage.NetworkEntity
		buf := bytes.NewReader(data)
		if err := jsonpb.Unmarshal(buf, &msg); err != nil {
			return err
		}
		return fn(&msg)
	}
	return nil
}

// AckKeysIndexed acknowledges the passed keys were indexed
func (s *storeImpl) AckKeysIndexed(keys ...string) error {
	return nil
}

// GetKeysToIndex returns the keys that need to be indexed
func (s *storeImpl) GetKeysToIndex() ([]string, error) {
	return nil, nil
}