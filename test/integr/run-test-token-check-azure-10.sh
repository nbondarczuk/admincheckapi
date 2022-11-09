#!/bin/bash

function run_test() {
	msg=$(curl -X POST -H "Content-Type: application/json" -d@testdata/cocacola-device-jwt-token.json http://localhost:1234/api/client/COCACOLA/admin/token 2>/dev/null)
	echo Result: $? $msg
}

n=${1:-10}

while true
do
	if test $n -gt 0
	then
		run_test
	else
		break
	fi
	let n=$((n - 1))
done	
