package storage

import (
	"encoding/json"
	"fmt"

	"github.com/LICODX/PoSSR-RNRCORE/internal/params"
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Store struct {
	db *leveldb.DB
}

// NewLevelDB creates a new Store instance
func NewLevelDB(path string) (*Store, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

// GenerateBlockBodyKey creates a key for block body
func GenerateBlockBodyKey(height uint64) []byte {
	return []byte(fmt.Sprintf("block-body-%d", height))
}

// SaveBlock saves a block to the database
func (s *Store) SaveBlock(block types.Block) error {
	// 1. Save Body
	bodyKey := GenerateBlockBodyKey(block.Header.Height)
	bodyData, _ := json.Marshal(block.Shards) // Simple JSON serialization for prototype
	if err := s.db.Put(bodyKey, bodyData, nil); err != nil {
		return err
	}

	// 2. Save Header (TODO: Separate bucket)
	headerKey := []byte(fmt.Sprintf("block-header-%d", block.Header.Height))
	headerData, _ := json.Marshal(block.Header)
	return s.db.Put(headerKey, headerData, nil)
}

// PruneOldBlocks dipanggil setiap kali blok baru ditambahkan
func (s *Store) PruneOldBlocks(currentHeight uint64) error {
	if currentHeight <= params.PruningWindow {
		return nil
	}

	// Target blok yang harus dihapus (N - 25)
	targetHeight := currentHeight - params.PruningWindow

	// Hapus BODY blok (transaksi raw 1 GB), tapi simpan HEADER
	key := GenerateBlockBodyKey(targetHeight)

	// LevelDB Delete (Fast I/O)
	err := s.db.Delete(key, nil)
	if err != nil {
		return err
	}

	// Lakukan CompactRange secara berkala untuk membebaskan disk space fisik
	if targetHeight%100 == 0 {
		s.db.CompactRange(util.Range{Start: nil, Limit: nil})
	}

	return nil
}

// GetDB returns the underlying LevelDB instance
func (s *Store) GetDB() *leveldb.DB {
	return s.db
}

// SaveTip saves the current chain tip
func (s *Store) SaveTip(tipData []byte) error {
	return s.db.Put([]byte("tip"), tipData, nil)
}

// GetTip loads the current chain tip
func (s *Store) GetTip() ([]byte, error) {
	return s.db.Get([]byte("tip"), nil)
}

// HasBlock checks if a block exists at the given height
func (s *Store) HasBlock(height uint64) bool {
	key := []byte(fmt.Sprintf("block-header-%d", height))
	_, err := s.db.Get(key, nil)
	return err == nil
}

// GetBlockHeaderByHeight retrieves a block header by its height
func (s *Store) GetBlockHeaderByHeight(height uint64) (*types.BlockHeader, error) {
	key := []byte(fmt.Sprintf("block-header-%d", height))
	data, err := s.db.Get(key, nil)
	if err != nil {
		return nil, fmt.Errorf("block header not found for height %d: %v", height, err)
	}

	var header types.BlockHeader
	if err := json.Unmarshal(data, &header); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block header: %v", err)
	}

	return &header, nil
}
