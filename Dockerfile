FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /timegobetter-api ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /timegobetter-api /timegobetter-api
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/timegobetter-api"]
