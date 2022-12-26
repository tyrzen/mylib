FROM golang:bullseye as BUILDER
RUN  mkdir /app
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/main.go

FROM alpine:latest as PRODUCT
COPY --from=BUILDER /app .
CMD ["./app"]