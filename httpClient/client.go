package httpClient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Jackmeng1985/ethParser"
	"sync"

	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

type RPCRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type RPCResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
}

// BLOCK_INTERVAL after merge 12 seconds per block
// According to ethstats.dev
const BLOCK_INTERVAL = 11
const BLOCK_FETCH_DELAY = 4

// New creates a new HttpClient instance
func New(logger ethParser.Logger, endpoint string) *HttpClient {
	return &HttpClient{
		log:      logger,
		endpoint: endpoint,
		txsFeeds: make(map[string]chan<- *ethParser.Transaction),
	}
}

type HttpClient struct {
	ctx context.Context
	log ethParser.Logger
	//
	endpoint           string
	currentBlockNumber uint64
	//
	txsFeeds     map[string]chan<- *ethParser.Transaction
	txsFeedsLock sync.RWMutex
}

func (e *HttpClient) loop() {

	timer := time.NewTimer(2 * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-e.ctx.Done():
			return
		case <-timer.C:
			var block ethParser.Block
			if err := e.request(&block, "eth_getBlockByNumber", ethParser.HexNumber(e.currentBlockNumber+1), true); err != nil {
				e.log.Warn("Cannot fetch eth_getBlockByNumber from endpoint err: %v", err)
				timer.Reset(2 * time.Second)
			} else {
				e.log.Info("Fetch block from endpoint blockNr: %v, txs: %d", block.Number, len(block.Transactions))
				e.currentBlockNumber = uint64(block.Number)
				delay := time.Unix(int64(block.Timestamp)+BLOCK_INTERVAL+BLOCK_FETCH_DELAY, 0).Sub(time.Now())
				timer.Reset(delay)
				//
				e.sendTxsFeeds(block.Transactions)
			}
		}
	}
}

func (e *HttpClient) request(result interface{}, method string, args ...interface{}) error {

	if result != nil && reflect.TypeOf(result).Kind() != reflect.Ptr {
		return fmt.Errorf("call result parameter must be pointer or nil interface: %v", result)
	}

	// Build RPC request
	request := RPCRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  args,
		ID:      1, // request ID
	}
	// Convert request to JSON
	requestData, err := json.Marshal(request)
	if err != nil {
		return err
	}
	// Send POST request to Ethereum node
	response, err := http.Post(e.endpoint, "application/json", bytes.NewBuffer(requestData))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// Parse JSON response
	var rpcResponse RPCResponse
	err = json.Unmarshal(body, &rpcResponse)
	if err != nil {
		e.log.Error("Error: %v,  body: %s", err, body)
		return err
	}
	return json.Unmarshal(rpcResponse.Result, result)
}

func (e *HttpClient) sendTxsFeeds(transactions ethParser.Transactions) {
	e.txsFeedsLock.RLock()
	defer e.txsFeedsLock.RUnlock()

	for _, transaction := range transactions {
		//incoming
		if sub, ok := e.txsFeeds[transaction.To]; ok {
			//e.log.Info(fmt.Sprintf("send feed tx.To: %s", transaction.To))
			sub <- transaction
		}
		//outgoing
		if sub, ok := e.txsFeeds[transaction.From]; ok {
			//e.log.Info(fmt.Sprintf("send feed tx.From: %s", transaction.From))
			sub <- transaction
		}
		//or call contract？ todo v2
		//if sub, ok := e.txsFeeds[transaction.Data]; ok {
		//	sub <- transaction
		//}
	}
}

func (e *HttpClient) Start(ctx context.Context) error {
	e.ctx = ctx
	//
	var bn ethParser.HexNumber
	if err := e.request(&bn, "eth_blockNumber"); err != nil {
		return fmt.Errorf("cannot fetch current block number err: %v", err)
	}
	//
	e.currentBlockNumber = uint64(bn)
	e.log.Info(fmt.Sprintf("Currnet block number: %d", e.currentBlockNumber))
	//
	go e.loop()
	return nil
}

func (e *HttpClient) SubscribeTransaction(address string, results chan<- *ethParser.Transaction) error {
	e.txsFeedsLock.Lock()
	defer e.txsFeedsLock.Unlock()

	if _, ok := e.txsFeeds[address]; ok {
		return fmt.Errorf("adress: %s alreay have a subscribe chan", address)
	}
	e.log.Info("Address %s has subscribe the transaction feed", address)
	e.txsFeeds[address] = results
	return nil
}
