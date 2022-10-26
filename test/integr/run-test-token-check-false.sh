#!/bin/bash

function run_test() {
	msg=$(curl -X POST -H "Content-Type: application/json" http://localhost:1234/api/client/purge 2>/dev/null)
	echo Result: $? $msg
	msg=$(curl -X POST -H "Content-Type: application/json" -d@testdata/jwt-token.json http://localhost:1234/api/client/XXX/admin/token 2>/dev/null)
	echo Result: $? $msg
}

n=${1:-1}

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

	
