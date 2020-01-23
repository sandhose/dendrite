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

# Generate the keys if they don't already exist.
if [ ! -f /etc/dendrite/server.key ] || \
   [ ! -f /etc/dendrite/server.crt ] || \
   [ ! -f /etc/dendrite/matrix_key.pem ]; then
    echo "Generating keys ..."

    rm -f \
      /etc/dendrite/server.key \
      /etc/dendrite/server.crt \
      /etc/dendrite/matrix_key.pem

    test -f /etc/dendrite/server.key || \
    openssl req -x509 -newkey rsa:4096 \
      -keyout /etc/dendrite/server.key \
      -out /etc/dendrite/server.crt \
      -days 3650 -nodes \
      -subj /CN=localhost

    test -f /etc/dendrite/matrix_key.pem || \
    /usr/local/bin/generate-keys \
      -private-key /etc/dendrite/matrix_key.pem
fi

# Start nginx
echo "Starting Riot"
nginx &

# Start dendrite
echo "Starting Dendrite"
cd /etc/dendrite && dendrite-monolith-server
