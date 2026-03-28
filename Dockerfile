FROM golang:1.25-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /grimoire ./cmd/grimoire

FROM gcr.io/distroless/static-debian12:latest
WORKDIR /app
COPY --from=build /grimoire /app/grimoire
ENV GRIMOIRE_DB_PATH=/data/grimoire.db
ENTRYPOINT ["/app/grimoire"]
