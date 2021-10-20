package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/stackrox/rox/pkg/utils"
)

const storeFile = `

package postgres

import (
	"bytes"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	ops "github.com/stackrox/rox/pkg/metrics"
	"database/sql"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/lib/pq"
	"github.com/stackrox/rox/pkg/set"
)

var (
	log = logging.LoggerForModule()

	table = "{{.Table}}"
)

type Store interface {
	Count() (int, error)
	Exists(id string) (bool, error)
	GetIDs() ([]string, error)
	Get(id string) (*storage.{{.Type}}, bool, error)
	GetMany(ids []string) ([]*storage.{{.Type}}, []int, error)
	{{- if .NoKeyField}}
	UpsertWithID(id string, obj *storage.{{.Type}}) error
	UpsertManyWithIDs(ids []string, objs []*storage.{{.Type}}) error
	{{- else }}
	Upsert(obj *storage.{{.Type}}) error
	UpsertMany(objs []*storage.{{.Type}}) error
	{{- end}}
	Delete(id string) error
	DeleteMany(ids []string) error
	{{- if .NoKeyField}}
	WalkAllWithID(fn func(id string, obj *storage.{{.Type}}) error) error
	{{- else }}
	Walk(fn func(obj *storage.{{.Type}}) error) error
	{{- end}}
	AckKeysIndexed(keys ...string) error
	GetKeysToIndex() ([]string, error)
}

type storeImpl struct {
	db *sql.DB

	countStmt *sql.Stmt
	existsStmt *sql.Stmt
	getIDsStmt *sql.Stmt
	getStmt *sql.Stmt
	getManyStmt *sql.Stmt
	upsertWithIDStmt *sql.Stmt
	upsertStmt *sql.Stmt
	deleteStmt *sql.Stmt
	deleteManyStmt *sql.Stmt
	walkStmt *sql.Stmt
	walkWithIDStmt *sql.Stmt
}

func alloc() proto.Message {
	return &storage.{{.Type}}{}
}

{{- if not .NoKeyField}}

func keyFunc(msg proto.Message) string {
	return msg.(*storage.{{.Type}}).{{.KeyFunc}}
}
{{- end}}

{{- if .UniqKeyFunc}}

func uniqKeyFunc(msg proto.Message) string {
	return msg.(*storage.{{.Type}}).{{.UniqKeyFunc}}
}
{{- end}}

func compileStmtOrPanic(db *sql.DB, query string) *sql.Stmt {
	vulnStmt, err := db.Prepare(query)
	if err != nil {
		panic(err)
	}
	return vulnStmt
}

const (
	createTableQuery = "create table if not exists {{.Table}} (id varchar primary key, value jsonb)"
	createIDIndexQuery = "create index if not exists {{.Table}}_id on {{.Table}} using hash ((id))"
)

// New returns a new Store instance using the provided sql instance.
func New(db *sql.DB) Store {
	globaldb.RegisterTable(table, "{{.Type}}")

	_, err := db.Exec(createTableQuery)
	if err != nil {
		panic("error creating table")
	}

	_, err = db.Exec(createIDIndexQuery)
	if err != nil {
		panic("error creating index")
	}

//	{{- if .UniqKeyFunc}}
//	return &storeImpl{
//		crud: generic.NewUniqueKeyCRUD(db, bucket, {{if .NoKeyField}}nil{{else}}keyFunc{{end}}, allocCluster, uniqKeyFunc, {{.TrackIndex}}),
//	}
//	{{- else}}
	return &storeImpl{
		db: db,

		countStmt: compileStmtOrPanic(db, "select count(*) from {{.Table}}"),
		existsStmt: compileStmtOrPanic(db, "select exists(select 1 from {{.Table}} where id = $1)"),
		getIDsStmt: compileStmtOrPanic(db, "select id from {{.Table}}"),
		getStmt: compileStmtOrPanic(db, "select value from {{.Table}} where id = $1"),
		getManyStmt: compileStmtOrPanic(db, "select value from {{.Table}} where id = ANY($1::text[])"),
		upsertStmt: compileStmtOrPanic(db, "insert into {{.Table}}(id, value) values($1, $2) on conflict(id) do update set value=$2"),
		deleteStmt: compileStmtOrPanic(db, "delete from {{.Table}} where id = $1"),
		deleteManyStmt: compileStmtOrPanic(db, "delete from {{.Table}} where id = ANY($1::text[])"),
		walkStmt: compileStmtOrPanic(db, "select value from {{.Table}}"),
		walkWithIDStmt: compileStmtOrPanic(db, "select id, value from {{.Table}}"),
	}
//	{{- end}}
}

// Count returns the number of objects in the store
func (s *storeImpl) Count() (int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Count, "{{.Type}}")

	row := s.countStmt.QueryRow()
	if err := row.Err(); err != nil {
		return 0, err
	}
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// Exists returns if the id exists in the store
func (s *storeImpl) Exists(id string) (bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Exists, "{{.Type}}")

	row := s.existsStmt.QueryRow(id)
	if err := row.Err(); err != nil {
		return false, nilNoRows(err)
	}
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, nilNoRows(err)
	}
	return exists, nil
}

// GetIDs returns all the IDs for the store
func (s *storeImpl) GetIDs() ([]string, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.GetAll, "{{.Type}}IDs")

	rows, err := s.getIDsStmt.Query()
	if err != nil {
		return nil, nilNoRows(err)
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

func nilNoRows(err error) error {
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

// Get returns the object, if it exists from the store
func (s *storeImpl) Get(id string) (*storage.{{.Type}}, bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Get, "{{.Type}}")

	row := s.getStmt.QueryRow(id)
	if err := row.Err(); err != nil {
		return nil, false, nilNoRows(err)
	}

	var data []byte
	if err := row.Scan(&data); err != nil {
		return nil, false, nilNoRows(err)
	}

	msg := alloc()
	buf := bytes.NewBuffer(data)
	if err := jsonpb.Unmarshal(buf, msg); err != nil {
		return nil, false, err
	}
	return msg.(*storage.{{.Type}}), true, nil
}

// GetMany returns the objects specified by the IDs or the index in the missing indices slice 
func (s *storeImpl) GetMany(ids []string) ([]*storage.{{.Type}}, []int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.GetMany, "{{.Type}}")

	rows, err := s.getManyStmt.Query(pq.Array(ids))
	if err != nil {
		if err == sql.ErrNoRows {
			missingIndices := make([]int, 0, len(ids))
			for i := range ids {
				missingIndices = append(missingIndices, i)
			}
			return nil, missingIndices, nil
		}
		return nil, nil, err
	}
	defer rows.Close()
	elems := make([]*storage.{{.Type}}, 0, len(ids))
	foundSet := set.NewStringSet()
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, nil, err
		}
		msg := alloc()
		buf := bytes.NewBuffer(data)
		if err := jsonpb.Unmarshal(buf, msg); err != nil {
			return nil, nil, err
		}
		elem := msg.(*storage.{{.Type}})
		foundSet.Add(elem.GetId())
		elems = append(elems, elem)
	}
	missingIndices := make([]int, 0, len(ids)-len(foundSet))
	for i, id := range ids {
		if !foundSet.Contains(id) {
			missingIndices = append(missingIndices, i)
		}
	}
	return elems, missingIndices, nil
}

{{- if .NoKeyField}}
// UpsertWithID inserts the object into the DB
func (s *storeImpl) UpsertWithID(id string, obj *storage.{{.Type}}) error {
	return upsert(id, obj)
}

// UpsertManyWithIDs batches objects into the DB
func (s *storeImpl) UpsertManyWithIDs(ids []string, objs []*storage.{{.Type}}) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.AddMany, "{{.Type}}")

	// txn? or partial? what is the impact of one not being upserted
	for i, id := range ids {
		if err := s.upsert(id, objs(i)); err != nil {
			return err
		}
	}
	return nil
}
{{- else}}

func (s *storeImpl) upsert(id string, obj *storage.{{.Type}}) error {
	value, err := (&jsonpb.Marshaler{
		EnumsAsInts:  true,
		EmitDefaults: true,
	}).MarshalToString(obj)
	if err != nil {
		return err
	}
	_, err = s.upsertStmt.Exec(id, value)
	return err
}


// Upsert inserts the object into the DB
func (s *storeImpl) Upsert(obj *storage.{{.Type}}) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Add, "{{.Type}}")
	return s.upsert(keyFunc(obj), obj)
}

// UpsertMany batches objects into the DB
func (s *storeImpl) UpsertMany(objs []*storage.{{.Type}}) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.AddMany, "{{.Type}}")

	// Txn? or all errors to be passed through?
	for _, obj := range objs {
		if err := s.upsert(keyFunc(obj), obj); err != nil {
			return err
		}
	}
	return nil
}
{{- end}}

// Delete removes the specified ID from the store
func (s *storeImpl) Delete(id string) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Remove, "{{.Type}}")

	if _, err := s.deleteStmt.Exec(id); err != nil {
		return err
	}
	return nil
}

// Delete removes the specified IDs from the store
func (s *storeImpl) DeleteMany(ids []string) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.RemoveMany, "{{.Type}}")

	if _, err := s.deleteManyStmt.Exec(pq.Array(ids)); err != nil {
		return err
	}
	return nil
}

{{- if .NoKeyField}}
// WalkAllWithID iterates over all of the objects in the store and applies the closure
func (s *storeImpl) WalkAllWithID(fn func(id string, obj *storage.{{.Type}}) error) error {

	panic("unimplemented")	
//return b.crud.WalkAllWithID(func(id []byte, msg proto.Message) error {
	rows, err := s.walkStmt.Query()
	if err != nil {
		return nilNoRows(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var data []byte
		if err := rows.Scan(&id, &data); err != nil {
			return err
		}
		msg := alloc()
		buf := bytes.NewBuffer(data)
		if err := jsonpb.Unmarshal(buf, msg); err != nil {
			return err
		}
		return fn(id, msg.(*storage.{{.Type}}))
	}
	return nil
}
{{- else}}

// Walk iterates over all of the objects in the store and applies the closure
func (s *storeImpl) Walk(fn func(obj *storage.{{.Type}}) error) error {
	rows, err := s.walkStmt.Query()
	if err != nil {
		return nilNoRows(err)
	}
	defer rows.Close()
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return err
		}
		msg := alloc()
		buf := bytes.NewBuffer(data)
		if err := jsonpb.Unmarshal(buf, msg); err != nil {
			return err
		}
		return fn(msg.(*storage.{{.Type}}))
	}
	return nil
}
{{- end}}

// AckKeysIndexed acknowledges the passed keys were indexed
func (s *storeImpl) AckKeysIndexed(keys ...string) error {
	return nil
}

// GetKeysToIndex returns the keys that need to be indexed
func (s *storeImpl) GetKeysToIndex() ([]string, error) {
	return nil, nil
}
`

