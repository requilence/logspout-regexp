# logspout-regexp

A minimalistic adapter for github.com/gliderlabs/logspout to match logs with regexp and notify in case match found
Currently supporting Telegram and StdErr as output
Follow the instructions in https://github.com/gliderlabs/logspout/tree/master/custom on how to build your own Logspout container with custom modules. Basically just copy the contents of the custom folder and include:

```go
package main

import (
  _ "github.com/requilence/logspout-regexp"
)
```

in modules.go.

Use by setting a docker environment variable `ROUTE_URIS=logstash://host:port` to the Logstash server.
The default protocol is UDP, but it is possible to change to TCP by adding ```+tcp``` after the logstash protocol when starting your container.

```bash
docker run --name="logspout" \
    --volume=/var/run/docker.sock:/var/run/docker.sock \
    -e "ROUTE_URIS=regexp+tg://bot?file=regexps.txt&throttle_seconds=600&hide_matched_string=1&chat=123&token=112233444:AAEfzA2_Q-FnUfasuib2_DAfdsn23jnK5s6QcQ" \
    logspout-regexp:latest
```

In your `logspout_regexps.txt` you can put the right regexps to match:
```
(?i)([A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,24})
part_of_sensitive_info
```
