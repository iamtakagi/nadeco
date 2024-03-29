FROM golang:1.20-alpine

WORKDIR /app
COPY go.sum go.mod ./
RUN go mod download 
COPY *.go .
RUN go build -o /nadeco

EXPOSE 53/udp

WORKDIR /
COPY config.yaml .
CMD ["/nadeco"]