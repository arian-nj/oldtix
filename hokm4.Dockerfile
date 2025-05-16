FROM golang:1.24-bookworm AS deps

WORKDIR /code

COPY ./backend/go.mod ./backend/go.sum ./

RUN go mod download

FROM golang:1.24-bookworm AS builder

WORKDIR /code
COPY --from=deps /go/pkg /go/pkg

COPY ./backend .

# RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./build/core_server ./cmd/core/.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./build/hokm4_server ./cmd/hokm4/.

FROM debian:bookworm-slim

WORKDIR /code


# Create a non-root user and group
RUN groupadd -r appuser && useradd -r -g appuser appuser

COPY --from=builder /code/build/hokm4_server .

# Change ownership of the application binary
RUN chown appuser:appuser ./hokm4_server

USER appuser

EXPOSE 4444

CMD [ "./hokm4_server" ]