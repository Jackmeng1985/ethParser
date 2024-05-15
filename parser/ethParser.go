package parser

import (
	"context"
	"github.com/Jackmeng1985/ethParser"
)

// config to hold dependencies for the sync service.
type config struct {
	log       ethParser.Logger
	db        ethParser.Database
	ethClient ethParser.EthClient
}

// New creates a new HttpClient instance
func New(ctx context.Context, opts ...Option) *EthParser {
	r := &EthParser{
		ctx:     ctx,
		config:  &config{},
		txsChan: make(chan *ethParser.Transaction),
	}
	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil
		}
	}
	return r
}

type EthParser struct {
	config *config
	//
	ctx     context.Context
	txsChan chan *ethParser.Transaction
}

func (e *EthParser) loop() {
	for {
		select {
		case <-e.ctx.Done():
			break
		case tx := <-e.txsChan:
			if err := e.config.db.AddTransaction(tx); err != nil {
				e.config.log.Error("Cannot insert a record to db err: %v", err)
			} else {
				//e.config.log.Info("Insert a txs record to DB From:%s To:%s", tx.From, tx.To)
			}

		}
	}
}

func (e *EthParser) Start() error {
	if err := e.config.ethClient.Start(e.ctx); err != nil {
		return err
	}
	go e.loop()
	return nil
}

func (e *EthParser) GetCurrentBlock() int {
	//TODO implement me
	panic("implement me")
}

func (e *EthParser) Subscribe(address string) bool {
	if err := e.config.ethClient.SubscribeTransaction(address, e.txsChan); err != nil {
		e.config.log.Warn("Cannot subscribe txs feed Err:%v", err)
		return false
	}
	return true
}

func (e *EthParser) GetTransactions(address string) (txs ethParser.Transactions) {
	var err error
	if txs, err = e.config.db.GetTransactionsByAddress(address); err != nil {
		e.config.log.Error("Cannot get txs from db err:%v, address: %s", err, address)
		return nil
	}
	return txs
}
