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

in modules.go. Or just clone this repo and use [logspout-module](https://github.com/requilence/logspout-regexp/tree/master/logspout-module) dir in this repo.

## How to use
Use by setting a docker environment variable: `ROUTE_URIS=regexp://bot?file=regexps.txt&hide_matched_string=1"`

The default transport is `stderr`, but this adapter mainly created to work with Telegram. You should put the right bot's token(you've got from [@BotFather](https://t.me/BotFather) and chat's id. Here is an example:
`ROUTE_URIS=regexp+tg://bot?file=regexps.txt&throttle_seconds=600&hide_matched_string=1&chat=123&token=112233444:AAEfzA2_Q-FnUfasuib2_DAfdsn23jnK5s6QcQ"`

Full command to run docker container with `logspout-regexp`. Please note, that first you must build it locally: `docker build -t logspout-regexp .`
```bash
docker run --name="logspout" \
    --volume=/var/run/docker.sock:/var/run/docker.sock \
    --mount type=bind,source=$(pwd)/regexps.txt,target=/regexps.txt \
    -e "ROUTE_URIS=regexp+tg://bot?file=regexps.txt&throttle_seconds=600&hide_matched_string=1&chat=123&token=112233444:AAEfzA2_Q-FnUfasuib2_DAfdsn23jnK5s6QcQ" \
    logspout-regexp:latest
```

In your `regexps.txt` you can put regexps to match container logs:
```
(?i)([A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,24})
part_of_sensitive_info
```
