#!/bin/sh

# Kick off postgres
export POSTGRES_DB=dendrite
export POSTGRES_USER=dendrite
export POSTGRES_PASSWORD=itsasecret
./usr/local/bin/docker-entrypoint.sh postgres &

# Wait for postgres to be ready
sleep 5
while ! pg_isready;
do
  sleep 1
done

# Start nginx
echo "Starting Riot"
nginx

# Wait forever
exec tail -f /dev/null
