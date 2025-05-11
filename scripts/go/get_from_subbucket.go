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
	return []byte(arg), nil
}

// parseHexArg parses a hex string argument (starting with "0x") into a byte slice.
func parseHexArg(arg string) ([]byte, error) {
	arg = strings.TrimPrefix(arg, "0x")
	return hex.DecodeString(arg)
}

func main() {
	// Get $HOME for the default
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to determine home directory: %v", err)
	}
	defaultPath := filepath.Join(home, ".neutrino", "data", "filters.db")

	dataPath := flag.String("data", defaultPath, "path to filters.db (default $HOME/.neutrino/data/filters.db)")

	// Custom usage
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] <topBucket> <subBucketStringOrHex> [keyHex]\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "Flags:")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nExample (all keys): %s header-index 0x5929\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Example (specific key): %s header-index 0x5929 0xabcdef...\n", os.Args[0])
	}
	flag.Parse()

	if flag.NArg() < 2 || flag.NArg() > 3 {
		flag.Usage()
		os.Exit(1)
	}

	topBucket := flag.Arg(0)
	subBucket, err := parseBucketArg(flag.Arg(1))
	if err != nil {
		log.Fatalf("Invalid subBucket: %v\n", err)
	}

	var targetKey []byte
	if flag.NArg() == 3 {
		targetKey, err = parseHexArg(flag.Arg(2))
		if err != nil {
			log.Fatalf("Invalid key: %v\n", err)
		}
	}

	db, err := bolt.Open(*dataPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		parent := tx.Bucket([]byte(topBucket))
		if parent == nil {
			return fmt.Errorf("bucket %s not found", topBucket)
		}
		b := parent.Bucket(subBucket)
		if b == nil {
			return fmt.Errorf("sub-bucket not found")
		}

		if len(targetKey) == 0 {
			// Print all keys/values
			return b.ForEach(func(k, v []byte) error {
				fmt.Printf("Key: %x  Value: %x\n", k, v)
				return nil
			})
		} else {
			// Print specific key
			val := b.Get(targetKey)
			if val == nil {
				fmt.Println("Key not found")
				return nil
			}
			fmt.Printf("Key: %x  Value: %x\n", targetKey, val)
			return nil
		}
	})
	if err != nil {
		log.Fatal(err)
	}
}
