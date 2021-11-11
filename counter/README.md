# counter

## start

- START 0x01 0x02 0x03

```
$ go run main.go
```

- GET 0x01

```
$ curl -s -v -XGET http://127.0.0.1:30001
0*   Trying 127.0.0.1:30001...
* Connected to 127.0.0.1 (127.0.0.1) port 30001 (#0)
> GET / HTTP/1.1
> Host: 127.0.0.1:30001
> User-Agent: curl/7.78.0
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Thu, 11 Nov 2021 15:25:47 GMT
< Content-Length: 1
< Content-Type: text/plain; charset=utf-8
<
{ [1 bytes data]
* Connection #0 to host 127.0.0.1 left intact
```

- PUT 0x02

```
$ curl -s -v -XPUT http://127.0.0.1:30002
*   Trying 127.0.0.1:30002...
* Connected to 127.0.0.1 (127.0.0.1) port 30002 (#0)
> PUT / HTTP/1.1
> Host: 127.0.0.1:30002
> User-Agent: curl/7.78.0
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 204 No Content
< Date: Thu, 11 Nov 2021 15:26:28 GMT
<
* Connection #0 to host 127.0.0.1 left intact
```

- GET 0x03

```
$ curl -s -v -XPUT http://127.0.0.1:30003
1*   Trying 127.0.0.1:30003...
* Connected to 127.0.0.1 (127.0.0.1) port 30003 (#0)
> GET / HTTP/1.1
> Host: 127.0.0.1:30003
> User-Agent: curl/7.78.0
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Thu, 11 Nov 2021 15:26:47 GMT
< Content-Length: 1
< Content-Type: text/plain; charset=utf-8
<
{ [1 bytes data]
* Connection #0 to host 127.0.0.1 left intact
```

## restart

- RESTART 0x01 0x02 0x03

```
$ go run main.go
```

- GET 0x02

```
$ curl -s -v -XGET http://127.0.0.1:30002
1*   Trying 127.0.0.1:30002...
* Connected to 127.0.0.1 (127.0.0.1) port 30002 (#0)
> GET / HTTP/1.1
> Host: 127.0.0.1:30002
> User-Agent: curl/7.78.0
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Thu, 11 Nov 2021 17:16:58 GMT
< Content-Length: 1
< Content-Type: text/plain; charset=utf-8
<
{ [1 bytes data]
* Connection #0 to host 127.0.0.1 left intact
```