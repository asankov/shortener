FROM golang:alpine AS builder

RUN mkdir /app

ADD . /app/

WORKDIR /app

RUN go build -o shortener ./cmd/shortener

FROM gcr.io/distroless/static-debian11:latest

COPY --from=builder /app/shortener .
COPY --from=builder /app/internal/ui/template /internal/ui/template
COPY --from=builder /app/static static

ENTRYPOINT [ "./shortener" ]
