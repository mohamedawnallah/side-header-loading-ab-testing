import hashlib
import argparse
import sys

def double_sha256(data: bytes) -> bytes:
    """Perform double SHA256 hashing on the given data."""
    return hashlib.sha256(hashlib.sha256(data).digest()).digest()

def main():
    parser = argparse.ArgumentParser(
        description="Calculate filter header using double SHA256 of filter + previous header."
    )

    # Positional arguments (required by default)
    parser.add_argument(
        "filter",
        type=str,
        help="Hex-encoded filter (little-endian input)"
    )

    parser.add_argument(
        "prev_header",
        type=str,
        help="Hex-encoded previous filter header (little-endian input)"
    )

    args = parser.parse_args()

    try:
        # Step 1: Convert hex filter to bytes.
        filter_bytes = bytes.fromhex(args.filter)

        # Step 2: Compute filter hash (double SHA256)
        filter_hash = double_sha256(filter_bytes)

        # Step 3: Convert previous header from hex to bytes (little-endian, as is)
        prev_header_bytes = bytes.fromhex(args.prev_header)

        # Step 4: Concatenate filter_hash + prev_header
        combined = filter_hash + prev_header_bytes

        # Step 5: Compute filter header
        filter_header = double_sha256(combined)

        # Step 6: Output the result
        print("Computed Filter Header:")
        print(" - Little-endian:", filter_header.hex())
        print(" - Big-endian   :", filter_header[::-1].hex())

    except ValueError as e:
        print("Error: Invalid hex input -", e)
        sys.exit(1)

if __name__ == "__main__":
    main()
