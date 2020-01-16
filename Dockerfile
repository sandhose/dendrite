FROM postgres:12.1-alpine AS postgres

#####

FROM vectorim/riot-web AS riot

#####

FROM docker.io/golang:1.13.6-alpine3.11 AS golang

ENV GO111MODULE=on

RUN apk --update --no-cache add openssl bash git
RUN go get -v github.com/matrix-org/dendrite/cmd/dendrite-monolith-server@p2p
RUN go get -v github.com/matrix-org/dendrite/cmd/generate-keys@p2p

##### 

FROM alpine:latest AS final

RUN apk --update --no-cache add openssl

VOLUME /var/lib/postgresql/data
VOLUME /etc/dendrite

COPY --from=riot . /
COPY --from=postgres . /
COPY --from=golang /go/bin/* /usr/local/bin/

ENV PGDATA=/var/lib/postgresql/data

ADD docker/p2p-entrypoint.sh /usr/local/bin
ADD docker/postgres/create_db.sh /docker-entrypoint-initdb.d/create_db.sh
ADD dendrite.yaml /etc/dendrite/dendrite.yaml

RUN chmod +x /docker-entrypoint-initdb.d/create_db.sh
RUN sed -i '3i\ \ \ \ application/wasm wasm\;' /etc/nginx/mime.types
RUN adduser --system nginx
RUN addgroup --system nginx
RUN rm -rf /usr/share/nginx/html && ln -s /app /usr/share/nginx/html

CMD sh /usr/local/bin/p2p-entrypoint.sh

EXPOSE 80
EXPOSE 8008
EXPOSE 8448
