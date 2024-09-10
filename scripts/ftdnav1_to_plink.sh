#!/bin/bash

# Directory containing the files
DIR="$(pwd)"

for file in "$DIR"/*.csv; do
    if [ -f "$file" ]; then
        # Create the output file name by appending "_converted" before the file extension
        output_file="${file%.csv}_converted.txt"

        echo "Converting $file to $output_file"

        # Run the conversion command with the new output file name
        ./terraseq convert --inFile "$file" --inFormat ftdnav1 --outFile "$output_file" --outFormat 23andme
    fi
done

for file in "$DIR"/*.txt; do
    if [ -f "$file" ]; then
        # Create the output file name by appending "_converted" before the file extension
        output_file="${file%.txt}_plink"

        echo "Converting $file to $output_file"

        # Run the conversion command with the new output file name
        ./plink --23file "$file" --make-bed --out "$output_file"
    fi
done

echo "All files have been converted"
