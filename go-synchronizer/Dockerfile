FROM golang:alpine

WORKDIR /app

COPY . /app

RUN apk add --no-cache alpine-sdk
RUN go build

CMD ["./go-synchronizer"]
