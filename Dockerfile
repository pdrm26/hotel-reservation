FROM golang:1.24.4-alpine

WORKDIR /app
COPY go.mod go.sum ./

# Download and install required Go  dependecies
RUN go mod download

COPY . .

# Build the Go app
RUN go build -o main .

EXPOSE 3000

CMD ["./main"]