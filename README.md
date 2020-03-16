# go-longsock
Keep a websocket connection alive for as long as possible

## How to build
```
go build
```

## How to use

### Server

Run locally
```
./go-longsock
```

Run on CloudFoundry
```
cf push
```

### Client

```
Usage: ./go-longsock [options] [server-url]
If server-url given, launch as client. Otherwise launch as server.
Options:
  -retry
    	Reconnect after disconnection (client only)
```

If you pass a server URL the program will automatically start in client mode.
Be aware that you need to use `ws://` or `wss://` for websockets instead of `http://` or `https://`
as protocol.


If you pass the `-retry=true` option the client will try to reconnect to the server
in case there are any errors and keep the connection alive.

### Examples

Connect to a local server without retry
```
./go-longsock ws://localhost:8080
```

Connect to a secure server on CloudFoundry with retry
```
./go-longsock -retry=true wss://go-longsock.my-cf.com/
```