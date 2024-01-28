# nadeco

## Configuration
```yaml
nameservers: ['8.8.8.8', '8.8.4.4']
records:
  - target: '127.0.0.1'
    values: ['local.example.dev.']
    type: 'A'
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