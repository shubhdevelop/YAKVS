#!/bin/bash

# Test script to demonstrate YAKVS commands
echo "Testing YAKVS Commands"
echo "======================"

# Build the project
echo "Building YAKVS..."
go build -o yakvs

if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

echo "Build successful!"
echo ""

# Test commands using normal command format
echo "Testing basic commands:"
echo ""

# Create a temporary input file with commands
cat > test_input.txt << EOF
SET key1 value1
GET key1
EXISTS key1
TTL key1
EXPIRE key1 3600
TTL key1
DEL key1
GET key1
EXISTS key1
exit
EOF

echo "Running commands from input file..."
./yakvs < test_input.txt

# Clean up
rm test_input.txt

echo ""
echo "Testing completed!"
