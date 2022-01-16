FROM golang:1.17

WORKDIR /app

COPY go.mod .
COPY go.sum .

ARG opts
RUN env ${opts} go mod download

COPY . .


RUN env ${opts} go build ./cmd/main.go

WORKDIR /app
ENTRYPOINT ["./main"]