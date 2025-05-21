#!/usr/bin/env python3

import struct, os

def add_headers_metadata(file_name, header_type):
    # Define the header bytes
    chain_type = 3  # TestNet4ID
    start_height = 0
    
    # Pack into binary format
    header = struct.pack('<BBL', chain_type, header_type, start_height)
    
    # Use home environment variable
    home_dir = os.environ.get('HOME')
    file_path = f'{home_dir}/.neutrino/{file_name}'
    
    # Read existing file
    with open(file_path, 'rb') as f:
        content = f.read()
    
    # Write new file with correct header
    with open(f'{file_path}.new', 'wb') as f:
        f.write(header + content)
    
    # Backup and replace
    os.rename(file_path, f'{file_path}.bak')
    os.rename(f'{file_path}.new', file_path)
    
    print(f'Metadata header inserted successfully. Original file backed up as {file_name}.bak')

# Add metadata headers to both files
add_headers_metadata('block_headers.bin', 0)
add_headers_metadata('reg_filter_headers.bin', 1)
