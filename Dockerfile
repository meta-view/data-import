FROM golang:1.12-alpine3.11
RUN apk --update add git build-base gcc unzip bash openssh sed icu-dev upx && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

WORKDIR /app
COPY . /app

RUN chmod +x update.sh && bash ./update.sh
RUN GOOS=linux GOARCH=amd64 go build -ldflags "-X main.VersionString=$(git describe --always --dirty --tags)" --tags "icu json1 fts5 secure_delete" -o bin/meta-view-service && /usr/bin/upx /app/bin/meta-view-service

FROM alpine:3.11
RUN mkdir /data
RUN apk --update add icu-libs && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

WORKDIR /
ENTRYPOINT ["/app/meta-view-service"]
ADD plugins /plugins
ADD templates /templates
ADD static /static

COPY --from=0 /app/bin/meta-view-service /app/meta-view-service