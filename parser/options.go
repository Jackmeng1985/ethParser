package parser

import (
	"github.com/Jackmeng1985/ethParser"
)

type Option func(e *EthParser) error

func WithLog(log ethParser.Logger) Option {
	return func(e *EthParser) error {
		e.config.log = log
		return nil
	}
}

func WithDatabase(db ethParser.Database) Option {
	return func(e *EthParser) error {
		e.config.db = db
		return nil
	}
}

func WithEthClient(client ethParser.EthClient) Option {
	return func(e *EthParser) error {
		e.config.ethClient = client
		return nil
	}
}
