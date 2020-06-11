<p align="center">
  <img src="img/logo.svg" alt="Logo botnet-ghost" width="200px">
</p>

# Run
⚠️ __Warning:__ This project are development
```go
go build main.go
./main.go <file template .html>
```

## C2 server
Server in GO with Sqlite and api, see endpoints:
* /some-string-random?gclid=base-64

* * For new client: `details of machine==ip address` (Response: `Command for run==id of db`)

* * For synchronize exists client: `id-database==response==status exited` (Response: `command for run==time for new request`)

## Response
The response command of server running within HTML in specific tag with template.
Ex:
```html
...
<img src='trump-idiot.jpg' class='<command-response>'/>
...
```

Encrypted: base64 [[See T1132 technique](https://attack.mitre.org/techniques/T1132/)]

Protocol: HTTP(s) [[See T1071 technique](https://attack.mitre.org/techniques/T1071/)]

* /victim (for attacker)

API for some frontend use
