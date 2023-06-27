# RFC3986

RFC 3986 URI Query Escape/Unescape inspired from "net/url" written in Go

## Installation

Use go get.

    go get github.com/colduction/rfc3986
  
## Differences

* Unlike the `"net/url"` standard package, it percent-encodes space character with binary octet **`00100000` (ABNF: %x20)** too.
* Current package is limited to only two functions:  
    1. `QueryEscape`
    2. `QueryUnescape`