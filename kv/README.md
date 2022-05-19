# kv

## start

- START 0x01 0x02 0x03

```
$ go run main.go
```

- GET 0x01 key

```
$ curl -s -v -XGET http://127.0.0.1:30001/key
*   Trying 127.0.0.1:30001...
* Connected to 127.0.0.1 (127.0.0.1) port 30001 (#0)
> GET /key HTTP/1.1
> Host: 127.0.0.1:30001
> User-Agent: curl/7.78.0
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 404 Not Found
< Content-Type: text/plain; charset=utf-8
< X-Content-Type-Options: nosniff
< Date: Tue, 23 Nov 2021 12:39:11 GMT
< Content-Length: 12
<
Fail to GET
* Connection #0 to host 127.0.0.1 left intact
```

- PUT 0x02 key value

```
$ curl -s -v -XPUT http://127.0.0.1:30002/key -d value
*   Trying 127.0.0.1:30002...
* Connected to 127.0.0.1 (127.0.0.1) port 30002 (#0)
> PUT /key HTTP/1.1
> Host: 127.0.0.1:30002
> User-Agent: curl/7.78.0
> Accept: */*
> Content-Length: 5
> Content-Type: application/x-www-form-urlencoded
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 204 No Content
< Date: Tue, 23 Nov 2021 12:42:53 GMT
<
* Connection #0 to host 127.0.0.1 left intact
```

- GET 0x03 key

```
$ curl -s -v -XGET http://127.0.0.1:30003/key
*   Trying 127.0.0.1:30003...
* Connected to 127.0.0.1 (127.0.0.1) port 30003 (#0)
> GET /key HTTP/1.1
> Host: 127.0.0.1:30003
> User-Agent: curl/7.78.0
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Tue, 23 Nov 2021 12:45:39 GMT
< Content-Length: 5
< Content-Type: text/plain; charset=utf-8
<
value* Connection #0 to host 127.0.0.1 left intact
```

## restart

- RESTART 0x01 0x02 0x03

```
$ go run main.go
```

- GET 0x02 key # wait leader election

```
$ curl -s -v -XGET http://127.0.0.1:30002/key
*   Trying 127.0.0.1:30002...
* Connected to 127.0.0.1 (127.0.0.1) port 30002 (#0)
> GET /key HTTP/1.1
> Host: 127.0.0.1:30002
> User-Agent: curl/7.78.0
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Tue, 23 Nov 2021 12:47:23 GMT
< Content-Length: 5
< Content-Type: text/plain; charset=utf-8
<
value* Connection #0 to host 127.0.0.1 left intact
```