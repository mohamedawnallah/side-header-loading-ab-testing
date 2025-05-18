def remove_bytes_from_end(file_path, n_bytes):
    """
    Remove n bytes from the end of a binary file by modifying it in place.
    
    Args:
        file_path (str): Path to the file
        n_bytes (int): Number of bytes to remove from the end
    """
    try:
        # Get the current file size
        with open(file_path, 'rb') as f:
            f.seek(0, 2)  # Seek to the end of the file
            file_size = f.tell()  # Get current position (file size)
        
        if n_bytes >= file_size:
            raise ValueError(f"Cannot remove {n_bytes} bytes: file is only {file_size} bytes long")
        
        # Calculate new size
        new_size = file_size - n_bytes
        
        # Truncate the file to the new size
        with open(file_path, 'r+b') as f:
            f.truncate(new_size)
            
        print(f"Successfully removed {n_bytes} bytes from {file_path}")
        print(f"Original size: {file_size} bytes, New size: {new_size} bytes")
        
    except Exception as e:
        print(f"Error: {e}")

# Example usage
if __name__ == "__main__":
    import sys
    
    if len(sys.argv) != 3:
        print("Usage: python script.py file_path n_bytes")
    else:
        file_path = sys.argv[1]
        try:
            n_bytes = int(sys.argv[2])
            remove_bytes_from_end(file_path, n_bytes)
        except ValueError as e:
            print(f"Error: {e}")
