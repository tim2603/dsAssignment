#!/bin/bash

# Function to stop all background processes
stop_processes() {
    echo "Stopping all processes..."
    kill -- -$$ # Sends a signal to the entire process group
}

# Trap the EXIT signal to call the stop_processes function
trap 'stop_processes' EXIT

# Start your programs in the background
go run ./grpc/master/master_main.go & 
for i in {1..5}
do
    go run ./grpc/worker/worker_main.go $i 5000$1 &
done

# Keep the script running
read -rp "Press Enter to stop all processes..."
