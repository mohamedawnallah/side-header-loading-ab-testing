package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btclog"
	"github.com/btcsuite/btcwallet/walletdb"
	_ "github.com/btcsuite/btcwallet/walletdb/bdb" // <-- ADD THIS!
	"github.com/lightninglabs/neutrino"
	"github.com/lightninglabs/neutrino/headerfs"
)

var shouldPersistToDisk = flag.Bool("persist_to_disk", false, "Persist the filter to the disk")
var shouldRescan = flag.Bool("rescan", false, "Perform blockchain rescan for script pubkey")
var scriptPubKeyHex = flag.String("scriptpubkey", "", "Hex-encoded script pubkey to scan for")

func setupLogger(logPath string, level btclog.Level) (btclog.Logger, *os.File) {
	// Open or create the log file
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Create a multi-writer to write to both stdout and the log file
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Set up the neutrino logger
	backend := btclog.NewBackend(multiWriter)
	logger := backend.Logger("NEUTRINO")
	logger.SetLevel(level)

	return logger, logFile
}

func main() {
	flag.Parse()
	home := os.Getenv("HOME")
	logPath := fmt.Sprintf("%s/.neutrino/neutrino.log", home)
	logger, logFile := setupLogger(logPath, btclog.LevelDebug)
	defer logFile.Close()

	neutrino.UseLogger(logger)

	// Create a context that will be canceled on program interruption
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up channel to catch interrupt signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println("\nReceived shutdown signal. Shutting down...")
		cancel()
	}()

	// Get home directory for data storage
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		return
	}
	neutrinoDataDir := filepath.Join(homeDir, ".neutrino", "data")

	if err := os.MkdirAll(neutrinoDataDir, 0700); err != nil {
		fmt.Printf("Error creating neutrino data directory: %v\n", err)
		return
	}

	filterDbPath := filepath.Join(neutrinoDataDir, "filters.db")
	db, err := walletdb.Create("bdb", filterDbPath, true, 5*time.Second, false)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Configure the Neutrino client
	config := neutrino.Config{
		DataDir:       neutrinoDataDir,
		Database:      db,
		ChainParams:   chaincfg.TestNet4Params,
		PersistToDisk: *shouldPersistToDisk,
	}

	// Create the chain service
	fmt.Println("Creating Neutrino chain service...")
	chainService, err := neutrino.NewChainService(config)
	if err != nil {
		fmt.Printf("Error creating chain service: %v\n", err)
		return
	}

	// Start the chain service
	fmt.Println("Starting Neutrino chain service...")
	if err := chainService.Start(); err != nil {
		fmt.Printf("Error starting chain service: %v\n", err)
		return
	}

	// Make sure to properly shut down the service when done
	defer func() {
		fmt.Println("Stopping chain service...")
		chainService.Stop()
		fmt.Println("Chain service stopped.")
	}()

	// Display initial sync info
	fmt.Println("Neutrino chain service started. Syncing headers...")

	// Create a ticker to periodically check sync status
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Main loop to keep program running and show sync status
syncLoop:
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			bestBlock, err := chainService.BestBlock()
			if err != nil {
				fmt.Printf("Error getting best block: %v\n", err)
				continue
			}

			fmt.Printf("Error getting best block: %v\n", err)

			if bestBlock.Height >= 81373 {
				fmt.Println("Chain is synchronized enough to begin scanning")
				break syncLoop
			}
		}
	}

	if *shouldRescan {
		if *scriptPubKeyHex == "" {
			fmt.Println("Please provide a script pubkey in hex using -scriptpubkey")
			return
		}
		fmt.Println("Starting blockchain rescan for script pubkey...")
		err = performRescan(ctx, chainService, *scriptPubKeyHex)
		if err != nil {
			fmt.Printf("Error during rescan: %v\n", err)
			return
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
			bestBlock, err := chainService.BestBlock()
			if err == nil {
				fmt.Printf("Current best block: %d, hash: %s\n",
					bestBlock.Height, bestBlock.Hash)
			}
		}
	}

}

func performRescan(ctx context.Context, chainService *neutrino.ChainService, scriptPubKeyHex string) error {
	// Replace this with your actual script pubkey
	// Example: This could be a P2PKH, P2SH, P2WPKH, etc.
	scriptPubKey, err := hex.DecodeString(scriptPubKeyHex)
	if err != nil {
		return fmt.Errorf("failed to decode script pubkey: %v", err)
	}

	// Define the starting block height for the rescan
	// In a real application, you might want to start from a specific block height
	// related to when the address was created or from a checkpoint
	startHeight := int64(0) // Start from genesis block

	// Get the starting block hash
	startHash, err := chainService.GetBlockHash(startHeight)
	if err != nil {
		return fmt.Errorf("failed to get starting block hash: %v", err)
	}

	startBlockStamp := &headerfs.BlockStamp{
		Height: int32(startHeight),
		Hash:   *startHash,
	}

	// Get current best block to know where to stop the rescan
	bestBlock, err := chainService.BestBlock()
	if err != nil {
		return fmt.Errorf("failed to get best block: %v", err)
	}

	endBlockStamp := &headerfs.BlockStamp{
		Height: bestBlock.Height,
		Hash:   bestBlock.Hash,
	}

	fmt.Printf("Starting rescan from block %d to %d for script pubkey: %x\n",
		startHeight, bestBlock.Height, scriptPubKey)

	// Set up notification handlers.
	ntfnHandlers := rpcclient.NotificationHandlers{
		OnFilteredBlockConnected:    func(height int32, header *wire.BlockHeader, txs []*btcutil.Tx) {},
		OnFilteredBlockDisconnected: func(height int32, header *wire.BlockHeader) {},
		OnRecvTx:                    func(tx *btcutil.Tx, details *btcjson.BlockDetails) {},
		OnRedeemingTx:               func(tx *btcutil.Tx, details *btcjson.BlockDetails) {},
	}

	// Create the set of watch criteria.
	watchList := []neutrino.RescanOption{
		neutrino.StartBlock(startBlockStamp),
		neutrino.EndBlock(endBlockStamp),
		neutrino.WatchAddrs(&btcutil.AddressPubKeyHash{}),
		neutrino.NotificationHandlers(ntfnHandlers),
	}

	// Create a done channel for synchronization
	doneChan := make(chan error, 1)

	// Start the rescan in a goroutine
	go func() {
		scan := neutrino.NewRescan(
			&neutrino.RescanChainSource{
				ChainService: chainService,
			},
			watchList...,
		)
		doneChan <- <-scan.Start()
	}()

	// Wait for the rescan to complete or be canceled.
	select {
	case err := <-doneChan:
		if err != nil {
			return fmt.Errorf("rescan error: %v", err)
		}
		fmt.Println("Rescan completed successfully")
	case <-ctx.Done():
		fmt.Println("Rescan was canceled")
	}

	return nil
}
