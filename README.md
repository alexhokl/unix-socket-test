To run server

```sh
go run main.go
```

To connect using `curl`

```sh
curl -s -N --unix-socket /tmp/test.sock http://localhost/
curl -s -N --unix-socket /tmp/test.sock http://localhost/health
```
