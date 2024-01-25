#!/bin/bash
cd ./grpc/worker/
./worker $(hostname) 10.166.0.7 50051 $(hostname -I) 50061
