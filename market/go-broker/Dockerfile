FROM golang:alpine

WORKDIR /app

COPY . /app

RUN apk add --no-cache alpine-sdk linux-headers
RUN go build

ENTRYPOINT ["./go-broker"]
