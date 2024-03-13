FROM golang:1.22.1-alpine3.19 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# Compilar o aplicativo Go para um binário estático.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sonarDockServer .

FROM alpine:latest  

WORKDIR /root/

COPY --from=builder /app/sonarDockServer .

EXPOSE 8080

CMD ["./sonarDockServer"]
