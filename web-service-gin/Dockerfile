FROM golang:alpine

RUN apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o binary

ENTRYPOINT ["/app/binary"]

