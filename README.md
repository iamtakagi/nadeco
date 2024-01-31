# nadeco

## Configuration
```yaml
nameservers: ['8.8.8.8', '8.8.4.4']
records:
  - target: '127.0.0.1'
    values: ['local.example.com.']
    type: 'A'
    ttl: 3600
  - target: '127.0.0.1'
    values: ['www.local.example.com.']
    type: 'CNAME'
    ttl: 3600
  - target: '74.125.202.108'
    values: ['smtp.gmail.com.']
    type: 'MX'
    ttl: 3600
```

## Installation

### Run on local
```sh
go mod download
go run .
```

### Run on container
```sh
docker-compose up -d --build
```