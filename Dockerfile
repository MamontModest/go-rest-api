FROM golang:alpine
WORKDIR /build
ADD go.mod .
COPY . .
RUN go build -o main src/server/main.go
CMD ["./main"]