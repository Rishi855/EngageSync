#!/bin/sh
# wait-for host:port -- command...
hostport="$1"
shift
host=$(echo $hostport | cut -d: -f1)
port=$(echo $hostport | cut -d: -f2)

until nc -z $host $port; do
  echo "Waiting for $host:$port..."
  sleep 1
done

exec "$@"
