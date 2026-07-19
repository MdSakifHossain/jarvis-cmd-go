APP_NAME=jarvis
BIN_DIR=bin
INSTALL_PATH=/usr/local/bin/$(APP_NAME)

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME)

install:
	make build
	sudo install -m 755 $(BIN_DIR)/$(APP_NAME) $(INSTALL_PATH)

clean:
	rm -rf $(BIN_DIR)