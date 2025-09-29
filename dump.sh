# The name of the file where all the code will be saved.
OUTPUT_FILE="kitkat_project_code.txt"

echo "🔍 Consolidating all Go source code into a single file..."

# Remove the output file if it already exists to start fresh.
rm -f "$OUTPUT_FILE"

# Find all files ending in .go, then loop through each one.
# The `find . -type f -name "*.go"` command ensures we only get files.
find . -type f -name "*.go" | sort | while read -r filepath; do
    
    # Print a clear header with the file path to the output file.
    echo "
# $filepath
" >> "$OUTPUT_FILE"
    
    # Append the actual content of the file.
    cat "$filepath" >> "$OUTPUT_FILE"
    
done

echo "✅ Success! All code has been written to $OUTPUT_FILE."