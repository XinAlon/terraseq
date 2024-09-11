#!/bin/bash

# Directory containing the files
DIR="$(pwd)"

merge_file="$DIR/merge_list.txt"
> merge_file

for file in "$DIR"/*.bed; do
    if [ -f "$file" ]; then
        base_file="${file%.bed}"  # Remove the .bed extension
        echo "$base_file" >> "$merge_file"
    fi
done

if [ -f "$merge_file" ]; then
    echo "Merging ..."
    ./plink --merge-list "$merge_file" --out "$DIR/merged_output"
else
    echo "List file $merge_file does not exist. No files to merge."
fi

echo "All files have been merged"
