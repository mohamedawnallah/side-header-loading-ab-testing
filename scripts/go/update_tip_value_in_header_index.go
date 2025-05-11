package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	bolt "go.etcd.io/bbolt"
)

// parseHexArg parses a hex string argument (starting with "0x") into a byte slice.
func parseHexArg(arg string) ([]byte, error) {
	arg = strings.TrimPrefix(arg, "0x")
	return hex.DecodeString(arg)
}

func main() {
	// Determine HOME for default path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to get home directory: %v", err)
	}
	defaultDBPath := filepath.Join(homeDir, ".neutrino", "data", "filters.db")

	dataPath := flag.String("data", defaultDBPath, "Path to filters.db (default $HOME/.neutrino/data/filters.db)")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] <key> <valueHex>\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "Flags:")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nExample: %s mykey 0x68656c6c6f\n", os.Args[0])
	}
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	key := []byte(flag.Arg(0))

	value, err := parseHexArg(flag.Arg(1))
	if err != nil {
		log.Fatalf("Invalid valueHex: %v\n", err)
	}

	db, err := bolt.Open(*dataPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("header-index"))
		if b == nil {
			return fmt.Errorf("bucket header-index not found")
		}
		val := b.Get(key)
		if val == nil {
			return fmt.Errorf("key not found")
		}
		if err := b.Put(key, value); err != nil {
			return fmt.Errorf("failed to update key: %v", err)
		}
		fmt.Println("Updated successfully")
		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
	}
}
