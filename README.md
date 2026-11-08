# sara

Very simple key value in memory storage, part of learn by project roadmap.

## Overview

Sara support base `redis` commands:

`PING`
If all fine, return PONG

`HELLO`
If all file, Sara say HI

`SET key value`
Set value of key

`SET key value EX ttl`
Set value for key with time to life. Janitor goroutine clean expired keys every 200ms

`GET key`
Return value of key, if key not exists return `(none)`

Get working with many keys, so it can be use like this
`GET key_1 key_2 key_3 ...`

`DEL key`
Delete key from storage

`EXISTS key`
Check key on exists

Commands `DEL` and `EXISTS` also work with many keys.

`DUMP`
Prints all keys and values

`BYE`/`QUIT` for close connection
