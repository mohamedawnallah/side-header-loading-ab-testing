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

// parseBucketArg parses an argument as either a hex string (starting with "0x") or as a plain string byte slice.
func parseBucketArg(arg string) ([]byte, error) {
	if strings.HasPrefix(arg, "0x") || strings.HasPrefix(arg, "0X") {
		return hex.DecodeString(strings.TrimPrefix(arg, "0x"))
	}
	// Use as-is (plain string to bytes)
	return []byte(arg), nil
}

// parseHexArg parses a hex string argument (starting with "0x") into a byte slice.
func parseHexArg(arg string) ([]byte, error) {
	arg = strings.TrimPrefix(arg, "0x")
	return hex.DecodeString(arg)
}

func main() {
	// Set up the data dir flag with default value $HOME/.neutrino/data
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to get home directory: %v", err)
	}
	defaultDBPath := filepath.Join(homeDir, ".neutrino", "data", "filters.db")
	dbPath := flag.String("db", defaultDBPath, "path to filters.db (default $HOME/.neutrino/data/filters.db)")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] <topBucketName> <subBucketHexOrString> <targetKey>\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nExample (hex sub-bucket): %s header-index 0x5929 0x592991c706fbb1af7438103ae6dc80e703343c82256c0c3112d940a100000000\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Example (plain sub-bucket): %s header-index subBucketName 0x592991c706fbb1af7438103ae6dc80e703343c82256c0c3112d940a100000000\n", os.Args[0])
	}

	flag.Parse()

	if flag.NArg() != 3 {
		flag.Usage()
		os.Exit(1)
	}

	topBucketName := flag.Arg(0)
	subBucket, err := parseBucketArg(flag.Arg(1))
	if err != nil {
		log.Fatalf("Invalid subBucket: %v\n", err)
	}
	targetKey, err := parseHexArg(flag.Arg(2))
	if err != nil {
		log.Fatalf("Invalid targetKey: %v\n", err)
	}

	db, err := bolt.Open(*dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		parent := tx.Bucket([]byte(topBucketName))
		if parent == nil {
			return fmt.Errorf("bucket %s not found", topBucketName)
		}
		b := parent.Bucket(subBucket)
		if b == nil {
			return fmt.Errorf("sub-bucket not found")
		}
		val := b.Get(targetKey)
		if val == nil {
			return fmt.Errorf("key not found")
		}
		if err := b.Delete(targetKey); err != nil {
			return fmt.Errorf("failed to delete key: %v", err)
		}
		fmt.Println("Deleted successfully")
		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
	}
}
