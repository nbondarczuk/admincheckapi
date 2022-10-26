#!/bin/bash

function crd() {
	echo "Host:$1"
	echo "Port:$2"
	echo "Run: curl -v -X POST -H Content-Type: application/json http://$1:$2/api/client/XXX/admin/group/ABC"
	curl -v -X POST -H "Content-Type: application/json" http://$1:$2/api/client/XXX/admin/group/ABC
	echo "Run: curl -v -X GET -H Content-Type: application/json http://$1:$2/api/client/XXX/admin/group"
	curl -v -X GET -H "Content-Type: application/json" http://$1:$2/api/client/XXX/admin/group
	echo "Run: curl -v -X DELETE -H Content-Type: application/json http://$1:$2/api/client/XXX/admin/group/ABC"
	curl -v -X DELETE -H "Content-Type: application/json" http://$1:$2/api/client/XXX/admin/group/ABC
}

echo Starting
env
echo Starting with setup
echo "ADMINCHECKAPI_HOST: ${ADMINCHECKAPI_HOST:-localhost}"
echo "       SLEEP_CYCLE: ${SLEEP_CYCLE:-60}"

while true
do
	ping -c 1 ${ADMINCHECKAPI_HOST:-localhost}
	crd ${ADMINCHECKAPI_HOST:-localhost} 1234
	sleep ${SLEEP_CYCLE:-60} 
done

