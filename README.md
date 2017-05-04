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
