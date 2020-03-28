# gomx

Http service to check if a domain has an MX record.

This is a research project, use at your own risk.

## Usage

```
$ ./gomx
Service listening on: 127.0.0.1:8000
```

### Options available

```
$ ./gomx --help
Usage of ./gomx:
  -address string
    	Address to bind to (default "127.0.0.1:8000")
```

## API

```
$ http --json -f POST http://127.0.0.1:8000/api/v1/check domain=iicloud.com

HTTP/1.1 200 OK
Content-Length: 16
Content-Type: text/plain; charset=utf-8
Date: Thu, 19 Mar 2020 11:16:16 GMT

{
    "status": false
}
```
