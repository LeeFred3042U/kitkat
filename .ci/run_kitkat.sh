#!/usr/bin/bash
set -e

echo "Testing kitkat..."

# Build kitkat binary
go build -o kitkat ./cmd/main.go

# Initialize empty repository for testing
mkdir -p test-repo
cd test-repo
git init

# Configure kitkat (config for this run)
../kitkat config --global user.name "testci"
../kitkat config --global user.email "testci@example.com"

echo "Kitkat setup completed "
