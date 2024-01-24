#!/bin/bash

# Function to stop all background processes
stop_processes() {
    echo "Stopping all processes..."
    kill -- -$$ # Sends a signal to the entire process group
}

# Trap the EXIT signal to call the stop_processes function
trap 'stop_processes' EXIT

AMOUNT_WORKERS=$1
# Start your programs in the background
cd ./grpc/master/
./master $1 & 
echo starting master with ${AMOUNT_WORKERS} workers
cd ../worker/
for (( i=1; i<=$AMOUNT_WORKERS; i++ ))
do
    echo starting worker $i at port 5006"$i"
    sleep 1s
    ./worker $i 50051 5006"$i" &
done

# Keep the script running
read -rp "Press Enter to stop all processes..."
