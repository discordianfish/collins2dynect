#!/bin/sh
[ -z "$INTERVAL" ] && INTERVAL=900

if [ -n "$COLLINS_PORT" ]
then
	COLLINS="-collins.url=$(echo $COLLINS_PORT|sed 's/^tcp/http/')/api"
fi


echo "Running every $INTERVAL seconds"
while true
do
	./collins2dynect "$COLLINS" "$@"
	sleep $INTERVAL
done
