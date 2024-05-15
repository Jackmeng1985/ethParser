package memDB

import (
	"github.com/Jackmeng1985/ethParser"
	"sync"
)

// InMemoryDB is an in-memory implementation of the Database interface
type InMemoryDB struct {
	mu           sync.RWMutex
	transactions map[string][]*ethParser.Transaction // Key is the address
}

func New() *InMemoryDB {
	return &InMemoryDB{
		transactions: make(map[string][]*ethParser.Transaction),
	}
}

func (db *InMemoryDB) AddTransaction(tx *ethParser.Transaction) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.transactions[tx.From] = append(db.transactions[tx.From], tx)
	db.transactions[tx.To] = append(db.transactions[tx.To], tx)

	return nil
}

func (db *InMemoryDB) GetTransactionsByAddress(address string) (ethParser.Transactions, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return append(db.transactions[address]), nil
}
