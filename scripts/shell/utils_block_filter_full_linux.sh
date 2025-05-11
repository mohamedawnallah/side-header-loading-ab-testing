#!/bin/bash
# utils_block_filter_full_linux.sh
# Linux-compatible, all files in $HOME/.neutrino/data by default

#-------------------------------
# Configurable data directory and file paths
#-------------------------------

# Data directory (override with NEUTRINO_DATA_DIR environment variable if desired)
NEUTRINO_DATA_DIR="${NEUTRINO_DATA_DIR:-$HOME/.neutrino/data}"

# Header files
BLOCK_HEADER_FILE="$NEUTRINO_DATA_DIR/block_headers.bin"
FILTER_HEADER_FILE="$NEUTRINO_DATA_DIR/reg_filter_headers.bin"

# Temporary files (always in data dir)
TMP1="$NEUTRINO_DATA_DIR/part1.bin"
TMP2="$NEUTRINO_DATA_DIR/part2.bin"
BLOCK_TMP_OUT="$NEUTRINO_DATA_DIR/block_headers_trimmed.bin"
FILTER_TMP_OUT="$NEUTRINO_DATA_DIR/reg_filter_headers_trimmed.bin"

#-------------------------------
# Print Nth Block Header
#-------------------------------

print_nth_block_header_from_start() {
    local n="$1"
    local skip=$((n - 1))
    dd if="$BLOCK_HEADER_FILE" bs=80 skip="$skip" count=1 2>/dev/null | xxd -p | tr -d '\n'; echo
}

print_nth_block_header_from_end() {
    local n="$1"
    local total
    total=$(stat -c %s "$BLOCK_HEADER_FILE")
    local count=$(( total / 80 ))
    local skip=$(( count - n ))
    dd if="$BLOCK_HEADER_FILE" bs=80 skip="$skip" count=1 2>/dev/null | xxd -p | tr -d '\n'; echo
}

#-------------------------------
# Remove Nth Block Header
#-------------------------------

remove_nth_block_header_from_start() {
    local n="$1"
    local before=$((n - 1))
    dd if="$BLOCK_HEADER_FILE" of="$TMP1" bs=80 count="$before" 2>/dev/null
    dd if="$BLOCK_HEADER_FILE" of="$TMP2" bs=80 skip="$n" 2>/dev/null
    cat "$TMP1" "$TMP2" > "$BLOCK_TMP_OUT"
    mv "$BLOCK_TMP_OUT" "$BLOCK_HEADER_FILE"
    rm -f "$TMP1" "$TMP2"
}

remove_nth_block_header_from_end() {
    local n="$1"
    local total
    total=$(stat -c %s "$BLOCK_HEADER_FILE")
    local count=$(( total / 80 ))
    local target=$((count - n + 1))
    remove_nth_block_header_from_start "$target"
}

#-------------------------------
# Print Nth Filter Header
#-------------------------------

print_nth_filter_header_from_start() {
    local n="$1"
    local skip=$((n - 1))
    dd if="$FILTER_HEADER_FILE" bs=32 skip="$skip" count=1 2>/dev/null | xxd -p | tr -d '\n'; echo
}

print_nth_filter_header_from_end() {
    local n="$1"
    local total
    total=$(stat -c %s "$FILTER_HEADER_FILE")
    local count=$(( total / 32 ))
    local skip=$(( count - n ))
    dd if="$FILTER_HEADER_FILE" bs=32 skip="$skip" count=1 2>/dev/null | xxd -p | tr -d '\n'; echo
}

#-------------------------------
# Remove Nth Filter Header
#-------------------------------

remove_nth_filter_header_from_start() {
    local n="$1"
    local before=$((n - 1))
    dd if="$FILTER_HEADER_FILE" of="$TMP1" bs=32 count="$before" 2>/dev/null
    dd if="$FILTER_HEADER_FILE" of="$TMP2" bs=32 skip="$n" 2>/dev/null
    cat "$TMP1" "$TMP2" > "$FILTER_TMP_OUT"
    mv "$FILTER_TMP_OUT" "$FILTER_HEADER_FILE"
    rm -f "$TMP1" "$TMP2"
}

remove_nth_filter_header_from_end() {
    local n="$1"
    local total
    total=$(stat -c %s "$FILTER_HEADER_FILE")
    local count=$(( total / 32 ))
    local target=$((count - n + 1))
    remove_nth_filter_header_from_start "$target"
}

#-------------------------------
# Example Usage (not executed)
#-------------------------------
# print_nth_block_header_from_start 3
# print_nth_block_header_from_end 2
# remove_nth_block_header_from_start 5
# remove_nth_block_header_from_end 1
# print_nth_filter_header_from_start 7
# print_nth_filter_header_from_end 2
# remove_nth_filter_header_from_start 4
# remove_nth_filter_header_from_end 3

#-------------------------------
# End of utils_block_filter_full_linux.sh