type properties struct {
	Type        string
	Table       string
	NoKeyField  bool
	KeyFunc     string
	UniqKeyFunc string
	Cache       bool
	TrackIndex  bool
}

func main() {
	c := &cobra.Command{
		Use: "generate store implementations",
	}

	var props properties
	c.Flags().StringVar(&props.Type, "type", "", "the (Go) name of the object")
	utils.Must(c.MarkFlagRequired("type"))

	c.Flags().StringVar(&props.Table, "table", "", "the logical table of the objects")
	utils.Must(c.MarkFlagRequired("table"))

	c.Flags().BoolVar(&props.NoKeyField, "no-key-field", false, "whether or not object contains key field. If no, then to key function is not applied on object")
	c.Flags().StringVar(&props.KeyFunc, "key-func", "GetId()", "the function on the object to retrieve the key")
	c.Flags().StringVar(&props.UniqKeyFunc, "uniq-key-func", "", "when set, unique key constraint is added on the object field retrieved by the function")

	c.RunE = func(*cobra.Command, []string) error {
		templateMap := map[string]interface{}{
			"Type":       props.Type,
			"Bucket":     props.Table,
			"NoKeyField": props.NoKeyField,
			"KeyFunc":    props.KeyFunc,
			//"UniqKeyFunc": props.UniqKeyFunc,
			"Table": props.Table,
		}

		t := template.Must(template.New("gen").Parse(autogenerated + storeFile))
		buf := bytes.NewBuffer(nil)
		if err := t.Execute(buf, templateMap); err != nil {
			return err
		}
		if err := ioutil.WriteFile("store.go", buf.Bytes(), 0644); err != nil {
			return err
		}
		return nil
	}
	if err := c.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}