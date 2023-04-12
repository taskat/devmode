FROM golang:1.20 AS builder
WORKDIR /src/devmode
COPY . /src/devmode
RUN go build -o dev_server ./main.go

FROM alpine:3.17.3 AS final
WORKDIR /app
COPY --from=builder /src/devmode/dev_server /bin/dev_server