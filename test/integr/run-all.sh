#!/bin/bash

n=${1:-1}

./run-test-system.sh $n
./run-test-create-read-delete.sh $n
./run-test-token-check-azure.sh $n
./run-test-token-check.sh $n

