A DDP Server Written in Go
==========================

# Example Usage
```
package main

import (
    "github.com/krishamoud/golang-ddp-server"
)


func main() {
  s := server.New()
  s.Method("hello", handler)
  s.Listen(":8080")
}

func handler(ctx server.MethodContext) {

  ctx.SendResult("Hello, world!")

  ctx.SendUpdated()
}
```

**todo**

1. tests
2. sockJs support
