FROM golang:1.25-alpine

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git curl

# Install Air
RUN GOPROXY=direct go install github.com/air-verse/air@latest

# Install postgres tools
RUN apk add --no-cache postgresql-client

COPY go.* ./
RUN go mod download

COPY . .

# Run Air by default
CMD ["air", "-c", ".air.toml"]
