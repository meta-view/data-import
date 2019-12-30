FROM golang:alpine3.8
RUN apk --update add git build-base gcc openssh icu-dev upx && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

WORKDIR /app
COPY . /app

RUN GOOS=linux GOARCH=amd64 go build --tags "icu json1 fts5 secure_delete" -o bin/meta-view-service && /usr/bin/upx /app/bin/meta-view-service

FROM alpine:3.8
RUN mkdir /data
RUN apk --update add icu-libs && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

WORKDIR /
ENTRYPOINT ["/app/meta-view-service"]
ADD plugins /plugins
ADD templates /templates
COPY --from=0 /app/bin/meta-view-service /app/meta-view-service