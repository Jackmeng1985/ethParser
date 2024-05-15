package ethParser

import "context"

type Parser interface {
	// GetCurrentBlock last parsed block
	GetCurrentBlock() int
	// Subscribe add address to observer
	Subscribe(address string) bool
	// GetTransactions list of inbound or outbound transactions for an address
	GetTransactions(address string) Transactions
}

type EthClient interface {
	Start(ctx context.Context) error
	SubscribeTransaction(address string, results chan<- *Transaction) error
}

// Database defines the interface for our transaction database
type Database interface {
	AddTransaction(tx *Transaction) error
	GetTransactionsByAddress(address string) (Transactions, error)
}

type Logger interface {
	Info(msg string, ctx ...interface{})
	Warn(msg string, ctx ...interface{})
	Error(msg string, ctx ...interface{})
}
