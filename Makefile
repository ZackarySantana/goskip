mock:
	go run github.com/vektra/mockery/v2@v2.51.1

skip:
	@if [ -z "$(EXAMPLE)" ]; then \
		echo "Error: EXAMPLE is not set. Please provide it as an environment variable (e.g., make run EXAMPLE=groups)."; \
		exit 1; \
	fi
	docker build -t goskip-dev goskip-image
	docker run --name goskip -dv ./examples/$(EXAMPLE)/skip.ts:/app/skip.ts -p 8080:8080 -p 8081:8081 goskip-dev

skip-stop:
	docker rm -f goskip
