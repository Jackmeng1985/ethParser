package main

import (
	"context"
	"github.com/Jackmeng1985/ethParser/httpClient"
	"github.com/Jackmeng1985/ethParser/log"
	"github.com/Jackmeng1985/ethParser/memDB"
	"github.com/Jackmeng1985/ethParser/parser"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	usdt    = "0xdac17f958d2ee523a2206206994597c13d831ec7"
	uniswap = "0x7a250d5630b4cf539739df2c5dacb4c659f2488d"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	logger := log.New(log.INFO)
	ethclient := httpClient.New(logger, "https://cloudflare-eth.com")
	db := memDB.New()

	ethParaser := parser.New(ctx,
		parser.WithDatabase(db),
		parser.WithLog(logger),
		parser.WithEthClient(ethclient),
	)

	if err := ethParaser.Start(); err != nil {
		logger.Error("Cannot start ethClient err: %v", err)
		return
	}

	ethParaser.Subscribe(usdt)
	ethParaser.Subscribe(uniswap)

	wg.Add(1)
	go func() {
		timer := time.NewTimer(1 * time.Minute)
		defer timer.Stop()
		defer wg.Done()
		//
		select {
		case <-ctx.Done():
			break
		case <-timer.C:
			//for _, tx := range txs {
			//	logger.Error("tx from: %s to:%s", tx.From, tx.To)
			//}
			logger.Info("usdt has txs: %d", len(ethParaser.GetTransactions(usdt)))
			logger.Info("uniswap has txs: %d", len(ethParaser.GetTransactions(uniswap)))
			timer.Reset(1 * time.Minute)
		}
	}()

	//shutdown func
	shutdown := func() {
		cancel()
		wg.Wait()
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigc)

	for {
		<-sigc
		logger.Warn("Got interrupt, shutting down...")
		shutdown()
		return
	}
}
