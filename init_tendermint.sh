#!/bin/bash

# Declare variables
file_name="config.toml"

# Give this script to read hidden folder or file ( eg .storage,.something )
shopt -s dotglob

# Go to tendermint config folder
cd .storage/tendermint/config/ || exit

# Shall we begin ?
# Edit create_empty_blocks to false
sed -i '' "s|create_empty_blocks = false|create_empty_blocks = true|g" "$file_name"

# Edit proxy_app to tcp://abci:26658
sed -i '' "s|proxy_app = \"kvstore\"|proxy_app = \"tcp://abci:26658\"|g" "$file_name"

# Edit wal_file to .storage/tendermint/data/cs.wal/wal
sed -i '' "s|wal_file = \"data/cs.wal/wal\"|wal_file = \".storage/tendermint/data/cs.wal/wal\"|g" "$file_name"

# Edit laddr to tcp://0.0.0.0:26657
sed -i '' "s|laddr = \"tcp://0.0.0.0:26657\"|laddr = \"tcp://0.0.0.0:26657\"|g" "$file_name"
