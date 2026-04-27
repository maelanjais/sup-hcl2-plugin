.PHONY: all build-plugin build-sup clean install

# Nom du module (à adapter si besoin)
MODULE = github.com/maelan/sup-hcl2-plugin

# Binaires de sortie
PLUGIN_BIN = sup-hcl2-parser
SUP_BIN = sup

all: build-plugin

# Construire le binaire plugin
build-plugin:
	@echo "==> Building HCL2 parser plugin..."
	go build -o $(PLUGIN_BIN) ./plugin/
	@echo "==> Built: $(PLUGIN_BIN)"

# Installer le plugin dans le PATH
install: build-plugin
	@echo "==> Installing $(PLUGIN_BIN) to /usr/local/bin..."
	cp $(PLUGIN_BIN) /usr/local/bin/$(PLUGIN_BIN)
	@echo "==> Installed!"

# Nettoyer
clean:
	rm -f $(PLUGIN_BIN)

# Initialiser les dépendances
deps:
	go mod tidy

# Tester le plugin (parse un exemple et affiche le YAML)
test: build-plugin
	@echo "==> Testing HCL2 parser with example..."
	@echo "--- HCL2 input ---"
	@cat example/Supfile.hcl
	@echo ""
	@echo "--- Plugin output (YAML) ---"
	@go run ./test/main.go
