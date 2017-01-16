## rrframework

A collection of modules to make backend programming easier.

### Modules
#### config
configuration file parser, supporting formats:
* json


```go
package main
import (
        "github.com/songtianyi/rrframework/config"
	"fmt"
)

fun main() {
	rc, err := rrconfig.LoadJsonConfigFromFile("config.json")
	if err != nil {
		panic(err)
	}
	v, err := rc.GetStringSlice("files.ufile")
	if err != nil {
		panic(err)
	}
	fmt.Println(v)
	
}
```

#### connector

	
