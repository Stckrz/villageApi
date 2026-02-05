# build stage

FROM golang:1.25.3-bookworm AS build
WORKDIR /src


# Install build deps for CGO + sqlite
RUN apt-get update && apt-get install -y \
    gcc \
    sqlite3 \
    libsqlite3-dev \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o api ./cmd/api

#runtime stage
FROM gcr.io/distroless/cc-debian12

USER nonroot:nonroot

WORKDIR /app

# where SQLite will live
VOLUME ["/data"]

COPY --from=build /src/api /app/api

EXPOSE 8080
CMD ["/app/api"]

