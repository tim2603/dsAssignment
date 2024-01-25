#!/bin/bash
cd ./grpc/worker/
CGO_ENABLED=0 go build
cd ../master/
CGO_ENABLED=0 go build