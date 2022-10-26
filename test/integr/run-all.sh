#!/bin/bash

n=${1:-1}

./run-test-system.sh $n
./run-test-create-read-delete.sh $n
./run-test-token-check-false.sh $n
./run-test-token-check-true.sh $n
