FROM golang:latest
WORKDIR /app
RUN go install github.com/air-verse/air@latest
COPY go.mod go.sum .env ./
RUN go mod download
COPY . .
CMD ["air", "-c", ".air.toml"]