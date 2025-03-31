
APP=reddit-pulse
BIN="bin"

.PHONY: clean
clean:
	@rm -f $(BIN)/$(APP)

.PHONY: bin
bin:
	@mkdir -p $(BIN)

.PHONY: build
build: bin
	@go build -o $(BIN)/$(APP)

.PHONY: run
run: build
	@./$(BIN)/$(APP)

.PHONY: test
test:
	@go run test -v ./...
