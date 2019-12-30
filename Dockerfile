FROM golang:alpine3.8
RUN apk --update add git build-base gcc openssh upx && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

WORKDIR /app
COPY . /app

RUN GOOS=linux GOARCH=amd64 go build -o bin/meta-view-service && /usr/bin/upx /app/bin/meta-view-service

FROM alpine:3.8
RUN mkdir /data
RUN apk --update add git openssh && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

WORKDIR /
ENTRYPOINT ["/app/meta-view-service"]
ADD plugins /app/
ADD templates /app/
COPY --from=0 /app/bin/meta-view-service /app/meta-view-service