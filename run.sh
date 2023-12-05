#!/bin/bash

cd go
go build -o build/async-bench
cd - > /dev/null

cd tokio
cargo build --release
cd - > /dev/null

cd zmq/build
cmake .. -DCMAKE_BUILD_TYPE=Release
make -j
cd - > /dev/null

echo
echo "#==============================================================#"
echo "#                              go                              #"
echo "#==============================================================#"
echo

time env GOMAXPROCS=4 ./go/build/async-bench

echo
echo "#==============================================================#"
echo "#                             tokio                            #"
echo "#==============================================================#"
echo

time env TOKIO_WORKER_THREADS=4 ./tokio/target/release/async-bench

echo
echo "#==============================================================#"
echo "#                              zmq                             #"
echo "#==============================================================#"
echo

time ./zmq/build/async-bench
