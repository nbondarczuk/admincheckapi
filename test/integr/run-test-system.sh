#!/bin/bash

function run_test() {
	msg=$(curl -X GET -H "Content-Type: application/json" http://localhost:1234/system/alive 2>/dev/null)
	echo Result: $? $msg
	msg=$(curl -X GET -H "Content-Type: application/json" http://localhost:1234/system/health 2>/dev/null)
	echo Result: $? $msg
	msg=$(curl -X GET -H "Content-Type: application/json" http://localhost:1234/system/stat 2>/dev/null)
	echo Result: $? $msg
	msg=$(curl -X GET -H "Content-Type: application/json" http://localhost:1234/system/version 2>/dev/null)
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
