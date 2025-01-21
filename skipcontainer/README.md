# Test Container

The `skipcontainer` package provides a easy way to start up a Skip service. It has some helpers to add Skip file(s) which is required for the container to run.

```go
skipFile, err := os.Open("skip.ts")
// handle err
defer skipFile.Close()

// Create a new Skip container with a single Skip file.
skipContainer, err := skipcontainer.Run(ctx, "lidtop/goskip", skipcontainer.WithSkipFile(skipFile))
// handle err
defer skipContainer.Terminate(ctx)

// Create a new Skip container with multiple Skip files.
skipContainer, err := skipcontainer.Run(ctx, "lidtop/goskip", skipcontainer.WithSkipFiles(
    skipcontainer.WithFiles(
        skipcontainer.File{
            Reader:            skipFile,
            ContainerFilePath: "/app/skip.ts",
        },
        skipcontainer.File{
            Reader:            otherFile,
            ContainerFilePath: "/app/other.ts",
        },
    ),
))

// Or create a new Skip container with a directory of Skip files.
skipContainer, err := skipcontainer.Run(ctx, "lidtop/goskip", skipcontainer.WithSkipFiles(
    skipcontainer.WithDirectory("/path/to/skip/files"),
))
```

The `SkipContainer` also exposes helpers to create clients to connect to it.

```go
controlClient, err := skip.NewControlClient(skipContainer.GetControlURL())
streamClient, err := skip.NewStreamClient(skipContainer.GetStreamURL())
```
