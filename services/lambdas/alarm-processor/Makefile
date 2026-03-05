OS_ENV_PREFIX := GOOS=linux GOARCH=amd64 CGO_ENABLED=0
ifeq ($(OS),Windows_NT)
OS_ENV_PREFIX := set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0&&
endif

build-AlarmProcessorFunction:
	$(OS_ENV_PREFIX) go build -o $(ARTIFACTS_DIR)/bootstrap ./cmd/handler
