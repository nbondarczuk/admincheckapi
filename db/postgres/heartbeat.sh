#!/bin/bash

SLEEP_TIME=1
DBNAME=argonadmindb
PORT=5432
HOST=postgresdb.admincheckapi-dev-default.svc.cluster.local
USERNAME=test
TIMEOUT=5

export PGPASSWORD=test

echo  $(date) Starting
while true
do
	msg=$(pg_isready --dbname=$DBNAME --host=$HOST --port=$PORT --timeout=$TIMEOUT --username=$USERNAME)
	echo $(date) $DBNAME $PORT $msg
	echo $(date) $DBNAME $PORT pgstat:
	pgstat -h $HOST -p $PORT -U $USERNAME -d $DBNAME 1 1
	sleep $SLEEP_TIME
done

