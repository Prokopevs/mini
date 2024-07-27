FROM golang:1.22.3 as builder

WORKDIR /app

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY cmd cmd
COPY internal internal
COPY db db
COPY docs docs

COPY Makefile .

RUN make build

FROM debian:stable

WORKDIR /app
USER nobody

COPY --from=builder /app/server /app/server

ENTRYPOINT ["/app/server"]