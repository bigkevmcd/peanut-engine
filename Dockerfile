FROM golang:latest AS build
MAINTAINER Kevin McDermott <bigkevmcd@gmail.com>
WORKDIR /go/src
COPY . /go/src
RUN CGO_ENABLED=0 GOOS=linux go build -a ./cmd/peanut-engine

FROM alpine
RUN apk add --update ca-certificates \
 && apk add --update -t deps bash git \
 && apk add --update bash git \
 && apk del --purge deps \
 && rm /var/cache/apk/*
WORKDIR /root/
COPY --from=build /go/src/peanut-engine .
EXPOSE 9001
ENTRYPOINT ["./peanut-engine"]
