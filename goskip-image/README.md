# lidtop/goskip

This image is a simple Skip application that requires a `skip.ts` file in the `/app` directory.

An example usage of running this image is:

```bash
docker run -v ./examples/$(EXAMPLE)/skip.ts:/app/skip.ts -p 8080:8080 -p 8081:8081 lidtop/goskip
```

## Docker Hub

This is published to Docker Hub as [lidtop/goskip](https://hub.docker.com/repository/docker/lidtop/goskip).
