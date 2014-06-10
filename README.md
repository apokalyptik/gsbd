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

* websocket support
* see if avoiding calling out to C gets us better concurrency in github.com/rjohnsondev/go-safe-browsing-api

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
