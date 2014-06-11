gsbd
====

Google Safe Browsing Daemon

### Dependencies

These are not automatically satisfied by "go get"

* https://github.com/dcjones/hat-trie

### Installation

1. install dependencies
2. go get github.com/apokalyptik/gsbd
3. $GOPATH/bin/gsbd -h

### TODO

* see if avoiding calling out to C gets us better concurrency in github.com/rjohnsondev/go-safe-browsing-api
* batch API performance enhancements (possibly remove JSON for marshalling and use 1 line per request/response)

### Performance

It's not exactly highly optimize, but it's fairly performant as is. 
This simple test on my macbook pro (from localhost to localhost, 
and with all my normal things running, browsers, irc, etc) comes 
out pretty well.  I've not worked with it under any signifigant 
diversity of request urls though...

```
boom -n 1000000 -c 100 http://127.0.0.1:8888/safe/foo.com
1000000 / 1000000 Boooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo! 100.00 %

Summary:
  Total:	56.2107 secs.
  Slowest:	0.3203 secs.
  Fastest:	0.0003 secs.
  Average:	0.0056 secs.
  Requests/sec:	17790.1976

Status code distribution:
  [200]	1000000 responses

Response time histogram:
  0.000 [1]	|
  0.032 [996356]	|∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.064 [2112]	|
  0.096 [159]	|
  0.128 [125]	|
  0.160 [452]	|
  0.192 [203]	|
  0.224 [172]	|
  0.256 [156]	|
  0.288 [162]	|
  0.320 [102]	|

Latency distribution:
  10% in 0.0032 secs.
  25% in 0.0039 secs.
  50% in 0.0048 secs.
  75% in 0.0060 secs.
  90% in 0.0077 secs.
  95% in 0.0093 secs.
  99% in 0.0153 secs.
```

### Websockets API

#### /sock

Used to submit requests. Send a url as plain text, get a text response 
of "" (safe) or the list that the site is found in

### HTTP API

#### GET /

Used as a service health check. Always returns "ready"

#### GET /uptodate/ 

Whether the index is up to date or not. Returns "true" or "false"

#### GET /safe/{URL}

Checks the safe browsing database for {URL}.  Returns an empty string 
if the url is not listed, or a string representing the list that it 
was found in.

#### POST /batch

Batch processing. Request body should be a json encoded array 
of strings each of which is a url to check. Returns an array 
of strings matching the requested strings indexes to the response.
Just like /safe/{URL} it's "", or the list per response
