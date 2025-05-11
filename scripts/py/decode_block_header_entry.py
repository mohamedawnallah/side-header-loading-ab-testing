import struct
import hashlib
import sys

# Use the following command to get the last block header entry in block headers binary in hex.
# tail -c 80 block_headers.bin | xxd -p | tr -d '\n'; echo

if len(sys.argv) != 2:
    print(f"Usage: {sys.argv[0]} <header_hex>")
    sys.exit(1)

header_hex = sys.argv[1].strip()
header = bytes.fromhex(header_hex)

if len(header) != 80:
    print("Error: Block header should be exactly 80 bytes (160 hex characters).")
    sys.exit(1)

version, = struct.unpack("<I", header[0:4])
prev_block = header[4:36][::-1].hex()
merkle_root = header[36:68][::-1].hex()
timestamp, = struct.unpack("<I", header[68:72])
bits, = struct.unpack("<I", header[72:76])
nonce, = struct.unpack("<I", header[76:80])

# Calculate block hash.
hash1 = hashlib.sha256(header).digest()
hash2 = hashlib.sha256(hash1).digest()

# Get the Big-Endian representation.
block_hash = hash2[::-1].hex()

print("Version:", version)
print("Previous Block Hash:", prev_block)
print("Merkle Root:", merkle_root)
print("Timestamp (UNIX):", timestamp)
print("Bits:", hex(bits))
print("Nonce:", nonce)
print("Block Hash:", block_hash)
