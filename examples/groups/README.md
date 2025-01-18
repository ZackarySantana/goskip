# Groups Example

This example copies Skip's groups and adapts it to Go.

## Skip's Example

-   [groups.ts](https://github.com/SkipLabs/skip/blob/main/skipruntime-ts/examples/groups.ts) defines the `skip` service.
-   [groups-server.ts](https://github.com/SkipLabs/skip/blob/main/skipruntime-ts/examples/groups-server.ts) is a layer over the `skip` service.
-   [groups-client.ts](https://github.com/SkipLabs/skip/blob/main/skipruntime-ts/examples/groups-client.ts) queries the server to change and read the state.

## goskip's Example

-   [skip](./skip.ts) defines the `skip` service.
-   [main.go](./main.go) queries the skip service to change and read the state.

### Running the Example

Before running, make sure you have the dependencies installed:

```bash
cd examples && bun install
```

Run two terminals, one with:

```bash
bun run examples/groups/skip.ts
```

And the other with:

```bash
go run examples/groups/main.go
```

You should see the output like:

```bash
~/goskip$ go run examples/groups/main.go
Received Event: init, Data: [{Key:1001 Values:[1]}]
Setting Carol to active
Received Event: update, Data: [{Key:1001 Values:[1 2]} {Key:1002 Values:[2]}]
Setting Alice to inactive
Received Event: update, Data: [{Key:1001 Values:[2]}]
Setting Eve as Bob's friend
Received Event: update, Data: [{Key:1001 Values:[2 3]} {Key:1002 Values:[2]}]
Removing Carol and adding Eve to group 2
Received Event: update, Data: [{Key:1002 Values:[3]}]
```

Running the `main.go` file again will result in no updates, as the objects are already in that state.

```bash
~/goskip$ go run examples/groups/main.go
Received Event: init, Data: [{Key:1001 Values:[2 3]} {Key:1002 Values:[3]} {Key:1002 Values:[3]}]
Setting Carol to active
Setting Alice to inactive
Setting Eve as Bob's friend
Removing Carol and adding Eve to group 2
```
