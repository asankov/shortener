FROM golang:alpine AS builder

RUN mkdir /app

ADD . /app/

WORKDIR /app

RUN go build -o shortener ./cmd/shortener

FROM gcr.io/distroless/static-debian11:latest

COPY --from=builder /app/shortener .

ENTRYPOINT [ "./shortener" ]
