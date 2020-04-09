<img src="img/logo.svg" alt="Logo botnet-ghost">

# Run
> This project are development
```bash
go build main.go
./main.go <file template .html>
```

## C2 server
Server in GO with Sqlite and api, see endpoints:
* /some-string-random?gclid=base-64
For client sent in base64: `id\response`
* /auth
API for some frontend use
