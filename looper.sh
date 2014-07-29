#!/bin/sh
[ -z "$INTERVAL" ] && INTERVAL=900

if [ -n "$COLLINS_PORT" ]
then
  COLLINS="-collins.url=$(echo $COLLINS_PORT|sed 's/^tcp/http/')/api"
fi


echo "Running every $INTERVAL seconds"
while true
do
  ./collins2dynect $COLLINS "$@" > /tmp/collins2dynect.log
  if [ "$?" -eq "0" ]
  then
    cat /tmp/collins2dynect.log
    [ -n "$PUSHGATEWAY" ] && ( tail -1 /tmp/collins2dynect.log; date +"last_run %s" ) | curl --data-binary @- "$PUSHGATEWAY"
  fi
  tail -1 /tmp/collins2dynect.log
  sleep $INTERVAL
done
