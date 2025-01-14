package standalone_storage

import (
	"github.com/Connor1996/badger"
	"github.com/pingcap-incubator/tinykv/kv/config"
	"github.com/pingcap-incubator/tinykv/kv/storage"
	"github.com/pingcap-incubator/tinykv/kv/util/engine_util"
	"github.com/pingcap-incubator/tinykv/proto/pkg/kvrpcpb"
)

// StandAloneStorage is an implementation of `Storage` for a single-node TinyKV instance. It does not
// communicate with other nodes and all data is stored locally.
type StandAloneStorage struct {
	// Your Data Here (1).
	db *badger.DB
}

func NewStandAloneStorage(conf *config.Config) *StandAloneStorage {
	// Your Code Here (1).
	return &StandAloneStorage{
		db: engine_util.CreateDB(conf.DBPath, false),
	}
}

func (s *StandAloneStorage) Start() error {
	// Your Code Here (1).
	return nil
}

func (s *StandAloneStorage) Stop() error {
	// Your Code Here (1).
	return s.db.Close()
}

func (s *StandAloneStorage) Reader(ctx *kvrpcpb.Context) (storage.StorageReader, error) {
	// Your Code Here (1).
	tx := s.db.NewTransaction(false)
	return &StandAloneStorageReader{
		tx,
	}, nil
}

func (s *StandAloneStorage) Write(ctx *kvrpcpb.Context, batch []storage.Modify) error {
	// Your Code Here (1).
	tx := s.db.NewTransaction(true)
	for _, m := range batch {
		switch m.Data.(type) {
		case storage.Put:
			put := m.Data.(storage.Put)
			if err := tx.Set(engine_util.KeyWithCF(put.Cf, put.Key), put.Value); err != nil {
				return err
			}
		case storage.Delete:
			delete := m.Data.(storage.Delete)
			if err := tx.Delete(engine_util.KeyWithCF(delete.Cf, delete.Key)); err != nil {
				return err
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

type StandAloneStorageReader struct {
	tx *badger.Txn
}

func (r *StandAloneStorageReader) GetCF(cf string, key []byte) ([]byte, error) {
	val, err := engine_util.GetCFFromTxn(r.tx, cf, key)
	if err == badger.ErrKeyNotFound {
		return nil, err
	}
	return val, err
}

func (r *StandAloneStorageReader) Close() {
	r.tx.Discard()
}

func (r *StandAloneStorageReader) IterCF(cf string) engine_util.DBIterator {
	return engine_util.NewCFIterator(cf, r.tx)
}
