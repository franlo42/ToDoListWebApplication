#!/bin/bash
host="$1"
shift
port="$1"
shift
cmd="$@"

until nc -z -v -w30 $host $port
do
  echo "Esperando que $host:$port esté listo..."
  sleep 1
done
echo "$host:$port está listo, ejecutando el comando..."
exec $cmd