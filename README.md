gsbd
====

Google Safe Browsing Daemon

### Installation

1. go get github.com/apokalyptik/gsbd
2. $GOPATH/bin/gsbd -h

### TODO

* batch API performance enhancements (possibly remove JSON for marshalling and use 1 line per request/response)

### Performance

Needs testing, also "fast" depends on exactly how fast you need.

HTTP is slower than websockets, obviously. The service is easily load balanced.

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
