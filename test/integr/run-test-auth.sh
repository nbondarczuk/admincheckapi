#!/bin/bash

CONTENT_TYPE="Content-Type: application/json"
DATA=testdata/cocacola-device-claims-with-secret.json

function run_test() {
	msg=$(curl -X POST -H $CONTENT_TYPE -d@$DATA http://localhost:1234/api/client/XXX/admin/auth/secret 2>/dev/null)
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
