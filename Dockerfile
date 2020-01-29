FROM postgres:12.1-alpine AS postgres

#####

FROM vectorim/riot-web AS riot

#####

FROM alpine:latest AS final

RUN apk --update --no-cache add openssl

VOLUME /var/lib/postgresql/data

COPY --from=riot . /
COPY --from=postgres . /

ENV PGDATA=/var/lib/postgresql/data

ADD docker/p2p-entrypoint.sh /usr/local/bin
ADD docker/postgres/create_db.sh /docker-entrypoint-initdb.d/create_db.sh

RUN chmod +x /docker-entrypoint-initdb.d/create_db.sh
RUN sed -i '3i\ \ \ \ application/wasm wasm\;' /etc/nginx/mime.types
RUN adduser --system nginx
RUN addgroup --system nginx
RUN rm -rf /usr/share/nginx/html && ln -s /app /usr/share/nginx/html

CMD sh /usr/local/bin/p2p-entrypoint.sh

EXPOSE 80
EXPOXE 5432
