BIN_DIR := bin
BINARY := $(BIN_DIR)/umami

.PHONY: build
build:
	mkdir -p $(BIN_DIR)
	go build -ldflags "-X github.com/yborunov/umami-cli/internal/cmd.version=$(VERSION) -X github.com/yborunov/umami-cli/internal/cmd.commit=$(COMMIT) -X github.com/yborunov/umami-cli/internal/cmd.date=$(DATE)" -o $(BINARY) ./cmd/umami

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
