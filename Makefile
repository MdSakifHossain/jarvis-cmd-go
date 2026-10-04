APP_NAME = jarvis

BIN_DIR = bin
INSTALL_PATH = /usr/local/bin

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME)

clean:
	rm -rf $(BIN_DIR)

install: build
	sudo install -m 755 $(BIN_DIR)/$(APP_NAME) $(INSTALL_PATH)/$(APP_NAME)
	@echo "Installation complete."
	@echo 
	@echo "Run 'jarvis --help' to get started."

uninstall:
	sudo rm -f $(INSTALL_PATH)/$(APP_NAME)
	@echo
	@echo "Jarvis has been uninstalled."
