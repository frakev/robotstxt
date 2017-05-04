# robotstxt
robots.txt parser

**Example:**

```
package main

import (
  rbt "github.com/frakev/robotstxt"
)

func main() {
  access, err := rbt.IsAllowed("http://www.google.com/search") // true or false, err if exists
 }
 ```
