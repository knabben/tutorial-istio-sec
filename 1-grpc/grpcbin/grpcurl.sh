#!/usr/bin/env bash

grpcurl -d "{\"f_string\": \"`date`\"}" -plaintext $1:9000 grpcbin.GRPCBin/DummyServerStream
