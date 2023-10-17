FROM golang:1.21.2

COPY ./ /app

RUN export GOPATH=/app

WORKDIR /app

RUN go mod download

RUN go build cmd/api/main.go

ENTRYPOINT [ "./main" ]