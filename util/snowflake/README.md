# Introduction

A wrapper for the package <https://github.com/bwmarrin/snowflake>

## Usage

Just import the package into the file to use it.

**Example Program:**

```go
package foo

import (
 "fmt"

 "go.tekoapis.com/tekone/library/util/snowflake"
)

func doSomething() {

    // generate snowflake ID
    id, err := snowflake.NextID()
    if err != nil {
        // handle err
    }

    // use generated id to do something
    ....
}
```
