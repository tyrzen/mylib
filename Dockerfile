FROM golang:1-alpine3.17 as builder
RUN  mkdir /myapp
ADD . /myapp
WORKDIR /myapp
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp ./cmd/main.go

FROM alpine:3.17 as product
COPY --from=builder /myapp .
ENTRYPOINT ["./myapp"]