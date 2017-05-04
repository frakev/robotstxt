# Robots.txt checker

- Pass an url as parameter
- Get and parse robots.txt url
- Check if user-agent is allowed to crawl the url pass as parameter
- Return true or false, error if exists (HTTP 403 Forbidden for example)

**Installation:**

```go
go get github.com/frakev/robotstxt
```
 
**Example:**

```go
package main

import (
  rbt "github.com/frakev/robotstxt"
)

func main() {
  access, err := rbt.IsAllowed("http://www.google.com/search") // true or false, err if exists
 }
```
 
**TO DO:**

- [ ] Pass user-agent as parameter
