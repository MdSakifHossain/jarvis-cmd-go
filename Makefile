APP_NAME = jarvis
COMPLETION_FILE = _jarvis

BIN_DIR = bin
GENERATED_DIR = generated

INSTALL_PATH = /usr/local/bin
ZSH_COMPLETION_PATH = $(HOME)/.oh-my-zsh/completions

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME)

completion:
	mkdir -p $(GENERATED_DIR)
	go run ./tools/generate-completions/main.go ./jarvis-schema.json
	mv $(COMPLETION_FILE) $(GENERATED_DIR)

clean:
	rm -rf $(BIN_DIR)
	rm -rf $(GENERATED_DIR)

install: build completion
	sudo install -m 755 $(BIN_DIR)/$(APP_NAME) $(INSTALL_PATH)/$(APP_NAME)
	mkdir -p $(ZSH_COMPLETION_PATH)
	cp $(GENERATED_DIR)/$(COMPLETION_FILE) $(ZSH_COMPLETION_PATH)
	@echo "\n\nInstallation complete.\nRun 'exec zsh' to reload completions."

uninstall:
	sudo rm -f $(INSTALL_PATH)/$(APP_NAME)
	rm -f $(ZSH_COMPLETION_PATH)/$(COMPLETION_FILE)
	@echo "\n\nJarvis has been uninstalled.\nRun 'exec zsh' to take effect."