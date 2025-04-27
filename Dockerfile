FROM golang:latest AS dev
WORKDIR /app
RUN go install github.com/air-verse/air@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["air", "-c", ".air.toml"]

FROM golang:latest AS prod
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o sinartimur-app ./cmd/app
EXPOSE 8080
ENTRYPOINT ["/app/sinartimur-app"]
