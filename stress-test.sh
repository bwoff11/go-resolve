#!/bin/bash

# Check if a port number is provided
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <Port Number>"
    exit 1
fi

DNS_PORT=$1

# List of sample domains to query
DOMAINS=("google.com" "facebook.com" "twitter.com" "github.com" "reddit.com" "youtube.com" "linkedin.com" "wikipedia.org" "amazon.com" "instagram.com")

# Infinite loop to continuously send DNS queries
while true; do
    for domain in "${DOMAINS[@]}"; do
        dig @localhost -p 1053 $domain
    done
    sleep 0.2
done
