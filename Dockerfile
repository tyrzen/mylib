FROM golang:1-alpine3.17 as builder
RUN  mkdir /mylib
ADD . /mylib
WORKDIR /mylib
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mylib ./cmd/main.go

FROM alpine:3.17 as product
COPY --from=builder /mylib .
ENTRYPOINT ["./mylib"]