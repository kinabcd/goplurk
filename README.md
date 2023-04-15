goplurk
======
A golang wrapper of [Plurk API 2.0](https://www.plurk.com/API)

Getting started
----
With Go module support, simply add the following import
```
import "github.com/kinabcd/goplurk"
```

Basic usage
----
```golang
package main

import (
	"github.com/kinabcd/goplurk"
)

var consumerToken = "..."
var consumerSecret = "..."
var accessToken = "..."
var accessSecret = "..."

func main() {
	client, _ := goplurk.NewClient(consumerToken, consumerSecret, accessToken, accessSecret)
	client.Timeline.PlurkAdd("says", "somecontent")
}

```

Advenced usage
----
see **example/** for more usages

Author
------

Kin Lo :: kinabcd@gmail.com :: @kinabcd
