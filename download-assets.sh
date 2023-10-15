#!/bin/bash

# URL of the indexed web directory
web_directory="https://eden.sh/versions/"

# File name to download
file_to_download="assets"

# Get the latest file name from the web directory
latest_file=$(curl -s "$web_directory" | grep -o "$file_to_download-.*\.tgz" | sort -V | tail -n 1 | sed 's/.*">//' )

# Check if a file was found
if [ -z "$latest_file" ]; then
    echo "No files found in the web directory."
    exit 1
fi

# Download the latest file
curl -O "${web_directory}${latest_file}"

echo "Downloaded $latest_file"

