#!/bin/bash
cd ./grpc/worker/
./worker $(hostname) localhost 50051 $(hostname -I) 50061
